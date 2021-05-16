package milestoneassets

import (
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
)

type AssetsRepositorier interface {
	saver.Saver
	getter.Getter
	deleter.Deleter
}