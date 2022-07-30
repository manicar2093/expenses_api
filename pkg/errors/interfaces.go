package errors

type HandleableError interface {
	StatusCode() int
}
