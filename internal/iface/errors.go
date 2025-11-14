package iface

type Errors interface {
	Unwrap() []error
}
