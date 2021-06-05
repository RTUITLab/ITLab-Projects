package content

import (
	"encoding/base64"
	"fmt"
	"time"
)

type Content struct {
	RepoID		uint64
	// Deprecated
	DownloadURL string 	`json:"download_url"`
	Content		string	`json:"content"`
	Commit		*Commit
}

func (c *Content) GetContent() ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(c.Content)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("Nil data")
	}

	return data, nil
}

func (c *Content) GetDate() time.Time {
	return c.Commit.Commit.Commiter.Date
}

type Commit struct {
	Commit 	*CommitInfo	`json:"commit"`
}

type CommitInfo struct {
	Commiter *Committer	`json:"committer"`
}

type Committer struct {
	Date	time.Time	`json:"date"`
}