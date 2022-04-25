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
	mux.HandleFunc("/RegisterHandler", RegisterHandler)
	mux.HandleFunc("/deleteUser", deleteUser)
	mux.HandleFunc("/updateUserinfo", updateUserinfo)
	mux.HandleFunc("/findDonors", findDonors)
	mux.HandleFunc("/listDonors", listDonors)
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

func addUser(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "addUser Page")
}

func deleteUser(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "deleteUser Page")
}

func updateUserinfo(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "updateUserinfo Page")
}

func findDonors(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "findDonors Page")
}

func listDonors(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "listDonors Page")
}
