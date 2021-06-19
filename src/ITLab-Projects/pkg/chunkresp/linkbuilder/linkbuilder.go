package linkbuilder

import (
	"fmt"
	"net/url"

	"github.com/ITLab-Projects/pkg/chunkresp"
)

type LinkBuilder struct {
	StartKey		string
	LimitKey		string
	links 			*chunkresp.Links
}

func New(
	StartKey,
	LimitKey 	string,
) *LinkBuilder {
	return &LinkBuilder{
		StartKey: StartKey,
		LimitKey: LimitKey,
	}
}

type Builder interface {
	Build(
		Chunck 	*chunkresp.ChunkResp,
		URL		*url.URL,
	) *chunkresp.Links
}

func (l *LinkBuilder) GetLinks(
) *chunkresp.Links {
	return l.links
}

func (l *LinkBuilder) BuildSelf(
	Chunck 	*chunkresp.ChunkResp,
	URL		*url.URL,
) *chunkresp.Links {
	if l.links == nil {
		l.links = chunkresp.NewLink()
	}
	values := URL.Query()

	values.Set(l.LimitKey, fmt.Sprint(Chunck.Limit)) // Because limit can be handled by endpoint

	return l.links.AddSelf(
		fmt.Sprintf(
			"%s?%s",
			URL.Path,
			chunkresp.ParseUrlQueryToHref(
				values,
			),
		),
	)
}

// If can't build prev return nil
func (l *LinkBuilder) BuildPrev(
	Chunck 	*chunkresp.ChunkResp,
	URL		*url.URL,
) *chunkresp.Links {
	if l.links == nil {
		l.links = chunkresp.NewLink()
	}

	if Chunck.Start == 0 {
		return nil
	}

	values := URL.Query()
	var prevStart int64
	if Chunck.Start >= Chunck.Limit {
		prevStart = Chunck.Start - Chunck.Limit
	} else {
		values.Set(l.LimitKey, fmt.Sprint(Chunck.Start))
		prevStart = 0
	}

	values.Set(l.StartKey, fmt.Sprint(prevStart))

	return l.links.AddPrev(
		fmt.Sprintf(
			"%s?%s",
			URL.Path,
			chunkresp.ParseUrlQueryToHref(
				values,
			),
		),
	)
}

// If don't have more return nil
func (l *LinkBuilder) BuildNext(
	Chunck 	*chunkresp.ChunkResp,
	URL		*url.URL,
) *chunkresp.Links {
	if l.links == nil {
		l.links = chunkresp.NewLink()
	}

	if !Chunck.HasMore {
		return nil
	}

	values := URL.Query()
	values.Set(l.LimitKey, fmt.Sprint(Chunck.Limit)) // Because limit can be handled by endpoint
	nextStart := Chunck.Start + Chunck.Count

	values.Set(l.StartKey, fmt.Sprint(nextStart))

	return l.links.AddNext(
		fmt.Sprintf(
			"%s?%s",
			URL.Path,
			chunkresp.ParseUrlQueryToHref(
				values,
			),
		),
	)
}

func (l *LinkBuilder) Build(
	Chunck 	*chunkresp.ChunkResp,
	URL		*url.URL,
) *chunkresp.Links {
	l.BuildSelf(
		Chunck,
		URL,
	)
	l.BuildPrev(
		Chunck,
		URL,
	)
	l.BuildNext(
		Chunck,
		URL,
	)

	return l.links
}