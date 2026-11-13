package Test_Main

type nativeText string
type nativeFlag bool
type nativeInts []int64
type nativeNumber float64
type nativeInteger int64

func NativeValue(name string) any {
	switch name {
	case "string":
		return "hello"
	case "emptyString":
		return ""
	case "namedString":
		return nativeText("named")
	case "true":
		return true
	case "false":
		return false
	case "namedBoolean":
		return nativeFlag(true)
	case "namedNumber":
		return nativeNumber(42.5)
	case "namedInteger":
		return nativeInteger(42)
	case "unsigned":
		return uint32(42)
	case "smallSigned":
		return int8(-42)
	case "float32":
		return float32(42.5)
	case "array":
		return []any{"hello", true, float64(42)}
	case "typedArray":
		return []string{"a", "b"}
	case "namedArray":
		return nativeInts{1, 2}
	case "nested":
		return []any{[]any{"nested", false}}
	case "emptyArray":
		return []any{}
	case "nilArray":
		return []any(nil)
	case "object":
		return map[string]any{"value": 42}
	case "null":
		return nil
	default:
		panic("Unknown native Foreign test value: " + name)
	}
}
