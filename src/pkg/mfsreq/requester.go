package mfsreq

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/ITLab-Projects/pkg/clientwrapper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Requester interface {
	DeleteFile(ID primitive.ObjectID) error
	GenerateDownloadLink(ID primitive.ObjectID) string
}

type Config struct {
	BaseURL		string
	TestMode	bool
}

type MFSRequester struct {
	baseURL		string

	clientWithWrap *clientwrapper.ClientWithWrap

	client *http.Client

	TestMode	bool
}

func New(cfg *Config) *MFSRequester {
	r := &MFSRequester{
		baseURL: cfg.BaseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
			},
		},
	}

	r.clientWithWrap = clientwrapper.New(r.client)
	r.TestMode = cfg.TestMode

	return r
}

func (r *MFSRequester) DeleteFile(ID primitive.ObjectID) error {
	if r.TestMode {
		return nil
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/delete/%s", r.baseURL, ID.Hex()), nil)
	if err != nil {
		return err
	}

	resp, err := r.clientWithWrap.Do(req)
	if _, ok := err.(net.Error); ok {
		return errors.Wrapf(NetError, "%v", err)
	} else if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return &UnexpectedCodeErr{
			Err:	errors.Wrapf(ErrUnexpectedCode, "%v", resp.StatusCode),
			Code: 	resp.StatusCode,
		}
	}
	return nil
}

func (r *MFSRequester) GenerateDownloadLink(ID primitive.ObjectID) string {
	return fmt.Sprintf("%s/download/%s", r.baseURL, ID.Hex())
}