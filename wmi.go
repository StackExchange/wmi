package wmi

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	ole "github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"
)

var l = log.New(os.Stdout, "", log.LstdFlags)

func init() {
	ole.CoInitializeEx(0, 0)
	// todo: determine when/if to call ole.CoUninitialize()
}

var (
	ErrInvalidEntityType = errors.New("wmi: invalid entity type")
)

func Query(query string, dst interface{}) error {
	dv := reflect.ValueOf(dst)
	if dv.Kind() != reflect.Ptr || dv.IsNil() {
		return ErrInvalidEntityType
	}
	dv = dv.Elem()
	mat, elemType := checkMultiArg(dv)
	if mat == multiArgTypeInvalid {
		return ErrInvalidEntityType
	}

	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return err
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer")
	if err != nil {
		return err
	}
	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, err := oleutil.CallMethod(service, "ExecQuery", query)
	if err != nil {
		return err
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	count, err := oleInt64(result, "Count")
	if err != nil {
		return err
	}

	var errFieldMismatch error
	for i := int64(0); i < count; i++ {
		err := func() error {
			// item is a SWbemObject, but really a Win32_Process
			itemRaw, err := oleutil.CallMethod(result, "ItemIndex", i)
			if err != nil {
				return err
			}
			item := itemRaw.ToIDispatch()
			defer item.Release()

			ev := reflect.New(elemType)
			if err = loadEntity(ev.Interface(), item); err != nil {
				if _, ok := err.(*ErrFieldMismatch); ok {
					// We continue loading entities even in the face of field mismatch errors.
					// If we encounter any other error, that other error is returned. Otherwise,
					// an ErrFieldMismatch is returned.
					errFieldMismatch = err
				} else {
					return err
				}
			}
			if mat != multiArgTypeStructPtr {
				ev = ev.Elem()
			}
			dv.Set(reflect.Append(dv, ev))
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return errFieldMismatch
}

// ErrFieldMismatch is returned when a field is to be loaded into a different
// type than the one it was stored from, or when a field is missing or
// unexported in the destination struct.
// StructType is the type of the struct pointed to by the destination argument.
type ErrFieldMismatch struct {
	StructType reflect.Type
	FieldName  string
	Reason     string
}

func (e *ErrFieldMismatch) Error() string {
	return fmt.Sprintf("wmi: cannot load field %q into a %q: %s",
		e.FieldName, e.StructType, e.Reason)
}

// loadEntity loads a SWbemObject into a struct pointer.
func loadEntity(dst interface{}, src *ole.IDispatch) error {
	v := reflect.ValueOf(dst).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		n := v.Type().Field(i).Name
		if !f.CanSet() {
			return &ErrFieldMismatch{
				StructType: f.Type(),
				FieldName:  n,
				Reason:     "CanSet() is false",
			}
		}
		prop, err := oleutil.GetProperty(src, n)
		if err != nil {
			return err
		}
		switch f.Kind() {
		case reflect.String:
			f.SetString(prop.ToString())
		default:
			l.Println("ignore:", n, f.Type())
		}
	}
	return nil
}

type multiArgType int

const (
	multiArgTypeInvalid multiArgType = iota
	multiArgTypeStruct
	multiArgTypeStructPtr
)

// checkMultiArg checks that v has type []S, []*S for some struct type S.
//
// It returns what category the slice's elements are, and the reflect.Type
// that represents S.
func checkMultiArg(v reflect.Value) (m multiArgType, elemType reflect.Type) {
	if v.Kind() != reflect.Slice {
		return multiArgTypeInvalid, nil
	}
	elemType = v.Type().Elem()
	switch elemType.Kind() {
	case reflect.Struct:
		return multiArgTypeStruct, elemType
	case reflect.Ptr:
		elemType = elemType.Elem()
		if elemType.Kind() == reflect.Struct {
			return multiArgTypeStructPtr, elemType
		}
	}
	return multiArgTypeInvalid, nil
}

func oleInt64(item *ole.IDispatch, prop string) (int64, error) {
	v, err := oleutil.GetProperty(item, prop)
	if err != nil {
		return 0, err
	}
	i := int64(v.Val)
	return i, nil
}
