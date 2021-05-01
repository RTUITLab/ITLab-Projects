package tags

import (
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
)

type Tager interface {
	saver.SaverWithDelete
	getter.Getter
	deleter.Deleter
}