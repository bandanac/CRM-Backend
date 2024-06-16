package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Customer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"`
}

var totalCustomers int = 5
var customers = map[string]Customer{
	"8f1f4cb1-aff0-4094-9398-bbf2dc89d6ad": {
		ID:        "8f1f4cb1-aff0-4094-9398-bbf2dc89d6ad",
		Name:      "Aria Bell",
		Role:      "Frontend Engineer",
		Email:     "aria_bell@xyz.com",
		Phone:     "+49-7001-3206995",
		Contacted: true,
	},
	"bc322825-8be4-406f-b79e-63c993750e7d": {
		ID:        "bc322825-8be4-406f-b79e-63c993750e7d",
		Name:      "Blake Miller",
		Role:      "Backend Engineer",
		Email:     "blake_miller@xyz.com",
		Phone:     "+49-4313-7492708",
		Contacted: true,
	},
	"eecb15b3-495d-4977-9801-56cdd93add94": {
		ID:        "eecb15b3-495d-4977-9801-56cdd93add94",
		Name:      "Caroline Sanders",
		Role:      "Full stack Engineer",
		Email:     "caroline_sanders@xyz.com",
		Phone:     "+49-2977-8324914",
		Contacted: false,
	},
	"4d48cb8d-45ac-4035-ac1c-cd31f7e05b3c": {
		ID:        "4d48cb8d-45ac-4035-ac1c-cd31f7e05b3c",
		Name:      "Dahlia Martin",
		Role:      "Mobile Developer",
		Email:     "dahlia_martin@xyz.com",
		Phone:     "+49-6858-4638667",
		Contacted: false,
	},
	"d145860c-977c-4466-9a01-b285b2070e00": {
		ID:        "d145860c-977c-4466-9a01-b285b2070e00",
		Name:      "Ella Walker",
		Role:      "Frontend Engineer",
		Email:     "ella_walker@xyz.com",
		Phone:     "+49-4430-6425473",
		Contacted: false,
	},
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	if _, ok := customers[id]; ok {
		json.NewEncoder(w).Encode(customers[id])
		w.WriteHeader(http.StatusOK) // 200 OK
	} else {
		w.WriteHeader(http.StatusNotFound) // 404 NOT FOUND
		w.Write([]byte("404 - Customer doesn't exists!"))
	}
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newEntry Customer

	if _, ok := customers[mux.Vars(r)["id"]]; ok {
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		w.Write([]byte("409 - Customer already exists!"))
		return
	}

	requestBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(requestBody, &newEntry)

	id := uuid.New().String()
	newEntry.ID = id
	customers[id] = newEntry

	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(newEntry)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var updateCustomer map[string]Customer
	requestBody, _ := io.ReadAll(r.Body)
	id := mux.Vars(r)["id"]

	if _, ok := customers[id]; ok {
		json.Unmarshal(requestBody, &updateCustomer)
		customers[id] = updateCustomer[id]
		w.WriteHeader(http.StatusOK) // 200 OK
		json.NewEncoder(w).Encode(customers[id])
	} else {
		w.WriteHeader(http.StatusNotFound) // 404 NOT FOUND
		w.Write([]byte("404 - Customer doesn't exists!"))
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	if _, ok := customers[id]; ok {
		delete(customers, id)
		w.WriteHeader(http.StatusOK) // 200 OK
		json.NewEncoder(w).Encode(customers)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404 NOT FOUND
		w.Write([]byte("404 - Customer doesn't exists!"))
	}
}

func main() {
	router := mux.NewRouter()

	fileServer := http.FileServer(http.Dir("./static"))

	// Home page
	router.HandleFunc("/", fileServer.ServeHTTP)

	// Get all customers
	router.HandleFunc("/customer", getCustomers).Methods("GET")
	// GET - http://localhost:3000/customer (success)

	// Get a single customer through id
	router.HandleFunc("/customer/{id}", getCustomer).Methods("GET")
	// GET - http://localhost:3000/customer/8f1f4cb1-aff0-4094-9398-bbf2dc89d6ad (success) - existing user
	// GET - http://localhost:3000/customer/8f1f4cb1-aff0-4094-9398-bbf2dc89d600 (error) - non existing user

	// Add a customer
	router.HandleFunc("/customer", addCustomer).Methods("POST")
	// POST - http://localhost:3000/customer (success) - non existing user

	// {
	// 	"name": "Faith Bell",
	// 	"role": "Cloud Engineer",
	// 	"email": "faith_bell@xyz.com",
	// 	"phone": "+49-6858-4638667",
	// 	"contacted": true
	// }

	// Update a customer
	router.HandleFunc("/customer/{id}", updateCustomer).Methods("PUT")
	// PUT - http://localhost:3000/customer/8f1f4cb1-aff0-4094-9398-bbf2dc89d600 (error) - non existing user
	// PUT - http://localhost:3000/customer/8f1f4cb1-aff0-4094-9398-bbf2dc89d6ad (success) - existing user

	// {"8f1f4cb1-aff0-4094-9398-bbf2dc89d6ad": {
	//     "id": "8f1f4cb1-aff0-4094-9398-bbf2dc89d6ad",
	//     "name": "Aria Bell",
	//     "role": "Mobile Developer",
	//     "email": "aria_bell@xyz.com",
	//     "phone": "+49-7001-3206111",
	//     "contacted": false
	// }}

	// Delete a customer
	router.HandleFunc("/customer/{id}", deleteCustomer).Methods("DELETE")
	// DELETE - http://localhost:3000/customer/8f1f4cb1-aff0-4094-9398-bbf2dc89d6ad (success) - existing user
	// DELETE - http://localhost:3000/customer/8f1f4cb1-aff0-4094-9398-bbf2dc89d600 (error) - non existing user

	fmt.Println("Server is starting on port 3000...")
	http.ListenAndServe(":3000", router)

}
