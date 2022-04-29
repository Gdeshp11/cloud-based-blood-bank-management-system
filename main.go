package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	mongodbEndpoint = "mongodb://localhost:27017"
)

type Post struct {
	ID            primitive.ObjectID `bson:"_id"`
	Username      string             `bson:"username"`
	Password      string             `bson:"password"`
	BloodType     string             `bson:"blood_type"`
	ContactNumber string             `bson:"contact_number"`
	Location      string             `bson:"location"`
	DonationCount uint16             `bson:"donation_count"`
	CreatedAt     time.Time          `bson:"created_at"`
	Tags          string             `bson:"tags"`
}

var ctx context.Context
var col *mongo.Collection
var tpl *template.Template

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

	tpl, _ = template.ParseGlob("static/*.html")

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("static")))
	mux.HandleFunc("/loginHandler", loginHandler)
	mux.HandleFunc("/registerHandler", registerHandler)
	mux.HandleFunc("/deleteUser", deleteUser)
	mux.HandleFunc("/updateUserinfo", updateUserinfo)
	mux.HandleFunc("/findDonors", findDonors)
	mux.HandleFunc("/listAllDonors", listDonors)
	mux.HandleFunc("/requestBlood", requestBlood)
	mux.HandleFunc("/makeDonation", makeDonation)
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

type loginPageData struct {
	Username string
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	// fmt.Fprintln(w, "loginHandler Page")
	req.ParseForm()
	Username := req.FormValue("username")
	password := req.FormValue("password")
	// fmt.Fprintln(w, "username:", username, "password:", password)
	data := loginPageData{Username}
	var hash string

	filter := bson.M{"username": Username}

	// find one document
	var p Post
	if err := col.FindOne(ctx, filter).Decode(&p); err != nil {
		// fmt.Fprintln(w, "user:", username, " is not registered!") // if the item does not exist write and error

	} else {
		fmt.Printf("post: %+v\n", p)
		// fmt.Fprintln(w, "hashed password of ", username, " : ", p.Password)
		hash = p.Password
	}

	ok := CheckPasswordHash(password, hash)
	if !ok {
		fmt.Fprintln(w, "username or password is incorrect!")
	} else {
		// fmt.Fprintln(w, "Login Successful!")
		// t = template.Must(template.ParseFiles(("static/update.html")))

		tpl.ExecuteTemplate(w, "splash.html", data)
	}

}

func registerHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "RegisterHandler Page")
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")
	contactNumber := req.FormValue("contact")
	bloodType := req.FormValue("bloodtype")
	location := req.FormValue("locations")
	fmt.Println(w, "username:", username, "password:", password, "bloodType: ", bloodType, "contactNumber:", contactNumber, "location:", location)

	//check if username is available
	filter := bson.M{"username": username}
	var p Post
	if err := col.FindOne(ctx, filter).Decode(&p); err == nil {
		fmt.Fprintln(w, "username:", username, " is not available, please try different username")
		return
	}

	hashedPassword, err := HashPassword(password)
	checkError(err)

	// Insert one
	res, err := col.InsertOne(ctx, &Post{
		ID:            primitive.NewObjectID(),
		Username:      username,
		Password:      hashedPassword,
		BloodType:     bloodType,
		ContactNumber: contactNumber,
		Location:      location,
		CreatedAt:     time.Now(),
		Tags:          "bloodDonors",
		DonationCount: 1,
	})

	checkError(err)

	if err == nil {
		fmt.Printf("inserted id: %s\n", res.InsertedID.(primitive.ObjectID).Hex())
		fmt.Fprintln(w, "user:", username, " is registered successfully!")
	}

}

func deleteUser(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintln(w, "deleteUser Page")
	username := req.FormValue("username")
	fmt.Fprintln(w, "username: ", username)
	res, err := col.DeleteMany(ctx, bson.M{"username": username})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Fprintln(w, "delete count: ", res.DeletedCount)
	}
}

func updateUserinfo(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintln(w, "updateUserinfo Page")
	fmt.Fprintln(w, "updateUserinfo Page")
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")
	contactNumber := req.FormValue("contactNumber")
	bloodType := req.FormValue("BloodType")
	location := req.FormValue("location")
	fmt.Fprintln(w, "password:", password, "contactNumber:", contactNumber)
	hashedPassword, err := HashPassword(password)
	checkError(err)

	filter := bson.D{{"username", username}}

	// Insert one
	res, err := col.UpdateOne(ctx, filter, &Post{
		Password:      hashedPassword,
		ContactNumber: contactNumber,
		BloodType:     bloodType,
		Location:      location,
		Tags:          "bloodDonors"})

	if err == nil {
		fmt.Println("update count: ", res.ModifiedCount)
	} else {
		fmt.Fprint(w, "Error:\n", err)
	}

}

func findDonors(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "findDonors Page")
}

func listDonors(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "listDonors Page")
}

func requestBlood(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	location := req.FormValue("locations")
	bloodType := req.FormValue("bloodtype")
	// fmt.Println("bloodType: ", bloodType, "location:", location)

	filter := bson.M{"blood_type": bloodType, "location": location}
	opts := options.FindOneAndUpdate().SetSort(bson.D{{"donation_count", -1}})
	update := bson.M{"$inc": bson.M{"eval": -1}}

	var donorInfo Post
	err := col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&donorInfo)
	if err != nil {
		fmt.Fprintln(w, "No results found for requested search criteria")

	} else {
		if donorInfo.DonationCount > 0 {
			fmt.Println("Updated Donation count after request blood: ", donorInfo.DonationCount)
			fmt.Fprintln(w, "Requested blood is available, please find details below:")
			fmt.Fprintln(w, "Contact:", donorInfo.ContactNumber, "Location:", donorInfo.Location)
		} else {
			fmt.Fprintln(w, "No results found for requested search criteria")
		}
	}

}

func makeDonation(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintln(w, "listDonors Page")
	req.ParseForm()
	username := req.FormValue("username")
	// opts := options.FindOneAndUpdate().SetSort(bson.D{{"donation_count", 1}})
	update := bson.M{"$inc": bson.M{"eval": +1}}

	var donorInfo Post

	err := col.FindOneAndUpdate(ctx, bson.M{"username": username}, update).Decode(&donorInfo)
	if err != nil {
		fmt.Fprintln(w, "Could not update donation count for user:", username)
	} else {
		fmt.Fprintln(w, "make donation request successful! \nupdated donation count for ", username, " : ", donorInfo.DonationCount)
	}

}
