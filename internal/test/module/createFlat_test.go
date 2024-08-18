package module

import (
	"avito2024/internal/domain"
	mock_service "avito2024/internal/mocks"
	"avito2024/internal/service"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateFlatService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_service.NewMockRepository(ctrl)
	s := service.NewService(mockRepo)

	bdError := domain.NewCustomError(domain.InvalidInputError())

	tests := []struct {
		name          string
		flat          domain.Flat
		mockBehavior  func(m *mock_service.MockRepository)
		expectedFlat  domain.Flat
		expectedError error
	}{
		{
			name: "Invalid HouseID",
			flat: domain.Flat{
				HouseID: -1,
				Price:   300000,
				Rooms:   3,
			},
			mockBehavior:  func(m *mock_service.MockRepository) {},
			expectedFlat:  domain.Flat{},
			expectedError: domain.NewCustomError(domain.InvalidInputError()),
		},
		{
			name: "Not existing HouseID",
			flat: domain.Flat{
				HouseID: 1337,
				Price:   300000,
				Rooms:   3,
			},
			mockBehavior: func(m *mock_service.MockRepository) {
				m.EXPECT().CreateFlatRepo(gomock.Any()).Return(domain.Flat{}, bdError)
			},
			expectedFlat:  domain.Flat{},
			expectedError: bdError,
		},
		{
			name: "Invalid Price",
			flat: domain.Flat{
				HouseID: 101,
				Price:   0,
				Rooms:   3,
			},
			mockBehavior:  func(m *mock_service.MockRepository) {},
			expectedFlat:  domain.Flat{},
			expectedError: domain.NewCustomError(domain.InvalidInputError()),
		},
		{
			name: "Invalid Rooms",
			flat: domain.Flat{
				HouseID: 101,
				Price:   300000,
				Rooms:   0,
			},
			mockBehavior:  func(m *mock_service.MockRepository) {},
			expectedFlat:  domain.Flat{},
			expectedError: domain.NewCustomError(domain.InvalidInputError()),
		},
		{
			name: "Valid Flat",
			flat: domain.Flat{
				HouseID: 101,
				Price:   300000,
				Rooms:   3,
			},
			mockBehavior: func(m *mock_service.MockRepository) {
				m.EXPECT().CreateFlatRepo(gomock.Any()).Return(domain.Flat{
					ID:      1,
					HouseID: 101,
					Price:   300000,
					Rooms:   3,
					Status:  "created",
				}, nil)
			},
			expectedFlat: domain.Flat{
				ID:      1,
				HouseID: 101,
				Price:   300000,
				Rooms:   3,
				Status:  "created",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRepo)

			createdFlat, err := s.CreateFlatService(tt.flat)

			errors.Is(tt.expectedError, err)
			assert.Equal(t, tt.expectedFlat, createdFlat)
		})
	}
}
