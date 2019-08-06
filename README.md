# Tutorial for building a RESTful HTTP API in Golang    

1. Install postgresql on your system. Follow this [link](https://www.postgresqltutorial.com/install-postgresql/) 
    2. Run `psql -U postgres` on the terminal. (NOTE: `postgres` is a default role automatically created, if it's not you need to create it. Also, my commands are for mac, but other OSs should be pretty similar)  
    3. Create a new database. `create database guide`, I'm naming it guide, you can name it anything.  
    4. Create a simple table with two columns `customer_id` and `customer_name` by running `create table customer(customer_id int,customer_name text);` after connecting the `guide` database (do `\c guide`)  
    5. Create a new package in our golang project for all the database interactions at `/pkg/db`.
    6. We need a third party library for better database handling. Use [govendor](https://github.com/kardianos/govendor) for dependency management hence install it this way:-  
        * At the project root- `govendor init`  
        * `govendor fetch github.com/lib/pq` 
        (Note: If we would have not used govendor for dependency management we could have installed using- `go get -u github.com/lib/pq` but using one is always a better idea)
    6. `init` method in golang is the method that runs first even before running `main` hence we will setup the database connection there. 
        ``` 
        func init() {
            psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
            fmt.Println("conection string: ", psqlInfo)
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
        ```


5. First of all, let's register all the handlers that would be performing the `CRUD` operations. We are using a very efficient third-party router called "gorilla mux".  
    ```
    router := mux.NewRouter()
	router.HandleFunc("/customer", addCustomer).Methods("POST")
	router.HandleFunc("/customer/{id}", updateCustomer).Methods("PUT")
	router.HandleFunc("/customer/{id}", deleteCustomer).Methods("DELETE")
	router.HandleFunc("/customer", fetchCustomers).Methods("GET")
	log.Fatal(http.ListenAndServe(":8192", router))
    ```

6. Now let's write the handler for adding a customer's details.  
    ```
    func addCustomer(w http.ResponseWriter, r *http.Request) {
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

    ```

    What it is doing is, first it is reading the request JSON data and  unmarshalling it into `customer` struct, then it is making `INSERT`  query to the database to add the data. Simple!  
    

7. The code for updating the user is as follows:-   
    ```
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
    ``` 

We pass the data to be updated in the `body` of the request and customer ID of the customer whose details are being updated is passed in the URL.  
(NOTE: later when you see the request you will understand it better, for now focus on the logic.)  
8. Let's try to `DELETE` a resource now.  

    
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
    

In the above code, we are passing the ID of the customer to be deleted from our records. The query for deletion the simple `DELETE` command.    

9. Now finally, let's try to fetch the details of all customer data in JSON format.  
```
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
```
In the above code, we are querying for all the customer records, accessing them one by one and appending to a slice and finally serializing them into JSON using the `Marshall` method.   
