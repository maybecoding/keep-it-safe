package vscfg

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Fn type for custom tag to string converter
type Fn func(tag reflect.StructTag) string

// Tag converts tag to string directly
func Tag(tagName string) Fn {
	return func(tag reflect.StructTag) string {
		return tag.Get(tagName)
	}
}

// Env converts tag to env value with name of tag
func Env(tagName string) Fn {
	return func(tag reflect.StructTag) string {
		name := tag.Get(tagName)
		if name != "" {
			return os.Getenv(name)
		}
		return ""
	}
}

// Flag converts tags for flag and flag usage to flag value
func Flag(tagFlag, tagUsage string) []Fn {
	flags := make(map[string]*string)
	once := sync.Once{}
	prepare := func(tag reflect.StructTag) string {
		flg := tag.Get(tagFlag)
		flgUsage := tag.Get(tagUsage)
		if flg == "" {
			return ""
		}
		for _, f := range strings.Split(flg, ",") {
			if f == "" {
				continue
			}
			v := new(string)
			flag.StringVar(v, f, "", flgUsage)
			// Берем именно исходный флаг, взятый из тега
			flags[flg] = v
		}
		return ""
	}

	getValues := func(tag reflect.StructTag) string {
		once.Do(func() {
			flag.Parse()
		})

		flg := tag.Get(tagFlag)

		if f, ok := flags[flg]; ok {
			return *f
		}
		return ""
	}

	return []Fn{prepare, getValues}
}

func FillByTags(s reflect.Value, fns ...Fn) error {
	for _, fn := range fns {
		err := walk(s, fn)
		if err != nil {
			return err
		}
	}
	return nil
}

// walk through struct fields without pointers
func walk(v reflect.Value, fn Fn) error {
	for i := 0; i < v.NumField(); i += 1 {
		fld := v.Field(i)
		fldT := v.Type().Field(i)

		// Если структура идем рекурсивно по ней
		if fld.Kind() == reflect.Struct {
			err := walk(fld, fn)
			if err != nil {
				return err
			}
			continue
		}

		// Проверяем наличие тега
		value := fn(fldT.Tag)
		if value == "" {
			continue
		}

		// Если слайс
		if fld.Kind() == reflect.Slice {
			values := strings.Split(value, ",")
			valuesLen := len(values)
			sliceV := reflect.MakeSlice(fld.Type(), valuesLen, valuesLen)
			for j, valuesItem := range values {
				err := fill(sliceV.Index(j).Addr(), valuesItem)
				if err != nil {
					return fmt.Errorf("error parse slice field %s #%d: %w", fldT.Name, j, err)
				}
			}
			fld.Set(sliceV)
			continue
		}

		// По всем остальным просто вызываем fill
		err := fill(fld.Addr(), value)
		if err != nil {
			return fmt.Errorf("error parse field %s: %w", fldT.Name, err)
		}
	}
	return nil
}

// fills pointer with parsed from data string value
func fill(vPtr reflect.Value, data string) error {
	v := vPtr.Elem()
	t := v.Type()

	// duration
	if t == reflect.TypeOf(time.Duration(0)) {
		durVal, err := time.ParseDuration(data)
		if err != nil {
			return fmt.Errorf("failed to parse \"%s\" as Duration", data)
		}
		v.SetInt(int64(durVal))
		return nil
	}

	// Проверяем строки и целые числа
	switch v.Kind() {
	case reflect.String:
		v.SetString(data)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(data, 10, t.Bits())
		if err != nil {
			return fmt.Errorf("failed to parse %s as int%d", data, t.Bits())
		}
		v.SetInt(intVal)
	case reflect.Bool:
		if data == "true" {
			v.SetBool(true)
		}
	default:
		return fmt.Errorf("failed to parse \"%s\", please implement handling of %v", data, v.Kind())

	}
	return nil
}
