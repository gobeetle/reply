package envelope

type MessageEnvelope[M ~string] struct {
	Message M `json:"message"`
}
