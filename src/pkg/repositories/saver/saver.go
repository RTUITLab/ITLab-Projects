package saver

import (
	"context"
)

type Saver interface {
	Save(interface{}) error
}

type SaverWithDelete interface {
	Saver
	SaveAndDeletedUnfind(context.Context, interface{}) error
}
