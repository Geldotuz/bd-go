package models

import "time"

type Detalle_venta struct {
	VentaID        string    `json:"venta_id"`
	ProductoCodigo string    `json:"producto_codigo"`
	Fecha          time.Time `json:"fecha"`
	Producto       *Producto `json:"producto,omitempty" gorm:"references:Codigo"`
}

type Venta struct {
	ID        string    `json:"id"`
	Fecha     time.Time `json:"fecha"`
	ClienteID string    `json:"cliente_id"`
	EmplID    string    `json:"empl_id"`
}
