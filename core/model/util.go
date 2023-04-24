package model

type Role int
type Estado byte
type CasoEstado byte
type FileType byte
type Order byte

const (
	DESC  Order = 1
	ASC = 2
)


const (
	Pendiente CasoEstado = iota
	EnEsperaDelFuncionario
	EnEsperaDelCliente 
	Resuelto
	NoResuelto	
)

const (
	RoleCliente Role = iota
	RoleFuncionario
	RoleClienteAdmin
	RoleFuncionarioAdmin
	RoleAdmin
)


const (
	Activo Estado = iota
	Inactivo
	Eliminado
)

const (
	XLSX FileType = iota
	PDF
	HTML
)



type Util interface {
	//devuelve true si el  rol es de cliente  cliente = 0
	IsCliente(rol int) bool
	//devuelve true si el rol es de funcionario funcionario = 1
	IsFuncionario(rol int) bool
	//devuelve true si el  rol es de cliente admin  cliente_admin = 2
	IsClienteAdmin(rol int) bool
	//devuelve true si el  rol es de funcionario admin  funcionario_admin = 3
	IsFuncionarioAdmin(rol int) bool
	//devuelve true si el rol es cliente o cliente admin
	IsClienteRol(rol int) bool
	//devuelve true si el rol es funcionario o funcionario admin
	IsFuncionarioRol(rol int) bool
	//devuelve true si es cliente admin o funcionario admin
	IsRolAdmin(rol int) bool
	//devuelve true si es cliente o funcionario
	IsUserRol(rol int) bool
	//devuele page count and offset
	PaginationValues(p int,of int)(page int,offset int)
	//is admin rol 4
	IsAdminFuncionario(rol int) bool
}

