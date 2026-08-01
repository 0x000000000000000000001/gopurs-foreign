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
