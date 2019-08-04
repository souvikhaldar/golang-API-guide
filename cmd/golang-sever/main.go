package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type customer struct {
	CustomerID   int32
	CustomerName string
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	var requestbody customer
	if err := json.NewDecoder(r.Body).Decode(&requestbody); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Println("Data recieved: ", requestbody)
}
func main() {
	http.HandleFunc("/customer", addCustomer)
	log.Fatal(http.ListenAndServe(":8192", nil))
}
