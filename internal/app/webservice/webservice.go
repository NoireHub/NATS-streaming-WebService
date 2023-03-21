package webservice

import (
	"database/sql"
	"net/http"

	"github.com/NoireHub/NATS-streaming-WebService/internal/app/store/sqlstore"
	_ "github.com/lib/pq"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
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

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
