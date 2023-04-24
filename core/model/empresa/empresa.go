package empresa

import (
	"context"
	"time"
)

type UserAddToArea struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AddUserRequestData struct {
	AreaId   int             `json:"area_id"`
	Users    []UserAddToArea `json:"users"`
	AreaName string          `json:"area_name"`
}

type Empresa struct {
	Id        int        `json:"id"`
	Slug      string     `json:"slug"`
	Nombre    string     `json:"nombre"`
	ParentId  *int       `json:"parent_id"`
	Telefono  *string    `json:"telefono"`
	Estado    int        `json:"estado"`
	CreatedOn time.Time  `json:"created_on"`
	UpdatedOn *time.Time `json:"updated_on"`
}

type EmpresaRepository interface {
	GetEmpresa(ctx context.Context, userId string, rol int) (res Empresa, err error)

	CreateProyecto(ctx context.Context, a *Proyecto) (err error)
	GetProyectosEmpresa(ctx context.Context, emId int) (res []Area, err error)

	GetProyectos(ctx context.Context, parentId int) (res []Area, err error)

	GetAreaByName(context.Context, string) (Area, error)
	GetProyectoByName(context.Context, string) (ProyectoDetail, error)

	GetAreasEmpresa(context.Context, int) ([]Area, error)
	GetAreasFuncionario(context.Context, string) ([]Area, error)


	StoreEmpresa(ctx context.Context, empresa *Empresa) (err error)
	GetAreasUser(ctx context.Context, userId string) (res []AreaUser, err error)

	GetProyectoFromUserArea(ctx context.Context, id string) (res []Area, err error)
	GetProyectosFuncionario(ctx context.Context, id string) (res []Area, err error)
	GetProyectosAdmin(ctx context.Context, id int) (res []Area, err error)


	// GetAreaUser(ctx context.Context,userId string)(res Area,err error)
	GetEmpresas(ctx context.Context, emId *int) (res []Empresa, err error)
	AddUserToArea(ctx context.Context, id string, a *AddUserRequestData) (err error)
	AreaChangeState(ctx context.Context, state int, id int) error

	StoreArea(ctx context.Context, area *Area) error
	// AddAreaProveedorArea(ctx context.Context,a *ProovedorArea)()

	GetFuncionariosByAreaId(ctx context.Context, areaId int) (res []UserArea, err error)
	GetClientesByAreaId(ctx context.Context, areaId int) (res []UserArea, err error)
}

type EmpresaUseCase interface {
	CreateProyecto(ctx context.Context, a *Proyecto) (err error)
	GetProyectos(ctx context.Context, parentId int) (res []Area, err error)

	GetAreas(ctx context.Context,emId int,id string,rol int) ([]Area, error)

	GetEmpresa(ctx context.Context, userId string, rol int) (res Empresa, err error)
	GetAreaByName(context.Context, string) (Area, error)
	GetProyectoByName(context.Context, string) (ProyectoDetail, error)

	StoreEmpresa(ctx context.Context, empresa *Empresa) (err error)

	GetAreasUser(ctx context.Context, userId string) (res []AreaUser, err error)

	GetAreasFromUser(ctx context.Context, userId string, emID, rol int) (res []Area, err error)
	// GetAreaUser(ctx context.Context,userId string)(res Area,err error)
	StoreArea(ctx context.Context, area *Area) error
	GetEmpresas(ctx context.Context, emId *int) (res []Empresa, err error)
	AddUserToArea(ctx context.Context, areaD *AddUserRequestData) (err error)
	AreaChangeState(ctx context.Context, state int, id int) error

	// GetUsersAreaByAreaId(ctx context.Context,areaId int)(res []UserArea,err error)

	GetFuncionariosByAreaId(ctx context.Context, areaId int) (res []UserArea, err error)
	GetClientesByAreaId(ctx context.Context, areaId int) (res []UserArea, err error)
}
