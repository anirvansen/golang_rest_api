package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"net/http"
	"reflect"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"createdOn"`
	UpdatedOn   string  `json:"updatedOn"`
	DeletedOn   string  `json:"deletedOn"`
}

func (p Product) ToJSON(writer io.Writer) error {
	e := json.NewEncoder(writer)
	return e.Encode(p)
}

func (p *Product) FromJSON(body io.Reader) error {
	d := json.NewDecoder(body)
	return d.Decode(p)
}


func getDbConnection(dbProperties map[string]string) (db *sql.DB) {
	db, err := sql.Open(dbProperties["dbDriver"], dbProperties["dbUser"]+":"+dbProperties["dbPass"]+"@/"+dbProperties["dbName"])
	if err != nil {
		panic(err.Error())
	}
	return db
}

type Products []Product

func (p Products) ToJSON(writer http.ResponseWriter) error {
	e := json.NewEncoder(writer)
	return e.Encode(p)
}


func GetProducts(dbProperties map[string]string) (Products,error) {
		var productList Products
		db := getDbConnection(dbProperties)

		defer db.Close()
		result,err := db.Query("SELECT * FROM products")
		if err != nil {
			return nil, err
		}

		var ID int
		var Name,Description,SKU,CreatedOn,UpdatedOn,DeletedOn string
		var Price float32
		for result.Next() {

			err := result.Scan(&ID,&Name,&Description,&Price,&SKU,&CreatedOn,&UpdatedOn,&DeletedOn)
			if err != nil {
				return nil,err
			}
			productList = append(productList,Product{
				ID:          ID,
				Name:        Name,
				Description: Description,
				Price:       Price,
				SKU:         SKU,
				CreatedOn:   CreatedOn,
				UpdatedOn:   UpdatedOn,
				DeletedOn:   DeletedOn,
			})
		}

	return productList,nil

}

func  GetProductById(id int,dbProperties map[string]string) (Product,error) {

	db := getDbConnection(dbProperties)
	defer db.Close()
	result,err := db.Query("SELECT * FROM products where Id=?",id)
	fmt.Println("error after selecting",err)
	if err != nil {
		fmt.Println("SELECT error",err)
		return Product{}, err
	}

	var ID int
	var Name,Description,SKU,CreatedOn,UpdatedOn,DeletedOn string
	var Price float32
	var selectProduct Product
	fmt.Println("came here")
	for result.Next() {
		fmt.Println("is it coming here")
		err := result.Scan(&ID,&Name,&Description,&Price,&SKU,&CreatedOn,&UpdatedOn,&DeletedOn)
		fmt.Println("error while fetching",err)
		if err != nil {
			return Product{},err
		}
		 selectProduct = Product{
			ID:          ID,
			Name:        Name,
			Description: Description,
			Price:       Price,
			SKU:         SKU,
			CreatedOn:   CreatedOn,
			UpdatedOn:   UpdatedOn,
			DeletedOn:   DeletedOn,
		}
	}
	return selectProduct,nil

}


func SaveProduct(p Product, dbProperties map[string]string) error {
	db := getDbConnection(dbProperties)
	defer db.Close()
	insertData, err := db.Prepare("INSERT INTO products(Name,Description,Price,SKU,CreatedOn,UpdatedOn,DeletedOn) VALUES (?,?,?,?,?,?,?)")

	if err != nil {
		return fmt.Errorf("Error inserting the values into products table")
	}

	insertData.Exec(p.Name,p.Description,p.Price,p.SKU,time.Now().UTC().String(),"","")
	return nil
}


func UpdateProduct(id int, p Product, dbProperties map[string]string) error {
	db := getDbConnection(dbProperties)
	defer db.Close()
	fmt.Println(p)
	e := reflect.ValueOf(&p).Elem()
	fieldNames := []string{}
	for i:= 0 ; i < e.NumField() ; i ++ {
		fieldNames = append(fieldNames,e.Type().Field(i).Name)
	}

	//Check If the product exist of not

	p, _ = GetProductById(id, dbProperties)
	if p == (Product{}) {
		return fmt.Errorf("Product is not found")
	}

	//updateRequiredForThisColumns := []string{}
	//for i := range fieldNames {
	//	if p.(fieldNames[i]) != "" {
	//		updateRequiredForThisColumns = append(updateRequiredForThisColumns,fieldNames[i])
	//	}
	//}
	//updateData, err := db.Prepare("INSERT INTO products(Name,Description,Price,SKU,CreatedOn,UpdatedOn,DeletedOn) VALUES (?,?,?,?,?,?,?)")
	//
	//if err != nil {
	//	return fmt.Errorf("Error inserting the values into products table")
	//}
	//
	//updateData.Exec(p.Name,p.Description,p.Price,p.SKU,time.Now().UTC().String(),"","")
	return nil
}



func DeleteProduct(id int,dbProperties map[string]string) error {
	db := getDbConnection(dbProperties)
	defer db.Close()
	p, _ := GetProductById(id, dbProperties)
	if p == (Product{}) {
		return fmt.Errorf("Product is not found")
	}


	deleteData, err := db.Prepare("DELETE FROM products WHERE id=?")

	if err != nil {
		return fmt.Errorf("Error deleting the record")
	}

	deleteData.Exec(id)
	return nil
}