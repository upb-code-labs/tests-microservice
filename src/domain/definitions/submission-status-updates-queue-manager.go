package definitions

import "github.com/upb-code-labs/tests-microservice/src/domain/dtos"

type SubmissionStatusUpdatesQueueManager interface {
	QueueUpdate(updateDTO *dtos.SubmissionStatusUpdateDTO) error
}
