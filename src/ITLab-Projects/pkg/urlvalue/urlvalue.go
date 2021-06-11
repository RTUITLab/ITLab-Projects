package urlvalue

import (
)

type UrlValue struct {
	Values map[string]interface{}
}

func New() *UrlValue {
	return &UrlValue{
		Values: map[string]interface{}{},
	}
}

// If value not exsist return nil
func (u *UrlValue) Get(value string) interface{} {
	v, find := u.Values[value]
	if !find {
		return nil
	}

	return v
}

// If value not find or not int return 0
func (u *UrlValue) GetInt(value string) int64 {
	intValue, ok := u.Get(value).(int64)
	if !ok {
		return 0
	}

	return intValue
}

// If value not find or not string return ""
func (u *UrlValue) GetString(value string) string {
	strValue, ok := u.Get(value).(string)
	if !ok {
		return ""
	}

	return strValue
}