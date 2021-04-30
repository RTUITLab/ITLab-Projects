package mfsreq

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mfsclient *http.Client

func init() {
	mfsclient =  &http.Client{
	 	Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
	}
}

type Requester interface {
	NewRequests(req *http.Request) Requests
	GenerateDownloadLink(ID primitive.ObjectID) string
}

type Requests interface {
	FileDeleter
}

type FileDeleter interface {
	DeleteFile(ID primitive.ObjectID) error
}

type Config struct {
	BaseURL		string
	TestMode	bool
}

type MFSRequester struct {
	baseURL		string

	TestMode	bool
}

type MFSRequests struct {
	baseURL 	string

	req *http.Request

	TestMode	bool
}

func (mfs *MFSRequests) DeleteFile(ID primitive.ObjectID) error {
	if mfs.TestMode {
		logrus.Info("TestMode activated")
		logrus.Info(mfs.req.Header)
		return nil
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/delete/%s", mfs.baseURL, ID.Hex()), nil)
	if err != nil {
		return err
	}

	req.Header = mfs.req.Header.Clone()

	resp, err := mfsclient.Do(req)
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


func New(cfg *Config) *MFSRequester {
	r := &MFSRequester{
		baseURL: cfg.BaseURL,
	}

	r.TestMode = cfg.TestMode

	return r
}

func (r *MFSRequester) NewRequests(req *http.Request) Requests {
	return &MFSRequests {
		baseURL: r.baseURL,
		TestMode: r.TestMode,
		req: req,
	}
}

func (r *MFSRequester) GenerateDownloadLink(ID primitive.ObjectID) string {
	return fmt.Sprintf("%s/download/%s", r.baseURL, ID.Hex())
}