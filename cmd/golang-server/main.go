package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "guide"
)

var dbDriver *sql.DB

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	var err error
	dbDriver, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = dbDriver.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to postgres!")
}

type customer struct {
	CustomerID   int32
	CustomerName string
}
type customerUpdate struct {
	CustomerName string
}

// Adding a new customer to the database
func addCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Adding a new customer")
	var requestbody customer
	if err := json.NewDecoder(r.Body).Decode(&requestbody); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("Data recieved: ", requestbody)
	if _, err := dbDriver.Exec("INSERT INTO customer(customer_id,customer_name) VALUES($1,$2)", requestbody.CustomerID, requestbody.CustomerName); err != nil {
		fmt.Println("Error in inserting to the database")
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintln(w, "Successfully inserted: ", requestbody)
}

// Update the details of a customer
func updateCustomer(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	idS := v["id"]
	id, _ := strconv.Atoi(idS)
	fmt.Println("Updating customer: ", id)
	var requestbody customerUpdate
	if err := json.NewDecoder(r.Body).Decode(&requestbody); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if _, err := dbDriver.Exec("UPDATE customer set customer_name=$1 where customer_id=$2", requestbody.CustomerName, id); err != nil {
		fmt.Println("Error in updating: ", err)
		http.Error(w, err.Error(), 500)
	}
	fmt.Fprintln(w, "Succesfully updated user")
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id := v["id"]
	fmt.Println("Deleting user: ", id)
	if _, err := dbDriver.Exec("DELETE FROM customer where customer_id=$1", id); err != nil {
		fmt.Println("Unable to delete the customer: ", err)
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintln(w, "Successfully deleted!")
}

func fetchCustomers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching all customers")
	rows, err := dbDriver.Query("SELECT * from customer")
	if err != nil {
		fmt.Println("Unable to read the table: ", err)
		http.Error(w, err.Error(), 500)
		return
	}
	var customers []customer
	defer rows.Close()
	for rows.Next() {
		var c customer
		if err := rows.Scan(&c.CustomerID, &c.CustomerName); err != nil {
			fmt.Println("Unable to scan")
		}
		customers = append(customers, c)
	}
	customerJSON, err := json.Marshal(customers)
	if err != nil {
		fmt.Println("Unable to marshall the data: ", err)
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("Customers: ", customers)
	fmt.Fprintln(w, string(customerJSON))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/customer", addCustomer).Methods("POST")
	router.HandleFunc("/customer/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customer/{id}", deleteCustomer).Methods("DELETE")
	router.HandleFunc("/customer", fetchCustomers).Methods("GET")
	log.Fatal(http.ListenAndServe(":8192", router))
}
