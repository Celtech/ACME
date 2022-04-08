package main

import (
	"baker-acme/internal/context"
	"baker-acme/web"
	con "context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"fmt"
)

var appContext *context.AppContext

func init() {
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
		appContext.Logger.Info(fmt.Sprintf("caught sig: %+v", sig))

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
	appContext.Logger.Info("attempting graceful shutdown")

	srv.Shutdown(con.Background())
	httpServerExitDone.Wait()

	appContext.Logger.Info("graceful shutdown complete")

	done <- struct{}{}
	os.Exit(<-returnCode)
}
