package counter

type Counter interface {
	Count() int64
	UpdateCount() (int64, error)
}