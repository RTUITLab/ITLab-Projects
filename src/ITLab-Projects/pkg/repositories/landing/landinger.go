package landing

import (
	"github.com/ITLab-Projects/pkg/repositories/agregate"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
)

type LandingRepositorier interface {
	saver.SaverWithDelUpdate
	deleter.Deleter
	getter.Getter
	agregate.Agregater
}