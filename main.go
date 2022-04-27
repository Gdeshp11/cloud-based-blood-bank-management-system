package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	mongodbEndpoint = "mongodb://172.17.0.2:27017" // Find this from the Mongo container
)

type Post struct {
	ID            primitive.ObjectID `bson:"_id"`
	Username      string             `bson:"username"`
	Password      string             `bson:"password"`
	BloodType     string             `bson:"blood_type"`
	ContactNumber string             `bson:"contact_number"`
	CreatedAt     time.Time          `bson:"created_at"`
	Tags         string              `bson:"tags"`
}

var ctx context.Context
var col *mongo.Collection

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	client, err := mongo.NewClient(
		options.Client().ApplyURI(mongodbEndpoint),
	)
	checkError(err)
	// Connect to mongo
	ctx = context.Background()
	err = client.Connect(ctx)
	checkError(err)
	// Disconnect
	defer client.Disconnect(ctx)

	// select collection from database
	col = client.Database("bloodBankDatabase").Collection("users")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("static")))
	mux.HandleFunc("/loginHandler", loginHandler)
	mux.HandleFunc("/RegisterHandler", RegisterHandler)
	mux.HandleFunc("/deleteUser", deleteUser)
	mux.HandleFunc("/updateUserinfo", updateUserinfo)
	mux.HandleFunc("/findDonors", findDonors)
	mux.HandleFunc("/listAllDonors", listDonors)
	log.Fatal(http.ListenAndServe(":8000", mux)) // Listens for curl communication of localhost
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "loginHandler Page")
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")
	fmt.Fprintln(w, "username:", username, "password:", password)


	var hash string

	filter := bson.M{"username": username}

	// find one document
	var p Post
	if err := col.FindOne(ctx, filter).Decode(&p); err != nil {
		fmt.Fprintln(w, "user:",username," is not registered!") // if the item does not exist write and error

	} else {
		fmt.Printf("post: %+v\n", p)
		fmt.Fprintln(w, "hashed password of ",username, " : ",p.Password)
		hash = p.Password
	}

	err := CheckPasswordHash(password,hash)
	if err == true {
		fmt.Fprintln(w, "username or password is incorrect!")
	}
	
}

func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "RegisterHandler Page")
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")
	bloodType := req.FormValue("bloodType")
	contactNumber := req.FormValue("contactNumber")
	fmt.Fprintln(w, "username:", username, "password:", password, "bloodType: ",bloodType, "contactNumber:",contactNumber)
	hashedPassword, err := HashPassword(password)
	checkError(err)

	// Insert one
	res, err := col.InsertOne(ctx, &Post{
		ID:            primitive.NewObjectID(),
		Username :     username,          
		Password :     hashedPassword,            
		BloodType :    bloodType,   
		ContactNumber: contactNumber, 
		CreatedAt:     time.Time,   
		Tag : "bloodDonors"} )
	checkError(err)

	if err == nil {
		fmt.Printf("inserted id: %s\n", res.InsertedID.(primitive.ObjectID).Hex())
	}

}

func deleteUser(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintln(w, "deleteUser Page")
	username := req.FormValue("username")
	res, err := col.DeleteMany(ctx, bson.M{"Username": username})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Fprintln(w, "delete count: ", res.DeletedUser)
	}
}

func updateUserinfo(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintln(w, "updateUserinfo Page")
		fmt.Fprintln(w, "updateUserinfo Page")
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")
	bloodType := req.FormValue("bloodType")
	contactNumber := req.FormValue("contactNumber")
	fmt.Fprintln(w, "username:", username, "password:", password, "bloodType: ",bloodType, "contactNumber:",contactNumber)
	hashedPassword, err := HashPassword(password)
	checkError(err)
	
	filter := bson.D{{"Username", username}}
	

	// Insert one
	res, err := col.UpdateOne(ctx, filter, &Post{
		Username :     username,          
		Password :     hashedPassword,            
		BloodType :    bloodType,   
		ContactNumber: contactNumber, 
		Tag : "bloodDonors" } )
	checkError(err)

	if err == nil {
		fmt.Printf("inserted id: %s\n", res.InsertedID.(primitive.ObjectID).Hex())
	}
	
}

func findDonors(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "findDonors Page")
}

func listDonors(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "listDonors Page")
}
