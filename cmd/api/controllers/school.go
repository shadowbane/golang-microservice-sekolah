package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/shadowbane/golang-microservice-sekolah/cmd/models"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/application"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
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

		//fmt.Printf("%#v\n", schools)

		response, _ := json.Marshal(formattedValues)
		_, err := w.Write(response)
		if err != nil {
			zap.S().Errorf(err.Error())
			writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}

		if app.Cfg.GetAppEnv() != "production" {
			PrintMemUsage()
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
			zap.S().Errorf("Error %d: %s", http.StatusNotFound, err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)

			if app.Cfg.GetAppEnv() != "production" {
				PrintMemUsage()
			}

			return
		}

		var formattedValues returnData
		formattedValues.Success = true
		formattedValues.Data = school

		//fmt.Printf("%#v\n", formattedValues)

		response, _ := json.Marshal(formattedValues)
		_, err := w.Write(response)
		if err != nil {
			zap.S().Errorf(err.Error())
			writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}

		if app.Cfg.GetAppEnv() != "production" {
			PrintMemUsage()
		}
	}
}

func SchoolCreate(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		var request schoolData

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
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
			zap.S().Errorf(err.Error())
			writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
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
			writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
			return
		}

		app.DB.Table("schools").Where("uuid = ?", p.ByName("id")).Update(&request)

		var formattedValues returnData
		formattedValues.Success = true
		formattedValues.Data = request

		response, _ := json.Marshal(formattedValues)
		_, err = w.Write(response)
		if err != nil {
			zap.S().Errorf(err.Error())
			writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

func SchoolDelete(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

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
			zap.S().Errorf(err.Error())
			writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}
}

// Writes the error response as a Standard API JSON response with a response code
func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errorCode)
	err := json.
		NewEncoder(w).Encode(&JsonErrorResponse{Error: &ApiError{Status: errorCode, Title: errorMsg}})

	if err != nil {
		zap.S().Fatalf(err.Error())
	}
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

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	alloc := bToMb(m.Alloc)
	total := bToMb(m.TotalAlloc)
	sys := bToMb(m.Sys)

	zap.S().Debugf("Alloc = %vMB, TotalAlloc = %vMB, Sys = %vMB", alloc, total, sys)

	//// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	//fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	//fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	//fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	//fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
