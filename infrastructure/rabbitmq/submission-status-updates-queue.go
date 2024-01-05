package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/upb-code-labs/tests-microservice/domain/dtos"
)

type SubmissionStatusUpdatesQueueMgr struct {
	Queue *amqp.Queue
}

// submissionStatusUpdatesQueueMgr Singleton struct
var submissionStatusUpdatesQueueMgr *SubmissionStatusUpdatesQueueMgr

// GetSubmissionStatusUpdatesQueueMgrInstance returns a pointer to the submission status updates queue
func GetSubmissionStatusUpdatesQueueMgrInstance() *SubmissionStatusUpdatesQueueMgr {
	if submissionStatusUpdatesQueueMgr == nil {
		submissionStatusUpdatesQueueMgr = &SubmissionStatusUpdatesQueueMgr{
			Queue: getSubmissionsStatusUpdatesQueue(),
		}
	}

	return submissionStatusUpdatesQueueMgr
}

// getSubmissionsStatusUpdatesQueue returns a pointer to the submission status updates queue
func getSubmissionsStatusUpdatesQueue() *amqp.Queue {
	if submissionStatusUpdatesQueueMgr == nil {
		ch := GetRabbitMQChannel()

		// Declare queue
		qName := "submission-status-updates"
		qDurable := true
		qAutoDelete := false
		qExclusive := false
		qNoWait := false
		qArgs := amqp.Table{}

		q, err := ch.QueueDeclare(
			qName,
			qDurable,
			qAutoDelete,
			qExclusive,
			qNoWait,
			qArgs,
		)

		if err != nil {
			log.Fatal(err.Error())
		}

		// Set fair dispatch
		maxPrefetchCount := 8 // Limit to 8 updates per worker
		err = ch.Qos(
			maxPrefetchCount,
			0,
			false,
		)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("RabbitMQ: Submission status updates queue declared")
		return &q
	}

	return submissionStatusUpdatesQueueMgr.Queue
}

// QueueUpdate queues a submission status update
func (qMgr *SubmissionStatusUpdatesQueueMgr) QueueUpdate(updateDTO *dtos.SubmissionStatusUpdateDTO) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := GetRabbitMQChannel()

	// Parse update to JSON
	json, err := json.Marshal(updateDTO)
	if err != nil {
		return err
	}

	// Publish update to queue
	msgExchange := ""
	msgRoutingKey := qMgr.Queue.Name
	msgMandatory := false
	msgImmediate := false

	err = ch.PublishWithContext(
		ctx,
		msgExchange,
		msgRoutingKey,
		msgMandatory,
		msgImmediate,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        json,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
