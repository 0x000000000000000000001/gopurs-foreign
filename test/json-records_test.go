package Foreign

import (
	"encoding/json"
	"math"
	"os"
	"reflect"
	"strconv"
	"testing"

	rt "gopurs/output/gopurs_runtime"
)

func TestNativeForeignReaders(t *testing.T) {
	type text string
	type flag bool
	for _, value := range []any{"hello", text("hello")} {
		if TagOf(rt.Any(value)).StrVal() != "String" || UnsafeReadString(value) != "hello" {
			t.Fatalf("string reader: %T", value)
		}
	}
	for _, value := range []any{true, flag(true)} {
		if TagOf(rt.Any(value)).StrVal() != "Boolean" || !UnsafeReadBoolean(value) {
			t.Fatalf("boolean reader: %T", value)
		}
	}
	if UnsafeReadString("") != "" || UnsafeReadBoolean(false) {
		t.Fatal("empty/false primitive changed")
	}
	type number float64
	type integer int64
	for _, value := range []any{int(42), int8(42), int16(42), int32(42), int64(42),
		uint(42), uint8(42), uint16(42), uint32(42), uint64(42),
		float32(42), float64(42), number(42), integer(42)} {
		if TagOf(rt.Any(value)).StrVal() != "Number" || UnsafeReadNumber(value) != 42 {
			t.Fatalf("number reader: %T", value)
		}
	}
	if !math.IsNaN(UnsafeReadNumber(number(math.NaN()))) ||
		!math.IsInf(UnsafeReadNumber(number(math.Inf(1))), 1) ||
		!math.IsInf(UnsafeReadNumber(float32(math.Inf(-1))), -1) {
		t.Fatal("non-finite Number changed")
	}
}

func TestReadArrayPreservesElementRepresentations(t *testing.T) {
	type strings []string
	for _, value := range []any{[]any{int64(1), "a", true}, []string{"a", "b"}, strings{"a", "b"}, []any{}, []any(nil)} {
		input := rt.Any(value)
		if !IsArray(input).BoolVal() {
			t.Fatalf("array guard: %T", value)
		}
		got := UnsafeReadArray(input)
		source := reflect.ValueOf(value)
		if len(got) != source.Len() {
			t.Fatalf("array length: %T", value)
		}
		for i := range got {
			if !reflect.DeepEqual(rt.ForeignValue(got[i]), source.Index(i).Interface()) {
				t.Fatalf("array element %d", i)
			}
		}
	}
	function := rt.Func(func(value rt.Value) rt.Value { return value })
	boxed := []rt.Value{{}, function, {Type: rt.TypeConstructor, IntVal: 123}}
	borrowed := UnsafeReadArray(rt.Array(boxed))
	if &borrowed[0] != &boxed[0] || borrowed[0].Type != 0 || borrowed[2] != boxed[2] ||
		rt.Apply(borrowed[1], rt.Int(42)).IntVal != 42 {
		t.Fatal("legacy array contents or borrowing changed")
	}
	child := map[string]any{"value": true}
	cycle := make([]any, 1)
	cycle[0] = cycle
	items := UnsafeReadArray(rt.Any([]any{child, child, cycle, function, rt.Value{}}))
	for _, index := range []int{0, 1} {
		if reflect.ValueOf(items[index].PtrVal()).Pointer() != reflect.ValueOf(child).Pointer() {
			t.Fatal("shared child changed")
		}
	}
	if reflect.ValueOf(items[2].PtrVal()).Pointer() != reflect.ValueOf(cycle).Pointer() ||
		rt.Apply(items[3], rt.Int(42)).IntVal != 42 || items[4].Type != 0 {
		t.Fatal("native array children were recursively converted")
	}
}

func TestJSONRecordRepresentations(t *testing.T) {
	values := []rt.Value{rt.Any(nil), rt.Str(""), rt.Bool(false), rt.Int(0), rt.Str("Ω"), rt.Array([]rt.Value{rt.Any(nil), rt.Bool(true)})}
	keys := []string{"null", "empty", "false", "zero", "text", "array"}
	native := []any{nil, "", false, int64(0), "Ω", []any{nil, true}}
	compact := []rt.Value{
		rt.RecordDict0(),
		rt.RecordDict1(keys[0], values[0]),
		rt.RecordDict2(keys[0], keys[1], values[0], values[1]),
		rt.RecordDict3(keys[0], keys[1], keys[2], values[0], values[1], values[2]),
		rt.RecordDict4(keys[0], keys[1], keys[2], keys[3], values[0], values[1], values[2], values[3]),
		rt.RecordDict5(keys[0], keys[1], keys[2], keys[3], keys[4], values[0], values[1], values[2], values[3], values[4]),
	}
	for size := 0; size <= len(keys); size++ {
		fields := make(map[string]rt.Value)
		expected := make(map[string]any)
		for i := 0; i < size; i++ {
			fields[keys[i]] = values[i]
			expected[keys[i]] = native[i]
		}
		layouts := map[string]rt.Value{
			"map":    rt.Record(fields),
			"data":   rt.RecordDict(keys[:size], values[:size]),
			"opaque": rt.Any(fields),
		}
		if size < len(compact) {
			layouts["compact"] = compact[size]
		}
		for name, value := range layouts {
			t.Run(name+"/"+strconv.Itoa(size), func(t *testing.T) {
				if got := UnboxForJSON(value); !reflect.DeepEqual(got, expected) {
					t.Fatalf("record conversion: got %#v, want %#v", got, expected)
				}
			})
		}
	}
}

func TestJSONNestedNullAndUndefined(t *testing.T) {
	nested := rt.RecordDict([]string{"present", "absent"}, []rt.Value{rt.Any(nil), {}})
	value := rt.RecordDict2("object", "array", nested, rt.Array([]rt.Value{nested, rt.Any(nil)}))
	encoded, err := json.Marshal(UnboxForJSON(value))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(encoded), `{"array":[{"present":null},null],"object":{"present":null}}`; got != want {
		t.Fatalf("JSON: got %s, want %s", got, want)
	}
	if !IsNull(rt.Any(nil)).BoolVal() || IsUndefined(rt.Any(nil)).BoolVal() || !IsUndefined(rt.Value{}).BoolVal() {
		t.Fatal("null and undefined must remain distinct")
	}
}

func TestJSONAfterNativeForeignTransport(t *testing.T) {
	reference, err := os.ReadFile("json-boundary-reference.json")
	if err != nil {
		t.Fatal(err)
	}
	var expected map[string]any
	if err := json.Unmarshal(reference, &expected); err != nil {
		t.Fatal(err)
	}
	leaf := rt.RecordDict3("absent", "null", "zero", rt.Value{}, rt.Any(nil), rt.Int(0))
	object := rt.RecordDict3("absent", "nested", "list", rt.Value{}, leaf, rt.Array([]rt.Value{{}, rt.Any(nil), leaf}))
	mixed := rt.Any([]any{rt.Value{}, nil, map[string]any{"absent": rt.Value{}, "null": nil}})
	for name, boxed := range map[string]rt.Value{"object": object, "mixed": mixed} {
		t.Run(name, func(t *testing.T) {
			native := rt.ForeignValue(boxed)
			// The callback adapter wraps its native argument with Any, exactly
			// as the generated Foreign -> String callback does for PostgreSQL.
			encoded, err := json.Marshal(UnboxForJSON(rt.Any(native)))
			if err != nil {
				t.Fatal(err)
			}
			var actual any
			if err := json.Unmarshal(encoded, &actual); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(actual, expected[name]) {
				t.Fatalf("JSON differs from JavaScript: got %s, want %#v", encoded, expected[name])
			}
		})
	}
	native := rt.ForeignValue(object).(map[string]any)
	undefined, present := native["absent"]
	if !present || undefined == nil || !rt.IsForeignUndefined(undefined) {
		t.Fatal("transport discarded an undefined property or changed it to null")
	}
	if _, present := native["missing"]; present {
		t.Fatal("transport added an absent property")
	}
	for _, value := range []rt.Value{rt.Any(undefined), rt.Box(undefined)} {
		if !IsUndefined(value).BoolVal() || IsNull(value).BoolVal() || TypeOf(value).StrVal() != "undefined" || TagOf(value).StrVal() != "Undefined" {
			t.Fatal("native callback/result roundtrip changed undefined semantics")
		}
	}
	if !IsNull(rt.Any(native["nested"].(map[string]any)["null"])).BoolVal() {
		t.Fatal("null did not survive the native boundary")
	}
}
