package caso

import (
	"context"
	"database/sql"
	"log"
	"soporte-go/core/model"
	"soporte-go/core/model/caso"
	// "time"

	"github.com/lib/pq"
)

type pgCasoDataRepository struct {
	Conn *sql.DB
}

func NewCasoDataRepository(sql *sql.DB) caso.CasoDataRepository {
	return &pgCasoDataRepository{
		Conn: sql,
	}
}

func (p *pgCasoDataRepository) GetCasosCountEstadoByProyecto(ctx context.Context,d caso.ChartFilterData) (res []caso.DataCaso, err error) {
	query := `select estado, COUNT(estado) from casos where area = any($1) group by estado;`
	res, err = p.fetchCasosDataCount(ctx, query, pq.Array(d.Ids))
	return
}

func (p *pgCasoDataRepository)GetCasosCountByProyecto(ctx context.Context,d caso.ChartFilterData) (res []caso.DataCaso, err error){
	query := `select p.nombre , COUNT(c.area) as area from casos as c
	 left join proyectos as p on c.area = p.id where area = any($1) group by p.nombre;`
	res, err = p.fetchCasosDataCount(ctx, query, pq.Array(d.Ids))
	return
}

func (p *pgCasoDataRepository) GetCasosEstadoByDate(ctx context.Context,d caso.ChartFilterData) (res []caso.DataCasoEstado,err error){
	query := `select  date_trunc($4,d.date)::date,
	count(se.id) filter (where estado = 0) as pendiente,
	count(se.id) filter (where estado = 3) as resuelto,
	count(se.id) filter (where estado = 4) as no_resuelto
	FROM
			(
			SELECT date_trunc('day', ($1::date - offs)) AS date
			FROM generate_series(0, ($1::date - $2::date), 1) AS offs
			) d
			LEFT OUTER JOIN casos se ON d.date = date_trunc('day', se.created_on) 
			and area = any($3)
			group by date_trunc($4, d.date)
		    order by date_trunc($4, d.date);`
	res,err = p.fetchCasosDataEstado(ctx,query,d.ToDate,d.FromDate,pq.Array(d.Ids),d.TypeDate)
	return
}

func(p *pgCasoDataRepository)GetCasosCountCreatedLast30Days(ctx context.Context, d caso.ChartFilterData) (res []caso.DataCaso, err error){
	query := `select  date_trunc($4,d.date),count(se.id)
	FROM
		(
		SELECT date_trunc('day', ($1::date - offs)) AS date
		FROM generate_series(0, ($1::date - $2::date), 1) AS offs
		) d
		LEFT OUTER JOIN casos se
		ON d.date = date_trunc('day', se.created_on)
		and area = any($3)	
		group by date_trunc($4, d.date)
		order by date_trunc($4, d.date);`
	  res,err =  p.fetchCasosDataCount(ctx,query,d.ToDate,d.FromDate,pq.Array(d.Ids),d.TypeDate)
	  return 
}

func (p *pgCasoDataRepository) fetchCasosDataCount(ctx context.Context, query string, args ...interface{}) (result []caso.DataCaso, err error) {
	rows, err := p.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func() {
		rows.Close()
	}()
	result = make([]caso.DataCaso, 0)
	for rows.Next() {
		t := caso.DataCaso{}
		err = rows.Scan(
			&t.Name,
			&t.Value,
		)
		// log.Println(t.Name)
		result = append(result, t)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return result, nil
}

func (p *pgCasoDataRepository) fetchCasosDataEstado(ctx context.Context, query string, args ...interface{}) (result []caso.DataCasoEstado, err error) {
	rows, err := p.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func() {
		rows.Close()
	}()
	result = make([]caso.DataCasoEstado, 0)
	for rows.Next() {
		t := caso.DataCasoEstado{}
		err = rows.Scan(
			&t.Name,
			&t.Pendiente,
			&t.Resuelto,
			&t.NoResuelto,
		)
		// log.Println(t.Name)
		result = append(result, t)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return result, nil
}


func (p *pgCasoDataRepository) GetEstado(e int) string {
	switch e {
	case int(model.Pendiente):
		return "Pendientes"
	case int(model.Resuelto):
		return "Resueltos"
	case int(model.NoResuelto):
		return "No Resueltos"
	}
	return "En curso"
}
