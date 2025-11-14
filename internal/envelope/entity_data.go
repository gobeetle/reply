package envelope

type DataEnvelope[T any] struct {
	Data T `json:"data"`
}
