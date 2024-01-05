# Environment

This document describes the required environment variables to run the micro-service.

| Name                          | Description                                     | Example                                  | Mandatory |
| ----------------------------- | ----------------------------------------------- | ---------------------------------------- |  -------- |
| `RABBIT_MQ_CONNECTION_STRING` | The connection string to the RabbitMQ instance. | `amqp://username:password@address:port/` | Yes       |
| `STATIC_FILES_MICROSERVICE_ADDRESS`        | The address of the static files micro-service.  | `http://domain:port`                  | Yes       |
| `TESTS_EXECUTION_DIRECTORY`      | The path (preferably absolute) to the directory where the tests will be executed. | `/tmp` | No       |
