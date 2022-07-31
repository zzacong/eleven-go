package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zzacong/eleven-go/go-bookstore/pkg/models"
	"github.com/zzacong/eleven-go/go-bookstore/pkg/utils"
	"gorm.io/gorm"
)

var NewBook models.Book

func GetBooks(w http.ResponseWriter, r *http.Request) {
	newBooks := models.GetAllBooks()
	res, _ := json.Marshal(newBooks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookId, err := strconv.ParseInt(params["id"], 0, 0)
	if err != nil {
		fmt.Println("Error while parsing bookId")
	}
	bookDetails, db := models.GetBookById(bookId)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	CreateBook := &models.Book{}
	utils.ParseBody(r, CreateBook)
	b := CreateBook.CreateBook()
	res, _ := json.Marshal(b)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bookId, err := strconv.ParseInt(params["id"], 0, 0)
	if err != nil {
		fmt.Println("error while parsing bookId")
	}
	_, db := models.DeleteBook(bookId)
	if db.RowsAffected < 1 {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	res, _ := json.Marshal(true)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updateBook = &models.Book{}
	utils.ParseBody(r, updateBook)
	params := mux.Vars(r)
	bookId, err := strconv.ParseInt(params["id"], 0, 0)
	if err != nil {
		fmt.Println("error while parsing bookId")
	}
	bookDetails, db := models.GetBookById(bookId)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if updateBook.Name != "" {
		bookDetails.Name = updateBook.Name
	}
	if updateBook.Author != "" {
		bookDetails.Author = updateBook.Author
	}
	if updateBook.Publication != "" {
		bookDetails.Publication = updateBook.Publication
	}
	db.Save(&bookDetails)
	res, _ := json.Marshal(bookDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
