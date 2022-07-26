package endpoint_http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/creent-production/cdk-go/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"

	hublibs "github.com/IndominusByte/gorilla-ws/internal/libs/hub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Resolve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type homesUsecaseIface interface {
}

func AddHomes(r *chi.Mux, uc homesUsecaseIface, redisCli *redis.Pool, h *hublibs.Hub) {
	r.Route("/homes", func(r chi.Router) {
		// protected route
		r.Group(func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
					if err := auth.ValidateJWT(r.Context(), redisCli, "jwtRequired"); err != nil {
						return
					}
					// Token is authenticated, pass it through
					next.ServeHTTP(rw, r)
				})
			})
			r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
				conn, err := upgrader.Upgrade(rw, r, nil)
				if err != nil {
					log.Println(err)
					return
				}

				_, claims, _ := jwtauth.FromContext(r.Context())
				sub, _ := strconv.Atoi(claims["sub"].(string))
				// TODO: handle user not found

				client := &hublibs.Client{
					Hub:       h,
					Conn:      conn,
					Send:      make(chan []byte, 256),
					Ctx:       r.Context(),
					TypeToken: "jwtRequired",
					UserId:    sub,
					RedisCli:  redisCli,
				}
				client.Hub.Register <- client

				go client.ReadPump()
				go client.WritePump()
			})
		})
		// public route
	})
}
