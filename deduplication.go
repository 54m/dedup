package dedup

import (
	"errors"
	"fmt"
	"reflect"
)

type Deduplication struct {
	SliceStr []string
	SliceInt []int
	Value    *reflect.Value
	Error    error
}

func NewDeduplication() *Deduplication {
	return new(Deduplication)
}

func (d *Deduplication) validation(args interface{}) Deduplication {
	switch reflect.TypeOf(args).Kind() {
	case reflect.Slice, reflect.Array:
		targetValue := reflect.ValueOf(args)
		d.Value = &targetValue
	default:
		d.Error = errors.New("invalid type")
	}
	return *d
}

func (d *Deduplication) duplicationStr() {
	encountered := make(map[string]struct{}, d.Value.Len())
	for i := 0; i < d.Value.Len(); i++ {
		element := d.Value.Index(i).String()
		if _, ok := encountered[element]; !ok {
			encountered[element] = struct{}{}
			d.SliceStr = append(d.SliceStr, element)
		}
	}
}

func (d *Deduplication) duplicationInt() {
	encountered := make(map[int]struct{}, d.Value.Len())
	for i := 0; i < d.Value.Len(); i++ {
		element := int(d.Value.Index(i).Int())
		if _, ok := encountered[element]; !ok {
			encountered[element] = struct{}{}
			d.SliceInt = append(d.SliceInt, element)
		}
	}
}

func (d Deduplication) Do(args interface{}) Deduplication {
	if d.validation(args).Error != nil {
		return d
	}
	switch d.Value.Type().Elem().Kind() {
	case reflect.String:
		d.duplicationStr()
	case reflect.Int:
		d.duplicationInt()
	default:
		d.Error = errors.New("Unknown type: " + d.Value.Type().Elem().String())
		fmt.Println(d.Value.Type().Elem())
	}
	return d
}

func (d Deduplication) errorCheck(typ reflect.Kind) (err error) {
	if d.Error != nil {
		err = d.Error
	} else if elem := d.Value.Type().Elem(); elem.Kind() != typ {
		err = errors.New("Invalid Type -> " + elem.String())
	} else if len(d.SliceStr) == 0 && len(d.SliceInt) == 0 {
		err = errors.New("0 Size slice: " + d.Value.Type().String())
	}
	return
}

func (d Deduplication) Int() ([]int, error) {
	defer d.clear()
	if err := d.errorCheck(reflect.Int); err != nil {
		return nil, err
	}
	return d.SliceInt, nil
}

func (d Deduplication) Str() ([]string, error) {
	defer d.clear()
	if err := d.errorCheck(reflect.String); err != nil {
		return nil, err
	}
	return d.SliceStr, nil
}

func (d *Deduplication) clear() {
	d.SliceInt = nil
	d.SliceStr = nil
	d.Value = nil
	d.Error = nil
}
