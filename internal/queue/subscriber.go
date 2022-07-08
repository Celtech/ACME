package queue

import (
	"encoding/json"
	"github.com/Celtech/ACME/internal/acme"
	"github.com/Celtech/ACME/web/model"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const SEQUENTIAL_WAIT_TIME = 30 // in seconds

func (q *QueueManager) Subscribe() {
	for {
		evt, err := q.extractEventFromQueue()

		if err != nil {
			log.Error(err.Error())
		} else {
			log.Infof("Working on queue event for request id %d of type %s attempt %d",
				evt.RequestId,
				evt.ChallengeType,
				evt.Attempt,
			)

			if err := acme.Run(evt.Domain, evt.ChallengeType); err != nil {
				handleCertificateError(evt, err)
			} else {
				updateRequest(evt, model.STATUS_ISSUED)
			}
		}

		time.Sleep(SEQUENTIAL_WAIT_TIME * time.Second)
	}
}

func (q *QueueManager) extractEventFromQueue() (QueueEvent, error) {
	result, err := q.client.BLPop(q.ctx, 0*time.Second, q.queue).Result()
	if err != nil {
		return QueueEvent{}, err
	}

	evt := QueueEvent{}
	err = json.Unmarshal([]byte(result[1]), &evt)
	if err != nil {
		return QueueEvent{}, err
	}

	return evt, nil
}

func handleCertificateError(params QueueEvent, err error) {
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
