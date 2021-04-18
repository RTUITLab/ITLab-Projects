package saver

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
)

type Save struct {
	collection *mongo.Collection
	save saveWithReplaceFunc
}

func (s *Save) Save(v interface{}) error {
	return s.save(v)
}

type saveWithReplaceFunc func(interface{}) error

func New(
	c *mongo.Collection,
	// Should be not ptr and slice
	// reccomend use type like model.Model{}
	Type interface{},
	fun saveWithReplaceFunc,
) Saver {
	s := &Save{
		collection: c,
	}

	t, err := getTypeOrPanic(Type)
	if err != nil {
		log.WithFields(
			log.Fields{
				"package": "saver",
				"func": "New",
				"err": err,
			},
		).Panic()
	}

	saveFunc := func(v interface{}) error {
		typeOfV := reflect.TypeOf(v)

		if typeOfV.AssignableTo(t) {
			return fun(v)
		} else if typeOfV.AssignableTo(reflect.PtrTo(t)) {
			v = reflect.ValueOf(v).Elem().Interface()
			return fun(v)
		} else if typeOfV.AssignableTo(reflect.SliceOf(t)) {
			slice := reflect.ValueOf(v)
			for i := 0; i < slice.Len(); i++ {
				value := slice.Index(i).Interface()
				if err := fun(value); err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf(
				"Unexpected Type %T, Expected %s or %s or %s",
				v, t, reflect.PtrTo(t), reflect.SliceOf(t),
			)
		}

		return nil
	}

	s.save = saveFunc

	return s
}

func getTypeOrPanic(Type interface{}) (reflect.Type, error) {
	t := reflect.TypeOf(Type)
	switch t.Kind() {
	case reflect.Struct:
		break
	default:
		return nil, fmt.Errorf("You give %s of %s expect %s", t.Kind(), t, reflect.Struct)
	}

	return t, nil
}