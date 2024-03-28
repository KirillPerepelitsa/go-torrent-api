package beencode

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
)

type keyValue struct {
	key   string
	value reflect.Value
}

type keyValueList []keyValue

func (a keyValueList) Len() int { return len(a) }

func (a keyValueList) Less(i, j int) bool { return a[i].key < a[j].key }

func (a keyValueList) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func Marshal(w io.Writer, data interface{}) (io.Writer, error) {
	return w, encode(w, reflect.ValueOf(data))
}

func encode(w io.Writer, v reflect.Value) (err error) {
	if !v.IsValid() {
		return errors.New("value is invalid")
	}
	switch v.Kind() {
	case reflect.String:
		_, err = fmt.Fprintf(w, "%d:%s", len(v.String()), v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err = fmt.Fprintf(w, "i%de", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		_, err = fmt.Fprintf(w, "i%de", v.Uint())
	case reflect.Array, reflect.Slice:
		_, err = fmt.Fprint(w, "l")
		if err != nil {
			return err
		}
		for i := 0; i < v.Len(); i++ {
			err := encode(w, v.Index(i))
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprint(w, "e")
		if err != nil {
			return err
		}
	case reflect.Map:
		key := v.Type().Key()
		if key.Kind() != reflect.String {
			return fmt.Errorf("map keys can be only string, %s provided", key.String())
		}
		_, err := fmt.Fprint(w, "d")
		if err != nil {
			return err
		}
		var kv keyValueList
		for _, key := range v.MapKeys() {
			kv = append(kv, keyValue{key: key.String(), value: v.MapIndex(key)})
		}
		err = writeSortedKeys(w, &kv)
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(w, "e")
		if err != nil {
			return err
		}
	case reflect.Struct:
		_, err := fmt.Fprint(w, "d")
		if err != nil {
			return err
		}
		var kv keyValueList
		for i := 0; i < v.NumField(); i++ {
			kv = append(kv, keyValue{key: v.Type().Field(i).Name, value: v.Field(i)})
		}
		err = writeSortedKeys(w, &kv)
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(w, "e")
		if err != nil {
			return err
		}
	case reflect.Interface:
		err = encode(w, v.Elem())
	default:
		err = fmt.Errorf("cannot marshal value of type %s", v.Type().String())
	}
	return
}

func writeSortedKeys(w io.Writer, list *keyValueList) (err error) {
	sort.Sort(list)
	for _, item := range *list {
		_, err = fmt.Fprintf(w, "%d:%s", len(item.key), item.key)
		err = encode(w, item.value)
		if err != nil {
			return err
		}
	}
	return
}
