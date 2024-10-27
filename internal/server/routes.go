package server

import (
	"encoding/gob"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes(auth *Authenticator) http.Handler {
	r := gin.Default()

	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("auth-session", store))

	r.GET("/", s.HomeHandler)
	r.GET("/ticket/:id", s.IsAuthenticated, s.TicketHandler)
	r.GET("/ticket", s.TicketRedirectHandler)
	r.GET("/login", s.LoginHandler(auth))
	r.GET("/callback", s.CallbackHandler(auth))
	r.GET("/logout", s.LogoutHandler)

	r.POST("/api/v1/ticket", MtmMiddleware(), s.GenerateTicketHandler)

	return r
}
