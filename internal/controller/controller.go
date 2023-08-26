package controller

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/erupshis/kode.git/internal/jsonmsg"
	"github.com/erupshis/kode.git/internal/logger"
	"github.com/erupshis/kode.git/internal/spellchecker"
	"github.com/erupshis/kode.git/internal/storage"
	"github.com/erupshis/kode.git/internal/user"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	storage storage.Manager
	logger  logger.BaseLogger
	users   user.Users
}

func Create(logger logger.BaseLogger, storage storage.Manager) *Controller {
	controller := &Controller{
		storage: storage,
		logger:  logger,
		users:   *user.Create(),
	}

	return controller
}

func (c *Controller) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Use(c.logger.LogHandler)
	r.Use(c.users.Auth)

	r.Get("/", c.getHandler)
	r.Post("/", c.postHandler)

	r.NotFound(http.NotFound)
	return r
}

func (c *Controller) getHandler(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("[Controller::getHandler] Handle Get request")

	username, _, _ := r.BasicAuth()
	var output jsonmsg.Output
	var err error
	output.Texts, err = c.storage.GetTexts(username)
	if err != nil {
		c.logger.Info("[Controller::postHandler] Error during storage access: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	buf, err := json.Marshal(&output)
	if err != nil {
		c.logger.Info("[Controller::postHandler] Error during data parsing: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		c.logger.Info("[Controller::postHandler] Error writing data in response body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
}

func (c *Controller) postHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	c.logger.Info("[Controller::postHandler] Handle JSON request with body: %s", buf.String())
	var input jsonmsg.Input
	if err := json.Unmarshal(buf.Bytes(), &input); err != nil {
		c.logger.Info("[Controller::postHandler] Error during JSON parsing: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	spelledText, err := spellchecker.Handle(input.Text)
	if err != nil {
		c.logger.Info("[Controller::postHandler] Failed to spell input data: %v. Added original text", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	username, _, _ := r.BasicAuth()
	c.storage.AddText(username, spelledText)
}
