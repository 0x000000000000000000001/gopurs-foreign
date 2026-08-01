package Foreign_Index

import (
	"gopurs/output/gopurs_runtime"
	"strconv"
)

func UnsafeReadPropImpl(f gopurs_runtime.Value, s gopurs_runtime.Value, key gopurs_runtime.Value, value gopurs_runtime.Value) gopurs_runtime.Value {
	if value.Type == gopurs_runtime.TypeAny && value.UnsafePtr == nil {
		return f
	}
	
	var propVal gopurs_runtime.Value
	found := false
	
	var kStr string
	isStr := false
	if key.Type == gopurs_runtime.TypeString {
		kStr = *(*string)(key.UnsafePtr)
		isStr = true
	} else if key.Type == gopurs_runtime.TypeInt {
		kStr = strconv.FormatInt(key.IntVal, 10)
		isStr = true
	}
	
	if value.Type == gopurs_runtime.TypeRecord && isStr {
		m := gopurs_runtime.RecordToMap(value)
		propVal, found = m[kStr]
	} else if value.Type == gopurs_runtime.TypeArray {
		arr := *(*[]gopurs_runtime.Value)(value.UnsafePtr)
		var idx int
		if key.Type == gopurs_runtime.TypeInt {
			idx = int(key.IntVal)
		} else if isStr {
			idx, _ = strconv.Atoi(kStr)
		}
		if idx >= 0 && idx < len(arr) {
			propVal = arr[idx]
			found = true
		}
	} else if value.Type == gopurs_runtime.TypeAny && value.UnsafePtr != nil && isStr {
		native := *(*any)(value.UnsafePtr)
		if m, ok := native.(map[string]interface{}); ok {
			if v, ok := m[kStr]; ok {
				propVal = gopurs_runtime.Box(v)
				found = true
			}
		} else if m, ok := native.(map[string]gopurs_runtime.Value); ok {
			if v, ok := m[kStr]; ok {
				propVal = v
				found = true
			}
		}
	}
	
	if !found {
		propVal = gopurs_runtime.Value{Type: 0}
	}
	
	return gopurs_runtime.Apply(s, propVal)
}

func UnsafeHasOwnProperty(prop gopurs_runtime.Value, value gopurs_runtime.Value) gopurs_runtime.Value {
	var kStr string
	isStr := false
	if prop.Type == gopurs_runtime.TypeString {
		kStr = *(*string)(prop.UnsafePtr)
		isStr = true
	} else if prop.Type == gopurs_runtime.TypeInt {
		kStr = strconv.FormatInt(prop.IntVal, 10)
		isStr = true
	}
	
	if value.Type == gopurs_runtime.TypeRecord && isStr {
		m := gopurs_runtime.RecordToMap(value)
		_, ok := m[kStr]
		return gopurs_runtime.Bool(ok)
	} else if value.Type == gopurs_runtime.TypeAny && value.UnsafePtr != nil && isStr {
		native := *(*any)(value.UnsafePtr)
		if m, ok := native.(map[string]interface{}); ok {
			_, ok := m[kStr]
			return gopurs_runtime.Bool(ok)
		} else if m, ok := native.(map[string]gopurs_runtime.Value); ok {
			_, ok := m[kStr]
			return gopurs_runtime.Bool(ok)
		}
	}
	return gopurs_runtime.Bool(false)
}

func UnsafeHasProperty(prop gopurs_runtime.Value, value gopurs_runtime.Value) gopurs_runtime.Value {
	return UnsafeHasOwnProperty(prop, value)
}
