package queue

import (
	"context"
	"encoding/json"
	"github.com/Celtech/ACME/config"
	"github.com/Celtech/ACME/internal/acme"
	"github.com/Celtech/ACME/internal/util"
	"github.com/Celtech/ACME/web/model"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const SEQUENTIAL_WAIT_TIME = 60 // in seconds

var certificateRequestRetryCounter *prometheus.CounterVec
var certificateRequestFailureCounter *prometheus.CounterVec
var certificateRequestIssuedCounter *prometheus.CounterVec

func (q *QueueManager) Subscribe() {
	setupMetrics()

	for {
		evt, err := q.extractEventFromQueue()

		if err != nil {
			log.Error(err.Error())
		} else {
			log.Infof("Working on %s queue event for request id %d of type %s attempt %d",
				evt.Type,
				evt.RequestId,
				evt.ChallengeType,
				evt.Attempt,
			)

			var err error = nil
			if evt.Type == EVENT_ISSUE {
				err = acme.Run(evt.Domain, evt.ChallengeType)
			} else if evt.Type == EVENT_RENEW {
				err = acme.Renew([]string{evt.Domain}, evt.ChallengeType)
			}

			if err != nil {
				handleCertificateError(evt, err)
			} else {
				issuedOn := time.Now()
				renews := issuedOn.Add(90 * (time.Hour * 24))
				certificateRequestIssuedCounter.WithLabelValues(evt.Domain, issuedOn.Format(time.RFC3339), renews.Format(time.RFC3339)).Inc()

				updateRequest(evt, model.STATUS_ISSUED)
				processPlugins(evt.Domain)
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
	rateLimit := config.GetConfig().GetInt("acme.retryLimit")
	if params.Attempt >= rateLimit {
		certificateRequestFailureCounter.WithLabelValues(params.Domain).Inc()

		log.Errorf("error issuing certificate for %s on attempt %d of %d. Max attempts reached, marking as failed.\r\n%v", params.Domain, params.Attempt, rateLimit, err)
		updateRequest(params, model.STATUS_ERROR)
	} else {
		certificateRequestRetryCounter.WithLabelValues(params.Domain).Inc()

		log.Errorf("error issuing certificate for %s on attempt %d of %d. Re-queueing.\r\n%v", params.Domain, params.Attempt, rateLimit, err)
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

func processPlugins(domain string) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "domain", domain)
	r := util.Retry(acme.RunPlugins, 5, 5*time.Second)
	err := r(ctx)

	if err != nil {
		log.Errorf("There was a executing plugins: %v", err)
	}
}

func setupMetrics() {
	certificateRequestRetryCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssl_certify_certificate_request_retry_count", // metric name
			Help: "Count of number of retried certificate requests.",
		},
		[]string{"domain"}, // labels
	)

	certificateRequestFailureCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssl_certify_certificate_request_fail_count", // metric name
			Help: "Count of number of failed certificate requests, this does not include retries.",
		},
		[]string{"domain"}, // labels
	)

	certificateRequestIssuedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ssl_certify_certificate_request_issued_count", // metric name
			Help: "Count of number of successfully issued certificate requests.",
		},
		[]string{"domain", "issued", "renews"}, // labels
	)

	prometheus.MustRegister(certificateRequestRetryCounter)
	prometheus.MustRegister(certificateRequestFailureCounter)
	prometheus.MustRegister(certificateRequestIssuedCounter)
}
