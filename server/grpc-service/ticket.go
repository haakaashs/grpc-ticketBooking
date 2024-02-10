package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/haakaashs/grpc-ticketBooking/log"
	"github.com/haakaashs/grpc-ticketBooking/protos/ticketBooking"
)

type server struct {
	receipt map[string]*ticketBooking.Receipt
	mu      sync.Mutex
	ticketBooking.TicketBookingServiceServer
}

type train struct {
	counterA int32
	counterB int32
	sectionA map[string]int32
	sectionB map[string]int32
}

var trainSeats *train

func NewTicketBookingServer() *server {
	trainSeats = &train{
		sectionA: make(map[string]int32, 10),
		sectionB: make(map[string]int32, 10),
	}
	return &server{
		receipt: make(map[string]*ticketBooking.Receipt),
	}
}

func (s *server) CreateReceipt(ctx context.Context, receipt *ticketBooking.Receipt) (*ticketBooking.Message, error) {
	funcDesc := "CreateReceipt"
	log.Generic.INFO("enter server " + funcDesc)

	s.mu.Lock()
	defer s.mu.Unlock()

	var err error

	if _, ok := s.receipt[receipt.User.EmailId]; ok {
		err := errors.New("receipt already exist")
		log.Generic.ERROR(err)
		return &ticketBooking.Message{}, err
	}

	// assuming that the train has 20 seats each section has 10 seats
	s.receipt[receipt.User.EmailId] = receipt
	if trainSeats.counterA < 10 {
		trainSeats.counterA++
		trainSeats.sectionA[receipt.User.EmailId] = trainSeats.counterA
	} else if trainSeats.counterA == 10 && trainSeats.counterB < 10 {
		trainSeats.counterB++
		trainSeats.sectionB[receipt.User.EmailId] = trainSeats.counterB
	} else {
		err = errors.New("seats are full")
		log.Generic.ERROR(err)
		return &ticketBooking.Message{}, err
	}

	log.Generic.INFO("exit server " + funcDesc)
	return &ticketBooking.Message{Message: "created successfully"}, nil
}

func (s *server) ReadReceiptByUserEmail(ctx context.Context, email *ticketBooking.Message) (*ticketBooking.Receipt, error) {
	funcDesc := "ReadReceiptByUserEmail"
	log.Generic.INFO("enter server " + funcDesc)

	s.mu.Lock()
	defer s.mu.Unlock()

	res, ok := s.receipt[email.Message]
	if !ok {
		err := errors.New("receipt does not exist")
		log.Generic.ERROR(err)
		return res, err
	}

	log.Generic.INFO("exit server " + funcDesc)
	return res, nil
}

func (s *server) ReadUsersAndSeatsBySection(ctx context.Context, section *ticketBooking.Message) (*ticketBooking.UsersSeats, error) {
	funcDesc := "ReadUsersAndSeatsBySection"
	log.Generic.INFO("enter server " + funcDesc)

	s.mu.Lock()
	defer s.mu.Unlock()

	seatInfo := make(map[string]int32)
	var err error

	if section.Message == "sectionA" {
		for key := range s.receipt {
			if res, ok := s.receipt[key]; ok {
				seatInfo[fmt.Sprintf("%s %s", res.User.FirstName, res.User.SecondName)] = trainSeats.sectionA[key]
			}
		}

	} else if section.Message == "sectionB" {
		for key := range s.receipt {
			if res, ok := s.receipt[key]; ok {
				seatInfo[fmt.Sprintf("%s %s", res.User.FirstName, res.User.SecondName)] = trainSeats.sectionB[key]
			}
		}
	} else {
		err = errors.New(section.Message + " does not exist")
		log.Generic.ERROR(err)
		return &ticketBooking.UsersSeats{}, err
	}

	log.Generic.INFO("exit server " + funcDesc)
	return &ticketBooking.UsersSeats{MapData: seatInfo}, nil
}

func (s *server) DeleteUserByEmail(ctx context.Context, email *ticketBooking.Message) (*ticketBooking.Message, error) {
	funcDesc := "DeleteUserByEmail"
	log.Generic.INFO("enter server " + funcDesc)

	s.mu.Lock()
	defer s.mu.Unlock()

	var err error

	if _, ok := trainSeats.sectionA[email.Message]; ok {
		delete(trainSeats.sectionA, email.Message)
	} else if _, ok := trainSeats.sectionB[email.Message]; ok {
		delete(trainSeats.sectionB, email.Message)
	} else {
		err = errors.New(email.Message + " does not exist")
		log.Generic.ERROR(err)
		return &ticketBooking.Message{}, err
	}

	log.Generic.INFO("exit server " + funcDesc)
	return &ticketBooking.Message{Message: "deleted successfully"}, nil
}

func (s *server) UpdateSeatByEmail(ctx context.Context, input *ticketBooking.UpdateSeat) (*ticketBooking.Message, error) {
	funcDesc := "DeleteUserByEmail"
	log.Generic.INFO("enter server " + funcDesc)

	s.mu.Lock()
	defer s.mu.Unlock()

	var err error

	if input.Section == "sectionA" {
		if input.SeatNo < 10 && checkSeatAvailable(trainSeats.sectionA, input.SeatNo) {
			s.DeleteUserByEmail(context.Background(), &ticketBooking.Message{Message: input.Email})
			trainSeats.sectionA[input.Email] = input.SeatNo
		} else {
			err = fmt.Errorf("seat number %d does not exist", input.SeatNo)
			log.Generic.ERROR(err)
			return &ticketBooking.Message{}, err
		}
	} else if input.Section == "sectionB" {
		if input.SeatNo <= 10 && checkSeatAvailable(trainSeats.sectionB, input.SeatNo) {
			s.DeleteUserByEmail(context.Background(), &ticketBooking.Message{Message: input.Email})
			trainSeats.sectionB[input.Email] = input.SeatNo
		} else {
			err = fmt.Errorf("seat number %d does not exist", input.SeatNo)
			log.Generic.ERROR(err)
			return &ticketBooking.Message{}, err
		}
	} else {
		err = errors.New("section does not exist")
		log.Generic.ERROR(err)
		return &ticketBooking.Message{}, err
	}

	log.Generic.INFO("exit server " + funcDesc)
	return &ticketBooking.Message{Message: "deleted successfully"}, nil
}

func checkSeatAvailable(data map[string]int32, seatNo int32) bool {
	for _, value := range data {
		if value == seatNo {
			return false
		}
	}
	return true
}
