package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Product struct {
	ProductName     string `json:"name"`
	ProductNickname string `json:"nickname"`
}

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "samit@57"
	dbname   = "postgres"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM product")
	if err != nil {
		log.Fatal(err)
	}

	var p1 []Product

	for rows.Next() {
		var product Product
		rows.Scan(&product.ProductName, &product.ProductNickname)
		p1 = append(p1, product)
	}

	productBytes, _ := json.MarshalIndent(p1, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(productBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var p Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO product (product_name, product_nickname) VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, p.ProductName, p.ProductNickname)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func main() {
	http.HandleFunc("/", GETHandler)
	http.HandleFunc("/insert", POSTHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
