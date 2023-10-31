package models

type Producto struct {
	Codigo   string           `json:"codigo"`
	Nombre   string           `json:"nombre"`
	Precio   float64          `json:"precio"`
	Detalles []*Detalle_venta `json:"detalles" gorm:"many2many:detalle_venta;"`
}
