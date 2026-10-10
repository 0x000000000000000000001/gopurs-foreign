package Foreign

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	rt "gopurs/output/gopurs_runtime"
)

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
