package dto

type Readlist[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
}

// NewReadlist creates a new Readlist from a slice of items (both pointer and non-pointer types supported)
func NewReadlist[T any](items []T) *Readlist[T] {
	return &Readlist[T]{
		Items: items,
		Total: len(items),
	}
}

// NewReadListFromPtr creates a new ReadList from a slice of pointers
func NewReadlistFromPtr[T any](items []*T) *Readlist[T] {
	result := make([]T, 0, len(items))
	for _, item := range items {
		if item != nil {
			result = append(result, *item)
		} else {
			var zero T
			result = append(result, zero)
		}
	}
	return NewReadlist(result)
}
