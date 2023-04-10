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

	CreateSubArea(ctx context.Context, a *SubArea) (err error)
	GetSubAreasEmpresa(ctx context.Context, emId int) (res []Area, err error)

	GetSubAreas(ctx context.Context, parentId int) (res []Area, err error)

	GetAreaByName(context.Context, string) (Area, error)

	GetAreasEmpresa(context.Context, int) ([]Area, error)

	StoreEmpresa(ctx context.Context, empresa *Empresa) (err error)
	GetAreasUser(ctx context.Context, userId string) (res []AreaUser, err error)

	GetAreasClienteAdmin(ctx context.Context, id string) (res []Area, err error)
	GetSubAreasFuncionario(ctx context.Context, emId int) (res []Area, err error)

	// GetAreaUser(ctx context.Context,userId string)(res Area,err error)
	GetEmpresas(ctx context.Context, emId *int) (res []Empresa, err error)
	AddUserToArea(ctx context.Context, id string, a *AddUserRequestData) (err error)
	AreaChangeState(ctx context.Context, state int, id int) error

	StoreArea(ctx context.Context, area *Area) error
	// AddAreaProveedorArea(ctx context.Context,a *ProovedorArea)()

	GetUsersAreaByAreaId(ctx context.Context,areaId int)(res []UserArea,err error)
}

type EmpresaUseCase interface {
	CreateSubArea(ctx context.Context, a *SubArea) (err error)
	GetSubAreas(ctx context.Context, parentId int) (res []Area, err error)

	GetAreasEmpresa(context.Context, int) ([]Area, error)

	GetEmpresa(ctx context.Context, userId string, rol int) (res Empresa, err error)
	GetAreaByName(context.Context, string) (Area, error)
	StoreEmpresa(ctx context.Context, empresa *Empresa) (err error)

	GetAreasUser(ctx context.Context, userId string) (res []AreaUser, err error)

	GetAreasFromUser(ctx context.Context, userId string, emID, rol int) (res []Area, err error)
	// GetAreaUser(ctx context.Context,userId string)(res Area,err error)
	StoreArea(ctx context.Context, area *Area) error
	GetEmpresas(ctx context.Context, emId *int) (res []Empresa, err error)
	AddUserToArea(ctx context.Context, areaD *AddUserRequestData) (err error)
	AreaChangeState(ctx context.Context, state int, id int) error

	GetUsersAreaByAreaId(ctx context.Context,areaId int)(res []UserArea,err error)
}
