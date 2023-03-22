package webservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NoireHub/NATS-streaming-WebService/internal/app/model"
	"github.com/NoireHub/NATS-streaming-WebService/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
)

type server struct {
	logger   *logrus.Logger
	router   *mux.Router
	store    store.Store
	stanConn stan.Conn
	stanSub  stan.Subscription
}

func NewServer(store store.Store, stanClusterID string, stanSubj string) (*server, error) {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	sc, err := stan.Connect(stanClusterID, "subscriber", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		return nil, err
	}
	s.stanConn = sc

	sub, err := sc.Subscribe(stanSubj, s.handleNatsMessage)
	if err != nil {
		return nil, err
	}
	s.stanSub = sub

	if err:= s.store.Order().GetCache(); err != nil {
		return nil, err
	}

	s.configureRouter()

	return s, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Methods().HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s.preflightHandler(w, r)
		})

}

func (s *server) handleNatsMessage(m *stan.Msg) {
	d := &model.Order{}
	fmt.Println(m)
	err := json.Unmarshal(m.Data, d)
	fmt.Println(d)
	fmt.Println(err)
	if err := s.store.Order().AddOrder(d); err != nil {
		s.logger.Info(err)
	}

	s.logger.Info(s.store.Order().GetOrderById("b563feb7b2b84b6test"))
}

func (s *server) preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "content-Type")
	w.Header().Add("Connection", "keep-alive")
	w.Header().Add("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Add("Access-Control-Max-Age", "86400")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(200)
}
