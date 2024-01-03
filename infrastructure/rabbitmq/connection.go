package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/upb-code-labs/tests-microservice/infrastructure"
)

var rabbitMQChannel *amqp.Channel

func ConnectToRabbitMQ() {
	// Stablish connection
	rabbitMQConnectionString := infrastructure.GetEnvironment().RabbitMQConnectionString
	conn, err := amqp.Dial(rabbitMQConnectionString)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Set channel
	log.Println("Connected to RabbitMQ")
	rabbitMQChannel = ch
}

func CloseRabbitMQConnection() {
	if rabbitMQChannel != nil {
		rabbitMQChannel.Close()
	}
}

func GetRabbitMQChannel() *amqp.Channel {
	if rabbitMQChannel == nil {
		ConnectToRabbitMQ()
	}

	return rabbitMQChannel
}
