package rabbitmq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/upb-code-labs/tests-microservice/application"
	"github.com/upb-code-labs/tests-microservice/domain/entities"
	"github.com/upb-code-labs/tests-microservice/utils"
)

type SubmissionQueueManager struct {
	Queue          *amqp.Queue
	MessageChannel <-chan amqp.Delivery
	UseCases       *application.SubmissionsUseCases
}

// submissionQueueManagerInstance Singleton struct
var submissionQueueManagerInstance *SubmissionQueueManager

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
		log.Println("RabbitMQ submissions queue declared / set")
		return &q
	}

	return submissionQueueManagerInstance.Queue
}

// GetSubmissionQueueManager returns a pointer to the singleton instance of SubmissionQueueManager
func GetSubmissionQueueManager() *SubmissionQueueManager {
	if submissionQueueManagerInstance == nil {
		submissionQueueManagerInstance = &SubmissionQueueManager{
			Queue:    GetRabbitMQSubmissionsQueue(),
			UseCases: &application.SubmissionsUseCases{},
		}
	}

	return submissionQueueManagerInstance
}

// ListenForSubmissions starts listening for submissions
func (manager *SubmissionQueueManager) ListenForSubmissions() {
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
func (manager *SubmissionQueueManager) processSubmissions() {
	for msg := range manager.MessageChannel {
		go manager.processSubmission(msg)
	}
}

func (manager *SubmissionQueueManager) processSubmission(msg amqp.Delivery) {
	// log.Printf("Received a message: %s\n", msg.Body)

	// Unmarshal message
	var submissionWork entities.SubmissionWork
	err := json.Unmarshal(msg.Body, &submissionWork)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Process submission
	runner, err := utils.GetTestRunnerByLanguageUUID(submissionWork.LanguageUUID)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = manager.UseCases.RunTests(&submissionWork, runner)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Acknowledge message
	// msg.Ack(false)
}
