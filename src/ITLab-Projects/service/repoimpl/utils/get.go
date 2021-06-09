package utils

import (
	"go.mongodb.org/mongo-driver/mongo"
	"context"

	"github.com/ITLab-Projects/pkg/repositories/getter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAndScanTo(
	ctx		context.Context,
	Get		getter.GetAllerFiltered,
	filter	interface{},
	value	interface{},
	options	...*options.FindOptions,
) error {
	return Get.GetAllFiltered(
		ctx,
		filter,
		func(c *mongo.Cursor) error {
			if c.RemainingBatchLength() == 0 {
				return mongo.ErrNoDocuments
			}
			return c.All(
				ctx,
				value,
			)
		},
		options...
	)
}

func GetFiltrSortFromToAndScan(
	ctx			context.Context,
	Get			getter.GetAllerFiltered,
	filter		interface{},
	sort		interface{},
	value		interface{},
	from, to	int64,
) error {
	opts := options.Find().
				SetSort(sort).
				SetSkip(from)
	if to != 0 {
		opts = opts.SetLimit(to)
	}
	
	return GetAndScanTo(
		ctx,
		Get,
		filter,
		value,
		opts,
	)
}
