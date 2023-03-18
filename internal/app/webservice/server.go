package webservice

import (
	"github.com/NoireHub/NATS-streaming-WebService/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	logger *logrus.Logger
	router *mux.Router
	store	store.Store
}