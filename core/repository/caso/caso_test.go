package caso

import (
	"context"
	"testing"
	"time"

	"soporte-go/core/model/caso"
	// "soporte-go/core/repository"

	// "soporte-go/core/repository/caso"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetCasosCliente(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	funcionario_id:= "b4312b12b11"
	rows := sqlmock.NewRows([]string{"id", "titulo", "cliente_name", "funcionario_name", "updated_on", "created_on",
		"estado", "prioridad", "client_id", "funcionario_id"}).AddRow("12121213","Nuevo caso","alejandro","jorge",time.Now(),
time.Now(),0,1,"13121231",funcionario_id).
		AddRow("121212132","Nuevo caso2","alejandro","jorge",time.Now(),
		time.Now(),0,1,"131212311",funcionario_id)
	query := "select id,titulo,cliente_name,funcionario_name,created_on,updated_on,prioridad,estado,client_id,funcionario_id from casos where client_id = \\$1 limit \\$2 offset \\$3"
	mock.ExpectQuery(query).WillReturnRows(rows)
	r := NewPgCasoRepository(db,context.TODO())
	cursor := caso.CasoQuery{
		Page: 1,
	}
	list,err := r.GetCasosCliente(context.TODO(),funcionario_id,&cursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}


func TestGetCasosFuncinario(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	funcionario_id:= "b4312b12b11"
	rows := sqlmock.NewRows([]string{"id", "titulo", "cliente_name", "funcionario_name", "updated_on", "created_on",
		"estado", "prioridad", "client_id", "funcionario_id"}).AddRow("12121213","Nuevo caso","alejandro","jorge",time.Now(),
    time.Now(),0,1,"13121231",funcionario_id).
		AddRow("121212132","Nuevo caso2","alejandro","jorge",time.Now(),
		time.Now(),0,1,"131212311",funcionario_id)
	query := "select id,titulo,cliente_name,funcionario_name,created_on,updated_on,prioridad,estado,client_id,funcionario_id from casos where funcionario_id = \\$1 limit \\$2 offset \\$3"
	mock.ExpectQuery(query).WillReturnRows(rows)
	r := NewPgCasoRepository(db,context.TODO())
	cursor := caso.CasoQuery{
		Page: 1,
	}
	list,err := r.GetCasosFuncionario(context.TODO(),funcionario_id,&cursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}