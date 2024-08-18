package module

import (
	"avito2024/internal/domain"
	mock_service "avito2024/internal/mocks"
	"avito2024/internal/service"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGetHouseFlatsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock_service.NewMockRepository(ctrl)
	s := service.NewService(mockRepo)

	client := domain.User{UserType: "client"}
	moderator := domain.User{UserType: "moderator"}
	flats := []domain.Flat{
		{ID: 1, HouseID: 101, Price: 300000, Rooms: 3, Status: "approved"},
		{ID: 2, HouseID: 101, Price: 250000, Rooms: 2, Status: "created"},
		{ID: 3, HouseID: 101, Price: 400000, Rooms: 4, Status: "on moderation"},
		{ID: 4, HouseID: 101, Price: 150000, Rooms: 1, Status: "declined"},
		{ID: 5, HouseID: 101, Price: 500000, Rooms: 5, Status: "approved"},
		{ID: 6, HouseID: 101, Price: 350000, Rooms: 3, Status: "created"},
	}
	approvedFlats := []domain.Flat{
		{ID: 1, HouseID: 101, Price: 300000, Rooms: 3, Status: "approved"},
		{ID: 5, HouseID: 101, Price: 500000, Rooms: 5, Status: "approved"},
	}
	bdError := domain.NewCustomError(domain.InternalError(fmt.Errorf("internal server error")))

	tests := []struct {
		name          string
		id            int
		user          domain.User
		mockBehavior  func(m *mock_service.MockRepository)
		expectedFlats []domain.Flat
		expectedError error
	}{
		{
			name: "Check if output is valid",
			id:   101,
			user: client,
			mockBehavior: func(m *mock_service.MockRepository) {
				m.EXPECT().GetHouseFlatsRepo(101).Return(flats, nil)
			},
			expectedFlats: approvedFlats,
			expectedError: nil,
		},
		{
			name:          "Invalid house ID",
			id:            -1,
			user:          client,
			mockBehavior:  func(m *mock_service.MockRepository) {},
			expectedFlats: nil,
			expectedError: bdError,
		},
		{
			name: "House doesnt exist",
			id:   1337,
			user: client,
			mockBehavior: func(m *mock_service.MockRepository) {
				m.EXPECT().GetHouseFlatsRepo(1337).Return(nil, domain.NewCustomError(domain.InternalError(fmt.Errorf("internal server error"))))
			},
			expectedFlats: nil,
			expectedError: bdError,
		},
		{
			name: "Client gets only approved flats",
			id:   101,
			user: client,
			mockBehavior: func(m *mock_service.MockRepository) {
				m.EXPECT().GetHouseFlatsRepo(101).Return(flats, nil)
			},
			expectedFlats: approvedFlats,
			expectedError: nil,
		},
		{
			name: "Moderator gets all flats",
			id:   101,
			user: moderator,
			mockBehavior: func(m *mock_service.MockRepository) {
				m.EXPECT().GetHouseFlatsRepo(101).Return(flats, nil)
			},
			expectedFlats: flats,
			expectedError: nil,
		},
		{
			name: "Empty flats",
			id:   101,
			user: client,
			mockBehavior: func(m *mock_service.MockRepository) {
				m.EXPECT().GetHouseFlatsRepo(101).Return([]domain.Flat{}, nil)
			},
			expectedFlats: []domain.Flat{},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(mockRepo)

			flats, err := s.GetHouseFlatsService(tt.id, tt.user)

			errors.Is(tt.expectedError, err)
			assert.Equal(t, tt.expectedFlats, flats)
		})
	}
}
