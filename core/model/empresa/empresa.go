package empresa

import (
	"context"
	"time"
)

type UserAddToArea struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type AddUserRequestData struct {
	AreaId int `json:"area_id"`
	Users []UserAddToArea `json:"users"`
	AreaName string `json:"area_name"`
}

type Empresa struct {
	Id          int        `json:"id"`
	Slug        string     `json:"slug"`
	Nombre      string     `json:"nombre"`
	// Descripcion *string    `json:"descripcion"`
	Telefono    *string    `json:"telefono"`
	Estado      int        `json:"estado"`
	CreatedOn   time.Time  `json:"created_on"`
	UpdatedOn   *time.Time `json:"updated_on"`
}


type EmpresaRepository interface {
	GetEmpresa(ctx context.Context,userId string,rol int) (res Empresa, err error)
	GetAreaByName(context.Context,string)(Area,error)
	GetAreasEmpresa(context.Context,int)([]Area,error)
	StoreEmpresa(ctx context.Context, empresa *Empresa) (err error)
	GetAreasUser(ctx context.Context,userId string)(res []AreaUser,err error)
	GetAreasUserAdmin(ctx context.Context,userId *string,rol *int)(res []Area,err error)
	// GetAreaUser(ctx context.Context,userId string)(res Area,err error)
	StoreArea(ctx context.Context,area *Area,rol *int)(id int,err error)
	GetEmpresas(ctx context.Context)(res []Empresa,err error)
	AddUserToArea(ctx context.Context,id *string,n *string,a *AddUserRequestData)(err error)
	AreaChangeState(ctx context.Context,state int,id int) (error)
}



type EmpresaUseCase interface {
	GetEmpresa(ctx context.Context,userId string,rol int) (res Empresa, err error)
	GetAreaByName(context.Context,string)(Area,error)
	StoreEmpresa(ctx context.Context, empresa *Empresa) (err error)
	GetAreasUser(ctx context.Context,userId string)(res []AreaUser,err error)
	GetAreasUserAdmin(ctx context.Context,userId *string,rol *int)(res []Area,err error)
	// GetAreaUser(ctx context.Context,userId string)(res Area,err error)
	StoreArea(ctx context.Context,area *Area,rol *int)(id int,err error)
	GetEmpresas(ctx context.Context)(res []Empresa,err error)
	AddUserToArea(ctx context.Context,areaD *AddUserRequestData)(err error)
	AreaChangeState(ctx context.Context,state int,id int) (error)
}