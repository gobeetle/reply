package iface

type DataCoder interface {
	StatusCoder     // StatusCode() int
	DataProvider // Data() any
}
