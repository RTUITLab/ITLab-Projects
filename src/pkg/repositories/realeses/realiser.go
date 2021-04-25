package realeses

import (
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
)

type Realeser interface {
	getter.GetOner
	saver.SaverWithDelete
	deleter.Deleter
}