package utils

import "fmt"

func ValCopy[T any](t *T) *T {
	if t == nil {
		return nil
	}
	return new(*t) // This actually makes a copy https://goplay.tools/snippet/GzO6ESjC-JE
}

func NilStrPtr(t string) *string {
	if t == "" {
		return nil
	}
	return &t
}

func FormatPtr[T any](p *T) string {
	if p == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%+v", *p)
}
