package main

import (
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//"golang.org/x/crypto/bcrypt"

	"log"
	"net/http"
	"time"
)

var client *mongo.Client

type Users struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId   string             `json:"UserId,omitempty" bson:"UserId,omitempty"`
	Name     string             `json:"Name,omitempty" bson:"Name,omitempty"`
	Email    string             `json:"Email",omitempty" bson:"Email,omitempty"`
	Password string             `json:"pwd,omitempty" bson:"pwd,omitempty"`
}
type Post struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId  string             `json:"UserId,omitempty" bson:"UserId,omitempty"`
	Caption string             `json:"caption,omitempty" bson:"caption,omitempty"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var user Users
	_ = json.NewDecoder(r.Body).Decode(&user)
	collection := client.Database("Instadup").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(w).Encode(result)

}
func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var post Post
	_ = json.NewDecoder(r.Body).Decode(&post)
	collection := client.Database("Instadup").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, post)
	json.NewEncoder(w).Encode(result)

}
func listusers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Users
	var userId string
	collection := client.Database("Instadup").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"userId": userId})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Users
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

func main() {

	http.HandleFunc("/users", createUser)
	http.HandleFunc("/post", createPost)
	http.HandleFunc("/list", listusers)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
