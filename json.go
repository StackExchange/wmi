package wmi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/mattn/go-ole"
	"github.com/mattn/go-ole/oleutil"
)

// LoadJSON loads JSON data into dst
func LoadJSON(data []byte, dst interface{}) error {
	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	if len(r.Error) > 0 {
		return fmt.Errorf(r.Error)
	}
	m := r.Response
	dv := reflect.ValueOf(dst)
	if dv.Kind() != reflect.Ptr || dv.IsNil() {
		return ErrInvalidEntityType
	}
	dv = dv.Elem()
	mat, elemType := checkMultiArg(dv)
	if mat == multiArgTypeInvalid {
		return ErrInvalidEntityType
	}
	var errFieldMismatch error
	for _, v := range m {
		ev := reflect.New(elemType)
		if err := loadMap(ev.Interface(), v); err != nil {
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
	}
	return errFieldMismatch
}

// loadMap loads a map[string]interface{} into a struct pointer.
func loadMap(dst interface{}, src map[string]interface{}) (errFieldMismatch error) {
	v := reflect.ValueOf(dst).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		isPtr := f.Kind() == reflect.Ptr
		if isPtr {
			ptr := reflect.New(f.Type().Elem())
			f.Set(ptr)
			f = f.Elem()
		}
		n := v.Type().Field(i).Name
		if !f.CanSet() {
			return &ErrFieldMismatch{
				StructType: f.Type(),
				FieldName:  n,
				Reason:     "CanSet() is false",
			}
		}
		val, present := src[n]
		if !present {
			errFieldMismatch = &ErrFieldMismatch{
				StructType: f.Type(),
				FieldName:  n,
				Reason:     "no such struct field",
			}
			continue
		}
		switch reflect.ValueOf(val).Kind() {
		case reflect.Int64:
			iv := val.(int64)
			switch f.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				f.SetInt(iv)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				f.SetUint(uint64(iv))
			default:
				return &ErrFieldMismatch{
					StructType: f.Type(),
					FieldName:  n,
					Reason:     "not an integer class",
				}
			}
		case reflect.Float64:
			iv := val.(float64)
			switch f.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				f.SetInt(int64(iv))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				f.SetUint(uint64(iv))
			case reflect.Float32, reflect.Float64:
				f.SetFloat(iv)
			default:
				return &ErrFieldMismatch{
					StructType: f.Type(),
					FieldName:  n,
					Reason:     "not a number class",
				}
			}
		case reflect.String:
			sv := val.(string)
			iv, err := strconv.ParseInt(sv, 10, 64)
			switch f.Kind() {
			case reflect.String:
				f.SetString(sv)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if err != nil {
					return err
				}
				f.SetInt(iv)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if err != nil {
					return err
				}
				f.SetUint(uint64(iv))
			case reflect.Struct:
				switch f.Type() {
				case timeType:
					if len(sv) == 25 {
						sv = sv[:22] + "0" + sv[22:]
					}
					t, err := time.Parse("20060102150405.000000-0700", sv)
					if err != nil {
						return err
					}
					f.Set(reflect.ValueOf(t))
				}
			}
		case reflect.Bool:
			bv := val.(bool)
			switch f.Kind() {
			case reflect.Bool:
				f.SetBool(bv)
			default:
				return &ErrFieldMismatch{
					StructType: f.Type(),
					FieldName:  n,
					Reason:     "not a bool",
				}
			}
		default:
			typeof := reflect.TypeOf(val)
			if isPtr && typeof == nil {
				break
			}
			return fmt.Errorf("wmi: could not unmarshal %v with type %v", n, typeof)
		}
	}
	return errFieldMismatch
}

// QueryGen executes query and returns a map with keys of the columns slice.
func QueryGen(query string, columns []string, connectServerArgs ...interface{}) ([]map[string]interface{}, error) {
	var res []map[string]interface{}
	ole.CoInitializeEx(0, 0)
	unknown, err := oleutil.CreateObject("WbemScripting.SWbemLocator")
	if err != nil {
		return nil, err
	}
	defer unknown.Release()

	wmi, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return nil, err
	}
	defer wmi.Release()

	// service is a SWbemServices
	serviceRaw, err := oleutil.CallMethod(wmi, "ConnectServer", connectServerArgs...)
	if err != nil {
		return nil, err
	}

	service := serviceRaw.ToIDispatch()
	defer service.Release()

	// result is a SWBemObjectSet
	resultRaw, err := oleutil.CallMethod(service, "ExecQuery", query)
	if err != nil {
		return nil, err
	}
	result := resultRaw.ToIDispatch()
	defer result.Release()

	count, err := oleInt64(result, "Count")
	if err != nil {
		return nil, err
	}

	for i := int64(0); i < count; i++ {
		// item is a SWbemObject, but really a Win32_Process
		itemRaw, err := oleutil.CallMethod(result, "ItemIndex", i)
		if err != nil {
			return nil, err
		}
		item := itemRaw.ToIDispatch()
		defer item.Release()
		m := make(map[string]interface{})
		for _, c := range columns {
			prop, err := oleutil.GetProperty(item, c)
			if err != nil {
				return nil, err
			}
			m[c] = prop.Value()
		}
		res = append(res, m)
	}
	return res, nil
}

type WmiQuery struct {
	Query     string
	Namespace string
}

type Response struct {
	Error    string `json:",omitempty"`
	Response []map[string]interface{}
}
