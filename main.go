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

type userData struct {
	ID            primitive.ObjectID `bson:"_id"`
	Username      string             `bson:"username"`
	Password      string             `bson:"password"`
	BloodType     string             `bson:"blood_type"`
	ContactNumber string             `bson:"contact_number"`
	Location      string             `bson:"location"`
	DonationCount uint16             `bson:"donation_count"`
	CreatedAt     time.Time          `bson:"created_at"`
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
	mux.HandleFunc("/listAllDonors", listAllDonors)
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
	username := req.FormValue("username")
	password := req.FormValue("password")
	// fmt.Fprintln(w, "username:", username, "password:", password)
	data := loginPageData{username}
	var hash string

	filter := bson.M{"username": username}

	// find one document
	var p userData
	if err := col.FindOne(ctx, filter).Decode(&p); err != nil {
		fmt.Fprintln(w, "user:", username, " is not registered") // if the item does not exist write and error
		return
	} else {
		fmt.Printf("userData: %+v\n", p)
		hash = p.Password
	}

	ok := CheckPasswordHash(password, hash)
	if !ok {
		fmt.Fprintln(w, "username or password is incorrect!")
	} else {
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
	var p userData
	if err := col.FindOne(ctx, filter).Decode(&p); err == nil {
		fmt.Fprintln(w, "username:", username, " is not available, please try different username")
		return
	}

	hashedPassword, err := HashPassword(password)
	checkError(err)

	// Insert one
	res, err := col.InsertOne(ctx, &userData{
		ID:            primitive.NewObjectID(),
		Username:      username,
		Password:      hashedPassword,
		BloodType:     bloodType,
		ContactNumber: contactNumber,
		Location:      location,
		CreatedAt:     time.Now(),
		DonationCount: 1,
	})

	if err == nil {
		fmt.Printf("inserted id: %s\n", res.InsertedID.(primitive.ObjectID).Hex())
		fmt.Fprintln(w, "user:", username, " is registered successfully!")
	} else {
		fmt.Fprintln(w, "Error in Registration, Please try again")
	}

}

func deleteUser(w http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	fmt.Fprintln(w, "username: ", username)
	res, err := col.DeleteMany(ctx, bson.M{"username": username})
	if err != nil {
		log.Fatal(err)
	} else if res.DeletedCount == 0 {
		fmt.Fprintln(w, "Account is already deleted!")
	} else if res.DeletedCount > 0 {
		fmt.Fprintln(w, "Account Deleted Successfully ")
	}
}

func updateUserinfo(w http.ResponseWriter, req *http.Request) {
	// fmt.Fprintln(w, "updateUserinfo Page")
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")
	contactNumber := req.FormValue("contactNumber")
	location := req.FormValue("locations")
	fmt.Println("username", username, "password:", password, "contactNumber:", contactNumber, "location:", location)
	hashedPassword, err := HashPassword(password)
	checkError(err)

	filter := bson.D{{"username", username}}

	res, err := col.UpdateOne(ctx, filter, &userData{
		Password:      hashedPassword,
		ContactNumber: contactNumber,
		Location:      location,
	})

	if err == nil {
		fmt.Println("update count: ", res.ModifiedCount)
		fmt.Fprintln(w, "Update Successful!")
	} else {
		fmt.Fprint(w, "Can't Update, Please try again")
	}

}

func listAllDonors(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "listDonors Page")
}

func requestBlood(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	location := req.FormValue("locations")
	bloodType := req.FormValue("bloodtype")
	// fmt.Println("bloodType: ", bloodType, "location:", location)

	filter := bson.M{"blood_type": bloodType, "location": location}
	// opts := options.FindOneAndUpdate().SetSort(bson.D{{"donation_count", -1}})
	// update := bson.M{"$inc": bson.M{"donation_count": -1}}

	var donorInfo userData
	// err := col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&donorInfo)
	curr, err := col.Find(ctx, filter)
	defer curr.Close(ctx)
	// err = curr.Decode(&donorInfo)
	if err != nil {
		fmt.Fprintln(w, "No results found for requested search criteria")
		return
	} else {
		fmt.Fprintln(w, "Blood Donors Available, Please find details below:")
		for curr.Next(ctx) {
			err := curr.Decode(&donorInfo)
			if err != nil {
				fmt.Fprintln(w, "Error in Decoding")
				return
			} else {
				fmt.Fprintln(w, "---------------------+---------------------")
				fmt.Fprintln(w, "Contact:", donorInfo.ContactNumber, "Location:", donorInfo.Location)
			}
		}
	}

	// if err != nil {
	// 	fmt.Fprintln(w, "No results found for requested search criteria")

	// } else {

	// 	if donorInfo.DonationCount > 0 {
	// 		fmt.Println("Updated Donation count after request blood: ", donorInfo.DonationCount)
	// 		fmt.Fprintln(w, "Blood Donors Available, Please find details below:")
	// 		fmt.Fprintln(w, "Contact:", donorInfo.ContactNumber, "Location:", donorInfo.Location)
	// 		fmt.Fprintln(w, donorInfo)
	// 	} else {
	// 		fmt.Fprintln(w, "No results found for requested search criteria")
	// 	}
	// }

}

func makeDonation(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprintln(w, "listDonors Page")
	req.ParseForm()
	username := req.FormValue("username")
	// opts := options.FindOneAndUpdate().SetSort(bson.D{{"donation_count", 1}})
	update := bson.M{"$inc": bson.M{"donation_count": +1}}

	var donorInfo userData

	err := col.FindOneAndUpdate(ctx, bson.M{"username": username}, update).Decode(&donorInfo)
	if err != nil {
		fmt.Fprintln(w, "Could not update donation count for user:", username)
	} else {
		fmt.Fprintln(w, "make donation request successful! \nupdated donation count for ", username, " : ", donorInfo.DonationCount)
	}

}
