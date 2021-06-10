package landing

import (
	"context"

	model "github.com/ITLab-Projects/pkg/models/landing"
	"github.com/ITLab-Projects/pkg/repositories/agregate"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LandingByType struct {
	saver.SaverWithDelUpdate
	deleter.Deleter
	getter.Getter
	agregate.Agregater
	model mgm.Model
}

func NewByType(

) *LandingByType {
	l := &LandingByType{}

	m := model.Landing{}

	l.model = &m

	l.Deleter = deleter.NewDeleteByType(
		&m,
	)
	
	l.SaverWithDelUpdate = saver.NewSaverWithDelUpdateByType(
		m,
		&m,
		l.save,
		l.buildFilter,
	)

	l.Getter = getter.NewGetByType(
		&m,
	)

	l.Agregater = agregate.NewByType(
		&m,
	)

	return l
}

func (l *LandingByType) save(
	ctx context.Context,
	v	interface{},
) error {
	landing := getPointer(v)

	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"repo_id": landing.RepoId}

	if _, err := mgm.Coll(l.model).ReplaceOne(
		ctx,
		filter,
		landing,
		opts,
	); err != nil {
		return err
	}
	return nil
}

func (l *LandingByType) buildFilter(v interface{}) interface{} {
	var ids []uint64

	switch ls := v.(type) {
	case []model.Landing:
		for _, l := range ls {
			ids = append(ids, l.RepoId)
		}
	case []*model.Landing:
		for _, l := range ls {
			ids = append(ids, l.RepoId)
		}
	}

	return bson.M{"repo_id": bson.M{"$nin": ids}}
}

func getPointer(v interface{}) *model.Landing {
	switch landing := v.(type) {
	case model.Landing:
		return &landing
	case *model.Landing:
		return landing
	}
	panic("Unexcpected type")
}