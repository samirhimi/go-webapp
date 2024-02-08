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

// type book struct {
// 	ID       string `json:"id,omitempty"`
// 	Title string `json:"book1,omitempty"`
// 	Author string `json:"book2,omitempty"`
// 	Quantity string `json:"book3,omitempty"`
// }

   type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
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
	router.HandleFunc("/", GetQuestion).Methods("GET")
	router.HandleFunc("/", SubmitAnswer).Methods("POST")

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

// Question 1
func GetQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode("What is your favorite programming language?")
    json.NewEncoder(w).Encode("Please enter the book informations")

}

// Answer 1
func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	var answer book
	_ = json.NewDecoder(r.Body).Decode(&answer)

	log.Printf("answer.Title: %v", answer.Title)
	log.Printf("answer.Author: %v", answer.Author)
	log.Printf("answer.Quantity: %v", answer.Quantity)

	collection := client.Database("surveyDB").Collection("answers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, bson.M{"Answer1": answer.Title, "Answer2": answer.Author, "Answer3": answer.Author})
	if err != nil {
		log.Fatalf("Error inserting answer: %v", err)
	}

	json.NewEncoder(w).Encode(result.InsertedID)
}