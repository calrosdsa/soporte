package caso

import "context"

type DataCaso struct {
	Name  string `json:"name"`
	Value int `json:"value"`
}

type DataCasoEstado struct {
	Name  string `json:"name"`
	Pendiente int `json:"pendiente"`
	Resuelto int `json:"resuelto"`
	NoResuelto int `json:"no_resuelto"`
}


type DataCasoResponse struct {
	ProyectosCasos []DataCaso `json:"proyectos_casos"`
	CasosLastMonth []DataCaso  `json:"casos_last_month"`
	ProyectosCasosEstado []DataCaso `json:"proyecto_casos_estado"`
}

type ChartFilterData struct {
	Ids []int `json:"ids"`
	FromDate string `json:"from_date"`
	ToDate string `json:"to_date"`
	TypeDate string  `json:"type_date"`
} 

type CasoDataRepository interface {
	GetCasosCountByProyecto(ctx context.Context, d ChartFilterData) (res []DataCaso, err error)
	GetCasosCountEstadoByProyecto(ctx context.Context, d ChartFilterData) (res []DataCaso, err error)
	GetCasosCountCreatedLast30Days(ctx context.Context,d ChartFilterData) (res []DataCaso, err error)
	GetCasosEstadoByDate(ctx context.Context,d ChartFilterData) (res []DataCasoEstado,err error)

	// GetCasosCountCreatedLast30Days(ctx context.Context, pId int) (res []DataCaso, err error)

}

type CasoDataUseCase interface {
	GetCasosCountCreatedLast30Days(ctx context.Context,d ChartFilterData) (res []DataCaso, err error)

	GetCasosData(ctx context.Context,idUser string) (res []DataCaso,res2 []DataCaso,res3 []DataCaso, err error)
	GetCasosCountEstadoByProyecto(ctx context.Context,d ChartFilterData) (res []DataCaso, err error)

	GetCasosEstadoByDate(ctx context.Context,d ChartFilterData) (res []DataCasoEstado,err error)
}

	// select  date_trunc('week',d.date)::date,count(se.id)
	// FROM
	// 		(
	// 		SELECT date_trunc('day', (current_date - offs)) AS date
	// 		FROM generate_series(0, 30, 1) AS offs
	// 		) d
	// 		LEFT OUTER JOIN casos se
	// 		ON d.date = date_trunc('day', se.created_on) 
	// 		group by date_trunc('week', d.date)
	// 	    order by date_trunc('week', d.date);



	// select  date_trunc('day',d.date)::date,
	// ("Estado") as pm,
	// count(se.id) filter (where estado = 0) as pendiente,
	// count(se.id) filter (where estado = 3) as resuelto,
	// count(se.id) filter (where estado = 4) as no_resuelto
	// FROM
	// 		(
	// 		SELECT date_trunc('day', (current_date - offs)) AS date
	// 		FROM generate_series(0, ('2023-05-11'::date- '2023-04-11'::date), 1) AS offs
	// 		) d
	// 		LEFT OUTER JOIN casos se
	// 		ON d.date = date_trunc('day', se.created_on) 
	// 		group by date_trunc('day', d.date)	
	// 	    order by date_trunc('day', d.date);

	


		
