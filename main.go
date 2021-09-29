package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//BOOK STRUCT
// Definuješ typ a data která passnes JSONEM do API

type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"` //Typ Struct pro autor níže
}

//AUTHOR
type Author struct {
	Firstname string `json:firstname`
	Lastname  string `json:lastname`
}

//KNIZKY SLICE - struct
var books []Book

// GET ALL BOOKS -funct Všechny vrací stejné hodnoty w http.ResponseWriter, r *http.Request
func getBooks(w http.ResponseWriter, r *http.Request) {
	// Určuji typ imputu na json
	w.Header().Set("Content-Type", "aplication/json")
	//Vypise seznam knih predavam W coz je pole vsech knih encode podle type books
	json.NewEncoder(w).Encode(books)
}

// Pro vypis knihy
func getBook(w http.ResponseWriter, r *http.Request) {
	// Určuji typ imputu na json
	w.Header().Set("Content-Type", "aplication/json")
	params := mux.Vars(r) // Toto příjmá parametr v současné době ID knihy
	// Prohledání knížek cyklem FOR
	for _, item := range books {
		// zjisteni jestli ID v JSON je shodné s ID v URL
		// POZOR!! CASE SENZITIVE id != ID
		if item.ID == params["ID"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Vytvoreni knihy
func createBook(w http.ResponseWriter, r *http.Request) {
	// Určuji typ imputu na json
	w.Header().Set("Content-Type", "aplication/json")
	var book Book
	// Příjmá vstup z PUT requestu - používá se DECODER
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) //convert Intu na string // Vytvoření random INT pro ID knihy // Nepouzivat v produkci
	books = append(books, book)                 // Zapíše aktuálně sezbíraná data do BOOKS global VAR
	json.NewEncoder(w).Encode(book)

}

// Update data v knize
func updateBook(w http.ResponseWriter, r *http.Request) {
	// Určuji typ imputu na json
	w.Header().Set("Content-Type", "aplication/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["ID"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			// Příjmá vstup z PUT requestu - používá se DECODER
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["ID"]      //convert Intu na string // Vytvoření random INT pro ID knihy // Nepouzivat v produkci
			books = append(books, book) // Zapíše aktuálně sezbíraná data do BOOKS global VAR
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

// Smazat knihu
func deleteBook(w http.ResponseWriter, r *http.Request) {
	// Určuji typ imputu na json
	w.Header().Set("Content-Type", "aplication/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["ID"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	//INIT ROUTER
	router := mux.NewRouter()

	//DATA
	// Imput dat vždy pomocí append - do Slice (pole) Books
	// vkladam jeden zaznam typu book - vyplnim jednotlive promenne definovane v TYPE
	// Autor se definuje pres pointer na type Autor - pokud odkazuji na typ pres pointer musím použít &
	books = append(books, Book{ID: "1", Isbn: "12345", Title: "Pycha a predsudek", Author: &Author{Firstname: "George", Lastname: "Nameless"}})
	books = append(books, Book{ID: "2", Isbn: "256", Title: "Honzikova Cesta", Author: &Author{Firstname: "Jirka", Lastname: "Horejsi"}})

	//ROUTE HANDLERS
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{ID}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{ID}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{ID}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))

}
