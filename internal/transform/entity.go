package transform

import (
	"github.com/gobeetle/reply/internal/iface"
	"github.com/gobeetle/reply/internal/response"
)

type ResponseTransformer func(source response.Response) (iface.ErrorCoder, error)

// Config controls reshaping of the response after strip, before marshal.
// Zero value: pass-through (or the global transformer when set).
type Config struct {
	// Transformer is an optional custom reshape function.
	// Leave nil for the default pass-through (or global transformer when set).
	Transformer ResponseTransformer

	// DisableGlobal skips the package-level global transformer when Transformer is nil.
	DisableGlobal bool
}

// Resolve builds the ResponseTransformOpt used by handlers.
func (c Config) Resolve() ResponseTransformOpt {
	opt := ResponseTransformOpt{}
	if c.Transformer != nil {
		opt.transformer = c.Transformer
	}
	if c.DisableGlobal {
		opt.without_global_transformer = true
	}
	return opt
}

type ResponseTransformOpt struct {
	transformer                ResponseTransformer
	without_global_transformer bool // default is false, meaning it will use the global transformer if available
}

func (t *ResponseTransformOpt) WithResponseTransformer(transformer ResponseTransformer) *ResponseTransformOpt {
	if transformer == nil {
		return t
	}
	t.transformer = transformer
	return t
}

func (t *ResponseTransformOpt) WithGlobalTransformer(global_transformer bool) *ResponseTransformOpt {
	t.without_global_transformer = !global_transformer
	return t
}

func (t *ResponseTransformOpt) Transform(source response.Response) (iface.ErrorCoder, error) {
	if t.transformer != nil {
		return t.transformer(source)
	} else if !t.without_global_transformer && GetGlobalTransformer() != nil {
		return GetGlobalTransformer()(source)
	}
	return &source, nil
}
