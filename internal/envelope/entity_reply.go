package envelope

type ReplyEnvelope struct {
	Code    int      `json:"code,omitempty"`
	Errors  []string `json:"errors,omitempty"`
	Message []string `json:"message,omitempty"`
	Data    any      `json:"data,omitempty"`
}
