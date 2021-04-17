package saver

type Saver interface {
	Save(interface{}) error
}

// TODO make a standart realization with reflect.Type