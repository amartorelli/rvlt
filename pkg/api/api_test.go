package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

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

func TestSetBirthdayHandler(t *testing.T) {
	// test valid request
	rr := httptest.NewRecorder()
	dobr := DOBRequest{DOB: "2017-04-23"}
	dobrb, _ := json.Marshal(dobr)

	req := httptest.NewRequest("PUT", "http://localhost:8080/hello/john", bytes.NewReader(dobrb))
	req.Header.Set("Content-type", "application/json")
	a.mux.ServeHTTP(rr, req)

	resp := rr.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Error("setting the birthday for a valid user should return NoContent (204)")
	}

	// test invalid content type header
	rr = httptest.NewRecorder()
	dobr = DOBRequest{DOB: "2017-04-23"}
	dobrb, _ = json.Marshal(dobr)

	req = httptest.NewRequest("PUT", "http://localhost:8080/hello/john", bytes.NewReader(dobrb))
	a.mux.ServeHTTP(rr, req)

	resp = rr.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Error("a valid request with content type not set to application/json should give a MethodNotAllowed(405)")
	}

	// invalid json
	rr = httptest.NewRecorder()

	rs := "{\"date\"}"
	req = httptest.NewRequest("PUT", "http://localhost:8080/hello/john", bytes.NewReader([]byte(rs)))
	req.Header.Set("Content-type", "application/json")
	a.mux.ServeHTTP(rr, req)

	resp = rr.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("an invalid json request should return BadRequest (400)")
	}

	// invalid user
	rr = httptest.NewRecorder()
	dobr = DOBRequest{DOB: "2017-04-43"}
	dobrb, _ = json.Marshal(dobr)

	req = httptest.NewRequest("PUT", "http://localhost:8080/hello/john05", bytes.NewReader(dobrb))
	req.Header.Set("Content-type", "application/json")
	a.mux.ServeHTTP(rr, req)

	resp = rr.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Error("an invalid user request should return BadRequest (400)")
	}
}

func TestRenderBirthdayMessage(t *testing.T) {
	u := model.User{Username: "john", DOB: "2014-03-03"}

	// valid user, two days until birthday
	today, _ := time.Parse("2006-01-02", "2016-03-01")
	msg, err := renderBirthdayMessage(u, today)
	if err != nil {
		t.Error("a valid user should render the message")
	}
	if msg != "Hello, john! Your birthday is in 2 day(s)" {
		t.Errorf("john's birthday should be in one day, message is: %s", msg)
	}

	// valid user, today is his birthday
	today, _ = time.Parse("2006-01-02", "2016-03-03")
	msg, err = renderBirthdayMessage(u, today)
	if err != nil {
		t.Error("a valid user should render the message")
	}
	if msg != "Hello, john! Happy birthday!" {
		t.Errorf("john's birthday is today, message is: %s", msg)
	}

	u = model.User{Username: "john", DOB: "2914-03-03"}

	// invalid user, birthday is in the future
	today, _ = time.Parse("2006-01-02", "2016-03-03")
	msg, err = renderBirthdayMessage(u, today)
	if err != errBirthdayInTheFuture {
		t.Error("a user with birthday in the future should return an error")
	}
	if msg != "" {
		t.Errorf("when a user is invalid the message should be empty, message is: %s", msg)
	}
}

func TestMain(m *testing.M) {
	db, _ = database.NewDatabase("memory")
	db.Store(model.User{Username: "john", DOB: "2014-12-03"})

	a, _ = NewHelloWorldAPI(":8080", db)
	a.initHandlers()

	os.Exit(m.Run())
}
