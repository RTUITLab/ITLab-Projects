package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/Kamva/mgm"

	"github.com/ITLab-Projects/pkg/repositories/estimates"
	"github.com/ITLab-Projects/pkg/repositories/functasks"
	"github.com/ITLab-Projects/pkg/repositories/issues"
	"github.com/ITLab-Projects/pkg/repositories/landing"
	"github.com/ITLab-Projects/pkg/repositories/milestones"
	"github.com/ITLab-Projects/pkg/repositories/realeses"
	"github.com/ITLab-Projects/pkg/repositories/repos"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type Repositories struct {
	Repo 		repos.ReposRepositorier
	Milestone 	milestones.Milestoner
	Realese 	realeses.Realeser
	FuncTask 	functasks.FuncTaskRepositorier
	Estimate 	estimates.EstimateRepositorier
	Issue		issues.Issuer
	Landing		landing.LandingRepositorier
}

type Config struct {
	DBURI string
}

func New(cfg *Config) (*Repositories, error) {
	client, err := mongo.NewClient(
		options.Client().
			ApplyURI(cfg.DBURI).
			SetMaxPoolSize(50).
			SetMaxConnIdleTime(0).
			SetLocalThreshold(10*time.Millisecond),
	)
	if err != nil {
		return nil, errors.New("Error on created client")
	}
	defer client.Disconnect(context.Background())
	

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		return nil, errors.New("Error on connection")
	}

	ctx, _ = context.WithTimeout(context.Background(), 60*time.Second)
	if err = client.Ping(ctx, nil); err != nil {
		return nil, errors.New("Error on ping")
	}
	constr, err := connstring.Parse(cfg.DBURI)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI: err: %v", err)
	}
	
	dbName := constr.Database
	
	if err := mgm.SetDefaultConfig(
		nil,
		dbName,
		options.Client().ApplyURI(cfg.DBURI),
	); err != nil {
		return nil, err
	}

	
	return &Repositories{
		Repo: repos.NewByType(),
		Milestone: milestones.NewByType(),
		Realese: realeses.NewByType(),
		Estimate: estimates.NewByType(),
		FuncTask: functasks.NewByType(),
		Issue: issues.NewByType(),
		Landing: landing.NewByType(),
	}, 
	nil
}