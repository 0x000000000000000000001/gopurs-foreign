package Foreign

import (
	"testing"

	rt "gopurs/output/gopurs_runtime"
)

func TestFunctionDataClassification(t *testing.T) {
	ordinary := rt.Func2(func(a, b rt.Value) rt.Value { return b })
	wrapped := rt.WithFunctionData(ordinary, "metadata")
	eleven := rt.Func11(func(a, b, c, d, e, f, g, h, i, j, k rt.Value) rt.Value { return k })
	for _, value := range []rt.Value{ordinary, wrapped, eleven, rt.Apply(wrapped, rt.Int(1)), rt.Box(wrapped.AnyVal())} {
		if got := TypeOf(value).StrVal(); got != "function" {
			t.Fatalf("typeOf tag %d = %q", value.Type, got)
		}
		if got := TagOf(value).StrVal(); got != "Function" {
			t.Fatalf("tagOf tag %d = %q", value.Type, got)
		}
	}
	if TypeOf(rt.Int(1)).StrVal() != "number" || TagOf(rt.RecordDict0()).StrVal() != "Object" {
		t.Fatal("non-function classification changed")
	}
}
