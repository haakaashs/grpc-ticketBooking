package grpcservice

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/haakaashs/grpc-ticketBooking/log"
	"github.com/haakaashs/grpc-ticketBooking/protos/ticketBooking"
)

func CreateReceiptRest(w http.ResponseWriter, r *http.Request) {
	funcDesc := "CreateReceiptRest"
	log.Generic.INFO("enter rest " + funcDesc)

	w.Header().Set("Content-Type", "application/json")
	var receipt ticketBooking.Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil || receipt.From != "London" || receipt.To != "France" || receipt.PricePaid != 20 {
		log.Generic.ERROR(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	res, err := CreateReceipt(&receipt)
	if err != nil {
		log.Generic.ERROR(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Generic.INFO("exit rest " + funcDesc)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func ReadReceiptByUserEmailRest(w http.ResponseWriter, r *http.Request) {
	funcDesc := "ReadReceiptByUserEmailRest"
	log.Generic.INFO("enter rest " + funcDesc)

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	email := vars["email"]

	if !ValidateEmail(email) {
		log.Generic.ERROR("invalid user email")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user email"})
		return
	}

	res, err := ReadReceiptByUserEmail(&ticketBooking.Message{Message: email})
	if err != nil {
		log.Generic.ERROR(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Generic.INFO("exit rest " + funcDesc)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func ReadUsersAndSeatsBySectionRest(w http.ResponseWriter, r *http.Request) {
	funcDesc := "ReadUsersAndSeatsBySectionRest"
	log.Generic.INFO("enter rest " + funcDesc)

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	section := vars["section"]

	if section != "sectionA" && section != "sectionB" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid section"})
		return
	}

	res, err := ReadUsersAndSeatsBySection(&ticketBooking.Message{Message: section})
	if err != nil {
		log.Generic.ERROR(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Generic.INFO("exit rest " + funcDesc)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func DeleteUserByEmailRest(w http.ResponseWriter, r *http.Request) {
	funcDesc := "DeleteUserByEmailRest"
	log.Generic.INFO("enter rest " + funcDesc)

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	email := vars["email"]

	if !ValidateEmail(email) {
		log.Generic.ERROR("invalid user email")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user email"})
		return
	}

	res, err := DeleteUserByEmail(&ticketBooking.Message{Message: email})
	if err != nil {
		log.Generic.ERROR(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Generic.INFO("exit rest " + funcDesc)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func UpdateSeatByEmailRest(w http.ResponseWriter, r *http.Request) {
	funcDesc := "UpdateSeatByEmailRest"
	log.Generic.INFO("enter rest " + funcDesc)

	w.Header().Set("Content-Type", "application/json")
	var updateSeat ticketBooking.UpdateSeat
	err := json.NewDecoder(r.Body).Decode(&updateSeat)
	if err != nil {
		log.Generic.ERROR(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	res, err := UpdateSeatByEmail(&updateSeat)
	if err != nil {
		log.Generic.ERROR(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	log.Generic.INFO("exit rest " + funcDesc)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func ValidateEmail(email string) bool {
	pattern := `^[\w\.-]+@[a-zA-Z\d\.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		log.Generic.ERROR(err)
		return false
	}
	return matched
}
