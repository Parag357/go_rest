package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // switch dialects to change b/w dbs
	"log"
	"net/http"
	"rest/api"
	"rest/datastore"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "go_inventory"
)

func main(){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	datastore := datastore.NewProductDataStore(db)
	ctrl := api.NewController(datastore)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/delete/{id}",ctrl.DeleteProd).Methods("DELETE")
	myRouter.HandleFunc("/get",ctrl.ListProd).Methods("GET")
	myRouter.HandleFunc("/create",ctrl.CreateProd).Methods("POST")
	myRouter.HandleFunc("/update/{id}",ctrl.UpdateProd).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080",myRouter))
}
