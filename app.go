package main

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"strings"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)
var url string


const dbString = "localhost:27017"
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/SecretThing/{key}", func(w http.ResponseWriter, r *http.Request) {
		DecorateWithLog(SecretThingEndpoint)(w, r)
	})
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	fmt.Println(os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func GetThing(thingKey string)(SecretThing, error) {
	uri := GetDbString()
    if uri == "" {
            fmt.Println("no connection string provided")
            os.Exit(1)
    }
	dbName := GetDbName()
    sess, err := mgo.Dial(uri)
    if err != nil {
            fmt.Printf("Can't connect to mongo, go error %v\n", err)
            os.Exit(1)
    }
    defer sess.Close()

    //sess.SetSafe(&mgo.Safe{})
    fmt.Println("About to hit db. " + thingKey)
    collection := sess.DB(dbName).C("SecretService")

    result := SecretThing{}

    collection.Find(bson.M{"key": thingKey}).One(&result)
    fmt.Println(result)
    return result,nil

}

func PutThing(thing *SecretThing) error {
	uri := GetDbString()
	dbName := GetDbName()
	fmt.Println(uri)
	fmt.Println(GetDbName())
	sess, _ := mgo.Dial(uri)
	defer sess.Close()
	//sess.SetSafe(&mgo.Safe{})
	collection := sess.DB(dbName).C("SecretService")

	collection.Insert(thing)
	return nil
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

	// url := "http://localhost:8091/"

	authString, err := GetAuthenticationString(w, r)
	if err != nil {
		return
	}

	encrypted, err := Encrypt(Hash(authString), thing.Value)
	if err != nil {
		return
	}

	thing.Value = encrypted

	PutThing(thing)
	return
}

func GetSecretThing(w http.ResponseWriter, r *http.Request) (err error) {
	vars := mux.Vars(r)
	key := vars["key"]
	thing, err := GetThing(key)//
	fmt.Println(err)
	if err != nil || thing.Value == nil {
		http.Error(w, "Request returned no results.", 404)
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
	vars := mux.Vars(r)
	key := vars["key"]

	uri := GetDbString()
	sess, _ := mgo.Dial(uri)
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})
	collection := sess.DB("secretservice").C("SecretService")

	collection.Remove(bson.M{"key": key})
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


func GetDbString() string{
	uri := os.Getenv("MONGO_URL")
    if uri == "" {
            fmt.Println("no connection string provided")
            return dbString
    }else{
    	return uri
    }
}

func GetDbName() string{
	name := os.Getenv("MONGODB_DATABASE")
	if name == ""{
		fmt.Println("no db name configured")
		return "secretservice"
	}else{
		return name
	}
}
