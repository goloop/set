package set

import (
	"context"
	"fmt"
	"hash"
	"reflect"
)

// toHash is a helper function that takes a reflect.Value and creates a
// string representation of it. This function uses a switch statement to
// handle different kinds of complex types like Struct, Array, Slice, Map,
// Ptr, Interface, and Func. For each kind, it recursively builds a string
// representation and joins them together. If the kind doesn't fall into one of
// these categories, it uses the built-in Sprintf function to create a string.
// This function is mainly used by 'toHash' function to create unique keys for
// complex objects in the Set.
func toHash(ctx context.Context, v reflect.Value, hash hash.Hash64) error {
	// If the context is done, return an error.
	if ctx == nil {
		ctx = context.Background()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Handle different kinds of complex types like Struct, Array, Slice, Map,
	// Ptr, Interface, and Func.
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			err := toHash(ctx, v.Field(i), hash)
			if err != nil {
				return err
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			err := toHash(ctx, v.Index(i), hash)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			err := toHash(ctx, k, hash)
			if err != nil {
				return err
			}
			err = toHash(ctx, v.MapIndex(k), hash)
			if err != nil {
				return err
			}
		}
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return toHash(ctx, reflect.ValueOf("nil"), hash)
		}
		return toHash(ctx, v.Elem(), hash)
	case reflect.Func:
		if v.IsNil() {
			return toHash(ctx, reflect.ValueOf("func:nil"), hash)
		}
		return toHash(ctx, reflect.ValueOf(v.Type().String()+" Value"), hash)
	default:
		_, err := hash.Write([]byte(fmt.Sprintf("%v", v)))
		return err
	}

	return nil
}
