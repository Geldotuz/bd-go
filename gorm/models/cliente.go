package models

type Cliente struct {
	ID        string      `json:"id"`
	Nombre    string      `json:"nombre"`
	Direccion string      `json:"direccion"`
	Empleados []*Empleado `json:"empleados" gorm:"many2many:venta;"`
}

func (Cliente) TableName() string {
	return "cliente"
}
