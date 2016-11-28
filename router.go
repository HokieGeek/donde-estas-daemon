package dondeestas

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func UpdatePersonHandler(log *log.Logger, w http.ResponseWriter, r *http.Request) {
	log.Println("UpdatePersonHandler()")

	var update Person

	if err := getJson(r, &update); err != nil {
		postJson(w, 422, err)
	} else {
		// TODO: Finish handling person update
		log.Printf("Received update for person with id: %d\n", update.Id)
		w.WriteHeader(http.StatusOK)
	}
}

type PersonDataRequest struct {
	Id int `json:"id"`
}

func PersonRequestHandler(log *log.Logger, w http.ResponseWriter, r *http.Request) {
	log.Println("PersonRequestHandler()")

	var req PersonDataRequest

	if err := getJson(r, &req); err != nil {
		postJson(w, 422, err)
	} else {
		// TODO: Finish handling person request
		log.Printf("Received request for person with id: %d\n", req.Id)
		w.WriteHeader(http.StatusOK)
	}
}

func getJson(r *http.Request, data interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return err
	}

	if err := r.Body.Close(); err != nil {
		return err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	return nil
}

func postJson(w http.ResponseWriter, httpStatus int, send interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(send); err != nil {
		panic(err)
	}
}

func New(log *log.Logger, port int) {
	mux := http.NewServeMux()

	// TODO: new
	// person
	// update

	mux.HandleFunc("/update",
		func(w http.ResponseWriter, r *http.Request) {
			UpdatePersonHandler(log, w, r)
		})

	mux.HandleFunc("/person",
		func(w http.ResponseWriter, r *http.Request) {
			PersonRequestHandler(log, w, r)
		})

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
