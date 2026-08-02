package Foreign

import (
	"gopurs/output/gopurs_runtime"
	"reflect"
)

func TypeOf(v gopurs_runtime.Value) gopurs_runtime.Value {
	switch v.Type {
	case gopurs_runtime.TypeBool: return gopurs_runtime.Str("boolean")
	case gopurs_runtime.TypeInt, gopurs_runtime.TypeFloat: return gopurs_runtime.Str("number")
	case gopurs_runtime.TypeString: return gopurs_runtime.Str("string")
	case gopurs_runtime.TypeFunc, gopurs_runtime.TypeFunc2, gopurs_runtime.TypeFunc3, 
	     gopurs_runtime.TypeFunc4, gopurs_runtime.TypeFunc5, gopurs_runtime.TypeFunc6,
	     gopurs_runtime.TypeFunc7, gopurs_runtime.TypeFunc8, gopurs_runtime.TypeFunc9,
	     gopurs_runtime.TypeFunc10: 
		return gopurs_runtime.Str("function")
	case 0: return gopurs_runtime.Str("undefined")
	default:
		return gopurs_runtime.Str("object")
	}
}

func TagOf(v gopurs_runtime.Value) gopurs_runtime.Value {
	switch v.Type {
	case gopurs_runtime.TypeBool: return gopurs_runtime.Str("Boolean")
	case gopurs_runtime.TypeInt, gopurs_runtime.TypeFloat: return gopurs_runtime.Str("Number")
	case gopurs_runtime.TypeString: return gopurs_runtime.Str("String")
	case gopurs_runtime.TypeArray: return gopurs_runtime.Str("Array")
	case gopurs_runtime.TypeFunc, gopurs_runtime.TypeFunc2, gopurs_runtime.TypeFunc3, 
	     gopurs_runtime.TypeFunc4, gopurs_runtime.TypeFunc5, gopurs_runtime.TypeFunc6,
	     gopurs_runtime.TypeFunc7, gopurs_runtime.TypeFunc8, gopurs_runtime.TypeFunc9,
	     gopurs_runtime.TypeFunc10: 
		return gopurs_runtime.Str("Function")
	case 0: return gopurs_runtime.Str("Undefined")
	default:
		if v.Type == gopurs_runtime.TypeAny {
			if v.UnsafePtr == nil {
				return gopurs_runtime.Str("Null")
			}
			val := *(*any)(v.UnsafePtr)
			if val == nil {
				return gopurs_runtime.Str("Null")
			}
			rt := reflect.TypeOf(val)
			if rt != nil && rt.Kind() == reflect.Slice {
				return gopurs_runtime.Str("Array")
			}
		}
		return gopurs_runtime.Str("Object")
	}
}

func IsNull(v gopurs_runtime.Value) gopurs_runtime.Value {
	if v.Type == gopurs_runtime.TypeAny {
		if v.UnsafePtr == nil {
			return gopurs_runtime.Bool(true)
		}
		val := *(*any)(v.UnsafePtr)
		if val == nil {
			return gopurs_runtime.Bool(true)
		}
	}
	return gopurs_runtime.Bool(false)
}

func IsUndefined(v gopurs_runtime.Value) gopurs_runtime.Value {
	if v.Type == 0 {
		return gopurs_runtime.Bool(true)
	}
	return gopurs_runtime.Bool(false)
}

func IsArray(v gopurs_runtime.Value) gopurs_runtime.Value {
	if v.Type == gopurs_runtime.TypeArray {
		return gopurs_runtime.Bool(true)
	}
	if v.Type == gopurs_runtime.TypeAny && v.UnsafePtr != nil {
		val := *(*any)(v.UnsafePtr)
		if val != nil && reflect.TypeOf(val).Kind() == reflect.Slice {
			return gopurs_runtime.Bool(true)
		}
	}
	return gopurs_runtime.Bool(false)
}

func UnboxForJSON(v interface{}) interface{} {
	return deepUnbox(v)
}

func deepUnbox(v interface{}) interface{} {
	if val, ok := v.(gopurs_runtime.Value); ok {
		switch val.Type {
		case gopurs_runtime.TypeInt:
			return val.IntVal
		case gopurs_runtime.TypeFloat:
			if val.UnsafePtr != nil {
				return *(*float64)(val.UnsafePtr)
			}
			return 0.0
		case gopurs_runtime.TypeString:
			if val.UnsafePtr != nil {
				return *(*string)(val.UnsafePtr)
			}
			return ""
		case gopurs_runtime.TypeBool:
			return val.IntVal != 0
		case gopurs_runtime.TypeArray:
			if val.UnsafePtr != nil {
				arr := *(*[]gopurs_runtime.Value)(val.UnsafePtr)
				res := make([]interface{}, len(arr))
				for i, x := range arr {
					res[i] = deepUnbox(x)
				}
				return res
			}
			return []interface{}{}
		case gopurs_runtime.TypeRecord, gopurs_runtime.TypeRecord0, gopurs_runtime.TypeRecord1, gopurs_runtime.TypeRecord2, gopurs_runtime.TypeRecord3, gopurs_runtime.TypeRecord4, gopurs_runtime.TypeRecord5:
			rec := gopurs_runtime.RecordToMap(val)
			res := make(map[string]interface{})
			for k, x := range rec {
				res[k] = deepUnbox(x)
			}
			return res
		case gopurs_runtime.TypeAny:
			if val.UnsafePtr == nil {
				return nil
			}
			if *(*any)(val.UnsafePtr) == nil {
				return nil
			}
			return deepUnbox(*(*any)(val.UnsafePtr))
		}
	}
	return v
}
