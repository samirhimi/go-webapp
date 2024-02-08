package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

   type book struct {

    ID       string `bson:"_id,omitempty"`
    Title    string `bson:"title"`
    Author   string `bson:"author"`
    Quantity int    `bson:"quantity"`

}


var client *mongo.Client

func main() {
	// Set up MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Get server port from environment variable or use default
	serverPort := getEnv("SERVER_PORT", "8080")
	// Get MongoDB URI from environment variable or use default
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, _ = mongo.Connect(ctx, clientOptions)

	router := mux.NewRouter()
	
	router.HandleFunc("/", Welcome).Methods("GET")
    
	router.HandleFunc("/books", GetBooks).Methods("GET")
	
	router.HandleFunc("/newbook", createBook).Methods("POST")
	
	log.Printf("Starting server on port %s...", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, router))
}

// Get environment variable or fallback to default

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Welcome message

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode("Welcome to the the BOOKSTORE")
}

// Getting the books

func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode("The list of the books is below")
}


// Create new book

func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook book
	_ = json.NewDecoder(r.Body).Decode(&newBook)
    
	log.Printf("book.Id: %v", newBook.ID)
	log.Printf("book.Title: %v", newBook.Title)
	log.Printf("book.Author: %v", newBook.Author)
	log.Printf("book.Quantity: %v", newBook.Quantity)

	collection := client.Database("booksDB").Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, bson.M{"Id": newBook.ID,"Title": newBook.Title, "Author": newBook.Author, "Quantity": newBook.Quantity})
	if err != nil {
		log.Fatalf("Error inserting Book: %v", err)
	}

	json.NewEncoder(w).Encode(result.InsertedID)
}