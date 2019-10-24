package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/amartorelli/rvlt/pkg/database"
	"github.com/amartorelli/rvlt/pkg/model"
	"github.com/amartorelli/rvlt/pkg/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// HelloWorldAPI is the API server structure
type HelloWorldAPI struct {
	server *http.Server
	mux    *mux.Router
	db     database.Database
}

// DOBResponse is the structure used to encode the response
type DOBResponse struct {
	Message string `json:"message"`
}

// DOBRequest is the structure used to decode the request
type DOBRequest struct {
	DOB string `json:"dateOfBirth"`
}

var (
	errRendering           = errors.New("error rendering message")
	errBirthdayInTheFuture = errors.New("birthday date is in the future")
)

// NewHelloWorldAPI creates a new instance of the HelloWorldAPI structure
func NewHelloWorldAPI(addr string, db database.Database) (*HelloWorldAPI, error) {
	mux := mux.NewRouter()
	s := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &HelloWorldAPI{server: s, mux: mux, db: db}, nil
}

func (a *HelloWorldAPI) initHandlers() {
	a.mux.Handle("/health", http.HandlerFunc(a.healthHandler)).Methods("GET")
	a.mux.Handle("/metrics", promhttp.Handler())
	a.mux.Handle("/hello/{username}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(a.getBirthdayHandler))).Methods("GET")
	a.mux.Handle("/hello/{username}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(a.setBirthdayHandler))).Methods("PUT").HeadersRegexp("Content-Type", "application/json")
}

// Start starts the http server
func (a *HelloWorldAPI) Start() error {
	a.initHandlers()

	log.Info("starting API")
	err := a.server.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop stops the API server
func (a *HelloWorldAPI) Stop() error {
	err := a.server.Shutdown(context.Background())
	if err != nil {
		return err
	}

	err = a.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (a *HelloWorldAPI) healthHandler(w http.ResponseWriter, req *http.Request) {

}

func (a *HelloWorldAPI) getBirthdayHandler(w http.ResponseWriter, req *http.Request) {
	hstart := time.Now()
	defer reqDuration.WithLabelValues("get-birthday").Observe(time.Since(hstart).Seconds())

	w.Header().Set("Content-type", "application/json")
	var r DOBResponse
	vars := mux.Vars(req)

	u, err := a.db.Get(vars["username"])
	if err == database.ErrUserNotFound {
		log.Error(err)
		reqMetric.WithLabelValues("get-birthday", strconv.Itoa(http.StatusNotFound)).Inc()
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(r)
		return
	} else if err != nil {
		log.Error(err)
		reqMetric.WithLabelValues("get-birthday", strconv.Itoa(http.StatusInternalServerError)).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(r)
		return
	}

	msg, err := renderBirthdayMessage(u, time.Now())
	if err != nil {
		log.Error(err)
		reqMetric.WithLabelValues("get-birthday", strconv.Itoa(http.StatusInternalServerError)).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(r)
		return
	}

	r.Message = msg
	reqMetric.WithLabelValues("get-birthday", strconv.Itoa(http.StatusOK)).Inc()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(r)
}

func (a *HelloWorldAPI) setBirthdayHandler(w http.ResponseWriter, req *http.Request) {
	hstart := time.Now()
	defer reqDuration.WithLabelValues("set-birthday").Observe(time.Since(hstart).Seconds())

	vars := mux.Vars(req)
	var r DOBRequest

	// get and check DOB
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		log.Error(err)
		reqMetric.WithLabelValues("set-birthday", strconv.Itoa(http.StatusBadRequest)).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := model.User{Username: vars["username"], DOB: r.DOB}

	// check user is valid
	v, err := u.IsValid()
	if !v {
		log.Error(err)
		reqMetric.WithLabelValues("set-birthday", strconv.Itoa(http.StatusBadRequest)).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.db.Store(u)
	if err != nil {
		log.Error(err)
		reqMetric.WithLabelValues("set-birthday", strconv.Itoa(http.StatusInternalServerError)).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reqMetric.WithLabelValues("set-birthday", strconv.Itoa(http.StatusNoContent)).Inc()
	w.WriteHeader(http.StatusNoContent)
}

func renderBirthdayMessage(u model.User, now time.Time) (string, error) {
	var msg string
	bd, err := time.Parse("2006-01-02", u.DOB)
	if err != nil {
		log.Debugf("+%v", bd)
		return "", errRendering
	}

	if bd.After(now) {
		log.Debugf("+%v", bd)
		return "", errBirthdayInTheFuture
	}

	days, err := utils.DaysUntilBirthday(bd, now)
	if err != nil {
		return "", err
	}

	if days == 0 {
		msg = fmt.Sprintf("Hello, %s! Happy birthday!", u.Username)
	} else {
		msg = fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", u.Username, days)
	}

	return msg, nil
}
