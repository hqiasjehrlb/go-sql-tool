package sqltool

import "reflect"

type ExtFieldMap map[string]any

// To get reflect value by key
func (m ExtFieldMap) ValueOf(key string) reflect.Value {
	r := reflect.ValueOf(m[key])
	if r.Kind() == reflect.Pointer {
		r = r.Elem()
	}
	return r
}

// To get integer value in map by key
func (m ExtFieldMap) GetInt(key string) *int64 {
	var i64 int64
	t := reflect.TypeOf(i64)
	r := m.ValueOf(key)
	if r.IsValid() && r.CanConvert(t) {
		i64 = r.Convert(t).Int()
		return &i64
	}
	return nil
}

// To get string value in map by key
func (m ExtFieldMap) GetString(key string) *string {
	var str string
	t := reflect.TypeOf(str)
	r := m.ValueOf(key)
	if r.IsValid() && r.CanConvert(t) {
		str = r.Convert(t).String()
		return &str
	}
	return nil
}

// To get float value in map by key
func (m ExtFieldMap) GetFloat(key string) *float64 {
	var f64 float64
	t := reflect.TypeOf(f64)
	r := m.ValueOf(key)
	if r.CanConvert(t) {
		f64 = r.Convert(t).Float()
		return &f64
	}
	return nil
}

// To get bool value in map by key
func (m ExtFieldMap) GetBool(key string) *bool {
	var b bool
	t := reflect.TypeOf(b)
	r := m.ValueOf(key)
	if r.CanConvert(t) {
		b = r.Convert(t).Bool()
		return &b
	}
	return nil
}
