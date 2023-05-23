package repository

import (
	// "log"
	"soporte-go/core/model"
)

type RoleHandler struct {
}


func NewUtil() model.Util {
	return &RoleHandler{}
}


func (r *RoleHandler) IsCliente(rol int) (b bool) {
	return int(model.RoleCliente) == rol
}

func (r *RoleHandler) IsClienteAdmin(rol int) (b bool) {
	return int(model.RoleClienteAdmin) == rol
}

func (r *RoleHandler) IsFuncionario(rol int) (b bool) {
	return int(model.RoleFuncionario) == rol
}

func (r *RoleHandler) IsFuncionarioAdmin(rol int) (b bool) {
	return int(model.RoleFuncionarioAdmin) == rol
}

func (r *RoleHandler) IsAdminFuncionario(rol int) (b bool) {
	return int(model.RoleAdmin) == rol
}

func (r *RoleHandler) IsFuncionarioRol(rol int) (b bool) {
	return int(model.RoleFuncionario) == rol || int(model.RoleFuncionarioAdmin) == rol
}

func (r *RoleHandler) IsClienteRol(rol int) (b bool) {
	return int(model.RoleCliente) == rol || int(model.RoleClienteAdmin) == rol
}

func (r *RoleHandler) IsRolAdmin(rol int) (b bool) {
	return int(model.RoleFuncionarioAdmin) == rol || int(model.RoleClienteAdmin) == rol
}

func (r *RoleHandler) IsUserRol(rol int) (b bool) {
	return int(model.RoleCliente) == rol || int(model.RoleFuncionario) == rol
}

func (r *RoleHandler) PaginationValues(p int,of int)(page int,offset int){
	// log.Println(of)
	if of == 0 {
		offset = 10
	}else{
		offset = of
	}
	if p == 1 || p == 0 {
		page = 0
	} else {
		page = p - 1
	}
	return
}
