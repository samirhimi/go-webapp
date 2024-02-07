package main


// Importer les packages nécessaires :


import (
    "context"
    "fmt"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)


// Établir une connexion à la base de données

ctx := context.TODO()
clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
client, err := mongo.Connect(ctx, clientOptions)
if err != nil {
    log.Fatal(err)
}
defer client.Disconnect(ctx)


// Définir la structure du document "livre" et la collection correspondante

type Livre struct {
    Titre   string
    Auteur  string
    Année   int
    ISBN    string
}

collection := client.Database("votre_base_de_données").Collection("livres")

// Implémenter les opérations CRUD 

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

// Récupérer des livres 

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


// Mettre à jour un livre 
 

func updateLivre(ctx context.Context, collection *mongo.Collection, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
    result, err := collection.UpdateOne(ctx, filter, update)
    if err != nil {
        return nil, err
    }

    return result, nil
}

// Supprimer un livre

func deleteLivre(ctx context.Context, collection *mongo.Collection, filter bson.M) (*mongo.DeleteResult, error) {
    result, err := collection.DeleteOne(ctx, filter)
    if err != nil {
        return nil, err
    }

    return result, nil
}
