package apperrors

type HandleableError interface {
	StatusCode() int
}
