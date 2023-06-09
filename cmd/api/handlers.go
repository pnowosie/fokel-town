package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pnowosie/fokeltown-merkle/internal"
)

func (app *application) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	response := apiVersion{
		Name:     ServiceName,
		Version:  app.version,
		UpTime:   time.Now().Unix() - app.startTime,
		TrieRoot: app.storage.Root().String(),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		app.logger.Error("json marshal returns", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func (app *application) GetUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	userid := params.ByName("id")
	if (&internal.UserData{Id: userid}).IsValid() == false {
		app.logger.Warn("invalid id", "id", userid)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	userData, err := app.storage.Get(userid)
	if err != nil {
		if err == internal.ErrNotFound {
			app.logger.Warn("user not found", "id", userid)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		// Unspecified error
		app.logger.Error("storage get returns", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(userData)
	if err != nil {
		app.logger.Error("json marshal returns", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func (app *application) PutUser(w http.ResponseWriter, r *http.Request) {
	user := internal.UserData{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.logger.Warn("json decode returns", "error", err)
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if user.IsValid() == false {
		app.logger.Warn("invalid user data", "user", user)
		http.Error(w, "invalid user data", http.StatusBadRequest)
		return
	}

	err = app.storage.Put(user.Id, user)
	if err != nil {
		if err == internal.ErrAlreadyExists {
			app.logger.Warn("user already exist", "user", user)
			http.Error(w, fmt.Sprintf("already exist /v0/user/%s", user.Id), http.StatusFound)
			return
		}
		// Unspecified error
		app.logger.Error("storage get returns", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/v0/user/%s", user.Id))
	w.WriteHeader(http.StatusCreated)
}

type apiVersion struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	UpTime   int64  `json:"uptime"`
	TrieRoot string `json:"root"`
}
