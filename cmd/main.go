package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

// @title News API 
// @version 1.0
// @description Sample Server for News App
// @host localhost:8000
// @basepath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

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

	hasher := hasher.NewSHA1Hasher("salt")

	usersRepos := repository.NewUsersRepository(db)
	tokensRepos := repository.NewTokensRepository(db)
	usersService := service.NewUsersService(usersRepos, hasher, tokensRepos, []byte("sample secret"))

	authorsRepos := repository.NewAuthorsRepository(db)
	articlesRepos := repository.NewArticlesRepository(db)
	articlesService := service.NewArticlesService(authorsRepos, articlesRepos)

	handler := rest.NewHandler(usersService, articlesService)

	// init & run server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.InitRoutes(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Info("SERVER STARTED")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<- quit

	log.Info("SERVER SHUTDOWN")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Errorf("error on server shutting: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Errorf("error on db closing: %s", err.Error())
	}
}
