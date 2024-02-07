
package main


// Import Packages


import (
    "context"
    "fmt"
    "log"
	"time"
	"encoding/json"
	"net/http"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gorilla/mux"
)


// Connecting to the MongoDB

ctx := context.TODO()
clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
client, err := mongo.Connect(ctx, clientOptions)
if err != nil {
    log.Fatal(err)
}
defer client.Disconnect(ctx)


// Defining the Document structure "Book"

type Livre struct {
    Titre   string
    Auteur  string
    Année   int
    ISBN    string
}

collection := client.Database("votre_base_de_données").Collection("livres")

// Implementing CRUD operations 

nouveauLivre := Livre{
    Titre:   "Le Seigneur des Anneaux",
    Auteur:  "J.R.R. Tolkien",
    Année:   1954,
    ISBN:    "123456789",
}

_, err = collection.InsertOne(ctx, nouveauLivre)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Livre inséré avec succès !")

// Getting Books 

func getLivres(ctx context.Context, collection *mongo.Collection) ([]Livre, error) {
    var livres []Livre

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    if err := cursor.All(ctx, &livres); err != nil {
        return nil, err
    }

    return livres, nil
}


// Updating Books
 

func updateLivre(ctx context.Context, collection *mongo.Collection, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
    result, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return nil, err
    }

    return result, nil
}

// Removing Books

func deleteLivre(ctx context.Context, collection *mongo.Collection, filter bson.M) (*mongo.DeleteResult, error) {
    result, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        return nil, err
    }

    return result, nil
}
