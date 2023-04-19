package caso

import (
	"bytes"
	"context"
	"soporte-go/core/model"
	"time"
)

type CasoQuery struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	Estado    string `json:"estado"`
	Prioridad string `json:"prioridad"`
	Order     string `json:"order"`
}

type CasosResponse struct {
	Casos   []Caso `json:"results"`
	Size    int    `json:"page_size"`
	Current int    `json:"current_page"`
}

type FinalizacionDetail struct {
	Id     string `json:"id"`
	Detail string `json:"detail"`
	Estado int    `json:"estado"`
}

type CasoReporteOptions struct {
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	Estados   []int   `json:"estados"`
	Areas     []int   `json:"areas"`
	//    Empresa int `json:"empresa"`
}

type Caso struct {
	Id                     string     `json:"id"`
	ClienteId              *string    `json:"client_id"`
	FuncionarioId          *string    `json:"funcionario_id"`
	SuperiorId             *string    `json:"superior_id"`
	Titulo                 string     `json:"titulo"`
	Descripcion            *string    `json:"descripcion"`
	DetallesDeFinalizacion *string    `json:"detalles_de_finalizacion,omitempty"`
	Empresa                *int       `json:"empresa"`
	Area                   *int       `json:"area"`
	CreatedOn              time.Time  `json:"created_on,omitempty"`
	UpdatedOn              *time.Time `json:"updated_on,omitempty"`
	FechaInicio            *time.Time `json:"fecha_inicio,omitempty"`
	FechaFin               *time.Time `json:"fecha_fin,omitempty"`
	Prioridad              *int       `json:"prioridad"`
	Estado                 *int       `json:"estado"`
	Status                 *int       `json:"status"`
	ClienteName            *string    `json:"cliente_name"`
	FuncionarioName        *string    `json:"funcionario_name,omitempty"`
	ClienteApellido        *string    `json:"cliente_apellido"`
	FuncionarioApellido    *string    `json:"funcionario_apellido,omitempty"`
	Rol                    *int       `json:"rol"`
	ProyectoName           *string    `json:"proyecto_namE"`
}

type CasoRepository interface {
	GetCasoCliente(ctx context.Context, id string) (Caso, error)
	GetCasoFuncionario(ctx context.Context, id string) (Caso, error)

	GetCasosCountCliente(ctx context.Context, id string) (int, error)
	GetCasosCountFuncionario(ctx context.Context, id string) (int, error)
	GetCasosCountbySuperiorId(ctx context.Context, id string) (int, error)
	GetCasosCount(ctx context.Context) (int, error)

	GetCasosFuncionario(ctx context.Context, id string, query *CasoQuery) ([]Caso, error)
	GetCasosCliente(ctx context.Context, id string, query *CasoQuery) ([]Caso, error)

	GetAllCasosUserFuncionario(ctx context.Context, id int, query *CasoQuery) ([]Caso, error)
	GetAllCasosUserCliente(ctx context.Context, id string, query *CasoQuery) ([]Caso, error)
	UpdateCaso(ctx context.Context, c *Caso) error
	AsignarFuncionario(ctx context.Context, id string, idF string) error
	FinalizarCaso(ctx context.Context, fD *FinalizacionDetail) error
	// UploadRecurso(ctx context.Context)
	GetCasosCliForReporte(ctx context.Context, options *CasoReporteOptions) ([]Caso, error)
	GetCasosFunForReporte(ctx context.Context, options *CasoReporteOptions) ([]Caso, error)

	CreateCasoCliente(ctx context.Context, cas *Caso, id string, emI int, rol int) (err error)
	CreateCasoFuncionario(ctx context.Context, cas *Caso, id string, emI int, rol int) (err error)
}

type CasoUseCase interface {
	GetCaso(ctx context.Context, id string, rol int) (res Caso, err error)
	GetCasosUser(ctx context.Context, id string, query *CasoQuery, rol int) (casos []Caso, size int, err error)
	GetAllCasosUser(ctx context.Context, id string, query *CasoQuery, rol int) ([]Caso, int, error)
	CreateCaso(ctx context.Context, caso *Caso, id string, emI int, rol int) (err error)
	UpdateCaso(ctx context.Context, c *Caso) error
	AsignarFuncionario(ctx context.Context, id string, idF string) error
	FinalizarCaso(ctx context.Context, fD *FinalizacionDetail) error

	GetReporteCasos(ctx context.Context, t model.FileType, options *CasoReporteOptions) (b bytes.Buffer, err error)
	// CerrarCaso(ctx context.Context,id string)(error)
}
