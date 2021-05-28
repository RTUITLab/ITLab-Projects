package saver

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"

	"reflect"
	"github.com/Kamva/mgm"
)

type SaveByType struct {
	_type 		mgm.Model
	save 		saveWithReplaceFunc
	t			reflect.Type
	tPtr		reflect.Type
}

func (s *SaveByType) Save(ctx context.Context, v interface{}) error {
	return s.save(ctx, v)
}

func NewSaverByType(
	// Should be not ptr and slice
	// reccomend use type like model.Model{}
	Type interface{},
	ModelType mgm.Model,
	fun saveWithReplaceFunc,
) *SaveByType {
	s := &SaveByType{
		_type: ModelType,
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
	s.tPtr = reflect.PtrTo(t)

	saveFunc := func(ctx context.Context, v interface{}) error {
		typeOfV := reflect.TypeOf(v)

		if typeOfV.AssignableTo(t) {
			return fun(ctx, v)
		} else if typeOfV.AssignableTo(s.tPtr) {
			return fun(ctx, v)
		} else if typeOfV.AssignableTo(reflect.SliceOf(t)) {
			slice := reflect.ValueOf(v)
			for i := 0; i < slice.Len(); i++ {
				value := slice.Index(i).Addr().Interface()
				if err := fun(ctx, value); err != nil {
					return err
				}
			}
		} else if typeOfV.AssignableTo(reflect.SliceOf(s.tPtr)) {
			slice := reflect.ValueOf(v)
			for i := 0; i < slice.Len(); i++ {
				value := slice.Index(i).Interface()
				if err := fun(ctx, value); err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf(
				"Unexpected Type %T, Expected %s or %s or %s or %s",
				v, t, reflect.PtrTo(t), reflect.SliceOf(t), reflect.SliceOf(s.tPtr),
			)
		}

		return nil
	}

	s.save = saveFunc

	return s
}

func (s *SaveByType) makeSliceOfValue(value reflect.Value) interface{} {
	if value.Type().AssignableTo(
		s.t,
	) {
		value = value.Addr()
	}

	if value.Type().AssignableTo(
		s.tPtr,
	) {
		slice := reflect.MakeSlice(s.t, 0, 1)
		slice = reflect.Append(slice, value)

		return slice.Interface()
	}

	return value.Interface()
}

type SaveWithDeleteByType struct {
	s *SaveByType
	f filterBuilder
}

func NewSaverWithDeleteByType(
	Type interface{},
	ModelType mgm.Model,
	saveFunc saveWithReplaceFunc,
	f filterBuilder,
) *SaveWithDeleteByType {
	s := &SaveWithDeleteByType{
		f: f,
	}

	s.s = NewSaverByType(
		Type,
		ModelType,
		saveFunc,
	)


	return s
}

func (swd *SaveWithDeleteByType) SaveAndDeletedUnfind(
	ctx context.Context, 
	v interface{},
) error {
	if err := swd.s.Save(ctx, v); err != nil {
		return err
	}

	value := reflect.ValueOf(v)
	if value.Type().AssignableTo(
		swd.s.t,
	) {
		value = value.Addr()
	}

	if value.Type().AssignableTo(
		swd.s.tPtr,
	) {
		value = reflect.MakeSlice(swd.s.tPtr, 0, 1)
		value = reflect.Append(value, reflect.ValueOf(v))
	}

	filter := swd.f(value.Interface())
	opts := options.Delete()

	if _, err := mgm.Coll(swd.s._type).DeleteMany(
		ctx,
		filter,
		opts,
	); err != nil {
		return err
	}
	
	return nil
}

func (swd *SaveWithDeleteByType) Save(ctx context.Context,v interface{}) error {
	return swd.s.Save(ctx, v)
}

type SaveWithUpdateByType struct {
	s *SaveByType
	f filterBuilder
	Saver
}

func NewSaveWithUpdateByType(
	Type interface{},
	ModelType mgm.Model,
	saveFunc saveWithReplaceFunc,
	f filterBuilder,
) *SaveWithUpdateByType {
	s := &SaveWithUpdateByType{
		f: f,
	}

	s.s = NewSaverByType(
		Type,
		ModelType,
		saveFunc,
	)

	s.Saver = s.s

	return s
}

func(swu *SaveWithUpdateByType) SaveAndUpdatenUnfind(
	ctx context.Context, 
	v interface{},	// value that we  save
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

	if _, err := mgm.Coll(swu.s._type).UpdateMany(
		ctx,
		filter,
		updateFilter,
		options.Update(),
	); err != nil {
		return err
	}

	return nil
}


type SaveWithDelUpdateByType struct {
	SaverWithDelUpdate
}

func NewSaverWithDelUpdateByType(
	Type interface{},
	ModelType	mgm.Model,
	saveFunc saveWithReplaceFunc,
	f filterBuilder,
) *SaveWithDelUpdateByType {
	s := &SaveWithDelUpdateByType{}

	type _SaveWithDelUpdate struct {
		*SaveWithDeleteByType
		*SaveWithUpdateByType
	}

	_s := &_SaveWithDelUpdate{
		SaveWithDeleteByType: NewSaverWithDeleteByType(
			Type,
			ModelType,
			saveFunc,
			f,
		),
		SaveWithUpdateByType: NewSaveWithUpdateByType(
			Type,
			ModelType,
			saveFunc,
			f,
		),
	}
	s.SaverWithDelUpdate = _s

	return s
}