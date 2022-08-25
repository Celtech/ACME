package util

import (
	"context"
	log "github.com/sirupsen/logrus"
	"reflect"
	"runtime"
	"time"
)

type Effector func(context.Context) error

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) error {
		for r := 0; ; r++ {
			err := effector(ctx)
			if err == nil || r >= retries {
				// Return when there is no error or the maximum amount
				// of retries is reached.
				return err
			}

			log.Warnf("Function call %s failed, retrying in %v", GetFunctionName(effector), delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}
