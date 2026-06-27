package set

import (
	"encoding/json"
	"sort"
	"testing"
)

func TestJSONRoundTrip(t *testing.T) {
	s := New(3, 1, 2, 5, 4)

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got Set[int]
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !got.Equal(s) {
		t.Fatalf("round-trip changed set: got %v, want %v",
			got.Elements(), s.Elements())
	}
}

// Marshalling produces a JSON array of exactly the elements (order aside).
func TestMarshalShape(t *testing.T) {
	s := New(1, 2, 3)
	data, err := s.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}

	var arr []int
	if err := json.Unmarshal(data, &arr); err != nil {
		t.Fatalf("result is not a JSON array: %v", err)
	}
	sort.Ints(arr)
	if len(arr) != 3 || arr[0] != 1 || arr[1] != 2 || arr[2] != 3 {
		t.Fatalf("marshalled array = %v, want [1 2 3]", arr)
	}
}

// Unmarshalling collapses duplicates present in the JSON source.
func TestUnmarshalDeduplicates(t *testing.T) {
	var s Set[int]
	if err := s.UnmarshalJSON([]byte(`[1,2,2,3,3,3]`)); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}
	if s.Len() != 3 {
		t.Fatalf("Len = %d, want 3 after dedup", s.Len())
	}
}

// Unmarshalling into a populated set replaces, not merges.
func TestUnmarshalReplaces(t *testing.T) {
	s := New(7, 8, 9)
	if err := s.UnmarshalJSON([]byte(`[1,2]`)); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}
	want := New(1, 2)
	if !s.Equal(want) {
		t.Fatalf("after unmarshal got %v, want %v", s.Elements(), want.Elements())
	}
}

func TestUnmarshalInvalid(t *testing.T) {
	var s Set[int]
	if err := s.UnmarshalJSON([]byte(`{"not":"an array"}`)); err == nil {
		t.Fatal("expected error on invalid JSON")
	}
}

// A set nested in a struct must marshal and unmarshal transparently.
func TestJSONNestedStruct(t *testing.T) {
	type box struct {
		Tags *Set[string] `json:"tags"`
	}
	in := box{Tags: New("go", "set", "go")}

	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var out box
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if out.Tags == nil || !out.Tags.Equal(New("go", "set")) {
		t.Fatalf("nested set round-trip failed: %v", out.Tags)
	}
}
