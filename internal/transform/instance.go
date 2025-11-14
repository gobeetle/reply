package transform

var (
	globalTransformer ResponseTransformer
)

// SetGlobalTransform sets a global transform
func SetGlobalTransformer(transformer ResponseTransformer) {
	globalTransformer = transformer
}

// GetGlobalTransform retrieves the global transform
func GetGlobalTransformer() ResponseTransformer {
	return globalTransformer
}
