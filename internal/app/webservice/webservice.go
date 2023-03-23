package webservice

import (
	"net/http"
	"os"

	"github.com/NoireHub/NATS-streaming-WebService/internal/app/store/sqlstore"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Start(config *Config) error {
	connSTR:= config.DatabaseURL
	if os.Getenv("HOST") != "" {
		connSTR += " host=" + os.Getenv("HOST")
	}
	db, err := newDB(connSTR)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	server, err := NewServer(store, config.StanClusterID, config.StanSubject)
	if err != nil {
		return err
	}
	server.logger.Info("Server is Up")

	defer server.stanSub.Close()
	defer server.stanConn.Close()

	return http.ListenAndServe(config.BindAddr, server)
}

func newDB(dbURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres",dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
