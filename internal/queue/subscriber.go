package queue

import (
	"baker-acme/internal/acme"
	"encoding/json"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (q *QueueManager) Subscribe() error {
	for {
		result, err := q.client.BLPop(q.ctx, 0*time.Second, q.queue).Result()

		if err != nil {
			log.Error(err.Error())
		} else {
			log.Infof("working on %v", result)

			params := map[string]interface{}{}

			err := json.NewDecoder(strings.NewReader(string(result[1]))).Decode(&params)

			if err != nil {
				log.Error(err.Error())
			} else {
				if err := acme.Run(params["Domain"].(string), params["ChallengeType"].(string)); err != nil {
					log.Errorf("error issuing certificate for %s\r\n%v", params["Domain"].(string), err)
					// TODO: Requeue
				} else {
					// TODO: Update database
				}
			}
		}

		time.Sleep(30 * time.Second)
	}
}
