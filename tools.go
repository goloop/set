package set

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

// toStr is a helper function that takes a reflect.Value and creates a
// string representation of it. This function uses a switch statement to
// handle different kinds of complex types like Struct, Array, Slice, Map,
// Ptr, Interface, and Func. For each kind, it recursively builds a string
// representation and joins them together. If the kind doesn't fall into one of
// these categories, it uses the built-in Sprintf function to create a string.
// This function is mainly used by 'toHash' function to create unique keys for
// complex objects in the Set.
func toStr(ctx context.Context, v reflect.Value) (string, error) {
	// If the context is nil, create a new one.
	if ctx == nil {
		ctx = context.Background()
	}

	// If the context is done, return an empty string.
	// Execute the context cancellation check block here,
	// because this method can be executed recursively.
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Create a string representation of the given reflect.Value.
	// This procedure performs a recursive call toStr.
	switch v.Kind() {
	case reflect.Struct:
		var r []string
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			name := t.Field(i).Name
			value, err := toStr(ctx, v.Field(i))
			if err != nil {
				return "", err
			}

			r = append(r, fmt.Sprintf("%s:%s", name, value))
		}
		return "{" + strings.Join(r, ", ") + "}", nil
	case reflect.Array, reflect.Slice:
		var elements []string
		for i := 0; i < v.Len(); i++ {
			value, err := toStr(ctx, v.Index(i))
			if err != nil {
				return "", err
			}

			elements = append(elements, value)
		}
		return "[" + strings.Join(elements, ", ") + "]", nil
	case reflect.Map:
		var r []string
		for _, k := range v.MapKeys() {
			v := v.MapIndex(k)

			kValue, err := toStr(ctx, k)
			if err != nil {
				return "", err
			}

			vValue, err := toStr(ctx, v)
			if err != nil {
				return "", err
			}

			r = append(r, fmt.Sprintf("%s:%s", kValue, vValue))
		}
		return "{" + strings.Join(r, ", ") + "}", nil
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return "nil", nil
		}

		return toStr(ctx, v.Elem())
	case reflect.Func:
		if v.IsNil() {
			return "func:nil", nil
		}

		return v.Type().String() + " Value", nil
	default:
		// pass
	}

	return fmt.Sprintf("%v", v), nil
}
