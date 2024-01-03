package infrastructure

import "github.com/kelseyhightower/envconfig"

type Environment struct {
	RabbitMQConnectionString       string `split_words:"true" default:"amqp://rabbitmq:rabbitmq@localhost:5672/"`
	StaticFilesMicroserviceAddress string `split_words:"true" default:"http://localhost:8081"`
}

var env *Environment

func GetEnvironment() *Environment {
	if env == nil {
		env = &Environment{}

		err := envconfig.Process("", env)
		if err != nil {
			panic(err)
		}
	}

	return env
}
