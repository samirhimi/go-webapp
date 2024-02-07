
package main


// Import Packages


import (
    "context"
    "fmt"
    "log"
	"time"
	"encoding/json"
    "errors"
	"net/http"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gorilla/mux"
)

type book struct{
    ID        string  `json:"id"`
    Title     string  `json:"title"`
    Author    string  `json:"author"`
    Quantity  int     `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}