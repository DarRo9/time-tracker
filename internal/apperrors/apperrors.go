package apperrors

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

type ExternalAPIError struct {
	Message string
}

func (e *ExternalAPIError) Error() string {
	return e.Message
}

type DuplicateKeyError struct {
	Message string
}

func (e *DuplicateKeyError) Error() string {
	return e.Message
}

type NoRowsAffectedError struct {
	Message string
}

func (e *NoRowsAffectedError) Error() string {
	return e.Message
}

type NoUserError struct {
	Message string
}

func (e *NoUserError) Error() string {
	return e.Message
}

type NoTaskError struct {
	Message string
}

func (e *NoTaskError) Error() string {
	return e.Message
}

type TaskAlreadyEndedError struct {
	Message string
}

func (e *TaskAlreadyEndedError) Error() string {
	return e.Message
}
