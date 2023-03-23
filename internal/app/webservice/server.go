package webservice

import (
	"encoding/json"
	"net/http"
	"os"
	"fmt"
	"text/template"

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
	natsURL:= ""

	if os.Getenv("HOST") != "" {
		natsURL = fmt.Sprintf("nats://%s:4222", os.Getenv("HOST"))
	}else{
		natsURL = "nats://localhost:4222"
	}

	sc, err := stan.Connect(stanClusterID, "subscriber", stan.NatsURL(natsURL))
	if err != nil {
		return nil, err
	}
	s.stanConn = sc

	sub, err := sc.Subscribe(stanSubj, s.handleNatsMessage)
	if err != nil {
		return nil, err
	}
	s.stanSub = sub

	if err := s.store.Order().GetCache(); err != nil {
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
	s.router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	s.router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.router.HandleFunc("/order/{order_uid}", s.handleGetOrder()).Methods("GET")

}

func (s *server) handleNatsMessage(m *stan.Msg) {
	order := model.Order{}
	if err := json.Unmarshal(m.Data, &order); err != nil {
		s.logger.Error(err)
		return
	}
	if err := s.store.Order().AddOrder(&order); err != nil {
		s.logger.Error(err)
		return
	}

	s.logger.Info("Added order: " + order.OrderUid)
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html; charset=utf-8")

		t, err := template.ParseFiles("./static/index.html")
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
		}

		if err := t.Execute(w, nil); err != nil {
			s.error(w, r, http.StatusNotFound, err)
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleGetOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		orderUID := mux.Vars(r)["order_uid"]

		order, err := s.store.Order().GetOrderById(orderUID)
		if err != nil {
			s.error(w, r, http.StatusNoContent, err)
		}

		s.respond(w, r, http.StatusOK, order)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error:": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
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
