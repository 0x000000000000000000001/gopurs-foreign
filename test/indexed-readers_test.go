package Foreign_Index

import (
	"reflect"
	"testing"

	rt "gopurs/output/gopurs_runtime"
)

func readIndexForTest(key, value rt.Value) rt.Value {
	return UnsafeReadPropImpl(rt.Str("failure"), rt.Func(func(v rt.Value) rt.Value { return v }), key, value)
}

func TestReadNativeArrayIndices(t *testing.T) {
	type texts []string
	for name, value := range map[string]rt.Value{
		"native": rt.Any([]any{"a", "b"}),
		"typed":  rt.Any([]string{"a", "b"}),
		"named":  rt.Any(texts{"a", "b"}),
		"boxed":  rt.Array([]rt.Value{rt.Str("a"), rt.Str("b")}),
	} {
		for _, key := range []rt.Value{rt.Int(1), rt.Str("1")} {
			if got := rt.ForeignValue(readIndexForTest(key, value)); got != "b" {
				t.Fatalf("%s index: got %#v", name, got)
			}
		}
		for _, key := range []rt.Value{rt.Int(-1), rt.Int(2), rt.Int(1 << 62)} {
			if got := readIndexForTest(key, value); got.Type != 0 {
				t.Fatalf("%s bounds: %+v", name, got)
			}
		}
	}
	for _, key := range []rt.Value{rt.Str("x"), rt.Str("1x"), rt.Str("01"), rt.Str("+1"), rt.Str("1.0"), rt.Str("-0"), rt.Bool(true)} {
		if got := readIndexForTest(key, rt.Any([]any{"a", "b"})); got.Type != 0 {
			t.Fatalf("non-index key: %+v", key)
		}
	}
	for _, value := range []any{[]any{}, []any(nil)} {
		if got := readIndexForTest(rt.Int(0), rt.Any(value)); got.Type != 0 {
			t.Fatal("empty array index", got)
		}
	}
	if got := readIndexForTest(rt.Str("value"), rt.Any(map[string]any{"value": "kept"})); rt.ForeignValue(got) != "kept" {
		t.Fatal("existing native object property changed")
	}
}

func TestReadNativeObjectPreservesChildren(t *testing.T) {
	child := []any{"shared"}
	cycle := make([]any, 1)
	cycle[0] = cycle
	objectCycle := map[string]any{}
	objectCycle["self"] = objectCycle
	input := rt.Any(map[string]any{"first": child, "second": child, "cycle": cycle, "objectCycle": objectCycle})
	for _, key := range []string{"first", "second"} {
		got := readIndexForTest(rt.Str(key), input)
		if got.Type != rt.TypeAny || reflect.ValueOf(got.PtrVal()).Pointer() != reflect.ValueOf(child).Pointer() {
			t.Fatal("shared native property identity", key)
		}
	}
	for key, expected := range map[string]any{"cycle": cycle, "objectCycle": objectCycle} {
		got := readIndexForTest(rt.Str(key), input)
		if got.Type != rt.TypeAny || reflect.ValueOf(got.PtrVal()).Pointer() != reflect.ValueOf(expected).Pointer() {
			t.Fatal("cyclic native property identity", key)
		}
	}
}

func TestReadNativeArrayPreservesChildren(t *testing.T) {
	child := map[string]any{"value": true}
	cycle := make([]any, 1)
	cycle[0] = cycle
	function := rt.Func(func(v rt.Value) rt.Value { return v })
	constructor := rt.Value{Type: rt.TypeConstructor, IntVal: 123}
	input := rt.Any([]any{child, child, cycle, function, rt.Value{}, constructor, nil})
	for _, index := range []int64{0, 1} {
		got := readIndexForTest(rt.Int(index), input)
		if reflect.ValueOf(got.PtrVal()).Pointer() != reflect.ValueOf(child).Pointer() {
			t.Fatal("shared child identity")
		}
	}
	if got := readIndexForTest(rt.Int(2), input); reflect.ValueOf(got.PtrVal()).Pointer() != reflect.ValueOf(cycle).Pointer() {
		t.Fatal("cycle identity")
	}
	if got := readIndexForTest(rt.Int(3), input); rt.Apply(got, rt.Int(42)).IntVal != 42 {
		t.Fatal("legacy function changed")
	}
	if got := readIndexForTest(rt.Int(4), input); got.Type != 0 {
		t.Fatal("undefined changed")
	}
	if got := readIndexForTest(rt.Int(5), input); got != constructor {
		t.Fatal("constructor changed")
	}
	if got := readIndexForTest(rt.Int(6), input); got.Type != rt.TypeAny || got.PtrVal() != nil {
		t.Fatal("null changed")
	}
}

func TestNativeArrayIndexMembership(t *testing.T) {
	input := rt.Any([]any{rt.Value{}, "present"})
	for _, key := range []rt.Value{rt.Int(0), rt.Int(1), rt.Str("1")} {
		if !UnsafeHasProperty(key, input).BoolVal() || !UnsafeHasOwnProperty(key, input).BoolVal() {
			t.Fatalf("existing native array index reported absent: %+v", key)
		}
	}
	for _, key := range []rt.Value{rt.Int(-1), rt.Int(2), rt.Str("01"), rt.Str("not-an-index")} {
		if UnsafeHasProperty(key, input).BoolVal() || UnsafeHasOwnProperty(key, input).BoolVal() {
			t.Fatalf("missing native array index reported present: %+v", key)
		}
	}
}
