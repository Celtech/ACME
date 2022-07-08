package cmd

import (
	"github.com/Celtech/ACME/internal/queue"
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

		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)
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
				// timeout after 30 seconds waiting for app to finish,
				// our application should Exit(1)
				returnCode <- 1
			case <-done:
				// if we got a message on done, we finished, so end app
				// our application should Exit(0)
				returnCode <- 0
			}
		}()

		srv := web.Serve(config.GetConfig())
		queue.QueueMgr.Subscribe()

		<-finishUP
		log.Info("attempting graceful shutdown")

		// 1 second less than force shutdown time
		ctx, cancel := ctx.WithTimeout(ctx.Background(), 29*time.Second)
		defer cancel()
		srv.Shutdown(ctx)

		log.Info("graceful shutdown complete")

		done <- struct{}{}
		os.Exit(<-returnCode)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
