package repositories

import (
	"github.com/ITLab-Projects/pkg/repositories/database"
	"context"
	"errors"
	"time"

	"github.com/ITLab-Projects/pkg/repositories/estimates"
	"github.com/ITLab-Projects/pkg/repositories/functasks"
	"github.com/ITLab-Projects/pkg/repositories/issues"
	"github.com/ITLab-Projects/pkg/repositories/milestones"
	"github.com/ITLab-Projects/pkg/repositories/realeses"
	"github.com/ITLab-Projects/pkg/repositories/repos"
	"github.com/ITLab-Projects/pkg/repositories/tags"
	"github.com/ITLab-Projects/pkg/repositories/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type Repositories struct {
	Repo repos.	ReposRepositorier
	Milestone 	milestones.Milestoner
	Realese 	realeses.Realeser
	FuncTask 	functasks.FuncTaskRepositorier
	Estimate 	estimates.EstimateRepositorier
	Tag			tags.Tager
	Issue		issues.Issuer
}

type Config struct {
	DBURI string
}

func New(cfg *Config) (*Repositories, error) {
	URI, err := utils.GetDBURIWithoutName(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	client, err := mongo.NewClient(
		options.Client().
			ApplyURI(URI).
			SetMaxPoolSize(50).
			SetMaxConnIdleTime(0).
			SetLocalThreshold(10*time.Millisecond),
	)
	if err != nil {
		return nil, errors.New("Error on created client")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		return nil, errors.New("Error on connection")
	}

	ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	if err = client.Ping(ctx, nil); err != nil {
		return nil, errors.New("Error on ping")
	}

	dbName, err := utils.GetDbNameByReg(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	db := client.Database(
		dbName, 
		options.Database(),
	)

	reposCollection := db.Collection("repos")
	milestoneCollection := db.Collection("milestones")
	realeseCollection := db.Collection("realese")
	estimeCollection := db.Collection("estimate")
	funcTaskCollection := db.Collection("functask")
	tagCollection := db.Collection("tags")
	issueCollection := db.Collection("issues")
	return &Repositories{
		Repo: repos.New(reposCollection),
		Milestone: milestones.New(milestoneCollection),
		Realese: realeses.New(realeseCollection),
		Estimate: estimates.New(estimeCollection),
		FuncTask: functasks.New(funcTaskCollection),
		Tag: tags.New(tagCollection),
		Issue: issues.New(issueCollection),
	}, 
	nil
}

func NewDB(cfg Config) (*database.DB, error) {
	URI, err := utils.GetDBURIWithoutName(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	client, err := mongo.NewClient(
		options.Client().
			ApplyURI(URI).
			SetMaxPoolSize(50).
			SetMaxConnIdleTime(0).
			SetLocalThreshold(10*time.Millisecond),
	)
	if err != nil {
		return nil, errors.New("Error on created client")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		return nil, errors.New("Error on connection")
	}

	ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	if err = client.Ping(ctx, nil); err != nil {
		return nil, errors.New("Error on ping")
	}

	dbName, err := utils.GetDbNameByReg(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	db := client.Database(
		dbName, 
		options.Database(),
	)

	return database.New(db), nil
}