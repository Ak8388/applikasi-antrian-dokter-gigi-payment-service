package manager

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/config"
	_ "github.com/lib/pq"
)

type InfraManager interface {
	DBConnection() *sql.DB
}

type infraManager struct {
	cfg *config.Config
	db  *sql.DB
}

func (infra *infraManager) opneConnect() error {

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", infra.cfg.DBUser, infra.cfg.DBPass, infra.cfg.DBName, infra.cfg.DBPort)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return errors.New("failed open your database" + err.Error())
	}

	if err := db.Ping(); err != nil {
		return errors.New("failed connect to your database")
	}

	infra.db = db

	return nil
}

func (infra *infraManager) DBConnection() *sql.DB {
	return infra.db
}

func NewInfraManager(cfg *config.Config) InfraManager {
	infra := &infraManager{cfg: cfg}

	if err := infra.opneConnect(); err != nil {
		panic(err)
	}

	return infra
}
