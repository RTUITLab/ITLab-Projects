package urlvalue

// This not concurency safety
type UrlValue map[string]interface{}

func New() *UrlValue {
	return &UrlValue{
	}
}

// If value not exsist return nil
func (u UrlValue) Get(value string) interface{} {
	v, find := u[value]
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

func (u UrlValue) Set(key string, value interface{}) {
	u[key] = value
}

func (u *UrlValue) SetInt(key string, value int) {
	u.Set(key, value)
}

func (u *UrlValue) SetString(key, value string) {
	u.Set(key, value)
}

func (u UrlValue) Map() map[string]interface{} {
	return u
}

func FromMap(m map[string]interface{}) *UrlValue {
	valuer := (UrlValue)(m)
	return &valuer
}