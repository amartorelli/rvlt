package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/amartorelli/rvlt/pkg/database"
	"github.com/amartorelli/rvlt/pkg/model"
)

var db database.Database
var a *HelloWorldAPI

func TestGetBirthdayHandler(t *testing.T) {
	// testing non existent user
	rr := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "http://localhost:8080/hello/zack", nil)
	a.mux.ServeHTTP(rr, req)

	resp := rr.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Error("looking for a non existing user should return NotFound (404)")
	}

	// testing non existent user
	rr = httptest.NewRecorder()

	req = httptest.NewRequest("GET", "http://localhost:8080/hello/john", nil)
	a.mux.ServeHTTP(rr, req)

	resp = rr.Result()
	if resp.StatusCode != http.StatusOK {
		t.Error("looking for an existing user should return Ok (200)")
	}
}

func TestGSetBirthdayHandler(t *testing.T) {
	// TODO
}

func TestMain(m *testing.M) {
	db, _ = database.NewDatabase("memory")
	db.Store(model.User{Username: "john", DOB: "2014-12-03"})

	a, _ = NewHelloWorldAPI(":8080", db)
	a.initHandlers()

	os.Exit(m.Run())
}
