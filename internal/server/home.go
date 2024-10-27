package server

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HomeHandler(c *gin.Context) {
	tmpl, err := template.ParseFiles("internal/templates/home.gohtml")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	amount, err := s.db.Count()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tmpl.Execute(c.Writer, amount)

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(http.StatusOK)
}
