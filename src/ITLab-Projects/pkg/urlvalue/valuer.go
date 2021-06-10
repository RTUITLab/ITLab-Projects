package urlvalue

type UrlValuer interface {
	Get(value string) interface{}
	GetInt(value string) int64
	GetString(value string) string
}