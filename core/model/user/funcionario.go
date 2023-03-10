package user

import "time"


type Funcionario struct {
	FuncionarioId string  `json:"cliente_id,omitempty"`
	Nombre *string  `json:"nombre,omitempty"`   
	Apellido *string `json:"apellido,omitempty"`   
	Email string  `json:"email,omitempty"`   
	Telefono *string `json:"telefono,omitempty"`
	Celular *string  `json:"celular,omitempty"`
	UserId *string  `json:"user_id,omitempty"`
	SuperiorId *string  `json:"superior_id,omitempty"`
	ProfilePhoto *string  `json:"profile_photo,omitempty"`
	Estado int  `json:"estado,omitempty"`
	EmpresaId int `json:"empresa_id,omitempty"`
	// IsAdmin int  `json:"is_admin,omitempty"`
	Areas []int  `json:"areas,omitempty"`
	CreatedOn time.Time `json:"created_on,omitempty"`
	UpdatedOn *time.Time `json:"updated_on,omitempty"`
	Rol          int        `json:"rol"`
}


type FuncionarioResponse struct{
	FuncionarioId string  `json:"funcionario_id,omitempty"`
	Nombre *string  `json:"nombre,omitempty"`   
	Email string  `json:"email,omitempty"`   
	EmpresaId int `json:"empresa_id,omitempty"`
	Estado int  `json:"estado,omitempty"`
	CreatedOn time.Time `json:"created_on,omitempty"`
	UserId *string  `json:"user_id,omitempty"`
}