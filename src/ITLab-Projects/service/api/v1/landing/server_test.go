package landing_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	model "github.com/ITLab-Projects/pkg/models/landing"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFunc_ServerTests(t *testing.T) {
	if err := RepoImp.LandingRepositoryImp.SaveAndDeleteUnfindLanding(
		context.Background(),
		[]model.Landing{
			{
				LandingCompact: model.LandingCompact{
					RepoId: 1,
					Title: "mock_1",
					Tags: []string{"Backend", "Web"},
				},
			},
			{
				LandingCompact: model.LandingCompact{
					RepoId: 2,
					Title: "mock_2",
					Tags: []string{"Web"},
				},
			},
			{
				LandingCompact: model.LandingCompact{
					RepoId: 3,
					Title: "mock_3",
					Tags: []string{"VR"},
				},
			},
		},
	); err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer Repositories.Landing.DeleteMany(
		context.Background(),
		bson.M{},
		nil,
	)

	testfunc_GetAllLanding_HTTP_WithOutParams(t)
	testfunc_GetAllLandings_HTTP_ByName(t)
	testfunc_GetAllLandings_HTTP_ByTag(t)
	testfunc_HTTP_GetByID(t)
	testfunc_HTTP_GetByID_NotFound(t)
}

func testfunc_GetAllLanding_HTTP_WithOutParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/landing", nil)
	resp := httptest.NewRecorder()

	Router.ServeHTTP(resp, req)

	if resp.Result().StatusCode != http.StatusOK {
		t.Log("Assert error")
		t.Log(resp.Result().StatusCode)
		t.FailNow()
	}

	var ls []*model.LandingCompact

	if err := json.NewDecoder(resp.Body).Decode(&ls); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(ls) != 3 {
		t.Log("assert err")
		t.FailNow()
	}

	for _, l := range ls {
		if !(l.Title == "mock_1" && l.RepoId == 1 || l.Title == "mock_2" && l.RepoId == 2 || l.Title == "mock_3" && l.RepoId == 3) {
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func testfunc_GetAllLandings_HTTP_ByName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/landing?name=mock_3", nil)
	resp := httptest.NewRecorder()

	Router.ServeHTTP(resp, req)

	if resp.Result().StatusCode != http.StatusOK {
		t.Log("Assert error")
		t.Log(resp.Result().StatusCode)
		t.FailNow()
	}

	var ls []*model.LandingCompact

	if err := json.NewDecoder(resp.Body).Decode(&ls); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(ls) != 1 {
		t.Log("Assert error")
		t.FailNow()
	}

	l := ls[0]

	if !(l.Title == "mock_3" && l.RepoId == 3) {
		t.Log("Assert error")
		t.FailNow()
	}
}

func testfunc_GetAllLandings_HTTP_ByTag(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/landing?tag=Web", nil)
	resp := httptest.NewRecorder()

	Router.ServeHTTP(resp, req)

	if resp.Result().StatusCode != http.StatusOK {
		t.Log("Assert error")
		t.Log(resp.Result().StatusCode)
		t.FailNow()
	}

	var ls []*model.LandingCompact

	if err := json.NewDecoder(resp.Body).Decode(&ls); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if len(ls) != 2 {
		t.Log("Assert error")
		t.FailNow()
	}

	
	for _, l := range ls {
		if !(l.Title == "mock_1" && l.RepoId == 1 || l.Title == "mock_2" && l.RepoId == 2) {
			t.Log("Assert error")
			t.FailNow()
		}
	}
}

func testfunc_HTTP_GetByID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/landing/1", nil)
	resp := httptest.NewRecorder()

	Router.ServeHTTP(resp, req)

	if resp.Result().StatusCode != http.StatusOK {
		t.Log("Assert error")
		t.Log(resp.Result().StatusCode)
		t.FailNow()
	}

	l := &model.Landing{}

	if err := json.NewDecoder(resp.Body).Decode(l); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !(l.RepoId == 1 && l.Title == "mock_1") {
		t.Log("assert error")
		t.FailNow()
	}
}

func testfunc_HTTP_GetByID_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/landing/123123", nil)
	resp := httptest.NewRecorder()

	Router.ServeHTTP(resp, req)

	if resp.Result().StatusCode != http.StatusNotFound {
		t.Log("Assert error")
		t.Log(resp.Result().StatusCode)
		t.FailNow()
	}
}