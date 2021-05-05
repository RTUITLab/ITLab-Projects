package saver

import (
	"context"
)

type Saver interface {
	Save(context.Context, interface{}) error
}

type SaverWithDelete interface {
	Saver
	SaveAndDeletedUnfind(context.Context, interface{}) error
}

type SaverWithUpdate interface {
	Saver
	SaveAndUpdatenUnfind(
		ctx context.Context, 
		filter interface{},	// value that we  
		updateFilter interface{},	// filter where you change field
	) error
}


type SaverWithDelUpdate interface {
	SaverWithDelete
	SaverWithUpdate
}