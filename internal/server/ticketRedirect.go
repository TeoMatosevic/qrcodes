package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s *Server) TicketRedirectHandler(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("id")

	fmt.Println("TicketRedirectPageHandler", id)

	if id == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/ticket/%s", id))
}
