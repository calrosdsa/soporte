package empresa

import "time"

type Area struct {
	Id          int     `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion *string   `json:"descripcion"`
	Estado      int       `json:"estado"`
	EmpresaId   int       `json:"empresa_id"`
	EmpresaCli   int       `json:"empresa_id_cli"`
	CreatedOn   time.Time `json:"created_on"`
	CreadorId   string    `json:"creador_id"`
}

type AreaUser struct {
	Id     int    `json:"id"`
	Name   string `json:"nombre"`
	Estado int    `json:"estado"`
}

type ProovedorArea struct {
	Id                 int    `json:"id"`
	EmpresaCliente     int    `json:"em_cliente"`
	EmpresaFuncionaria int    `json:"em_funcionaria"`
	AreaId             int    `json:"area_id"`
	AreaName           string `json:"area_name"`
	Estado             *int   `json:"estado"`
}
