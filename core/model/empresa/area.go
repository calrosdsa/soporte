package empresa

import "time"

type Area struct {
	Id          int       `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion *string   `json:"descripcion"`
	Estado      int       `json:"estado"`
	EmpresaId   int       `json:"empresa_id"`
	CreatedOn   time.Time `json:"created_on"`
	CreadorId   string    `json:"creador_id"`
}

type AreaUser struct {
	Id     int `json:"id"`
	Name   string `json:"nombre"`
	Estado int `json:"estado"`
}
