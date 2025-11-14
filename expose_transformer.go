package reply

import (
	"github.com/gobeetle/reply/internal/transform"
)

type (
	ResponseTransformer  = transform.ResponseTransformer
	ResponseTransformOpt = transform.ResponseTransformOpt
)

// SetGlobalTransform sets a global transform
func SetGlobalTransform(t ResponseTransformer) {
	transform.SetGlobalTransformer(t)
}

// GetGlobalTransform retrieves the global transform
func GetGlobalTransform() ResponseTransformer {
	return transform.GetGlobalTransformer()
}
