package repositories

import (
	"github.com/ITLab-Projects/pkg/repositories/tags"
	"github.com/ITLab-Projects/pkg/repositories/estimates"
	"github.com/ITLab-Projects/pkg/repositories/functasks"
	"github.com/ITLab-Projects/pkg/repositories/realeses"
	"github.com/ITLab-Projects/pkg/repositories/milestones"
	"errors"
	"github.com/ITLab-Projects/pkg/repositories/repos"
	"github.com/ITLab-Projects/pkg/repositories/utils"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
)


type Repositories struct {
	Repo repos.	ReposRepositorier
	Milestone 	milestones.Milestoner
	Realese 	realeses.Realeser
	FuncTask 	functasks.FuncTaskRepositorier
	Estimate 	estimates.EstimateRepositorier
	Tag			tags.Tager
}

type Config struct {
	DBURI string
}

func New(cfg *Config) (*Repositories, error) {
	URI, err := utils.GetDBURIWithoutName(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
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

	reposCollection := client.Database(dbName).Collection("repos")
	milestoneCollection := client.Database(dbName).Collection("milestones")
	realeseCollection := client.Database(dbName).Collection("realese")
	estimeCollection := client.Database(dbName).Collection("estimate")
	funcTaskCollection := client.Database(dbName).Collection("functask")
	tagCollection := client.Database(dbName).Collection("tags")
	return &Repositories{
		Repo: repos.New(reposCollection),
		Milestone: milestones.New(milestoneCollection),
		Realese: realeses.New(realeseCollection),
		Estimate: estimates.New(estimeCollection),
		FuncTask: functasks.New(funcTaskCollection),
		Tag: tags.New(tagCollection),
	}, 
	nil
}