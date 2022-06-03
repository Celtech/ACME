package main

import (
	"baker-acme/config"
	"baker-acme/internal/queue"
	"baker-acme/web"
	"baker-acme/web/service/database"
	ctx "context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var conf *viper.Viper

func init() {
	conf = config.Init(os.Getenv("ACME_ENV"))

	log.SetFormatter(&log.TextFormatter{
		DisableColors: !conf.GetBool("services.logger.color"),
		FullTimestamp: true,
	})

	database.Init()
	queue.QueueMgr = queue.NewQueue(conf.GetString("redis.name"))
}

// @title           ACME API
// @version         1.0
// @description     Issue dynamic SSL certificates for 1+n domains.
// @description.markdown
// @termsOfService  https://chargeover.com/terms-of-service

// @contact.name   API Support
// @contact.url    https://git.rykelabs.com/rykelabs/acme-server/-/issues/new
// @contact.email  tim@chargeover.com

// @license.name  MIT
// @license.url   https://git.rykelabs.com/rykelabs/acme-server/-/blob/main/LICENSE.md

// @host      acme.chargeover.com:9022
// @BasePath  /api/v1

// @accept	json
// @produce json

// @schemes https

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description					Description for what is this security definition being used

func main() {
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

	srv := web.Serve(conf)
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
}
