package transform

import (
	"github.com/gobeetle/reply/internal/iface"
	"github.com/gobeetle/reply/internal/response"
)

type ResponseTransformer func(source response.Response) (iface.ErrorCoder, error)

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
