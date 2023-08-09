package set

import (
	"context"
	"hash/fnv"
	"reflect"
	"testing"
)

// TestToStr tests toStr function.
func TestToStr(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  uint64
	}{
		{
			name:  "Pointer",
			input: new(int),
			want:  12638135523509116079,
		},
		{
			name:  "NilPointer",
			input: (*int)(nil),
			want:  2397808468787316396,
		},
		{
			name:  "Interface",
			input: (interface{})(new(int)),
			want:  12638135523509116079,
		},
		{
			name:  "Func",
			input: func() {},
			want:  852608543138426317,
		},
		{
			name:  "NilFunc",
			input: (func())(nil),
			want:  5584826337234219198,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash := fnv.New64a()
			err := toHash(nil, reflect.ValueOf(test.input), hash)
			if err != nil {
				t.Errorf("toStr(%s) no error was expected", test.name)
			}

			if got := hash.Sum64(); got != test.want {
				t.Errorf("toStr(%s) = %d, want %d", test.name, got, test.want)
			}
		})
	}
}

// TestToStrWithContext tests toStr function with context.
func TestToStrWithContext(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  uint64
	}{
		{
			name:  "Map",
			input: map[int]string{1: "one"},
			want:  14695981039346656037,
		},
		{
			name:  "Pointer",
			input: new(int),
			want:  14695981039346656037,
		},
		{
			name:  "NilPointer",
			input: (*int)(nil),
			want:  14695981039346656037,
		},
		{
			name:  "Interface",
			input: (interface{})(new(int)),
			want:  14695981039346656037,
		},
		{
			name:  "Func",
			input: func() {},
			want:  14695981039346656037,
		},
		{
			name:  "NilFunc",
			input: (func())(nil),
			want:  14695981039346656037,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			cancel()
			hash := fnv.New64a()
			err := toHash(ctx, reflect.ValueOf(test.input), hash)
			if err == nil {
				t.Errorf("toStr(%s) error was expected", test.name)
			}

			if got := hash.Sum64(); got != test.want {
				t.Errorf("toStr(%s) = %d, want %d", test.name, got, test.want)
			}
		})
	}
}
