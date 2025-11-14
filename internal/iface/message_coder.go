package iface

// ==============================

type MessageCoder interface {
	StatusCoder     // StatusCode() int
	MessageProvider // Message() []string
}
