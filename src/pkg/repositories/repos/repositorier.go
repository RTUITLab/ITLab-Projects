package repos

import (
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/counter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
)

type ReposRepositorier interface {
	saver.SaverWithDelUpdate
	counter.Counter
	getter.Getter
}