package chunkresp

import "net/url"

type Chuncker interface {
	ChunckWriter
	SetItems(
		list interface{},
	)
}

type ChunckWriter interface {
	WriteHasMore(b bool)
	WriteLimit(limit int64)
	WriteCount(count int64)
	WriteTotalResult(res int64)
	WriteStart(start int64)
	WriteLinks(links *Links)
}

type ChunkResp struct {
	HasMore		bool	`json:"has_more"`
	Limit		int64	`json:"limit"`
	Count		int64	`json:"count"`
	TotalResult	int64	`json:"total_result"`
	Links		*Links	`json:"links"`
	Start		int64	`json:"start"`
}

func (c *ChunkResp) WriteHasMore(b bool) {
	c.HasMore = b
}

func (c *ChunkResp) WriteLimit(limit int64) {
	c.Limit = limit
}

func (c *ChunkResp) WriteCount(count int64) {
	c.Count = count
}

func (c *ChunkResp) WriteTotalResult(res int64) {
	c.TotalResult = res
}

func (c *ChunkResp) WriteStart(start int64) {
	c.Start = start
}

func (c *ChunkResp) WriteLinks(links *Links) {
	c.Links = links
}


type Link struct {
	Rel		string		`json:"rel"`
	Href	string		`json:"href"`
}

type Links []*Link

func NewLink() *Links {
	return &Links{}
}

func (l *Links) AddSelf(href string) *Links {
	return l.add(
		"self",
		href,
	)
}

func (l *Links) AddPrev(
	href	string,
) *Links {
	return l.add(
		"prev",
		href,
	)
}

func (l *Links) AddNext(
	href	string,
) *Links {
	return l.add(
		"next",
		href,
	)
}

func (l *Links) add(
	rel		string,
	href 	string,
) *Links {
	*l = AddLink(
		*l,
		&Link{
			Rel: rel,
			Href: href,
		},
	)

	return l
}

func AddLink(l []*Link, link *Link) []*Link {
	l = append(l, link)
	return l
}

func ParseUrlQueryToHref(
	values url.Values,
) string {
	var href string
	var i int
	for k, v := range values {
		href += k
		href += "="+massOfStringToString(v)
		if !(i == len(values) - 1) {
			href += "&"
		}
		i++
	}

	return href
}

func massOfStringToString(mass []string) string {
	var str string
	for i := 0; i < len(mass); i++ {
		str += mass[i]
		if !(i == len(mass) - 1) {
			str += "+"
		}
	}
	return str
}