package routes

import (
	"github.com/gorilla/mux"
	"github.com/zzacong/eleven-go/go-bookstore/pkg/controllers"
)

var RegisterBookstoreRoutes = func(router *mux.Router) {
	router.HandleFunc("/book", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/book", controllers.GetBooks).Methods("GET")
	router.HandleFunc("/book/{id}", controllers.GetBookById).Methods("GET")
	router.HandleFunc("/book/{id}", controllers.UpdateBook).Methods("PUT")
	router.HandleFunc("/book/{id}", controllers.DeleteBook).Methods("DELETE")
}
