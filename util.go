package main

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

func genId() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}

func sub(i, j int) int {
	return i - j
}

func FormToStruct(ptr interface{}, vals map[string][]string, form string) (map[string]string, bool) {
	errors := make(map[string]string)
	formToStruct(ptr, vals, "", errors, form)
	return errors, len(errors) == 0
}

func formToStruct(ptr interface{}, vals map[string][]string, start string, errors map[string]string, form string) {
	var strct reflect.Value
	if reflect.TypeOf(ptr) == reflect.TypeOf(reflect.Value{}) {
		strct = ptr.(reflect.Value)
	} else {
		strct = reflect.ValueOf(ptr).Elem()
	}
	strctType := strct.Type()
	for i := 0; i < strct.NumField(); i++ {
		fld := strct.Field(i)
		name := ToLowerFirst(strctType.Field(i).Name)
		if ok, v := GetVal(start+name, vals); ok || fld.Kind() == reflect.Struct {
			required := strctType.Field(i).Tag.Get("required")
			if fld.Kind() != reflect.Struct && v == "" && (strings.Contains(required, "all") || (form != "" && strings.Contains(required, form))) {
				errors[form+"."+start+name] = parseErrorName(name) + " is required"
			}
			switch fld.Kind() {
			case reflect.String:
				fld.SetString(v)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				in, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					errors[form+"."+start+name] = parseErrorName(name) + " must be a number"
				}
				fld.SetInt(in)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				u, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					errors[form+"."+start+name] = parseErrorName(name) + " must be a number"
				}
				fld.SetUint(u)
			case reflect.Float32, reflect.Float64:
				f, err := strconv.ParseFloat(v, 64)
				if err != nil {
					errors[form+"."+start+name] = parseErrorName(name) + " must be a number"
				}
				fld.SetFloat(f)
			case reflect.Bool:
				b, err := strconv.ParseBool(v)
				if err != nil {
					errors[form+"."+start+name] = parseErrorName(name) + " must be either true or false"
				}
				fld.SetBool(b)
			case reflect.Slice:
				ss := reflect.MakeSlice(fld.Type(), 0, 0)
				fld.Set(genSlice(ss, v, start, name, form, errors))
			case reflect.Struct:
				st := reflect.Indirect(fld)
				formToStruct(st, vals, start+name+".", errors, form)
				fld.Set(st)
			}
		}
	}
}

func genSlice(sl reflect.Value, val, start, name, form string, errors map[string]string) reflect.Value {
	vs := strings.Split(val, ",")
	for _, v := range vs {
		switch sl.Type().String() {
		case "[]string":
			sl = reflect.Append(sl, reflect.ValueOf(v))
		case "[]int":
			in, err := strconv.ParseInt(v, 10, 0)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of numbers"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(int(in)))
		case "[]int8":
			in, err := strconv.ParseInt(v, 10, 8)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of numbers"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(int8(in)))
		case "[]int16":
			in, err := strconv.ParseInt(v, 10, 16)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of numbers"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(int16(in)))
		case "[]int32":
			in, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of numbers"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(int32(in)))
		case "[]int64":
			in, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of numbers"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(int64(in)))
		case "[]uint":
			in, err := strconv.ParseUint(v, 10, 0)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of numbers"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(uint(in)))
		case "[]uint8":
			in, err := strconv.ParseUint(v, 10, 8)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of numbers"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(uint8(in)))
		case "[]uint16":
			in, err := strconv.ParseUint(v, 10, 16)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of numbers"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(uint16(in)))
		case "[]uint32":
			in, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of numbers"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(uint32(in)))
		case "[]uint64":
			in, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of numbers"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(uint64(in)))
		case "[]float32":
			in, err := strconv.ParseFloat(v, 32)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of numbers"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(float32(in)))
		case "[]float64":
			in, err := strconv.ParseFloat(v, 64)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of numbers"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(float64(in)))
		case "[]bool":
			b, err := strconv.ParseBool(v)
			if err != nil {
				errors[form+"."+start+name] = parseErrorName(name) + " must be a list of either true or false"
				break
			}
			sl = reflect.Append(sl, reflect.ValueOf(b))
		}
	}
	return sl
}

func GetVal(key string, v map[string][]string) (bool, string) {
	if v == nil {
		return false, ""
	}
	vs, ok := v[key]
	if !ok || len(vs) == 0 {
		return false, ""
	}
	return true, vs[0]
}

func ToLowerFirst(s string) string {
	return strings.ToLower(string(s[0])) + s[1:len(s)]
}

func ToUpperFirst(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:len(s)]
}

func parseErrorName(s string) string {
	for i := 0; i < len(s); i++ {
		if 'A' <= s[i] && s[i] <= 'Z' {
			s = s[:i] + " " + s[i:]
			i++
		}
	}
	return ToUpperFirst(s)
}
