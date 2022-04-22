package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var RWLock sync.RWMutex

func main() {
	// db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("static")))
	mux.HandleFunc("/loginHandler", loginHandler)
	mux.HandleFunc("/Register", Register)
	mux.HandleFunc("/RegisterHandler", RegisterHandler)
	// mux.HandleFunc("/addUser", addUser)
	// mux.HandleFunc("/deleteUser", deleteUser)
	// mux.HandleFunc("/updateUserinfo", updateUserinfo)
	// mux.HandleFunc("/findDonors", findDonors)
	// mux.HandleFunc("listDonors", listDonors)
	log.Fatal(http.ListenAndServe(":8000", mux)) // Listens for curl communication of localhost
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "loginHandler Page")
	req.ParseForm()
	username := req.FormValue("username")
	password := req.FormValue("password")
	fmt.Fprintln(w, "username:", username, "password:", password)
}

// type database map[string]dollars // database of items with their polar prices
func Register(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Register Page")
}

func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "RegisterHandler Page")
}
