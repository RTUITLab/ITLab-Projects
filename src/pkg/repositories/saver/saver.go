package saver

type Saver interface {
	Save(interface{}) error
}