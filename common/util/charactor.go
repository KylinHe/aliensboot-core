package util

import (
	"fmt"
	"github.com/KylinHe/aliensboot-core/log"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

var (
	stringType = reflect.TypeOf("")
)

func StringArray2Int32Array(array []string) []int32 {
	var result []int32
	for _, value := range array {
		new_value, _ := strconv.Atoi(value)
		result = append(result, int32(new_value))
	}
	return result
}

func ContainsInt32(value int32, array []int32) bool {
	if array == nil || len(array) == 0 {
		return false
	}
	for _, member := range array {
		if member == value {
			return true
		}
	}
	return false
}

func ContainsString(value string, array []string) bool {
	if array == nil || len(array) == 0 {
		return false
	}
	for _, member := range array {
		if member == value {
			return true
		}
	}
	return false
}

func StringToInt(value string) int {
	result, _ := strconv.Atoi(value)
	return result
}

func StringToInt32(value string) int32 {
	result, _ := strconv.Atoi(value)
	return int32(result)
}

func StringToInt64(value string) int64 {
	result, _ := strconv.ParseInt(value, 10, 64)
	return result
}

func StringToFloat32(value string) float32 {
	result, _ := strconv.ParseFloat(value, 64)
	return float32(result)
}

func StringToFloat64(value string) float64 {
	result, _ := strconv.ParseFloat(value, 64)
	return result
}

func Int32ToString(value int32) string {
	return strconv.Itoa(int(value))
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func FirstToUpper(str string) string {
	length := len(str)
	if length == 0 {
		return str
	}
	return strings.ToUpper(str[0:1]) + str[1:length]
	//temp := strings.Split(str, "_")
	//var upperStr string
	//for y := 0; y < len(temp); y++ {
	//	vv := []rune(temp[y])
	//	if y != 0 {
	//		for i := 0; i < len(vv); i++ {
	//			if i == 0 {
	//				vv[i] -= 32
	//				upperStr += string(vv[i]) // + string(vv[i+1])
	//			} else {
	//				upperStr += string(vv[i])
	//			}
	//		}
	//	}
	//}
	//return temp[0] + upperStr
}

func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//类型转换  byte slice to string
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//func ToStr(value interface{}, args ...int) (s string) {
//	switch v := value.(type) {
//	case bool:
//		s = strconv.FormatBool(v)
//	case float32:
//		s = strconv.FormatFloat(float64(v), 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 32))
//	case float64:
//		s = strconv.FormatFloat(v, 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 64))
//	case int:
//		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
//	case int8:
//		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
//	case int16:
//		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
//	case int32:
//		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
//	case int64:
//		s = strconv.FormatInt(v, argInt(args).Get(0, 10))
//	case uint:
//		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
//	case uint8:
//		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
//	case uint16:
//		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
//	case uint32:
//		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
//	case uint64:
//		s = strconv.FormatUint(v, argInt(args).Get(0, 10))
//	case string:
//		s = v
//	case []byte:
//		s = string(v)
//	default:
//		s = fmt.Sprintf("%v", v)
//	}
//	return s
//}

func Int(v interface{}) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case uint64:
		return int64(n)
	case int:
		return int64(n)
	case int32:
		return int64(n)
	case uint32:
		return int64(n)
	case int16:
		return int64(n)
	case uint16:
		return int64(n)
	case int8:
		return int64(n)
	case uint8:
		return int64(n)
	case float32:
		return int64(n)
	case float64:
		return int64(n)
	default:
		log.Panicf("Int: can not convert: %T %v", v, v)
		return 0
	}
}

func Float(v interface{}) float64 {
	switch n := v.(type) {
	case int64:
		return float64(n)
	case uint64:
		return float64(n)
	case int:
		return float64(n)
	case int32:
		return float64(n)
	case uint32:
		return float64(n)
	case int16:
		return float64(n)
	case uint16:
		return float64(n)
	case int8:
		return float64(n)
	case uint8:
		return float64(n)
	case float32:
		return float64(n)
	case float64:
		return n
	default:
		log.Panicf("Float: can not convert: %T %v", v, v)
		return 0
	}
}

func IntTuple(v interface{}) []int64 {
	if t, ok := v.([]int64); ok {
		return t
	}
	if t, ok := v.([]int); ok {
		ret := make([]int64, len(t))
		for i, v := range t {
			ret[i] = int64(v)
		}
		return ret
	}

	if t, ok := v.([]interface{}); ok {
		ret := make([]int64, len(t))
		for i, v := range t {
			ret[i] = Int(v)
		}
		return ret
	}
	log.Panicf("IntTuple: can not convert: %T %v", v, v)
	return nil
}

func FloatTuple(v interface{}) []float64 {
	if t, ok := v.([]float64); ok {
		return t
	}
	if t, ok := v.([]float32); ok {
		ret := make([]float64, len(t))
		for i, v := range t {
			ret[i] = float64(v)
		}
		return ret
	}

	if t, ok := v.([]interface{}); ok {
		ret := make([]float64, len(t))
		for i, v := range t {
			ret[i] = Float(v)
		}
		return ret
	}

	log.Panicf("FloatTuple: can not convert: %T %v", v, v)
	return nil
}

func String(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	val := reflect.ValueOf(v)
	return val.Convert(stringType).Interface().(string)
}

func MapStringAnything(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}

	if m, ok := v.(map[interface{}]interface{}); ok {
		m2 := make(map[string]interface{}, len(m))
		for k, v := range m {
			m2[k.(string)] = v
		}
		return m2
	}

	log.Panicf("MapStringAnything: can not convert: %T %v", v, v)
	return nil
}

// try to convert value to target type, panic if fail
func Convert(val interface{}, targetType reflect.Type) reflect.Value {
	value := reflect.ValueOf(val)
	if !value.IsValid() {
		return reflect.Zero(targetType)
	}

	valType := value.Type()
	if valType.ConvertibleTo(targetType) {
		return value.Convert(targetType)
	}

	//fmt.Printf("Value type is %v, emptyInterfaceType is %v, equals %v\n", valType, emptyInterfaceType, valType == emptyInterfaceType)
	interfaceVal := value.Interface()

	switch realVal := interfaceVal.(type) {
	case float64:
		return reflect.ValueOf(realVal).Convert(targetType)
	case []interface{}:
		// val is of type []interface{}, try to convert to typ
		sliceSize := len(realVal)
		targetSlice := reflect.MakeSlice(targetType, 0, sliceSize)
		elemType := targetType.Elem()
		for i := 0; i < sliceSize; i++ {
			targetSlice = reflect.Append(targetSlice, Convert(value.Index(i), elemType))
		}
		return targetSlice
	}

	panic(fmt.Errorf("convert from type %v to %v failed: %v", valType, targetType, value))
}
