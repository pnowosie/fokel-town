package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/pnowosie/fokeltown-merkle/internal"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetUser(t *testing.T) {
	existingUser := internal.UserData{Id: "00000f"}
	tests := map[string]struct {
		userid       string
		expectedCode int
	}{
		"finding existing user": {
			userid:       existingUser.Id,
			expectedCode: http.StatusOK,
		},
		"not finding non-existing user": {
			userid:       "11111f",
			expectedCode: http.StatusNotFound,
		},
		"invalid user id": {
			userid:       "abefgh",
			expectedCode: http.StatusBadRequest,
		},
	}

	// Setup
	merkle := &internal.MerkleTrie{}
	err := merkle.Put(existingUser.Id, existingUser)
	if err != nil {
		t.Fatal(err)
	}
	app := newApp(hclog.Default(), merkle)
	mux := app.Routes()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, newGetRequest(t, test.userid))
			rs := rr.Result()

			assert.Equal(t, test.expectedCode, rs.StatusCode)
		})
	}
}

func TestHandler_PutUser(t *testing.T) {
	existingUser := internal.UserData{Id: "00000f"}
	tests := map[string]struct {
		userid       string
		expectedCode int
	}{
		"already existing user": {
			userid:       existingUser.Id,
			expectedCode: http.StatusFound,
		},
		"new user": {
			userid:       "11111f",
			expectedCode: http.StatusCreated,
		},
		"invalid user id": {
			userid:       "abefgh",
			expectedCode: http.StatusBadRequest,
		},
	}

	// Setup
	merkle := &internal.MerkleTrie{}
	err := merkle.Put(existingUser.Id, existingUser)
	if err != nil {
		t.Fatal(err)
	}
	app := newApp(hclog.Default(), merkle)
	mux := app.Routes()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, newPutRequest(t, internal.UserData{Id: test.userid}))
			rs := rr.Result()

			assert.Equal(t, test.expectedCode, rs.StatusCode)
		})
	}
}

func TestHappyPath(t *testing.T) {
	// Setup
	merkle := &internal.MerkleTrie{}
	logger := hclog.Default()
	logger.SetLevel(hclog.Debug)
	app := newApp(logger, merkle)
	mux := app.Routes()

	user := internal.UserData{Id: "00000f", FirstName: "Joe", LastName: "Doe"}

	// Get not (yet) existing user
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, newGetRequest(t, user.Id))
	rs := rr.Result()
	assert.Equal(t, http.StatusNotFound, rs.StatusCode)

	// Put user
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, newPutRequest(t, user))
	rs = rr.Result()
	assert.Equal(t, http.StatusCreated, rs.StatusCode)

	// Get just created user
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, newGetRequest(t, user.Id))
	rs = rr.Result()
	assert.Equal(t, http.StatusOK, rs.StatusCode)
	defer rs.Body.Close()
	var data internal.UserData
	err := json.NewDecoder(rs.Body).Decode(&data)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, user, data)

	// Put already existing user
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, newPutRequest(t, user))
	rs = rr.Result()
	assert.Equal(t, http.StatusFound, rs.StatusCode)
}

func newGetRequest(t *testing.T, userid string) *http.Request {
	req, err := http.NewRequest("GET", "/v0/user/"+userid, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func newPutRequest(t *testing.T, data internal.UserData) *http.Request {
	bs, _ := json.Marshal(data)
	req, err := http.NewRequest("PUT", "/v0/user", bytes.NewReader(bs))
	if err != nil {
		t.Fatal(err)
	}
	return req
}
