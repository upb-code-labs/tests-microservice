package main

import (
	"github.com/upb-code-labs/tests-microservice/src/infrastructure/rabbitmq"
)

func main() {
	// Setup rabbitmq connection
	rabbitmq.ConnectToRabbitMQ()
	defer rabbitmq.CloseRabbitMQConnection()

	// Start listening to submissions queue
	submissionsQueueManager := rabbitmq.GetSubmissionQueueMgr()
	submissionsQueueManager.ListenForSubmissions()

	// Block forever
	var forever chan bool
	<-forever
}
