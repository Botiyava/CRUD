package service

import (
	"errors"
	"testing"

	"github.com/Botiyava/todo-app"
	mock_repository "github.com/Botiyava/todo-app/pkg/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthService_CreateUser(t *testing.T) {
	type args struct {
		user todo.User
	}
	type mockBehavior func(s *mock_repository.MockAuthorization, user todo.User)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
		expectedID   int
	}{
		{
			name: "Ok",
			args: args{
				user: todo.User{
					Name:     "Test",
					Username: "test",
					Password: "password",
				},
			},
			mockBehavior: func(s *mock_repository.MockAuthorization, user todo.User) {
				user.Password = generatePasswordHash(user.Password)
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			wantErr:    false,
			expectedID: 1,
		},
		{
			name: "Error",
			args: args{
				user: todo.User{
					Name:     "Test",
					Username: "test",
					Password: "password",
				},
			},
			mockBehavior: func(s *mock_repository.MockAuthorization, user todo.User) {
				user.Password = generatePasswordHash(user.Password)
				s.EXPECT().CreateUser(user).Return(0, errors.New("repository error"))
			},
			wantErr:    true,
			expectedID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_repository.NewMockAuthorization(c)
			tt.mockBehavior(repo, tt.args.user)

			s := NewAuthService(repo)

			id, err := s.CreateUser(tt.args.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}

func TestGeneratePasswordHash(t *testing.T) {
	password := "password"
	hash1 := generatePasswordHash(password)
	hash2 := generatePasswordHash(password)

	assert.Equal(t, hash1, hash2)
	assert.NotEqual(t, password, hash1)
}
