package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zzacong/eleven-go/go-bookstore/pkg/routes"
	_ "gorm.io/driver/mysql"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterBookstoreRoutes(r)
	http.Handle("/", r)
	fmt.Println("Starting server at http://localhost:8000")
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", r))
}
