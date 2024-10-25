package database

import (
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

var (
	appURL = os.Getenv("APP_URL")
)

type TicketData struct {
	Vatin     string `json:"vatin" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type Ticket struct {
	ID        string `json:"id"`
	Vatin     string `json:"vatin"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
}

func NewTicket(data TicketData) Ticket {
	return Ticket{
		ID:        uuid.New().String(),
		Vatin:     data.Vatin,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}

func (t Ticket) GenerateQRCode() []byte {
	var png []byte
	png, err := qrcode.Encode(appURL+"/"+t.ID, qrcode.Medium, 256)

	if err != nil {
		panic(err)
	}

	return png
}
