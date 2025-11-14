package marshal

var (
	globalMarshaler ResponseMarshaler
)

// SetGlobalMarshaler sets a global marshaler
func SetGlobalMarshaler(marshaler ResponseMarshaler) {
	globalMarshaler = marshaler
}

// GetGlobalMarshaler retrieves the global marshaler
func GetGlobalMarshaler() ResponseMarshaler {
	return globalMarshaler
}
