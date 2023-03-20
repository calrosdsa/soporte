package caso

import (
	"context"
	"time"
)

type CasoQuery struct {
	Page      int    `json:"page"`
	Estado    string `json:"estado"`
	Prioridad string `json:"prioridad"`
}

type CasosResponse struct {
	Casos   []Caso `json:"results"`
	Size    int    `json:"page_size"`
	Current int    `json:"current_page"`
}

type Caso struct {
	Id                     string     `json:"id"`
	ClienteId              *string    `json:"client_id"`
	FuncionarioId          *string    `json:"funcionario_id"`
	SuperiorId             *string    `json:"superior_id"`
	Titulo                 *string    `json:"titulo"`
	Descripcion            *string    `json:"descripcion"`
	DetallesDeFinalizacion *string    `json:"detalles_de_finalizacion,omitempty"`
	Empresa                *int       `json:"empresa"`
	Area                   *int       `json:"area"`
	ClienteName            *string    `json:"cliente_name"`
	FuncionarioName        *string    `json:"funcionario_name,omitempty"`
	CreatedOn              time.Time  `json:"created_on,omitempty"`
	UpdatedOn              *time.Time `json:"updated_on,omitempty"`
	FechaInicio            *time.Time `json:"fecha_inicio,omitempty"`
	FechaFin               *time.Time `json:"fecha_fin,omitempty"`
	Prioridad              *int       `json:"prioridad"`
	Estado                 *int       `json:"estado"`
}

type CasoRepository interface {
	GetCaso(ctx context.Context, id string) (Caso, error)
	GetCasosFuncionario(ctx context.Context, id *string, query *CasoQuery) ([]Caso, int, error)
	GetCasosCliente(ctx context.Context, id *string, query *CasoQuery) ([]Caso, int, error)
	GetAllCasosUser(ctx context.Context, id string, query *CasoQuery) ([]Caso, int, error)
	StoreCaso(ctx context.Context, cas *Caso, id string, emI int) (idCaso string, err error)
	UpdateCaso(ctx context.Context, columns []string, values ...interface{}) error
	AsignarFuncionario(ctx context.Context,id string,idF string) error
	// UploadRecurso(ctx context.Context)
}

type CasoUseCase interface {
	GetCaso(ctx context.Context, id string) (res Caso, err error)
	GetCasosUser(ctx context.Context, id *string, query *CasoQuery,rol *int) (casos []Caso, size int, err error)
	GetAllCasosUser(ctx context.Context, id string, query *CasoQuery) ([]Caso, int, error)
	StoreCaso(ctx context.Context, caso *Caso, id string, emI int) (idCaso string, err error)
	UpdateCaso(ctx context.Context, columns []string, values ...interface{}) error
	AsignarFuncionario(ctx context.Context,id string,idF string) error

}
