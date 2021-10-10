package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

var client *mongo.Client

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"Username,omitempty" bson:"Username,omitempty"`
	Password string             `json:"Password,omitempty" bson:"Password,omitempty"`
	Email    string             `json:"Email,omitempty" bson:"Email,omitempty"`
}

type Post struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption   string             `json:"Caption,omitempty" bson:"Caption,omitempty"`
	URL       string             `json:"url,omitempty" bson:"url,omitempty"`
	Timestamp string             `json:"Timestamp,omitempty" bson:"Timestamp,omitempty"`
	UserId    string             `json:"userid,omitempty" bson:"userid,omitempty"`
}

type LastId struct {
	Lastid     string `json:"lastid,omitempty" bson:"lastid,omitempty"`
	PagedPosts []Post `json:"pagedposts,omitempty" bson:"pagedposts,omitempty"`
}

func Encrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

func Decrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user User
	_ = json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("Insta").Collection("User")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	key, err := GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	ciphertext, err := Encrypt(key, []byte(user.Password))
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(ciphertext)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}

func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	// params := mux.Vars(request)
	// id, _ := primitive.ObjectIDFromHex(params["id"])
	id_p := request.URL.Query().Get("id")
	id, _ := primitive.ObjectIDFromHex(id_p)
	var user User
	collection := client.Database("Insta").Collection("User")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func GetPostEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	// params := mux.Vars(request)
	// id, _ := primitive.ObjectIDFromHex(params["id"])
	id_p := request.URL.Query().Get("id")
	id, _ := primitive.ObjectIDFromHex(id_p)
	if id.IsZero() {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`post not found`))
		return
	}
	var post Post
	collection := client.Database("Insta").Collection("Post")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Post{ID: id}).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(post)
}

func GetPostsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	// params := mux.Vars(request)
	// id, _ := params["userid"]
	id := request.URL.Query().Get("id")
	limit_p := request.URL.Query().Get("limit")
	limit, _ := strconv.Atoi(limit_p)
	i := 0
	var posts []Post

	collection := client.Database("Insta").Collection("Post")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	var lastid string
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) && i < limit {
		var post Post
		cursor.Decode(&post)
		if post.UserId == id {
			posts = append(posts, post)
			lastid = post.ID.String()
		}
		i += 1
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	lastobj := LastId{PagedPosts: posts, Lastid: lastid}
	ret, err := json.Marshal(lastobj)
	fmt.Fprintf(response, string(ret))
}

func CreatePostEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var post Post
	_ = json.NewDecoder(request.Body).Decode(&post)
	collection := client.Database("Insta").Collection("Post")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, post)
	json.NewEncoder(response).Encode(result)
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	godotenv.Load(".env")
	clientOptions := options.Client().ApplyURI(os.Getenv("link"))
	client, _ = mongo.Connect(ctx, clientOptions)
	http.HandleFunc("/users/", GetUserEndpoint)
	http.HandleFunc("/users", CreateUserEndpoint)
	http.HandleFunc("/posts", CreatePostEndpoint)
	http.HandleFunc("/posts/", GetPostEndpoint)
	http.HandleFunc("/posts/users/", GetPostsEndpoint)
	log.Fatal(http.ListenAndServe(":12345", nil))
}
