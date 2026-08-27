package Foreign

import (
	"gopurs/output/gopurs_runtime"
	"reflect"
)

var undefinedForJSON = struct{}{}

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
		if v.Type == gopurs_runtime.TypeAny && v.UnsafePtr != nil {
			val := *(*any)(v.UnsafePtr)
			if val != nil {
				rt := reflect.TypeOf(val)
				if rt != nil {
				if rt.Name() == "Date" {
					return gopurs_runtime.Str("Date")
				}

					switch rt.Kind() {
					case reflect.String:
						return gopurs_runtime.Str("string")
					case reflect.Bool:
						return gopurs_runtime.Str("boolean")
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					     reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
					     reflect.Float32, reflect.Float64:
						return gopurs_runtime.Str("number")
					}
				}
			}
		}
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
			if rt != nil {
				if rt.Kind() == reflect.Slice {
					return gopurs_runtime.Str("Array")
				}
				if rt.Kind() == reflect.String {
					return gopurs_runtime.Str("String")
				}
				if rt.Kind() == reflect.Bool {
					return gopurs_runtime.Str("Boolean")
				}
				if rt.Name() == "Date" {
					return gopurs_runtime.Str("Date")
				}

				switch rt.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					 reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
					 reflect.Float32, reflect.Float64:
					return gopurs_runtime.Str("Number")
				}
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
	return foreignDeepUnbox(v)
}

func foreignDeepUnbox(v interface{}) interface{} {
	if val, ok := v.(gopurs_runtime.Value); ok {
		switch val.Type {
		case gopurs_runtime.TypeInt:
			return val.IntVal
		case gopurs_runtime.TypeFloat:
			return val.FloatVal()
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
					res[i] = foreignDeepUnbox(x)
				}
				return res
			}
			return []interface{}{}
		case gopurs_runtime.TypeRecord, gopurs_runtime.TypeRecord0, gopurs_runtime.TypeRecord1, gopurs_runtime.TypeRecord2, gopurs_runtime.TypeRecord3, gopurs_runtime.TypeRecord4, gopurs_runtime.TypeRecord5:
			rec := gopurs_runtime.RecordToMap(val)
			res := make(map[string]interface{})
			for k, x := range rec {
				unboxed := foreignDeepUnbox(x)
				if unboxed != undefinedForJSON {
					res[k] = unboxed
				}
			}
			return res
		case gopurs_runtime.TypeAny:
			if val.UnsafePtr == nil {
				return nil
			}
			if *(*any)(val.UnsafePtr) == nil {
				return nil
			}
			return foreignDeepUnbox(*(*any)(val.UnsafePtr))
		case 0:
			return undefinedForJSON
		default:
			return nil
		}
	}
	if valSlice, ok := v.([]gopurs_runtime.Value); ok {
		res := make([]interface{}, len(valSlice))
		for i, x := range valSlice {
			res[i] = foreignDeepUnbox(x)
		}
		return res
	}
	if mapGopurs, ok := v.(map[string]gopurs_runtime.Value); ok {
		res := make(map[string]interface{})
		for k, x := range mapGopurs {
			unboxed := foreignDeepUnbox(x)
			if unboxed != undefinedForJSON {
				res[k] = unboxed
			}
		}
		return res
	}
	if mapRaw, ok := v.(map[string]interface{}); ok {
		res := make(map[string]interface{})
		for k, x := range mapRaw {
			unboxed := foreignDeepUnbox(x)
			if unboxed != undefinedForJSON {
				res[k] = unboxed
			}
		}
		return res
	}
	if arrRaw, ok := v.([]interface{}); ok {
		res := make([]interface{}, len(arrRaw))
		for i, x := range arrRaw {
			res[i] = foreignDeepUnbox(x)
		}
		return res
	}
	return v
}
