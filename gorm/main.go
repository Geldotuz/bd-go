package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/connection"
	"github.com/gorilla/mux"
)

func errorFalta(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func paginaInicio(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("se logro")
}

/* func handlerRequest(dataBase connection.Database) {
	http.HandleFunc("/registrar", dataBase.RegistrarProducto)
} */

func main() {
	r := mux.NewRouter()

	database, err := connection.NewDataBase("localhost", "root", "root", "gorm-martin", "5432")
	database.RegisterRoutes(r)
	errorFalta(err)
	//database.RegistrarProducto()
	//database.ConsultarProducto()
	//database.ConsultarProductoCodigo("002")
	//database.UpdateProducto()
	//database.DeleteProducto("001")

	fmt.Println("OK")

	r.HandleFunc("/", paginaInicio).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}
	srv.ListenAndServe()
}
