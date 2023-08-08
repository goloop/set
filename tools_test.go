package set

import (
	"context"
	"reflect"
	"testing"
)

// TestToStr tests toStr function.
func TestToStr(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name:  "Pointer",
			input: new(int),
			want:  "0",
		},
		{
			name:  "NilPointer",
			input: (*int)(nil),
			want:  "nil",
		},
		{
			name:  "Interface",
			input: (interface{})(new(int)),
			want:  "0",
		},
		{
			name:  "Func",
			input: func() {},
			want:  "func() Value",
		},
		{
			name:  "NilFunc",
			input: (func())(nil),
			want:  "func:nil",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := toStr(nil, reflect.ValueOf(test.input))
			if err != nil {
				t.Errorf("toStr(%s) no error was expected", test.name)
			}

			if got != test.want {
				t.Errorf("toStr(%s) = %s, want %s", test.name, got, test.want)
			}
		})
	}
}

// TestToStrWithContext tests toStr function with context.
func TestToStrWithContext(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name:  "Map",
			input: map[int]string{1: "one"},
			want:  "0",
		},
		{
			name:  "Pointer",
			input: new(int),
			want:  "0",
		},
		{
			name:  "NilPointer",
			input: (*int)(nil),
			want:  "nil",
		},
		{
			name:  "Interface",
			input: (interface{})(new(int)),
			want:  "0",
		},
		{
			name:  "Func",
			input: func() {},
			want:  "func() Value",
		},
		{
			name:  "NilFunc",
			input: (func())(nil),
			want:  "func:nil",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			cancel()
			_, err := toStr(ctx, reflect.ValueOf(test.input))
			if err == nil {
				t.Errorf("toStr(%s) error was expected", test.name)
			}
		})
	}
}
