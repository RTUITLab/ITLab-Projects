package functasks

import (
	"github.com/ITLab-Projects/pkg/repositories/getter"
	"github.com/ITLab-Projects/pkg/repositories/saver"
)

type FuncTaskRepositorier interface {
	saver.Saver
	getter.Getter
	Delete(MilestoneID uint64) error
}