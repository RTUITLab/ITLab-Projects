package issues

import (
	"github.com/ITLab-Projects/pkg/repositories/agregate"
	"github.com/ITLab-Projects/pkg/repositories/counter"
	"github.com/ITLab-Projects/pkg/repositories/deleter"
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
)

type Issuer interface {
	saver.SaverWithDelUpdate
	getter.Getter
	deleter.Deleter
	counter.Counter
	agregate.Agregater
}