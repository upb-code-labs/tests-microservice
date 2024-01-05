package definitions

import "github.com/upb-code-labs/tests-microservice/domain/dtos"

type SubmissionStatusUpdatesQueueManager interface {
	QueueUpdate(updateDTO *dtos.SubmissionStatusUpdateDTO) error
}
