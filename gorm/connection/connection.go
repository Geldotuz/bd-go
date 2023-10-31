package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Database struct {
	*gorm.DB
}

func NewDataBase(host, username, password, dbname, port string) (Database, error) {
	const dsn = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable"
	con, err := gorm.Open(postgres.Open(fmt.Sprintf(dsn, host, username, password, dbname, port)),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
			Logger:         logger.Default.LogMode(logger.Info),
		})
	return Database{
		DB: con,
	}, err
}

func (db *Database) RegistrarProducto(w http.ResponseWriter, r *http.Request) {
	mouse := models.Producto{
		Codigo: "005",
		Nombre: "Audifono",
		Precio: 1000,
	}
	teclado := models.Producto{
		Codigo: "002",
		Nombre: "Teclado",
		Precio: 312,
	}
	microfono := models.Producto{
		Codigo: "003",
		Nombre: "Microfono",
		Precio: 450,
	}
	rs := db.Create([]models.Producto{mouse, teclado, microfono})
	//rs := db.Create(mouse)
	if rs.Error != nil {
		log.Println(rs.Error)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mouse)
}

func (db *Database) ListSales(w http.ResponseWriter, r *http.Request) {
	var sales []models.Detalle_venta

	wp := r.Context().Value("withPreload").(bool)

	fmt.Printf("AQUI %v \n", wp)

	ndb := db.DB
	if wp {
		ndb = ndb.Preload("Producto")
	}

	if err := ndb.Find(&sales).Error; err != nil {
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sales)
}

func (db *Database) ConsultarProducto(w http.ResponseWriter, r *http.Request) {

	var productos []models.Producto
	criteria := r.FormValue("criteria")
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")

	var order string

	if criteria == "lower" {
		order = "ASC"
	} else if criteria == "bigger" {
		order = "DESC"
	}

	var limitInt int
	var err error
	var offsetInt int

	/* withPreloadBool, err := strconv.ParseBool(withPreload)
	if err != nil {
		panic(err.Error())
		return
	} */

	/* if withPreloadBool == true {
		cnx := r.Context().Value(withPreloadBool)
	} */

	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if offset != "" {
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	ndb := db.Order(fmt.Sprintf("codigo %s", order))
	if limitInt > 0 {
		ndb = ndb.Limit(limitInt)
	}

	if offsetInt > 0 {
		ndb = ndb.Offset(offsetInt)
	}
	if err := ndb.Find(&productos).Error; err != nil {
		log.Println(err.Error())
		return
	}

	for _, p := range productos {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(p)
	}

}

func (db *Database) ConsultarProductoCodigo(w http.ResponseWriter, r *http.Request) {
	var producto models.Producto
	//codigo := r.FormValue("codigo")

	codigo := mux.Vars(r)["id"]

	rs := db.Where("codigo = ?", codigo).Limit(1).Find(&producto)

	if rs.Error != nil {
		fmt.Println(rs.Error)
		return
	}

	fmt.Println("PRODUCTO", producto)

	if rs.RowsAffected == 0 {
		fmt.Println("no se encontro el registro: ", codigo)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(producto)
}

func (db *Database) UpdateProducto(w http.ResponseWriter, r *http.Request) {
	var producto = models.Producto{Nombre: "Pantalla2", Precio: 50}
	//codigo := r.FormValue("codigo")

	user_id := mux.Vars(r)["id"]

	rs := db.Where("codigo = ?", user_id).Updates(&producto)

	if rs.Error != nil {
		log.Fatal(rs.Error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(producto)
}

func (db *Database) DeleteProducto(w http.ResponseWriter, r *http.Request) {
	var producto models.Producto
	//codigo := r.FormValue("codigo")

	user_id := mux.Vars(r)["id"]

	rs := db.Where("codigo = ?", user_id).Delete(&producto)
	if rs.Error != nil {
		log.Println(rs.Error)
		return
	}
	if rs.RowsAffected == 0 {
		log.Println("No se elimino nada")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(producto)
}

func (db *Database) RegistrarClientes(w http.ResponseWriter, r *http.Request) {
	var cliente models.Cliente

	err := json.NewDecoder(r.Body).Decode(&cliente)
	if err != nil {
		log.Fatalln("Error decoding")
		return
	}

	ct := db.Create(cliente)
	if ct.Error != nil {
		log.Println(ct.Error)
		return
	}
}

func (db *Database) RegistrarEmpleados(w http.ResponseWriter, r *http.Request) {
	var empleado models.Empleado

	err := json.NewDecoder(r.Body).Decode(&empleado)
	if err != nil {
		log.Fatalln("Error decoding")
		return
	}

	ct := db.Create(empleado)
	if ct.Error != nil {
		log.Println(ct.Error)
		return
	}
}

type SalePayLoad struct {
	ClientID   string `json:"client_id"`
	EmployeeID string `json:"employee_id"`
}

func (db *Database) RegistrarVentas(w http.ResponseWriter, r *http.Request) {
	var producto models.Producto
	var cliente models.Cliente
	var empleado models.Empleado
	var payload SalePayLoad
	// por payload que traiga el ID del empleado y el ID del cliente
	idProducto := mux.Vars(r)["idProducto"] // 001

	rs := db.First(&producto, "codigo = ?", idProducto)

	if rs.Error != nil {
		fmt.Print(rs.Error.Error())
		return
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Fatalln("Error decoding")
		return
	}

	rs = db.First(&cliente, "id = ?", payload.ClientID)
	if rs.Error != nil {
		fmt.Print(rs.Error.Error())
		return
	}
	rs = db.First(&empleado, "id = ?", payload.EmployeeID)
	if rs.Error != nil {
		fmt.Print(rs.Error.Error())
		return
	}

	CrearVenta := models.Venta{
		ID:        "4",
		Fecha:     time.Now(),
		ClienteID: payload.ClientID,
		EmplID:    payload.EmployeeID,
	}

	ven := db.Create(CrearVenta)
	if ven.Error != nil {
		log.Println(ven.Error)
	}
	CrearDetalle := models.Detalle_venta{
		VentaID:        CrearVenta.ID,
		ProductoCodigo: producto.Codigo, // idProducto
		Fecha:          time.Now(),
	}
	det := db.Create(CrearDetalle)
	if det.Error != nil {
		log.Println(det.Error)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CrearVenta)
	json.NewEncoder(w).Encode(CrearDetalle)

}

func (db *Database) ClienteCodigo(w http.ResponseWriter, r *http.Request) {
	var cliente models.Cliente

	id := mux.Vars(r)["id"]

	rs := db.DB.Preload("Empleados").First(&cliente, "id = ?", id)

	if rs.Error != nil {
		fmt.Println(rs.Error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cliente)
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//query string withPreload (true,false)
		// SI es true inject en el context el ("withPreload",true)
		// EN el handler si withPreload del contexto es true. Agregar el Preload
		withPreload := r.FormValue("withPreload")

		ctx := r.Context()
		ctx = context.WithValue(ctx, "withPreload", false)
		if withPreload == "true" {
			ctx = context.WithValue(ctx, "withPreload", true)
		}

		log.Println("middleware", r.URL)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (db *Database) RegisterRoutes(r *mux.Router) {
	r.Use(Middleware)
	r.HandleFunc("/registrar", db.RegistrarProducto).Methods(http.MethodPost)
	r.HandleFunc("/consultar", db.ConsultarProducto).Methods(http.MethodGet)
	r.HandleFunc("/consultarC/{id}", db.ConsultarProductoCodigo).Methods(http.MethodGet)
	r.HandleFunc("/delete/{id}", db.DeleteProducto).Methods(http.MethodDelete)
	r.HandleFunc("/update/{id}", db.UpdateProducto).Methods(http.MethodPut)
	r.HandleFunc("/registrarClientes", db.RegistrarClientes).Methods(http.MethodPost)
	r.HandleFunc("/registrarEmpleados", db.RegistrarEmpleados).Methods(http.MethodPost)
	r.HandleFunc("/registrarVentas/{idProducto}", db.RegistrarVentas).Methods(http.MethodPost)
	r.HandleFunc("/listsales", db.ListSales).Methods(http.MethodGet)
	r.HandleFunc("/clienteC/{id}", db.ClienteCodigo).Methods(http.MethodGet)

	// endpoint para listar cliente por codigo x
	// crea un cliente Hugo
	// crea 2 empleados
	// crea 2 ventas, una con cada empleado nuevo creado
}
