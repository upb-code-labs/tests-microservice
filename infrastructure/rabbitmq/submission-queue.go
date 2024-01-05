package rabbitmq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/upb-code-labs/tests-microservice/application"
	"github.com/upb-code-labs/tests-microservice/domain/dtos"
	"github.com/upb-code-labs/tests-microservice/domain/entities"
	"github.com/upb-code-labs/tests-microservice/utils"
)

type SubmissionQueueMgr struct {
	Queue          *amqp.Queue
	MessageChannel <-chan amqp.Delivery
	UseCases       *application.SubmissionsUseCases
}

// submissionQueueManagerInstance Singleton struct
var submissionQueueManagerInstance *SubmissionQueueMgr

// GetRabbitMQSubmissionsQueue returns a pointer to the submissions queue
func GetRabbitMQSubmissionsQueue() *amqp.Queue {
	if submissionQueueManagerInstance == nil {
		ch := GetRabbitMQChannel()

		// Declare queue
		qName := "submissions"
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
		maxPrefetchCount := 4 // Limit to 4 submissions per worker
		err = ch.Qos(
			maxPrefetchCount,
			0,
			false,
		)

		if err != nil {
			log.Fatal(err.Error())
		}

		// Set queue
		log.Println("RabbitMQ: Submissions queue declared")
		return &q
	}

	return submissionQueueManagerInstance.Queue
}

// GetSubmissionQueueMgr returns a pointer to the singleton instance of SubmissionQueueMgr
func GetSubmissionQueueMgr() *SubmissionQueueMgr {
	if submissionQueueManagerInstance == nil {
		submissionQueueManagerInstance = &SubmissionQueueMgr{
			Queue:    GetRabbitMQSubmissionsQueue(),
			UseCases: &application.SubmissionsUseCases{},
		}
	}

	return submissionQueueManagerInstance
}

// ListenForSubmissions starts listening for submissions
func (manager *SubmissionQueueMgr) ListenForSubmissions() {
	ch := GetRabbitMQChannel()

	// Set consumer
	qName := manager.Queue.Name
	qConsumer := ""   // DEFAULT value so the server will generate a unique name
	qAutoAck := false // Manual ack to implement fair dispatch
	qExclusive := false
	qNoLocal := false
	qNoWait := false
	qArgs := amqp.Table{}

	msgs, err := ch.Consume(
		qName,
		qConsumer,
		qAutoAck,
		qExclusive,
		qNoLocal,
		qNoWait,
		qArgs,
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	// Set message channel
	manager.MessageChannel = msgs

	// Start processing submissions
	manager.processSubmissions()
}

// processSubmissions infinite loop to process received submissions
func (manager *SubmissionQueueMgr) processSubmissions() {
	for msg := range manager.MessageChannel {
		go manager.processSubmission(msg)
	}
}

func (manager *SubmissionQueueMgr) processSubmission(msg amqp.Delivery) {
	// ACK message after processing
	defer msg.Ack(false)

	// Unmarshal message
	var submissionWork entities.SubmissionWork
	err := json.Unmarshal(msg.Body, &submissionWork)
	if err != nil {
		log.Println("[ERROR] Failed to unmarshal submission work", err.Error())
		return
	}

	// Send submission status update
	statusDTO := &dtos.SubmissionStatusUpdateDTO{
		SubmissionUUID:   submissionWork.SubmissionUUID,
		TestsPassed:      false,
		TestsOutput:      "",
		SubmissionStatus: "running",
	}

	err = manager.sendStatusUpdate(statusDTO)
	if err != nil {
		log.Println("[ERROR] Failed to send submission status update", err.Error())
		return
	}

	// Process submission
	runner, err := utils.GetTestRunnerByLanguageUUID(submissionWork.LanguageUUID)
	if err != nil {
		manager.sendStatusUpdate(utils.GetSubmissionStatusUpdateDTOFromErrorMessage(
			submissionWork.SubmissionUUID,
			"[ERROR] We couldn't find a test runner for the language you submitted",
		))
		return
	}

	result := manager.UseCases.RunTests(&submissionWork, runner)

	// Send submission status update
	err = manager.sendStatusUpdate(result.ToSubmissionStatusUpdateDTO("ready"))
	if err != nil {
		log.Println("[ERROR] Failed to send submission status update", err.Error())
		return
	}
}

func (manager *SubmissionQueueMgr) sendStatusUpdate(dto *dtos.SubmissionStatusUpdateDTO) error {
	submissionStatusQueueMgr := GetSubmissionStatusUpdatesQueueMgrInstance()
	return submissionStatusQueueMgr.QueueUpdate(dto)
}
