package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Our User Struct
type User struct {
	gorm.Model
	Name  string
	Email string
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=docker sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()
	var users []User
	db.Find(&users)
	fmt.Println("{}", users)

	json.NewEncoder(w).Encode(users)
}

func newUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=docker sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	db.Create(&User{Name: name, Email: email})
	fmt.Fprintf(w, "New User Successfully Created")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=docker sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	fmt.Println("user to delete", name)

	var user User
	db.Where("name = ?", name).Find(&user)
	db.Delete(&user)

	fmt.Fprintf(w, "Successfully Deleted User")

}

func updateUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=docker sslmode=disable")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	var user User
	db.Where("name = ?", name).Find(&user)

	user.Email = email

	db.Save(&user)
	fmt.Fprintf(w, "Successfully Updated User")

}

func custom(w http.ResponseWriter, r *http.Request) {
	var p User
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		panic("failed to decode")
	}
	fmt.Println(p.Name)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/users", allUsers).Methods("GET")
	myRouter.HandleFunc("/user/{name}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{name}/{email}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/user/{name}/{email}", newUser).Methods("POST")
	myRouter.HandleFunc("/user/custom", custom).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func initialMigration() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=docker sslmode=disable")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{})
}

func main() {
	fmt.Println("Go ORM Tutorial")

	initialMigration()
	// Handle Subsequent requests
	handleRequests()
}
