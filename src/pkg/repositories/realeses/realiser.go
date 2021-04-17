package realeses

import (
	"github.com/ITLab-Projects/pkg/repositories/saver"
	"github.com/ITLab-Projects/pkg/repositories/getter"
)

type Realeser interface {
	getter.GetOner
	saver.Saver
}