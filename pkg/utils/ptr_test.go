package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValCopy(t *testing.T) {
	var a *int
	var b *int

	a = new(2)

	b = ValCopy(a)

	*a = 3

	assert.NotEqual(t, 3, *b)
}

func TestNilStrPtr(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "happy__non-empty",
			args: args{t: "hello"},
			want: new("hello"),
		},
		{
			name: "happy__empty",
			args: args{t: ""},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NilStrPtr(tt.args.t), "NilStrPtr(%v)", tt.args.t)
		})
	}
}

func TestFormatPtr(t *testing.T) {
	var a *int
	assert.Equal(t, "<nil>", FormatPtr(a))

	a = new(2)
	assert.Equal(t, "2", FormatPtr(a))

	var b *string
	assert.Equal(t, "<nil>", FormatPtr(b))

	b = new("hello")
	assert.Equal(t, "hello", FormatPtr(b))
}
