package server

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"qrcodes/internal/database"
)

func (s *Server) GenerateTicketHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	var ticket database.TicketData

	err := c.BindJSON(&ticket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticketModel := database.NewTicket(ticket)
	count, err := s.db.AmountByVatin(ticketModel.Vatin)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count >= 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Too many tickets for this VATIN"})
		return
	}

	err = s.db.Insert(ticketModel)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	png := ticketModel.GenerateQRCode()

	buffer := new(bytes.Buffer)
	buffer.Write(png)

	c.Header("Content-Type", "image/png")
	c.Header("Content-Length", strconv.Itoa(buffer.Len()))
	c.Writer.Write(buffer.Bytes())
}
