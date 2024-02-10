package main

import (
	"net/http"

	"github.com/gorilla/mux"
	grpcservice "github.com/haakaashs/grpc-ticketBooking/client/grpc-service"
	"github.com/haakaashs/grpc-ticketBooking/log"
	"github.com/haakaashs/grpc-ticketBooking/protos/ticketBooking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcPort = ":50001"
	restPort = ":8080"
)

func main() {
	conn, err := grpc.Dial(grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Generic.ERROR(err)
		log.Generic.FATAL("connection to server failed")
	}

	defer conn.Close()

	log.Generic.INFO("client listening at port " + grpcPort)

	grpcservice.Client = ticketBooking.NewTicketBookingServiceClient(conn)
	// we can directly hit the server using client

	// or we can indirectly hit the server using client using rest API
	server := mux.NewRouter()

	server.HandleFunc("/receipt", grpcservice.CreateReceiptRest).Methods("POST")
	server.HandleFunc("/receipt/{email}", grpcservice.ReadReceiptByUserEmailRest).Methods("GET")
	server.HandleFunc("/users/{section}", grpcservice.ReadUsersAndSeatsBySectionRest).Methods("GET")
	server.HandleFunc("/user/{email}", grpcservice.DeleteUserByEmailRest).Methods("DELETE")
	server.HandleFunc("/user", grpcservice.UpdateSeatByEmailRest).Methods("PUT")

	log.Generic.INFO("rest server started listening at port " + restPort)
	if err = http.ListenAndServe(restPort, server); err != nil {
		log.Generic.ERROR(err)
		log.Generic.FATAL("error starting rest web server at port " + restPort)
	}
}
