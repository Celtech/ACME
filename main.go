package main

import (
	"baker-acme/internal/context"
	"baker-acme/web"
	con "context"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var appContext *context.AppContext

func init() {
	var colorLogOutput bool
	if colorLogs := os.Getenv("BAKER_COLOR_LOGS"); len(colorLogs) == 0 {
		colorLogOutput = true
	} else {
		colorLogOutput, _ = strconv.ParseBool(colorLogs)
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: !colorLogOutput,
		FullTimestamp: true,
	})

	appContext = context.NewAppContext(nil)
}

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

	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv := web.StartServer(appContext, httpServerExitDone)

	<-finishUP
	log.Info("attempting graceful shutdown")

	srv.Shutdown(con.Background())
	httpServerExitDone.Wait()

	log.Info("graceful shutdown complete")

	done <- struct{}{}
	os.Exit(<-returnCode)
}
