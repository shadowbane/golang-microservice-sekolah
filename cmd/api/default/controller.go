package _default

import (
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/application"
	"github.com/shadowbane/golang-microservice-sekolah/pkg/server"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

type myJSON struct {
	Id    uint32
	Email string
}

type WhatsAppMessage struct {
	Phone   []string `json:"phone"`
	Message *string  `json:"message"`
	Token   *string  `json:"token"`
}

func Index(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		jsondat := []myJSON{
			{Id: 1, Email: "admin@btp.ac.id"},
			{Id: 2, Email: "yolandra@iteba.ac.id"},
			{Id: 4, Email: "syaif@iteba.ac.id"},
			{Id: 5, Email: "aji@btp.ac.id"},
			{Id: 6, Email: "anita@btp.ac.id"},
			{Id: 7, Email: "gloria@btp.ac.id"},
			{Id: 8, Email: "sandy@btp.ac.id"},
			{Id: 11, Email: "jufri@yayasanvitka.id"},
			{Id: 12, Email: "idnbryan7910@gmail.com"},
			{Id: 13, Email: "rendylee21@gmail.com"},
		}

		response, _ := json.Marshal(jsondat)
		w.Write(response)
	}
}

func Store(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer r.Body.Close()

		// set return header as json
		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // catch unwanted fields
		t := WhatsAppMessage{}

		// anonymous struct type: handy for one-time use
		//t := struct {
		//	Type *string `json:"type"` // pointer so we can test for field absence
		//}{}

		err := decoder.Decode(&t)
		if err != nil {
			server.SendHttpResp(w, err.Error(), 422)
			return
		}

		if t.Token == nil {
			server.SendHttpResp(w, "missing field 'type' from JSON request", 422)
			return
		}

		if t.Message == nil {
			server.SendHttpResp(w, "missing field 'message' from JSON request", 422)
			return
		}

		if t.Phone == nil {
			server.SendHttpResp(w, "missing field 'phone' from JSON request", 422)
			return
		}

		// optional extra check
		if decoder.More() {
			server.SendHttpResp(w, "extraneous data after JSON request", 422)
			return
		}

		if len(t.Phone) == 0 {
			server.SendHttpResp(w, "No Number to Send Message", 422)
			return
		}

		url := "https://sawit.wablas.com/api/send-message"
		var numbers = strings.Join(t.Phone, ", ")
		jsonMessage := map[string]interface{}{
			"phone":   numbers,
			"message": t.Message,
		}

		bytesRepresentation, err := json.Marshal(jsonMessage)
		if err != nil {
			zap.S().Error(err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytesRepresentation))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", *t.Token)
		//req.Header.Set("Authorization", "")

		// create http client
		client := http.Client{}
		resp, _ := client.Do(req)

		defer resp.Body.Close()
		if err != nil {
			zap.S().Error(err.Error())
		}

		if resp.StatusCode != 200 {
			body, _ := ioutil.ReadAll(resp.Body)

			w.WriteHeader(resp.StatusCode)
			w.Write(body)
			zap.S().Error(string(body))
		} else {
			server.SendHttpResp(w, "Message Sent", resp.StatusCode)
		}

		return
	}
}
