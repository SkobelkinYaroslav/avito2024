package module

import (
	"avito2024/internal/domain"
	mock_service "avito2024/internal/mocks"
	"avito2024/internal/service"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGetHouseFlatsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_service.NewMockRepository(ctrl) // Создаем мок репозитория
	s := service.NewService(mockRepo)                // Инициализируем сервис

	tests := []struct {
		name          string
		id            int
		user          domain.User
		mockBehavior  func(m *mock_service.MockRepository)
		expectedFlats []domain.Flat
		expectedError error
	}{
		{
			name: "Invalid house ID",
			id:   -1,
			user: domain.User{},
			mockBehavior: func(m *mock_service.MockRepository) {
			},
			expectedFlats: nil,
			expectedError: domain.NewCustomError(domain.InvalidInputError()),
		},
		{
			name: "Repository error",
			id:   1,
			user: domain.User{},
			mockBehavior: func(m *mock_service.MockRepository) {
				m.EXPECT().GetHouseFlatsRepo(1).Return(nil, domain.NewCustomError(domain.InternalError(fmt.Errorf("internal server error"))))
			},
			expectedFlats: nil,
			expectedError: domain.NewCustomError(domain.InternalError(fmt.Errorf("internal server error"))),
		},
		{
			name: "Non-client user gets all flats",
			id:   1,
			user: domain.User{UserType: "admin"},
			mockBehavior: func(m *mock_service.MockRepository) {
				m.EXPECT().GetHouseFlatsRepo(1).Return([]domain.Flat{
					{ID: 1, Status: "approved"},
					{ID: 2, Status: "pending"},
				}, nil)
			},
			expectedFlats: []domain.Flat{
				{ID: 1, Status: "approved"},
				{ID: 2, Status: "pending"},
			},
			expectedError: nil,
		},
		{
			name: "Client user gets only approved flats",
			id:   1,
			user: domain.User{UserType: "client"},
			mockBehavior: func(m *mock_service.MockRepository) {
				m.EXPECT().GetHouseFlatsRepo(1).Return([]domain.Flat{
					{ID: 1, Status: "approved"},
					{ID: 2, Status: "pending"},
				}, nil)
			},
			expectedFlats: []domain.Flat{
				{ID: 1, Status: "approved"},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRepo)

			flats, err := s.GetHouseFlatsService(tt.id, tt.user)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedFlats, flats)
		})
	}
}
