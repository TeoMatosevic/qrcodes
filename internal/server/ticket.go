package server

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketPageData struct {
	Status    string `json:"status"`
	Vatin     string `json:"vatin"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	CreatedAt string `json:"createdAt"`
}

func (s *Server) TicketHandler(ctx *gin.Context) {
	tmpl, err := template.ParseFiles("internal/templates/ticket.gohtml")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	ticket, err := s.db.GetTicket(id)

	if err != nil {
		t := TicketPageData{
			Status:    http.StatusText(http.StatusNotFound),
			Vatin:     "",
			FirstName: "",
			LastName:  "",
			CreatedAt: "",
		}

		tmpl.Execute(ctx.Writer, t)

		ctx.Writer.Header().Set("Content-Type", "text/html")
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	t := TicketPageData{
		Status:    http.StatusText(http.StatusOK),
		Vatin:     ticket.Vatin,
		FirstName: ticket.FirstName,
		LastName:  ticket.LastName,
		CreatedAt: ticket.CreatedAt,
	}

	tmpl.Execute(ctx.Writer, t)

	ctx.Writer.Header().Set("Content-Type", "text/html")
	ctx.Writer.WriteHeader(http.StatusOK)
}
