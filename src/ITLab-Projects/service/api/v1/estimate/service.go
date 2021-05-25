package estimate

import (
	"context"
	"net/http"

	"github.com/ITLab-Projects/pkg/models/estimate"
)

type Service interface {
	AddEstimate(
		ctx context.Context,
		est *estimate.EstimateFile,
	) error

	DeleteEstimate(
		ctx 			context.Context,
		MilestoneID		uint64,
		// For delete in microfileserver
		r				*http.Request,
	) error
}