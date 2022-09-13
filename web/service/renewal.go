package service

import (
	"github.com/Celtech/ACME/internal/queue"
	"github.com/Celtech/ACME/web/model"
	log "github.com/sirupsen/logrus"
	"time"
)

func ProcessRenewals() {
	for {
		requestModel := new(model.Request)
		res, _ := requestModel.GetAllExpiringSoon()

		if len(res) > 0 {
			log.Infof("Found %d certificates to renew, attempting renewal", len(res))
		} else {
			log.Infof("No certificates need renewal, next check at %s", time.Now().AddDate(0, 0, 1))
		}

		for _, req := range res {
			requestModel := req

			evt := queue.QueueEvent{
				RequestId:     requestModel.Id,
				Domain:        requestModel.Domain,
				ChallengeType: requestModel.ChallengeType,
				Type:          queue.EVENT_RENEW,
				Attempt:       1,
				CreatedAt:     time.Now(),
			}

			if err := queue.QueueMgr.Publish(evt); err != nil {
				log.Errorf("error publishing certificate renewal request for domain %s to queue, %v", requestModel.Domain, err)
			} else {
				log.Infof("Sent certificate renewal request for %s to queue", requestModel.Domain)
			}
		}

		time.Sleep(time.Hour * 24)
	}
}
