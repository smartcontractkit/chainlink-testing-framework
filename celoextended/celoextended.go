package celoextended

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
)

// ConvertType copy of github.com/ethereum/go-ethereum/accounts/abi method
func ConvertType(in interface{}, proto interface{}) interface{} {
	protoType := reflect.TypeOf(proto)
	if reflect.TypeOf(in).ConvertibleTo(protoType) {
		return reflect.ValueOf(in).Convert(protoType).Interface()
	}
	// Use set as a last ditch effort
	if err := set(reflect.ValueOf(proto), reflect.ValueOf(in)); err != nil {
		panic(err)
	}
	return proto
}

func set(dst, src reflect.Value) error {
	dstType, srcType := dst.Type(), src.Type()
	switch {
	case dstType.Kind() == reflect.Interface && dst.Elem().IsValid():
		return set(dst.Elem(), src)
	case dstType.Kind() == reflect.Ptr && dstType.Elem() != reflect.TypeOf(big.Int{}):
		return set(dst.Elem(), src)
	case srcType.AssignableTo(dstType) && dst.CanSet():
		dst.Set(src)
	case dstType.Kind() == reflect.Slice && srcType.Kind() == reflect.Slice && dst.CanSet():
		return setSlice(dst, src)
	case dstType.Kind() == reflect.Array:
		return setArray(dst, src)
	case dstType.Kind() == reflect.Struct:
		return setStruct(dst, src)
	default:
		return fmt.Errorf("abi: cannot unmarshal %v in to %v", src.Type(), dst.Type())
	}
	return nil
}

// setSlice attempts to assign src to dst when slices are not assignable by default
// e.g. src: [][]byte -> dst: [][15]byte
// setSlice ignores if we cannot copy all of src' elements.
func setSlice(dst, src reflect.Value) error {
	slice := reflect.MakeSlice(dst.Type(), src.Len(), src.Len())
	for i := 0; i < src.Len(); i++ {
		if src.Index(i).Kind() == reflect.Struct {
			if err := set(slice.Index(i), src.Index(i)); err != nil {
				return err
			}
		} else {
			// e.g. [][32]uint8 to []celoextended.Hash
			if err := set(slice.Index(i), src.Index(i)); err != nil {
				return err
			}
		}
	}
	if dst.CanSet() {
		dst.Set(slice)
		return nil
	}
	return errors.New("Cannot set slice, destination not settable")
}

func setArray(dst, src reflect.Value) error {
	array := reflect.New(dst.Type()).Elem()
	min := src.Len()
	if src.Len() > dst.Len() {
		min = dst.Len()
	}
	for i := 0; i < min; i++ {
		if err := set(array.Index(i), src.Index(i)); err != nil {
			return err
		}
	}
	if dst.CanSet() {
		dst.Set(array)
		return nil
	}
	return errors.New("Cannot set array, destination not settable")
}

func setStruct(dst, src reflect.Value) error {
	for i := 0; i < src.NumField(); i++ {
		srcField := src.Field(i)
		dstField := dst.Field(i)
		if !dstField.IsValid() || !srcField.IsValid() {
			return fmt.Errorf("Could not find src field: %v value: %v in destination", srcField.Type().Name(), srcField)
		}
		if err := set(dstField, srcField); err != nil {
			return err
		}
	}
	return nil
}
