package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/xopxe23/news-server/internal/config"
	"github.com/xopxe23/news-server/internal/repository"
	"github.com/xopxe23/news-server/internal/service"
	"github.com/xopxe23/news-server/internal/transport/rest"
	"github.com/xopxe23/news-server/pkg/database"
	hasher "github.com/xopxe23/news-server/pkg/hash"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("config: %+v\n", cfg)

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	hasher := hasher.NewSHA1Hasher("salt")

	usersRepos := repository.NewUsersRepository(db)
	tokensRepos := repository.NewTokensRepository(db)
	usersService := service.NewUsersService(usersRepos, hasher, tokensRepos, []byte("sample secret"))

	handler := rest.NewHandler(usersService)

	// init & run server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.InitRoutes(),
	}
	log.Info("SERVER STARTED")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
