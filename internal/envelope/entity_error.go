package envelope


type ErrorEnvelope[E error] struct {
	Error  E   `json:"error"`
	Errors []E `json:"errors"`
}
