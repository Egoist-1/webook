package integration

type Result[T any] struct {
	Msg  string
	Code int
	Data T
}
