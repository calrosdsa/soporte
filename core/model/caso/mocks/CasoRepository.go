package mocks

import (
	context "context"
	domain "soporte-go/core/model/caso"

	mock "github.com/stretchr/testify/mock"
)

type CasoRepository struct {
	mock.Mock
}

func (_m *CasoRepository) GetCaso(ctx context.Context, id string) (domain.Caso, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.Caso
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Caso); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.Caso)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *CasoRepository) GetCasosUser(ctx context.Context, id *string, query *domain.CasoQuery, rol *int) (res []domain.Caso,count int,err error) {
	return
}
func (_m *CasoRepository) GetAllCasosUser(ctx context.Context, id string, query *domain.CasoQuery) (res []domain.Caso,count int,err error) {
	return
}
func (_m *CasoRepository) StoreCaso(ctx context.Context, cas *domain.Caso, id string, emI int) (idCaso string, err error) {
	return
}
func (_m *CasoRepository) UpdateCaso(ctx context.Context, columns []string, values ...interface{}) (err error) {
	return
}
func (_m *CasoRepository) AsignarFuncionario(ctx context.Context, id string, idF string) (err error) {
	return
}
