package empresa

import "time"

type Area struct {
	Id     int    `json:"id"`
	Nombre string `json:"nombre"`
	// Descripcion *string   `json:"descripcion"`
	Estado    int `json:"estado"`
	EmpresaId int `json:"empresa_id"`
	// EmpresaCli   int       `json:"empresa_id_cli"`
	CreatedOn time.Time `json:"created_on"`
	CreadorId string    `json:"creador_id"`
}

type ProyectoDuration struct {
	Id         int       `json:"id"`
	ProyectoId int       `json:"proyecto_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}

type ProyectoDetail struct {
	Id       int    `json:"id"`
	Nombre   string `json:"nombre"`
	ParentId int    `json:"parent_id"`
	// Descripcion *string   `json:"descripcion"`
	Estado    int `json:"estado"`
	EmpresaId int `json:"empresa_id"`
	// EmpresaCli   int       `json:"empresa_id_cli"`
	EmpresaParentId  int                `json:"empresa_parent_id"`
	CreatedOn        time.Time          `json:"created_on"`
	CreadorId        string             `json:"creador_id"`
	AreaName string `json:"area_name"`
	EmpresaName string `json:"empresa_name"`
	ProyectoDuration []ProyectoDuration `json:"proyecto_duration"`
}

type Proyecto struct {
	Id       int    `json:"id"`
	Nombre   string `json:"nombre"`
	ParentId int    `json:"parent_id"`
	// Descripcion *string   `json:"descripcion"`
	Estado    int `json:"estado"`
	EmpresaId int `json:"empresa_id"`
	// EmpresaCli   int       `json:"empresa_id_cli"`
	EmpresaParentId int       `json:"empresa_parent_id"`
	CreatedOn       time.Time `json:"created_on"`
	CreadorId       string    `json:"creador_id"`
	Start           string    `json:"start"`
	End             string    `json:"end"`
}

type AreaUser struct {
	Id     int    `json:"id"`
	Name   string `json:"nombre"`
	Estado int    `json:"estado"`
}

type UserArea struct {
	UserId   string  `json:"user_id"`
	Estado   int     `json:"estado"`
	Name     string  `json:"nombre"`
	Apellido *string `json:"apellido"`
	Photo    *string `json:"profile_photo"`
}

type ProovedorArea struct {
	Id                 int    `json:"id"`
	EmpresaCliente     int    `json:"em_cliente"`
	EmpresaFuncionaria int    `json:"em_funcionaria"`
	AreaId             int    `json:"area_id"`
	AreaName           string `json:"area_name"`
	Estado             *int   `json:"estado"`
}
