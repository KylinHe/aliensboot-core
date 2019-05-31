/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/8/14
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func VisitTag(meta interface{}, tag string, callback func(fieldName string, tagValue string)) {
	dataValue := reflect.ValueOf(meta).Elem()
	dataType := reflect.TypeOf(meta).Elem()
	for i := 0; i < dataValue.NumField(); i++ {
		//field := dataValue.Field(i)
		fieldType := dataType.Field(i)
		tagValue := fieldType.Tag.Get(tag)

		callback(fieldType.Name, tagValue)
	}
}

func GetReflectValue(value reflect.Value) (interface{}, error) {
	data := value.Interface()
	if !IsStructType(value.Kind()) {
		return data, nil
	}
	if timeValue, ok := data.(time.Time); ok {
		//时间对象转时间戳
		return timeValue.Unix(), nil
	} else if bytesValue, ok := data.([]byte); ok {
		//字节数组转string
		return Bytes2Str(bytesValue), nil
	} else {
		//其他对象转json
		var jsonData, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}

		//TODO 支持protobuf
		return jsonData, nil
	}
}

//
func SetReflectValue(value reflect.Value, s string) error {
	if IsStructType(value.Kind()) {
		switch value.Interface().(type) {
		case time.Time:
			val, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(time.Unix(val, 0)))
		case []byte:
			value.SetBytes([]byte(s))
		default:
			var data = reflect.New(value.Type()).Interface()
			var err = json.Unmarshal([]byte(s), data)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(data).Elem())
		}
		return nil
	}

	switch value.Kind() {
	case reflect.String:
		value.SetString(s)
	case reflect.Bool:
		val, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		value.SetBool(val)
	case reflect.Float32:
		val, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return err
		}
		value.SetFloat(val)
	case reflect.Float64:
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		value.SetFloat(val)
	case reflect.Int, reflect.Int32:
		val, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return err
		}
		value.SetInt(val)
	case reflect.Int8:
		val, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return err
		}
		value.SetInt(val)
	case reflect.Int16:
		val, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return err
		}
		value.SetInt(val)
	case reflect.Int64:
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		value.SetInt(val)
	case reflect.Uint, reflect.Uint32:
		val, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return err
		}
		value.SetUint(val)
	case reflect.Uint8:
		val, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return err
		}
		value.SetUint(val)
	case reflect.Uint16:
		val, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return err
		}
		value.SetUint(val)
	case reflect.Uint64:
		val, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return err
		}
		value.SetUint(val)
	default:
		return fmt.Errorf("unkown-type :%v", reflect.TypeOf(value))
	}
	return nil
}

//是否扩展类型
func IsStructType(k reflect.Kind) bool {
	if k >= reflect.Array && k != reflect.String {
		return true
	}
	return false
}
