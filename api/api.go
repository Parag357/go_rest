package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io/ioutil"
	"net/http"
	"rest/model"
	//"rest/datastore"
	"strings"
)

type Controller struct {
	datastore model.Datastore
}

func NewController(datastore model.Datastore) Controller{
	return Controller {
		datastore: datastore,
	}
}

func (ctrl Controller) CreateProd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // to send json response
	msg:=make(map[string]string)
	jsn, _ := ioutil.ReadAll(r.Body)
	data := &model.Product{}
	json.Unmarshal(jsn,data)
	err := ValidateForCreate(data)
	if err != nil{
		w.WriteHeader(400)
		msg["error"]=err.Error()
		json.NewEncoder(w).Encode(msg)
	} else{
		err := ctrl.datastore.Create(data)
		if err != nil { // to check if create causes an error
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint"){ // to check if create causes an integrity error
				w.WriteHeader(400)
				msg["error"]="name already exists"
				json.NewEncoder(w).Encode(msg)
			}
		}else{
			w.WriteHeader(201)
			msg["error"]="created successfully"
			json.NewEncoder(w).Encode(msg)
		}
	}
}
func (ctrl Controller) DeleteProd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // to send json response
	msg:=make(map[string]string)
	id := mux.Vars(r)["id"]
	data := &model.Product{}
	ctrl.datastore.Delete(data, id)
	w.WriteHeader(200)
	msg["msg"]="deleted successfully"
	json.NewEncoder(w).Encode(msg)
}

func (ctrl Controller) ListProd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // to send json response
	msg:=make(map[string]string)
	params :=r.URL.Query()
	var prod  = ctrl.datastore.GetCategorisedProducts(params)
	if len(prod) == 0 {
		w.WriteHeader(404)
		msg["error"]="invalid category"
		json.NewEncoder(w).Encode(msg)
	}else{
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(prod)
	}
}

func (ctrl Controller) UpdateProd(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	msg:=make(map[string]string)
	id := mux.Vars(r)["id"]
	data := &model.Product{}
	err:=ctrl.datastore.GetProductForUpdate("id = ?", id,data)
	if err != nil{
		w.WriteHeader(404)
		msg["error"]="product is not available"
		json.NewEncoder(w).Encode(msg)
	}else{
		jsn, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(jsn,data)
		if data.Price < 0{
			w.WriteHeader(400)
			msg["error"]="price is invalid"
			json.NewEncoder(w).Encode(msg)
		}else{
			err = ctrl.datastore.Save(data)
			if err != nil{ // to check if create causes an error
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint"){ // to check if create causes an integrity error
					w.WriteHeader(400)
					msg["error"]="name already exists"
					json.NewEncoder(w).Encode(msg)
				}
			}else{
				w.WriteHeader(201)
				msg["error"]="updated successfully"
				json.NewEncoder(w).Encode(msg)
			}
		}
	}
}

func ValidateForCreate(data *model.Product) (err error){
	if data.Name==""{
		return errors.New("name is missing")
	}else if data.Price <= 0{
		return errors.New("price is missing or invalid")
	}else if data.CategoryId == 0{
		return errors.New("category is missing")
	}else{
		return nil
	}
}
