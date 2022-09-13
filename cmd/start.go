package cmd

import (
	"github.com/Celtech/ACME/internal/queue"
	"github.com/Celtech/ACME/web/database"
	"github.com/Celtech/ACME/web/database/migration"
	"github.com/Celtech/ACME/web/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	ctx "context"
	"github.com/Celtech/ACME/config"
	"github.com/Celtech/ACME/web"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the web server",
	Long:  "Start the web server",
	Run: func(cmd *cobra.Command, args []string) {
		var returnCode = make(chan int)
		var finishUP = make(chan struct{})
		var done = make(chan struct{})
		var gracefulStop = make(chan os.Signal)

		signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)
		go func() {
			// wait for our os signal to stop the app
			// on the graceful stop channel
			// this goroutine will block until we get an OS signal
			sig := <-gracefulStop
			log.Infof("caught sig: %+v", sig)

			// send message on "finish up" channel to tell the app to
			// gracefully shutdown
			finishUP <- struct{}{}

			// wait for word back if we finished or not
			select {
			case <-time.After(30 * time.Second):
				log.Info("graceful shutdown timed out, forcefully exiting")

				// timeout after 30 seconds waiting for app to finish,
				// our application should Exit(1)
				returnCode <- 1
			case <-done:
				log.Info("graceful shutdown complete")

				// if we got a message on done, we finished, so end app
				// our application should Exit(0)
				returnCode <- 0
			}
		}()

		log.Infof("SSL Certify v%s built %s", Version, Date)

		con := database.Init()
		migration.RunMigrations()
		queue.QueueMgr = queue.NewQueue(config.GetConfig().GetString("redis.name"))

		srv := web.Serve(config.GetConfig())
		go queue.QueueMgr.Subscribe()
		go service.ProcessRenewals()

		<-finishUP
		log.Info("attempting graceful shutdown")

		// 1 second less than force shutdown time
		ctx, cancel := ctx.WithTimeout(ctx.Background(), 29*time.Second)
		defer cancel()

		// Shut down queue
		log.Info("gracefully closing redis")
		err := queue.QueueMgr.Close()
		if err != nil {
			log.Fatalf("graceful shutdown of redis failed: %+v", err)
		}

		// Shut down database
		log.Info("gracefully closing database connection")
		if db, err := con.DB(); err == nil {
			// We can't log here so if the close command errors, we can't really
			// do anything about it. We have to let it exit forcefully.
			err := db.Close()
			if err != nil {
				log.Fatalf("graceful shutdown of mariadb failed: %+v", err)
			}
		}

		// Shut down web server
		log.Info("gracefully closing web server")
		err = srv.Shutdown(ctx)
		if err != nil {
			log.Fatalf("graceful shutdown of web server failed: %+v", err)
		}

		done <- struct{}{}
		os.Exit(<-returnCode)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
