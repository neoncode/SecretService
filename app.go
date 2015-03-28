package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/neoncode/NoSQLDataAccess"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Set up configuration
	url := os.Getenv("DBPATH")
	if url == "" {
		url = "http://localhost:8091/"
	}

	_ = url
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/SecretThing/{key}", func(w http.ResponseWriter, r *http.Request) {
		DecorateWithLog(SecretThingEndpoint)(w, r)
	})
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func SecretThingEndpoint(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Secret thing endpoint was hit.")
	var err error
	switch r.Method { 
	case "POST":
		err = PutOrPostSecretThing(w, r)
	case "PUT":
		err = PutOrPostSecretThing(w, r)
	case "DELETE":
		err = DeleteSecretThing(w, r)
	case "GET":
		err = GetSecretThing(w, r)
	}
	return err
}

func GetAuthenticationString(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)
	fmt.Println("made it past getting the header")
	if len(auth) != 2 || auth[0] != "Basic" {
		http.Error(w, "bad syntax", http.StatusBadRequest) //This is a good strategy for handling errors
		return nil, errors.New("Bad Syntax")
	}

	payload, err := base64.StdEncoding.DecodeString(auth[1])
	return payload, err
}

func GetSecretThingFromRequest(r *http.Request) (*SecretThing, error) {
	vars := mux.Vars(r)
	key := vars["key"]
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	secretThing := SecretThing{key, bytes}
	return &secretThing, err
}

func PutOrPostSecretThing(w http.ResponseWriter, r *http.Request) (err error) {
	thing, err := GetSecretThingFromRequest(r)
	if err != nil {
		return
	}

	url := "http://localhost:8091/"

	authString, err := GetAuthenticationString(w, r)
	if err != nil {
		return
	}

	encrypted, err := Encrypt(Hash(authString), thing.Value)
	if err != nil {
		return
	}

	thing.Value = encrypted

	couchbase := DataAccess.GetCouchbaseDAL(url, "default", "SecretThing")

	err = couchbase.Set(thing.Key, thing)
	return
}

func GetSecretThing(w http.ResponseWriter, r *http.Request) (err error) {
	url := "http://localhost:8091/"
	thing := new(SecretThing)
	vars := mux.Vars(r)
	key := vars["key"]
	couchbase := DataAccess.GetCouchbaseDAL(url, "default", "SecretThing")
	err = couchbase.Get(key, thing)
	fmt.Println(err)
	if err != nil {
		return
	}

	authString, err := GetAuthenticationString(w, r)
	if err != nil {
		return
	}

	decrypted, err := Decrypt(Hash(authString), thing.Value)
	if err != nil {
		return
	}

	w.Write(decrypted)
	return
}

func DeleteSecretThing(w http.ResponseWriter, r *http.Request) (err error) {
	url := "http://localhost:8091/"
	vars := mux.Vars(r)
	key := vars["key"]
	couchbase := DataAccess.GetCouchbaseDAL(url, "default", "SecretThing")
	err = couchbase.Remove(key)
	return
}

type appHandler func(http.ResponseWriter, *http.Request) error

func DecorateWithLog(fn appHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}
