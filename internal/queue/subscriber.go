package queue

import (
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

			// params := map[string]interface{}{}

			// err := json.NewDecoder(strings.NewReader(string(result[1]))).Decode(&params)

			if err != nil {
				log.Error(err.Error())
			} else {
				//log.Info(result[1])
			}
		}

		time.Sleep(10 * time.Second)
	}
}
