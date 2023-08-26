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

type controller struct {
	storage storage.Manager
	logger  logger.BaseLogger
	users   user.Users
}

func Create(logger logger.BaseLogger, storage storage.Manager) *controller {
	controller := &controller{
		storage: storage,
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
	r.Post("/", c.postHandler)

	r.NotFound(http.NotFound)
	return r
}

func (c *controller) getHandler(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("[controller::getHandler] Handle Get request")

	user, _, _ := r.BasicAuth()
	var output jsonmsg.Output
	output.Texts = c.storage.GetTexts(user)

	buf, err := json.Marshal(&output)
	if err != nil {
		c.logger.Info("[controller::postHandler] Error during data parsing")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		c.logger.Info("[controller::postHandler] Error writing data in response body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
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

	spelledText, err := spellchecker.Handle(input.Text)
	if err != nil {
		c.logger.Info("[controller::postHandler] Failed to spell input data: %v. Added original text", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user, _, _ := r.BasicAuth()
	c.storage.AddText(user, spelledText)
}
