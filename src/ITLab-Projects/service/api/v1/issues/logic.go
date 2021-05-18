package issues

import (
	"context"
	"strings"

	"github.com/ITLab-Projects/pkg/models/milestone"
	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type service struct {
	repository 	Repository
	logger		log.Logger
}

func New(
	repository 	Repository,
	logger		log.Logger,
) *service {
	return &service{
		repository: repository,
		logger: logger,
	}
}

func (s *service) GetIssues(ctx context.Context, 
	start int64, count int64, 
	name string, tag string,
) ([]*milestone.IssuesWithMilestoneID, error) {
	filter, err := s.buildFilterForGetIssues(
		ctx,
		name,
		tag,
	)
	if err != nil {
		return nil, err
	}

	is, err := s.repository.GetFiltrSortedFromToIssues(
		ctx,
		filter,
		bson.D{ {"createdat", -1}, {"deleted", 1}},
		start,
		count,
	)
	if err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return nil, err
	}

	return is, nil
}

func (s *service) GetLabels(
	ctx context.Context,
) ([]interface{}, error) {
	labels, err := s.repository.GetLabelsNameFromOpenIssues(
		ctx,
	)
	if err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return nil, err
	}

	return labels, nil
}

func (s *service) buildFilterForGetIssues(
	ctx 		context.Context,
	name, tag 	string,
) (interface{}, error) {
	filter := bson.M{
		"state": "open",
	}

	if tag != "" {
		s.buildFilterByLabelTags(
			ctx,
			&filter,
			tag,
		)
	}

	if name != "" {
		if err := s.buildFilterByName(
			ctx,
			&filter,
			name,
		); err != nil {
			return nil, err
		}
	}

	return filter, nil
}

func (s *service) buildFilterByLabelTags(
	ctx 	context.Context,
	filter 	*bson.M,
	tag		string,
) {
	tags := strings.Split(tag, " ")

	(map[string]interface{})(*filter)["labels.name"] = bson.M{"$in": tags}
}

func(s *service) buildFilterByName(
	ctx		context.Context,
	filter	*bson.M,
	name	string,
) error {
	type IDs struct {
		ID 	uint64		`bson:"id"`
	}

	var reposID	[]*IDs
	if err := s.repository.GetReposAndScanTo(
		ctx,
		bson.M{
			"name": bson.M{"$regex": name, "$options": "-i"},
		},
		&reposID,
		options.Find(),
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return err
	}

	var milestonesID []IDs
	if err := s.repository.GetMilestonesAndScanTo(
		ctx,
		bson.M{
			"title": bson.M{"$regex": name, "$options": "-i"},
		},
		&milestonesID,
		options.Find(),
	); err == mongo.ErrNoDocuments {
		// Pass
	} else if err != nil {
		return err
	}

	var array bson.A
	array = append(
		array,
		bson.D {
			{"title", bson.M{"$regex": name, "$options": "-i"}},
		},
		bson.D{
			{"description", bson.M{"$regex": name, "$options": "-i"}},
		},
	)

	if len(reposID) > 0 {
		var ids []uint64
		for _, id := range reposID {
			ids = append(ids, id.ID)
		}

		array = append(
			array,
			bson.D {
				{"repo_id", bson.M{"$in": ids}},
			},
		)
	}

	if len(milestonesID) > 0 {
		var ids []uint64
		for _, id := range milestonesID {
			ids = append(ids, id.ID)
		}
		array = append(
			array, 
			bson.D {
				{"milestone_id", bson.M{"$in": ids}},
			},
		)
	}

	(map[string]interface{})(*filter)["$or"] = array

	return nil
}
