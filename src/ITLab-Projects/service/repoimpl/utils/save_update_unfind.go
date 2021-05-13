package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"context"

	"github.com/ITLab-Projects/pkg/repositories/saver"
)

func SaveAndSetDeletedUnfind(
	ctx context.Context,
	Updater saver.SaverWithUpdate,
	values interface{},
) error {
	return Updater.SaveAndUpdatenUnfind(
		ctx,
		values,
		bson.M{"$set": bson.M{"deleted": true}},
	)
}