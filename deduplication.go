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

func (u *Deduplication) Validation(args interface{}) Deduplication {
	switch reflect.TypeOf(args).Kind() {
	case reflect.Slice, reflect.Array:
		targetValue := reflect.ValueOf(args)
		u.Value = &targetValue
	default:
		u.Error = errors.New("invalid type")
	}

	return *u
}

func (u *Deduplication) duplicationStr() {
	encountered := make(map[string]struct{}, u.Value.Len())
	for i := 0; i < u.Value.Len(); i++ {
		element := u.Value.Index(i).String()
		if _, ok := encountered[element]; !ok {
			encountered[element] = struct{}{}
			u.SliceStr = append(u.SliceStr, element)
		}
	}
}

func (u *Deduplication) duplicationInt() {
	encountered := make(map[int]struct{}, u.Value.Len())
	for i := 0; i < u.Value.Len(); i++ {
		element := int(u.Value.Index(i).Int())
		if _, ok := encountered[element]; !ok {
			encountered[element] = struct{}{}
			u.SliceInt = append(u.SliceInt, element)
		}
	}
}

func (u Deduplication) Do(args interface{}) Deduplication {
	if u.Validation(args).Error != nil {
		return u
	}
	switch u.Value.Type().Elem().Kind() {
	case reflect.String:
		u.duplicationStr()
	case reflect.Int:
		u.duplicationInt()
	default:
		fmt.Println(u.Value.Type().Elem())
	}
	return u
}

func (u Deduplication) errorCheck(typ reflect.Kind) (err error) {
	//fmt.Println(u.SliceInt)
	if elem := u.Value.Type().Elem(); elem.Kind() != typ {
		err = errors.New("Invalid Type -> " + elem.String())
	} else if len(u.SliceStr) == 0 && len(u.SliceInt) == 0 {
		err = errors.New("0 Size slice: " + u.Value.Type().String())
	}
	return
}

func (u Deduplication) Int() (results []int, err error) {
	defer u.clear()
	if err = u.errorCheck(reflect.Int); err != nil {
		return
	}
	results = u.SliceInt
	return
}

func (u Deduplication) Str() (results []string, err error) {
	defer u.clear()
	if err = u.errorCheck(reflect.String); err != nil {
		return
	}
	results = u.SliceStr
	return
}

func (u *Deduplication) clear() {
	u.SliceInt = nil
	u.SliceStr = nil
	u.Value = nil
	
	u.Error = nil
}
