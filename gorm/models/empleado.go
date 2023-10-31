package models

type Empleado struct {
	ID       string `json:"id" gorm:"primarykey"`
	Nombre   string `json:"nombre"`
	Telefono string `json:"telefono"`
}

func (Empleado) TableName() string {
	return "empleado"
}
