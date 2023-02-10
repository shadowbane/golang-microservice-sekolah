package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shadowbane/golang-microservice-sekolah/cmd/models"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/application"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/server"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type JsonErrorResponse struct {
	Error *ApiError `json:"error"`
}

type ApiError struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}

type returnData struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type schoolData struct {
	UUID          string    `json:"uuid" gorm:"primary_key"`
	Name          string    `json:"name"`
	KodeProvinsi  string    `json:"kode_provinsi"`
	KodeKabKota   string    `json:"kode_kab_kota"`
	KodeKecamatan string    `json:"kode_kecamatan"`
	NPSN          string    `json:"npsn"`
	Bentuk        string    `json:"bentuk"`
	Status        string    `json:"status"`
	AlamatJalan   string    `json:"alamat_jalan"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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

		var schools []schoolData
		app.DB.
			Select("uuid, name, kode_provinsi, kode_kab_kota, kode_kecamatan, npsn, bentuk, status, alamat_jalan, created_at, updated_at").
			Table("schools").
			Scan(&schools)

		var formattedValues returnData
		formattedValues.Success = true
		formattedValues.Data = schools

		response, _ := json.Marshal(formattedValues)
		_, err := w.Write(response)
		if err != nil {
			server.SendAndLogError(w, "Internal Server Error", err, http.StatusInternalServerError)

			return
		}
	}
}

func SchoolShow(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(r.Body)

		w.Header().Set("Content-Type", "application/json")

		var school schoolData

		if err := app.DB.Table("schools").Where("uuid = ?", p.ByName("id")).First(&school).Error; err != nil {
			server.SendAndLogError(w, "Not Found", err, http.StatusNotFound)

			return
		}

		var formattedValues returnData
		formattedValues.Success = true
		formattedValues.Data = school

		//fmt.Printf("%#v\n", formattedValues)

		response, _ := json.Marshal(formattedValues)
		_, err := w.Write(response)
		if err != nil {
			server.SendAndLogError(w, "Internal Server Error", err, http.StatusInternalServerError)

			return
		}
	}
}

func SchoolCreate(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		var request schoolData

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			server.SendAndLogError(w, "Unprocessible Entity", err, http.StatusUnprocessableEntity)

			return
		}

		request.UUID = strings.ToUpper(uuid.NewString())

		app.DB.Table("schools").Create(&request)

		var formattedValues returnData
		formattedValues.Success = true
		formattedValues.Data = request

		response, _ := json.Marshal(formattedValues)
		_, err = w.Write(response)
		if err != nil {
			server.SendAndLogError(w, "Internal Server Error", err, http.StatusInternalServerError)

			return
		}
	}
}

func SchoolUpdate(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		// decode update request
		var request schoolData
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			server.SendAndLogError(w, "Unprocessible Entity", err, http.StatusUnprocessableEntity)

			return
		}

		request.UUID = p.ByName("id")
		app.DB.Table("schools").Where("uuid = ?", request.UUID).Update(&request)

		var formattedValues returnData
		formattedValues.Success = true
		formattedValues.Data = request

		response, _ := json.Marshal(formattedValues)
		_, err = w.Write(response)
		if err != nil {
			server.SendAndLogError(w, "Internal Server Error", err, http.StatusInternalServerError)

			return
		}
	}
}

func SchoolDelete(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		var school schoolData
		if err := app.DB.Table("schools").Where("uuid = ?", p.ByName("id")).First(&school).Error; err != nil {
			server.SendAndLogError(w, "Not Found", err, http.StatusNotFound)

			return
		}

		app.DB.Where("uuid = ?", p.ByName("id")).Delete(models.School{})

		formattedValues := map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"message": "School deleted",
			},
		}

		response, _ := json.Marshal(formattedValues)
		_, err := w.Write(response)
		if err != nil {
			server.SendAndLogError(w, "Internal Server Error", err, http.StatusInternalServerError)

			return
		}
	}
}
