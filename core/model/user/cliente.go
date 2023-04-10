package user

import (
	"time"
)

type Cliente struct {
	Id           string     `json:"id,omitempty"`
	Nombre       *string    `json:"nombre,omitempty"`
	Apellido     *string    `json:"apellido,omitempty"`
	Email        string     `json:"email,omitempty"`
	Telefono     *string    `json:"telefono,omitempty"`
	Celular      *string    `json:"celular,omitempty"`
	UserId       *string    `json:"user_id,omitempty"`
	SuperiorId   *string    `json:"superior_id,omitempty"`
	ProfilePhoto *string    `json:"profile_photo,omitempty"`
	Estado       int        `json:"estado,omitempty"`
	EmpresaId    int        `json:"empresa_id,omitempty"`
	IsAdmin      bool       `json:"is_admin,omitempty"`
	Areas        []int      `json:"areas,omitempty"`
	CreatedOn    time.Time  `json:"created_on,omitempty"`
	UpdatedOn    *time.Time `json:"updated_on,omitempty"`
	Rol          int        `json:"rol"`
}

type UserAuth struct {
	Id        string `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	Estado    int    `json:"estado"`
	Rol       int    `json:"rol"`
	EmpresaId int    `json:"empresa_id,omitempty"`
	Username  string `json:"username"`
}

type Cliente2 struct {
	ClientId     string  `json:"cliente_id,omitempty"`
	Nombre       *string `json:"nombre,omitempty"`
	Apellido     *string `json:"apellido,omitempty"`
	Email        string  `json:"email,omitempty"`
	ProfilePhoto *string `json:"profile_photo,omitempty"`
	Estado       int     `json:"estado"`
	IsAdmin      bool    `json:"is_admin"`
}

type ClienteResponse struct {
	ClientId  string  `json:"cliente_id,omitempty"`
	Nombre    *string `json:"nombre,omitempty"`
	Email     string  `json:"email,omitempty"`
	EmpresaId int     `json:"empresa_id,omitempty"`
	UserId    *string `json:"user_id,omitempty"`
	Rol       int     `json:"rol"`
}

type 	UserShortInfo struct {
	Id        string     `json:"id"`
	Nombre    string     `json:"nombre"`
	Apellido  *string     `json:"apellido"`
	Pendiente bool       `json:"pendiente"`
	IsAdmin   bool       `json:"is_admin"`
	Email     *string    `json:"email"`
	Photo     *string    `json:"profile_photo"`
	Estado    int        `json:"estado"`
	DateTime  *time.Time `json:"created_on"`
}

type UserForList struct {
	Id        string     `json:"id"`
	Nombre    string     `json:"nombre"`
	Apellido  *string     `json:"apellido"`
	Photo     *string    `json:"profile_photo"`
	// Estado    int        `json:"estado"`
}

type UserArea struct {
	Id     string `json:"id"`
	Nombre string `json:"nombre"`
	Estado int    `json:"estado"`
	// AreaId    int    `json:"area_id"`
}

// client_id uuid DEFAULT uuid_generate_v4 (),
// nombre VARCHAR,
// apellido VARCHAR,
// celular VARCHAR,
// email VARCHAR ( 50 ) UNIQUE NOT NULL,
// superior_id VARCHAR,
// empresa_id INT NOT NULL,
// telefono VARCHAR,
// created_on TIMESTAMP NOT NULL,
// updated_on TIMESTAMP,
// user_id VARCHAR NOT NULL,
// is_admin INT DEFAULT 0,
// areas INT[],
// estado INT DEFAULT 0,
// profile_photo TEXT,
// PRIMARY KEY (client_id)
