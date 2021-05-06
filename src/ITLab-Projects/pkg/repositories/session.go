package repositories

import (
	"context"

	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetMongoSessionContext(ctx context.Context) (mongo.SessionContext, error) {
	_, client, _, err := mgm.DefaultConfigs()
	if err != nil {
		return nil, err
	}
	
	sess, err := client.StartSession()
	if err != nil {
		return nil, err
	}

	return mongo.NewSessionContext(
		ctx,
		sess,
	), nil

}