package urlvalue

type UrlValuer interface {
	Get(value string) interface{}
	GetInt(value string) int64
	GetString(value string) string
}

type MapValuer interface {
	UrlValuer
	Map() map[string]interface{}
	Set(key string, value interface{})
	SetInt(key string, value int)
	SetString(key, value string)
}