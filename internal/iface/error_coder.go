package iface

type ErrorCoder interface {
	StatusCoder // StatusCode() int
	error       // Error() string
}
