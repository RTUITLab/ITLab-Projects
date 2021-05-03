package database

import (
	"github.com/ITLab-Projects/pkg/repositories/estimates"
	"github.com/ITLab-Projects/pkg/repositories/functasks"
	"github.com/ITLab-Projects/pkg/repositories/issues"
	"github.com/ITLab-Projects/pkg/repositories/milestones"
	"github.com/ITLab-Projects/pkg/repositories/realeses"
	"github.com/ITLab-Projects/pkg/repositories/repos"
	"github.com/ITLab-Projects/pkg/repositories/tags"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	db *mongo.Database
}

func New(db *mongo.Database) *DB {
	return &DB{
		db,
	}
}

func (d *DB)TagsRepository() tags.Tager {
	return tags.New(
		d.db.Collection("tags"),
	)
}

func (d *DB)EstimateRepository() estimates.EstimateRepositorier {
	return estimates.New(
		d.db.Collection("estimate"),
	)
}

func (d *DB)FuncTaskRepository() functasks.FuncTaskRepositorier {
	return functasks.New(
		d.db.Collection("functask"),
	)
}

func (d *DB)IssuesRepository() issues.Issuer {
	return issues.New(
		d.db.Collection("issues"),
	)
}

func (d *DB)MilestoneRepository() milestones.Milestoner {
	return milestones.New(
		d.db.Collection("milestone"),
	)
}

func (d *DB)RealeseRepository() realeses.Realeser {
	return realeses.New(
		d.db.Collection("realese"),
	)
}

func (d *DB)ReposRepository() repos.ReposRepositorier {
	return repos.New(
		d.db.Collection("repos"),
	)
}