package grpcservice

import (
	"context"
	"time"

	"github.com/haakaashs/grpc-ticketBooking/log"
	"github.com/haakaashs/grpc-ticketBooking/protos/ticketBooking"
)

var Client ticketBooking.TicketBookingServiceClient

func CreateReceipt(input *ticketBooking.Receipt) (*ticketBooking.Message, error) {
	funcDesc := "CreateReceipt"
	log.Generic.INFO("enter client " + funcDesc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	res, err := Client.CreateReceipt(ctx, input)

	if err != nil {
		log.Generic.ERROR(err)
		return &ticketBooking.Message{}, err
	}

	log.Generic.INFO("exit client " + funcDesc)
	return res, nil
}

func ReadReceiptByUserEmail(input *ticketBooking.Message) (*ticketBooking.Receipt, error) {
	funcDesc := "ReadReceiptByUserEmail"
	log.Generic.INFO("enter client " + funcDesc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	res, err := Client.ReadReceiptByUserEmail(ctx, input)

	if err != nil {
		log.Generic.ERROR(err)
		return &ticketBooking.Receipt{}, err
	}

	log.Generic.INFO("exit client " + funcDesc)
	return res, nil
}

func ReadUsersAndSeatsBySection(input *ticketBooking.Message) (*ticketBooking.UsersSeats, error) {
	funcDesc := "ReadUsersAndSeatsBySection"
	log.Generic.INFO("enter client " + funcDesc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	res, err := Client.ReadUsersAndSeatsBySection(ctx, input)

	if err != nil {
		log.Generic.ERROR(err)
		return &ticketBooking.UsersSeats{}, err
	}

	log.Generic.INFO("exit client " + funcDesc)
	return res, nil
}

func DeleteUserByEmail(input *ticketBooking.Message) (*ticketBooking.Message, error) {
	funcDesc := "DeleteUserByEmail"
	log.Generic.INFO("enter client " + funcDesc)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	res, err := Client.DeleteUserByEmail(ctx, input)

	if err != nil {
		log.Generic.ERROR(err)
		return &ticketBooking.Message{}, err
	}

	log.Generic.INFO("exit client " + funcDesc)
	return res, nil
}

func UpdateSeatByEmail(input *ticketBooking.UpdateSeat) (*ticketBooking.Message, error) {
	funcDesc := "UpdateSeatByEmail"
	log.Generic.INFO("enter client " + funcDesc)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// defer func() {
	// 	if cancel != nil {
	// 		cancel()
	// 	}
	// }()
	ctx := context.Background()
	res, err := Client.UpdateSeatByEmail(ctx, input)

	if err != nil {
		log.Generic.ERROR(err)
		return &ticketBooking.Message{}, err
	}

	log.Generic.INFO("exit client " + funcDesc)
	return res, nil
}
