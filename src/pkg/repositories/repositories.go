package repositories

import (
	"errors"
	"github.com/ITLab-Projects/pkg/repositories/repos"
	"github.com/ITLab-Projects/pkg/repositories/utils"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
)


type Repositories struct {
	repos.ReposRepositorier
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
	
	return &Repositories{
		repos.New(reposCollection),
	}, 
	nil
}