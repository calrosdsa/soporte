package ucases

import (
	"context"
	domain "soporte-go/core/model/caso"
	"soporte-go/core/model/caso/mocks"
	// ucase "soporte-go/domain/ucases"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCaso(t *testing.T) {
	mockArticleRepo := new(mocks.CasoRepository)
	mockArticle := domain.Caso{
		Id: "2939219312mfmmf",
	}
	// mockAuthor := domain.Author{
	// 	ID:   1,
	// 	Name: "Iman Tumorang",
	// }

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetCaso", mock.Anything, mock.AnythingOfType("string")).Return(mockArticle, nil).Once()
		// mockAuthorrepo := new(mocks.CasoRepository)
		// mockAuthorrepo.On("GetCaso", mock.Anything, mock.AnythingOfType("string")).Return(mockAuthor, nil)
		u := NewCasoUseCase(mockArticleRepo, time.Second*2)

		a, err := u.GetCaso(context.TODO(), mockArticle.Id)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockArticleRepo.AssertExpectations(t)
		// mockAuthorrepo.AssertExpectations(t)
	})
	// t.Run("error-failed", func(t *testing.T) {
	// 	mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Article{}, errors.New("Unexpected")).Once()

	// 	mockAuthorrepo := new(mocks.AuthorRepository)
	// 	u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

	// 	a, err := u.GetByID(context.TODO(), mockArticle.ID)

	// 	assert.Error(t, err)
	// 	assert.Equal(t, domain.Article{}, a)

	// 	mockArticleRepo.AssertExpectations(t)
	// 	mockAuthorrepo.AssertExpectations(t)
	// })
}