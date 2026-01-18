package utils

func Ptr[T any](t T) *T {
	return &t
}

func ValCopy[T any](t *T) *T {
	if t == nil {
		return nil
	}
	return Ptr(*t) // This actually makes a copy https://goplay.tools/snippet/ipMDVGHhgOU
}

func NilStrPtr(t string) *string {
	if t == "" {
		return nil
	}
	return &t
}
