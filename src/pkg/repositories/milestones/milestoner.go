package milestones

import (
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/counter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
)

type Milestoner interface {
	saver.SaverWithDelete
	counter.Counter
	getter.Getter
}