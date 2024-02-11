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
	Title    string `bson:"Title"`
	Author   string `bson:"Author"`
	Quantity int    `bson:"Quantity"`
}

var client *mongo.Client

func main() {
	// Set up MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Get server port from environment variable or use default
	serverPort := getEnv("SERVER_PORT", "80")
	// Get MongoDB URI from environment variable or use default
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/", Welcome).Methods("GET")
	router.HandleFunc("/books", GetBooks).Methods("GET")
	router.HandleFunc("/newbook", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", DeleteBook).Methods("DELETE")
	log.Println("Connected to MongoDB")
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
	json.NewEncoder(w).Encode("Welcome to the Book Store!")
}

// Get the books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database("booksDB").Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Error finding Books: %v", err)
	}

	var books []book
	if err = cursor.All(ctx, &books); err != nil {
		log.Fatalf("Error decoding Books: %v", err)
	}
	json.NewEncoder(w).Encode("The list of the books is below")
	json.NewEncoder(w).Encode(books)
}

// Create new book
func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook book
	json.NewDecoder(r.Body).Decode(&newBook)
	collection := client.Database("booksDB").Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, newBook)
	if err != nil {
		log.Fatalf("Error inserting book: %v", err)
	}
	json.NewEncoder(w).Encode(newBook)
	log.Printf("Book created: %v", newBook)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Book created successfully")
}

// Delete book by ID

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	collection := client.Database("booksDB").Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Fatalf("Error deleting book: %v", err)
	}

	json.NewEncoder(w).Encode("Book deleted successfully")
}
