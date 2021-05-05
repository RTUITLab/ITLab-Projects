package saver

import (
	"context"
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Save struct {
	collection *mongo.Collection
	save 		saveWithReplaceFunc
	// type of elems
	t			reflect.Type
}

// TODO make check func before save

func (s *Save) Save(ctx context.Context, v interface{}) error {
	return s.save(ctx, v)
}

type saveWithReplaceFunc func(context.Context, interface{}) error

func NewSaver(
	c *mongo.Collection,
	// Should be not ptr and slice
	// reccomend use type like model.Model{}
	Type interface{},
	fun saveWithReplaceFunc,
) *Save {
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
	s.t = t

	saveFunc := func(ctx context.Context, v interface{}) error {
		typeOfV := reflect.TypeOf(v)

		if typeOfV.AssignableTo(t) {
			return fun(ctx, v)
		} else if typeOfV.AssignableTo(reflect.PtrTo(t)) {
			v = reflect.ValueOf(v).Elem().Interface()
			return fun(ctx, v)
		} else if typeOfV.AssignableTo(reflect.SliceOf(t)) {
			slice := reflect.ValueOf(v)
			for i := 0; i < slice.Len(); i++ {
				value := slice.Index(i).Interface()
				if err := fun(ctx, value); err != nil {
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

func (s *Save) makeSliceOfValue(value reflect.Value) interface{} {
	if value.Type().AssignableTo(
		reflect.PtrTo(s.t),
	) {
		value = value.Elem()
	}

	if value.Type().AssignableTo(
		s.t,
	) {
		slice := reflect.MakeSlice(s.t, 0, 1)
		slice = reflect.Append(slice, value)

		return slice.Interface()
	}

	return value.Interface()
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


type SaveWithDelete struct {
	s *Save
	f filterBuilder
}

type filterBuilder func(interface{}) (interface{})

func NewSaverWithDelete(
	collection *mongo.Collection,
	Type interface{},
	saveFunc saveWithReplaceFunc,
	f filterBuilder,
) *SaveWithDelete {
	s := &SaveWithDelete{
		f: f,
	}

	s.s = NewSaver(
		collection,
		Type,
		saveFunc,
	)


	return s
}

func (swd *SaveWithDelete) SaveAndDeletedUnfind(ctx context.Context, v interface{}) error {
	if err := swd.s.Save(ctx, v); err != nil {
		return err
	}

	value := reflect.ValueOf(v)
	if value.Type().AssignableTo(
		reflect.PtrTo(swd.s.t),
	) {
		value = value.Elem()
	}

	if value.Type().AssignableTo(
		swd.s.t,
	) {
		value = reflect.MakeSlice(swd.s.t, 0, 1)
		value = reflect.Append(value, reflect.ValueOf(v))
	}

	filter := swd.f(value.Interface())
	opts := options.Delete()

	if _, err := swd.s.collection.DeleteMany(
		ctx,
		filter,
		opts,
	); err != nil {
		return err
	}
	
	return nil
}

func (swd *SaveWithDelete) Save(ctx context.Context,v interface{}) error {
	return swd.s.Save(ctx, v)
}

type SaveWithUpdate struct {
	s *Save
	f filterBuilder
	Saver
}

func NewSaveWithUpdate(
	collection *mongo.Collection,
	Type interface{},
	saveFunc saveWithReplaceFunc,
	f filterBuilder,
) *SaveWithUpdate {
	s := &SaveWithUpdate{
		f: f,
	}

	s.s = NewSaver(
		collection,
		Type,
		saveFunc,
	)

	s.Saver = s.s

	return s
}

func(swu *SaveWithUpdate) SaveAndUpdatenUnfind(
	ctx context.Context, 
	v interface{},	// value that we  
	updateFilter interface{},	// filter where you change field
) error {
	if err := swu.s.Save(ctx, v); err != nil {
		return err
	}

	filter := swu.f(
		swu.s.makeSliceOfValue(
			reflect.ValueOf(v),
		),
	)

	if _, err := swu.s.collection.UpdateMany(
		ctx,
		filter,
		updateFilter,
		options.Update(),
	); err != nil {
		return err
	}

	return nil
}


type SaveWithDelUpdate struct {
	SaverWithDelUpdate
}

func NewSaverWithDelUpdate(
	collection *mongo.Collection,
	Type interface{},
	saveFunc saveWithReplaceFunc,
	f filterBuilder,
) *SaveWithDelUpdate {
	s := &SaveWithDelUpdate{}

	type _SaveWithDelUpdate struct {
		*SaveWithDelete
		*SaveWithUpdate
	}

	_s := &_SaveWithDelUpdate{
		SaveWithDelete: NewSaverWithDelete(
			collection,
			Type,
			saveFunc,
			f,
		),
		SaveWithUpdate: NewSaveWithUpdate(
			collection,
			Type,
			saveFunc,
			f,
		),
	}
	s.SaverWithDelUpdate = _s

	return s
}


