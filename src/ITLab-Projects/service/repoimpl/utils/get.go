package utils

import (
	"github.com/ITLab-Projects/pkg/conextvalue/chunck"
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ITLab-Projects/pkg/chunkresp"
	"github.com/ITLab-Projects/pkg/repositories/counter"
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

type ChucnkRepositorier interface {
	getter.GetAllerFiltered
	counter.Counter
}

// GetChuncked
func GetChuncked(
	ctx				context.Context,
	GetChunck		ChucnkRepositorier,
	filter,
	sort,
	value			interface{},
	from,
	to				int64,
) error {
	var Chunck chunkresp.ChunckWriter
	if _chunck, err := chunck.GetChunckFromContext(
		ctx,
	); err == chunck.ErrChunckNotFound {
		return err
	} else {
		Chunck = _chunck
	}
	
	Chunck.WriteStart(from)
	opts := options.Find().
				    SetSort(sort).
					SetSkip(from)
	if to != 0 {
		Chunck.WriteLimit(to)
		opts = opts.SetLimit(to)
	}

	var count int
	if err := GetChunck.GetAllFiltered(
		ctx,
		filter,
		func(c *mongo.Cursor) error {
			count = c.RemainingBatchLength()
			return c.All(
				ctx,
				value,
			)
		},
		opts,
	); err != nil {
		return err
	}

	Chunck.WriteCount(int64(count))

	total, err := GetChunck.UpdateCount()
	if err != nil {
		return err
	}
	Chunck.WriteTotalResult(total)

	if total - from - int64(count) <= 0 {
		Chunck.WriteHasMore(false)
	} else {
		Chunck.WriteHasMore(true)
	}

	return nil
}
