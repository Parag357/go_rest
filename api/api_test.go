package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"rest/mocks"
	"rest/model"
	"strconv"
	"testing"
	"time"
)

func TestCreateFailureWithNoName(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/create",ctrl.CreateProd).Methods("POST")

	prod := &model.Product{
		Name: "",
		Price: 34.0,
		CategoryId: 1,
	}
	jprod, _ := json.Marshal(prod)
	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jprod))
	resp := httptest.NewRecorder()
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code, "Bad Request is expected")
}

func TestCreateFailureWithNoPrice(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/create",ctrl.CreateProd).Methods("POST")

	prod := &model.Product{
		Name: "prod11",
		Price: 0,
		CategoryId: 1,
	}
	jprod, _ := json.Marshal(prod)
	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jprod))
	resp := httptest.NewRecorder()
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code, "Bad Request is expected")
}
func TestCreateFailureWithInvalidPrice(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/create",ctrl.CreateProd).Methods("POST")

	prod := &model.Product{
		Name: "prod11",
		Price: -89,
		CategoryId: 1,
	}
	jprod, _ := json.Marshal(prod)
	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jprod))
	resp := httptest.NewRecorder()
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code, "Bad Request is expected")
}

func TestCreateFailureWithNoCategory(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/create",ctrl.CreateProd).Methods("POST")

	prod := &model.Product{
		Name: "prod11",
		Price: 45,
		CategoryId: 0,
	}
	jprod, _ := json.Marshal(prod)
	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jprod))
	resp := httptest.NewRecorder()
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code, "Bad Request is expected")
}

func TestCreateFailureWithDuplicateName(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	prod := &model.Product{
		Name: "prod1",
		Price: 34,
		CategoryId: 1,
	}
	mockDatastore.EXPECT().Create(prod).Return(errors.New("duplicate key value violates unique constraint"))
	jprod, _ := json.Marshal(prod)
	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jprod))
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/create",ctrl.CreateProd).Methods("POST")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code, "Bad Request is expected")
}

func TestCreateSuccess(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	prod := &model.Product{
		Name: "prod100",
		Price: 34,
		CategoryId: 1,
	}
	mockDatastore.EXPECT().Create(prod).Return(nil)
	jprod, _ := json.Marshal(prod)
	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(jprod))
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/create",ctrl.CreateProd).Methods("POST")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 201, resp.Code, "Created successfully is expected")
}

func TestDeleteFailureWithInvalidID(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	//prod := &model.Product{}
	i:= strconv.Itoa(4200) // 2nd arg of delete is string not int
	mockDatastore.EXPECT().Delete(&model.Product{},i)
	//jprod, _ := json.Marshal(prod)
	req, _ := http.NewRequest("DELETE", "/delete/4200",nil)
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true) // need to specify router because server is not running while creaing tests
	myRouter.HandleFunc("/delete/{id}",ctrl.DeleteProd).Methods("DELETE")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code, "deleted successfully is expexted")
}

func TestDeleteSuccessWithValidID(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	//prod := &model.Product{}
	i:= strconv.Itoa(2) // 2nd arg of delete is string not int
	mockDatastore.EXPECT().Delete(&model.Product{},i)
	//jprod, _ := json.Marshal(prod)
	req, _ := http.NewRequest("DELETE", "/delete/2",nil)
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true) // need to specify router because server is not running while creaing tests
	myRouter.HandleFunc("/delete/{id}",ctrl.DeleteProd).Methods("DELETE")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code, "deleted successfully is expexted")
}

func TestUpdateFailureWithInvalidId(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	prod := &model.Product{}
	query:="id = ?"
	i:= strconv.Itoa(20000)
	mockDatastore.EXPECT().GetProductForUpdate(query,i,prod).Return(errors.New("product is not available"))
	newprod := &model.Product{
		Name: "prod32",
		Price: 55,
		CategoryId: 2,
	}
	jprod, _ := json.Marshal(newprod)
	req, _ := http.NewRequest("PUT", "/update/20000", bytes.NewBuffer(jprod))
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/update/{id}",ctrl.UpdateProd).Methods("PUT")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code, "Not Found is expected")
}
func TestUpdateFailureWithInvalidNewPrice(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	prod := &model.Product{}
	query:="id = ?"
	i:= strconv.Itoa(2)
	mockDatastore.EXPECT().GetProductForUpdate(query,i,prod).Return(nil)
	newprod := &model.Product{
		Name: "prod32",
		Price: -55,
		CategoryId: 2,
	}
	jprod, _ := json.Marshal(newprod)
	req, _ := http.NewRequest("PUT", "/update/2", bytes.NewBuffer(jprod))
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/update/{id}",ctrl.UpdateProd).Methods("PUT")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code, "Bad Request is expected")
}

func TestUpdateFailureWithDuplicateNewName(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	prod := &model.Product{}
	query:="id = ?"
	i:= strconv.Itoa(2)
	mockDatastore.EXPECT().GetProductForUpdate(query,i,prod).Return(nil)
	newprod := &model.Product{
		Name: "prod32",
		Price: 55,
		CategoryId: 2,
	}
	mockDatastore.EXPECT().Save(newprod).Return(errors.New("duplicate key value violates unique constraint"))
	jprod, _ := json.Marshal(newprod)
	req, _ := http.NewRequest("PUT", "/update/2", bytes.NewBuffer(jprod))
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/update/{id}",ctrl.UpdateProd).Methods("PUT")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 400, resp.Code, "Bad Request is expected")
}
func TestUpdateSuccess(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	prod := &model.Product{}
	query:="id = ?"
	i:= strconv.Itoa(2)
	mockDatastore.EXPECT().GetProductForUpdate(query,i,prod).Return(nil)
	newprod := &model.Product{
		Name: "prod32",
		Price: 55,
		CategoryId: 2,
	}
	mockDatastore.EXPECT().Save(newprod).Return(nil)
	jprod, _ := json.Marshal(newprod)
	req, _ := http.NewRequest("PUT", "/update/2", bytes.NewBuffer(jprod))
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/update/{id}",ctrl.UpdateProd).Methods("PUT")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 201, resp.Code, "Updated successfully is expected")
}

func TestGetFailureWithWrongCatgoryId(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	req, _ := http.NewRequest("GET", "/get", nil)
	q := req.URL.Query()
	q.Add("categoryId", "30")
	req.URL.RawQuery = q.Encode()
	mockDatastore.EXPECT().GetCategorisedProducts(q).Return([]model.Product{})
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/get",ctrl.ListProd).Methods("GET")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 404, resp.Code, "Not Found is expected")
}

func TestGetSuccessWithNoValues(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	req, _ := http.NewRequest("GET", "/get", nil)
	q := req.URL.Query()
	mockDatastore.EXPECT().GetCategorisedProducts(q).Return([]model.Product{{3, "prod120", 100, time.Time{}, 3}})
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/get",ctrl.ListProd).Methods("GET")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code, "OK is expected")
}


func TestGetSuccessWithIdOnly(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	req, _ := http.NewRequest("GET", "/get", nil)
	q := req.URL.Query()
	q.Add("categoryId", "3")
	req.URL.RawQuery = q.Encode()
	mockDatastore.EXPECT().GetCategorisedProducts(q).Return([]model.Product{{3, "prod120", 100, time.Time{}, 3}})
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/get",ctrl.ListProd).Methods("GET")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code, "OK is expected")
}

func TestGetSuccessWithSortOnly(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	req, _ := http.NewRequest("GET", "/get", nil)
	q := req.URL.Query()
	q.Add("sort", "price")
	req.URL.RawQuery = q.Encode()
	mockDatastore.EXPECT().GetCategorisedProducts(q).Return([]model.Product{{3, "prod120", 100, time.Time{}, 3}})
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/get",ctrl.ListProd).Methods("GET")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code, "OK is expected")
}

func TestGetSuccessWithSortAndOrder(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	req, _ := http.NewRequest("GET", "/get", nil)
	q := req.URL.Query()
	q.Add("sort", "price")
	q.Add("order", "desc")
	req.URL.RawQuery = q.Encode()
	mockDatastore.EXPECT().GetCategorisedProducts(q).Return([]model.Product{{3, "prod120", 100, time.Time{}, 3}})
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/get",ctrl.ListProd).Methods("GET")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code, "OK is expected")
}

func TestGetSuccessWithIdAndSort(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	req, _ := http.NewRequest("GET", "/get", nil)
	q := req.URL.Query()
	q.Add("sort", "price")
	q.Add("categoryId", "3")
	req.URL.RawQuery = q.Encode()
	mockDatastore.EXPECT().GetCategorisedProducts(q).Return([]model.Product{{3, "prod120", 100, time.Time{}, 3}})
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/get",ctrl.ListProd).Methods("GET")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code, "OK is expected")
}

func TestGetSuccessWithAllValues(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDatastore := mocks.NewMockDatastore(mockCtrl)
	ctrl := NewController(mockDatastore)
	req, _ := http.NewRequest("GET", "/get", nil)
	q := req.URL.Query()
	q.Add("categoryId", "3")
	q.Add("sort", "price")
	q.Add("order", "desc")
	req.URL.RawQuery = q.Encode()
	mockDatastore.EXPECT().GetCategorisedProducts(q).Return([]model.Product{{3, "prod120", 100, time.Time{}, 3}})
	resp := httptest.NewRecorder()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/get",ctrl.ListProd).Methods("GET")
	myRouter.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code, "OK is expected")
}




