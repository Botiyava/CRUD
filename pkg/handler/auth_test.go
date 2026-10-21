package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Botiyava/todo-app"
	"github.com/Botiyava/todo-app/pkg/service"
	mock_service "github.com/Botiyava/todo-app/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            todo.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name":"Test","username":"test","password":"password"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Invalid Input",
			inputBody: `{"username":"test"}`,
			inputUser: todo.User{},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"Key: 'User.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"name":"Test","username":"test","password":"password"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(0, errors.New("service failure"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			tt.mockBehavior(auth, tt.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(tt.inputBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, w.Body.String())
		})
	}
}
