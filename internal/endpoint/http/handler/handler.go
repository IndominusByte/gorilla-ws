package handler_http

import (
	"net/http"

	"github.com/IndominusByte/gorilla-ws/internal/config"
	endpoint_http "github.com/IndominusByte/gorilla-ws/internal/endpoint/http"
	hublibs "github.com/IndominusByte/gorilla-ws/internal/libs/hub"
	homesrepo "github.com/IndominusByte/gorilla-ws/internal/repo/homes"
	homesusecase "github.com/IndominusByte/gorilla-ws/internal/usecase/homes"
	"github.com/creent-production/cdk-go/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	Router *chi.Mux
	// Db config can be added here
	db       *sqlx.DB
	redisCli *redis.Pool
	cfg      *config.Config
}

func CreateNewServer(db *sqlx.DB, redisCli *redis.Pool, cfg *config.Config) *Server {
	s := &Server{db: db, redisCli: redisCli, cfg: cfg}
	s.Router = chi.NewRouter()
	return s
}

func CustomVerifier(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return jwtauth.Verify(ja, jwtauth.TokenFromHeader, jwtauth.TokenFromCookie, jwtauth.TokenFromQuery)
}

func (s *Server) MountHandlers() error {
	// jwt
	publicKey, privateKey := auth.DecodeRSA(s.cfg.JWT.PublicKey, s.cfg.JWT.PrivateKey)
	TokenAuthRS256 := jwtauth.New(s.cfg.JWT.Algorithm, privateKey, publicKey)
	s.Router.Use(CustomVerifier(TokenAuthRS256))

	// middleware stack
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	s.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), //The url pointing to API definition
	))

	hub := hublibs.NewHub()
	go hub.Run()

	// you can insert your behaviors here
	homesRepo, err := homesrepo.New(s.db)
	if err != nil {
		return err
	}
	homesUsecase := homesusecase.NewHomesUsecase(homesRepo)
	endpoint_http.AddHomes(s.Router, homesUsecase, s.redisCli, hub)

	return nil
}
