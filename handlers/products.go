package handlers

import (
	"fmt"
	"github.com/anirvansen/golang_rest_api/data"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type products struct {
	logger *log.Logger
	dbProperties map[string]string
}

func (p *products) GetProducts(writer http.ResponseWriter, request *http.Request) {
	p.logger.Println("Handle Get Requests")
	pList,err :=data.GetProducts(p.dbProperties)

	if err != nil {
		p.logger.Println("Error while fetching records from database")
		http.Error(writer,"No records found",http.StatusNoContent)
		return
	}
	p.logger.Println("Got result from database")
	err_her := pList.ToJSON(writer)
	if err_her != nil {
		http.Error(writer,"Unable to marshall json",http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *products) GetProductById(writer http.ResponseWriter, request *http.Request) {
	p.logger.Println("Handle Get Requests, GetProductById")
	vars := mux.Vars(request)
	id,_ := strconv.Atoi(vars["id"])
	p_info,err := data.GetProductById(id,p.dbProperties)
	fmt.Println(p_info)
	if err != nil{
		http.Error(writer,"Not able to find the product",http.StatusBadRequest)
		return
	}

	if p_info == (data.Product{}) {
		http.Error(writer,"Not able to find the product",http.StatusBadRequest)
		return
	}

	err_her := p_info.ToJSON(writer)
	if err_her != nil {
		http.Error(writer,"Unable to marshall json",http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *products) SaveProduct(writer http.ResponseWriter, request *http.Request) {
	p.logger.Println("Handle post requests")
	new_product := &data.Product{}
	err := new_product.FromJSON(request.Body)
	if err != nil {
		http.Error(writer,"Unable to un-marshall json",http.StatusInternalServerError)
		return
	}

	save_err := data.SaveProduct(*new_product,p.dbProperties)
	if save_err != nil {
		http.Error(writer,"Unable to write into table,Please try again",http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Successfully saved it to database"))

}

func (p *products) UpdateProductById(writer http.ResponseWriter, request *http.Request) {
	p.logger.Println("Handle PUT requests")

	vars := mux.Vars(request)
	id,_ := strconv.Atoi(vars["id"])
	updated_product := &data.Product{}
	err := updated_product.FromJSON(request.Body)
	if err != nil {
		http.Error(writer,"Unable to un-marshall json",http.StatusInternalServerError)
		return
	}

	update_err := data.UpdateProduct(id,*updated_product,p.dbProperties)
	if update_err != nil {
		http.Error(writer,update_err.Error(),http.StatusBadRequest)
		return
	}


	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Successfully updated the record"))
}

func (p *products) DeleteProductById(writer http.ResponseWriter, request *http.Request) {
	p.logger.Println("Handle Delete requests")

	vars := mux.Vars(request)
	id,_ := strconv.Atoi(vars["id"])
	err := data.DeleteProduct(id,p.dbProperties)

	if err != nil {
		http.Error(writer,err.Error(),http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Successfully deleted the record"))

}

func ProductHandler(logger *log.Logger, properties map[string]string) *products {
	return &products{logger: logger,dbProperties: properties}
}