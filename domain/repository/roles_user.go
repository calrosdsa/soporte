package repository

import (
	"soporte-go/core/model"
)

type RoleHandler struct{
	
}


func NewRolUser() model.RolUser{
	return RoleHandler{}
}

func (r *RoleHandler)IsCliente(rol *int) (b bool){
	return
}

func (r *RoleHandler)IsClienteAdmin(rol *int) (b bool){
	return
}

func (r *RoleHandler)IsFuncionario(rol *int) (b bool){
	return
}

func (r *RoleHandler)IsFuncionarioAdmin(rol *int) (b bool){
	return
}

func (r *RoleHandler)IsFuncionarioRol(rol *int) (b bool){
	return
}

func (r *RoleHandler)IsClienteRol(rol *int) (b bool){
	return
}

func (r *RoleHandler)IsRolAdmin(rol *int) (b bool){
	return
}

func (r *RoleHandler)IsUserRol(rol *int) (b bool){
	return
}