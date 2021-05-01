package typechecker

import (
	log "github.com/sirupsen/logrus"
	"fmt"
	"reflect"
)

type TypeChecker func(interface{}) error

// func NewWithPtrAndSlice(t reflect.Type) TypeChecker {
// 	if t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 	}

// 	return func(i interface{}) error {
// 		switch
// 	}
// }

func NewSingle(t reflect.Type) TypeChecker {
	switch t.Kind() {
	case reflect.Struct:
		break
	default:
		log.WithFields(log.Fields{
			"package": "typechecker",
			"func": "NewSingle",
			"err": fmt.Sprintf("You give %s of %s expect %s", t.Kind(), t, reflect.Struct),
		}).Panic()
	}

	slice := reflect.SliceOf(t)
	ptrToSlice := reflect.PtrTo(slice)

	return func(i interface{}) error {
		var err error = nil

		if !reflect.TypeOf(i).AssignableTo(ptrToSlice) {
			err = fmt.Errorf("Uknown type: %T Expected: %s", i, ptrToSlice)
		}

		return err
	}
}

func NewSingleByInterface(t interface{}) TypeChecker {
	return NewSingle(reflect.TypeOf(t))
}