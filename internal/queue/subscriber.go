package queue

import (
	"baker-acme/internal/acme"
	"baker-acme/web/model"
	"encoding/json"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const SEQUENTIAL_WAIT_TIME = 30 // in seconds

func (q *QueueManager) Subscribe() error {
	for {
		result, err := q.client.BLPop(q.ctx, 0*time.Second, q.queue).Result()

		if err != nil {
			log.Error(err.Error())
		} else {
			log.Infof("working on %v", result)

			params := QueueEvent{}
			err := json.Unmarshal([]byte(result[1]), &params)
			if err != nil {
				log.Error(err.Error())
			} else {
				if err := acme.Run(params.Domain, params.ChallengeType); err != nil {
					if params.Attempt >= 3 {
						log.Errorf("error issuing certificate for %s on attempt %d. Max attempts reached, marking as failed.\r\n%v", params.Domain, params.Attempt, err)
						updateRequest(params, model.STATUS_ERROR)
					} else {
						log.Errorf("error issuing certificate for %s on attempt %d of 3. Re-queueing.\r\n%v", params.Domain, params.Attempt, err)
						params.Attempt++
						if err := QueueMgr.Publish(params); err != nil {
							log.Errorf("error publishing certificate request for domain %s to queue, %v", params.Domain, err)
						}
					}
				} else {
					updateRequest(params, model.STATUS_ISSUED)
				}
			}
		}

		time.Sleep(SEQUENTIAL_WAIT_TIME * time.Second)
	}
}

func updateRequest(params QueueEvent, status string) {
	requestId := params.RequestId
	var requestModel = new(model.Request)
	err := requestModel.GetByID(strconv.Itoa(requestId))
	if err != nil {
		log.Errorf("error fetching request %d\r\n%v", requestId, err)
	} else {
		requestModel.Status = status
		err := requestModel.Update()
		if err != nil {
			log.Errorf("error updating request %d to status %s\r\n%v", requestId, status, err)
		}
	}
}
