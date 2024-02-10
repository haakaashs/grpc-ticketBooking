package test

import (
	"context"
	"log"
	"testing"

	"github.com/go-playground/assert"
	"github.com/haakaashs/grpc-ticketBooking/protos/ticketBooking"
	service "github.com/haakaashs/grpc-ticketBooking/server/grpc-service"
)

func TestCreateReceipt(t *testing.T) {
	tests := []struct {
		Desc     string
		Input    *ticketBooking.Receipt
		Expected interface{}
		Ctx      context.Context
	}{
		{
			Desc: "success test case",
			Input: &ticketBooking.Receipt{
				User: &ticketBooking.User{
					EmailId: "haakaash@gmail.com",
				}},
			Expected: &ticketBooking.Message{Message: "created successfully"},
			Ctx:      context.Background(),
		},
		{
			Desc: "failure test case",
			Input: &ticketBooking.Receipt{
				User: &ticketBooking.User{
					EmailId: "haakaash@gmail.com",
				}},
			Expected: "receipt already exist",
			Ctx:      context.Background(),
		},
	}
	serverObj := service.NewTicketBookingServer()
	for _, tt := range tests {
		log.Println(tt.Desc)
		t.Run(tt.Desc, func(t *testing.T) {
			res, err := serverObj.CreateReceipt(tt.Ctx, tt.Input)

			if err != nil {
				assert.Equal(t, err.Error(), tt.Expected.(string))
				return
			}
			assert.Equal(t, res.Message, tt.Expected.(*ticketBooking.Message).Message)

		})
	}
}

func TestReadReceiptByUserEmail(t *testing.T) {
	tests := []struct {
		Desc     string
		Input    *ticketBooking.Message
		Expected interface{}
		Ctx      context.Context
	}{
		{
			Desc:     "failure test case",
			Input:    &ticketBooking.Message{Message: "haakaash@gmail.com"},
			Expected: "receipt does not exist",
			Ctx:      context.Background(),
		},
		{
			Desc:  "success test case",
			Input: &ticketBooking.Message{Message: "haakaash@gmail.com"},
			Expected: &ticketBooking.Receipt{
				User: &ticketBooking.User{
					EmailId: "haakaash@gmail.com",
				}},
			Ctx: context.Background(),
		},
	}
	serverObj := service.NewTicketBookingServer()
	for _, tt := range tests {
		log.Println(tt.Desc)
		t.Run(tt.Desc, func(t *testing.T) {
			if tt.Desc == "success test case" {
				serverObj.CreateReceipt(tt.Ctx, &ticketBooking.Receipt{
					User: &ticketBooking.User{
						EmailId: "haakaash@gmail.com",
					}})
			}
			res, err := serverObj.ReadReceiptByUserEmail(tt.Ctx, tt.Input)

			if err != nil {
				assert.Equal(t, err.Error(), tt.Expected.(string))
				return
			}
			assert.Equal(t, res.User.EmailId, tt.Expected.(*ticketBooking.Receipt).User.EmailId)

		})
	}
}
