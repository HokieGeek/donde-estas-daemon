package dondeestas

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type PersonDataRequest struct {
	Ids []string `json:"ids"`
}

type PersonDataResponse struct {
	People []Person `json:"people"`
}

func PersonRequestHandler(log *log.Logger, db *dbclient, w http.ResponseWriter, r *http.Request) {
	var req PersonDataRequest

	if err := getJson(r.Body, &req); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), http.StatusUnprocessableEntity)
	} else {
		log.Printf("Received request for people with ids: %v\n", req.Ids)

		var resp PersonDataResponse
		resp.People = make([]Person, 0)

		for _, id := range req.Ids {
			if person, err := (*db).Get(id); err == nil { // TODO: this pointer dereference
				resp.People = append(resp.People, *person)
			}
		}

		if len(resp.People) == len(req.Ids) {
			postJson(w, http.StatusOK, resp)
		} else {
			postJson(w, http.StatusPartialContent, resp)
		}
	}
}

func UpdatePersonHandler(log *log.Logger, db *dbclient, w http.ResponseWriter, r *http.Request) {
	var update Person

	if err := getJson(r.Body, &update); err != nil {
		http.Error(w, fmt.Sprintf("%s\n", err), http.StatusUnprocessableEntity)
	} else {
		log.Printf("Received update for person with id: %d\n", update.Id)

		var err error
		if (*db).Exists(update.Id) {
			err = (*db).Update(update)
		} else {
			err = (*db).Create(update)
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("%s\n", err), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	}
}

func getJson(requestBody io.ReadCloser, data interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(requestBody, 1048576))
	if err != nil {
		return err
	}

	if err := requestBody.Close(); err != nil {
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

func New(log *log.Logger, port int, db *dbclient) {
	http.HandleFunc("/person",
		func(w http.ResponseWriter, r *http.Request) {
			PersonRequestHandler(log, db, w, r)
		})

	http.HandleFunc("/update",
		func(w http.ResponseWriter, r *http.Request) {
			UpdatePersonHandler(log, db, w, r)
		})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
