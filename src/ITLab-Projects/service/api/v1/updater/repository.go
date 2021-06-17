package updater

import (
	"context"
)

type Repository interface {
	RealeseRepository
	RepoRepository
	IssueRepository
	LandingRepository
	MilestoneRepository
}

type RepoRepository interface{
	SaveReposAndSetDeletedUnfind(
		ctx context.Context,
		repos interface{},
	) error
}

type MilestoneRepository interface {
	SaveMilestonesAndSetDeletedUnfind(
		ctx context.Context,
		ms interface{},
	) error
}

type IssueRepository interface {
	SaveIssuesAndSetDeletedUnfind(
		ctx context.Context,
		is 	interface{},
	) error
}

type RealeseRepository interface {
	SaveRealeses(
		ctx context.Context,
		rs interface{},
	) error
}

type LandingRepository interface {
	SaveAndDeleteUnfindLanding(
		ctx context.Context,
		ls interface{},
	) error
}