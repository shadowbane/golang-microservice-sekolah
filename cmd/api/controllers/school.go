package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shadowbane/golang-microservice-sekolah/cmd/models"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/application"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type JsonErrorResponse struct {
	Error *ApiError `json:"error"`
}

type ApiError struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}

type CreateSchoolRequest struct {
	Name          string `json:"name" binding:"required"`
	KodeProvinsi  string `json:"kode_provinsi" binding:"required"`
	KodeKabKota   string `json:"kode_kab_kota" binding:"required"`
	KodeKecamatan string `json:"kode_kecamatan" binding:"required"`
	Bentuk        string `json:"bentuk" binding:"required"`
	Status        string `json:"status" binding:"required"`
	AlamatJalan   string `json:"alamat_jalan"`
}

func SchoolIndex(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(r.Body)

		w.Header().Set("Content-Type", "application/json")

		var schools []models.School
		app.DB.Find(&schools)

		//fmt.Printf("%#v\n", schools)

		response, _ := json.Marshal(schools)
		w.Write(response)
	}
}

func SchoolCreate(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var request models.School

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
			return
		}

		request.UUID, _ = uuid.NewUUID()

		app.DB.Create(&request)

		response, _ := json.Marshal(request)
		w.Write(response)
	}
}

// Writes the error response as a Standard API JSON response with a response code
func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	json.
		NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: errorCode, Title: errorMsg}})
}

//Populates a model from the params in the Handler
func populateModelFromHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params, model interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, model); err != nil {
		return err
	}
	return nil
}
