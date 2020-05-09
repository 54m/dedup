package dedup

import (
	"fmt"
	"reflect"

	"golang.org/x/xerrors"
)

// Deduplication duplicate exclusion
//   └── Supports type: []float32, []float64, []int, []int64, []uint, []uint64, []string
type Deduplication struct {
	SliceFloat32 []float32
	SliceFloat64 []float64
	SliceInt     []int
	SliceInt64   []int64
	SliceUint    []uint
	SliceUint64  []uint64
	SliceStr     []string

	Value *reflect.Value
	Error error
}

// NewDeduplication constructor
func NewDeduplication() *Deduplication {
	return new(Deduplication)
}

func (d *Deduplication) validation(args interface{}) *Deduplication {
	switch reflect.TypeOf(args).Kind() {
	case reflect.Slice, reflect.Array:
		targetValue := reflect.ValueOf(args)
		d.Value = &targetValue
	default:
		d.Error = xerrors.New("invalid type")
	}
	return d
}

func (d *Deduplication) duplicationFloat32() {
	encountered := make(map[float32]struct{}, d.Value.Len())
	for i := 0; i < d.Value.Len(); i++ {
		element := float32(d.Value.Index(i).Float())
		if _, ok := encountered[element]; !ok {
			encountered[element] = struct{}{}
			d.SliceFloat32 = append(d.SliceFloat32, element)
		}
	}
}

func (d *Deduplication) duplicationFloat64() {
	encountered := make(map[float64]struct{}, d.Value.Len())
	for i := 0; i < d.Value.Len(); i++ {
		element := d.Value.Index(i).Float()
		if _, ok := encountered[element]; !ok {
			encountered[element] = struct{}{}
			d.SliceFloat64 = append(d.SliceFloat64, element)
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

func (d *Deduplication) duplicationInt64() {
	encountered := make(map[int64]struct{}, d.Value.Len())
	for i := 0; i < d.Value.Len(); i++ {
		element := d.Value.Index(i).Int()
		if _, ok := encountered[element]; !ok {
			encountered[element] = struct{}{}
			d.SliceInt64 = append(d.SliceInt64, element)
		}
	}
}

func (d *Deduplication) duplicationUint() {
	encountered := make(map[uint]struct{}, d.Value.Len())
	for i := 0; i < d.Value.Len(); i++ {
		element := uint(d.Value.Index(i).Uint())
		if _, ok := encountered[element]; !ok {
			encountered[element] = struct{}{}
			d.SliceUint = append(d.SliceUint, element)
		}
	}
}

func (d *Deduplication) duplicationUint64() {
	encountered := make(map[uint64]struct{}, d.Value.Len())
	for i := 0; i < d.Value.Len(); i++ {
		element := d.Value.Index(i).Uint()
		if _, ok := encountered[element]; !ok {
			encountered[element] = struct{}{}
			d.SliceUint64 = append(d.SliceUint64, element)
		}
	}
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

// Do divide processing by type
func (d *Deduplication) Do(args interface{}) *Deduplication {
	if d.validation(args).Error != nil {
		return d
	}
	switch d.Value.Type().Elem().Kind() {
	case reflect.Float32:
		d.duplicationFloat32()
	case reflect.Float64:
		d.duplicationFloat64()
	case reflect.Int:
		d.duplicationInt()
	case reflect.Int64:
		d.duplicationInt64()
	case reflect.Uint:
		d.duplicationUint()
	case reflect.Uint64:
		d.duplicationUint64()
	case reflect.String:
		d.duplicationStr()
	default:
		d.Error = xerrors.New("Unknown type: " + d.Value.Type().Elem().String())
		fmt.Println(d.Value.Type().Elem())
	}
	return d
}

func (d *Deduplication) errorCheck(typ reflect.Kind) (err error) {
	element := d.Value.Type().Elem()
	switch {
	case d.Error != nil:
		err = d.Error
	case element.Kind() != typ:
		err = xerrors.Errorf("Invalid Type -> %s", element.String())
	case len(d.SliceFloat32) == 0 &&
		len(d.SliceFloat64) == 0 &&
		len(d.SliceInt) == 0 &&
		len(d.SliceInt64) == 0 &&
		len(d.SliceUint) == 0 &&
		len(d.SliceUint64) == 0 &&
		len(d.SliceStr) == 0:
		err = xerrors.Errorf("0 Size slice: %s", d.Value.Type().String())
	}
	if d.Error != nil {
		err = d.Error
	} else if elem := d.Value.Type().Elem(); elem.Kind() != typ {
		err = xerrors.Errorf("Invalid Type -> %s", elem.String())
	} else if len(d.SliceStr) == 0 && len(d.SliceInt) == 0 {
		err = xerrors.Errorf("0 Size slice: %s", d.Value.Type().String())
	}
	return
}

// Float32 returns the deduplicated slice
func (d *Deduplication) Float32() ([]float32, error) {
	defer d.clear()
	if err := d.errorCheck(reflect.Int); err != nil {
		return nil, err
	}
	return d.SliceFloat32, nil
}

// Float64 returns the deduplicated slice
func (d *Deduplication) Float64() ([]float64, error) {
	defer d.clear()
	if err := d.errorCheck(reflect.Int64); err != nil {
		return nil, err
	}
	return d.SliceFloat64, nil
}

// Int returns the deduplicated slice
func (d *Deduplication) Int() ([]int, error) {
	defer d.clear()
	if err := d.errorCheck(reflect.Int); err != nil {
		return nil, err
	}
	return d.SliceInt, nil
}

// Int64 returns the deduplicated slice
func (d *Deduplication) Int64() ([]int64, error) {
	defer d.clear()
	if err := d.errorCheck(reflect.Int64); err != nil {
		return nil, err
	}
	return d.SliceInt64, nil
}

// Uint returns the deduplicated slice
func (d *Deduplication) Uint() ([]uint, error) {
	defer d.clear()
	if err := d.errorCheck(reflect.Uint); err != nil {
		return nil, err
	}
	return d.SliceUint, nil
}

// Uint64 returns the deduplicated slice
func (d *Deduplication) Uint64() ([]uint64, error) {
	defer d.clear()
	if err := d.errorCheck(reflect.Uint64); err != nil {
		return nil, err
	}
	return d.SliceUint64, nil
}

// Str returns the deduplicated slice
func (d *Deduplication) Str() ([]string, error) {
	defer d.clear()
	if err := d.errorCheck(reflect.String); err != nil {
		return nil, err
	}
	return d.SliceStr, nil
}

func (d *Deduplication) clear() {
	d.SliceFloat32 = nil
	d.SliceFloat64 = nil
	d.SliceInt = nil
	d.SliceInt64 = nil
	d.SliceUint = nil
	d.SliceUint64 = nil
	d.SliceStr = nil
	d.Value = nil
	d.Error = nil
}
