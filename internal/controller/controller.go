package controller

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/erupshis/kode.git/internal/config"
	"github.com/erupshis/kode.git/internal/jsonmsg"
	"github.com/erupshis/kode.git/internal/logger"
	"github.com/erupshis/kode.git/internal/storage"
	"github.com/erupshis/kode.git/internal/user"
	"github.com/go-chi/chi/v5"
)

type controller struct {
	config  config.Config
	storage storage.Storage
	logger  logger.BaseLogger
	users   user.User
}

func Create(config config.Config, logger logger.BaseLogger) *controller {
	controller := &controller{
		config:  config,
		storage: *storage.Create(),
		logger:  logger,
		users:   *user.Create(),
	}

	return controller
}

func (c *controller) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Use(c.logger.LogHandler)
	r.Use(c.users.Auth)
	r.Get("/", c.getHandler)

	//TODO: need to wrap data.
	r.Post("/", c.postHandler)

	r.NotFound(http.NotFound)
	return r
}

func (c *controller) getHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (c *controller) postHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	c.logger.Info("[controller::postHandler] Handle JSON request with body: %s", buf.String())
	var input jsonmsg.Input
	if err := json.Unmarshal(buf.Bytes(), &input); err != nil {
		c.logger.Info("[controller::postHandler] Error during JSON parsing")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, _, _ := r.BasicAuth()
	//TODO: need Ya.Speller package.
	c.storage.AddText(user, input.Text)

	w.WriteHeader(http.StatusOK)
}
