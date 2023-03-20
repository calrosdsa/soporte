package model

type Role int	
const (
	RoleCliente Role = iota
	RoleFuncionario
	RoleClienteAdmin
	RoleFuncionarioAdmin
)

type Estado byte
const (
	Activo Estado = iota
	Inactivo 
	Eliminado
)

type RolUser interface {
	//devuelve true si el  rol es de cliente  cliente = 0
	IsCliente(rol *int) bool
	//devuelve true si el rol es de funcionario funcionario = 1
	IsFuncionario(rol *int) bool
	//devuelve true si el  rol es de cliente admin  cliente_admin = 2
	IsClienteAdmin(rol *int) bool
	//devuelve true si el  rol es de funcionario admin  funcionario_admin = 3
	IsFuncionarioAdmin(rol *int) bool
	//devuelve true si el rol es cliente o cliente admin
	IsClienteRol(rol *int) bool
	//devuelve true si el rol es funcionario o funcionario admin
	IsFuncionarioRol(rol *int) bool
	//devuelve true si es cliente admin o funcionario admin
	IsRolAdmin(rol *int) bool
	//devuelve true si es cliente o funcionario 
	IsUserRol(rol *int) bool
}

func IsClienteRol(rol *int) bool{
	//retorna true si el rol es cliente = 0 o cliente admin = 2
	return *rol == int(RoleCliente) || *rol == int(RoleClienteAdmin)
}