// Code generated by MockGen. DO NOT EDIT.
// Source: ./projectPostService_interface.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=./projectPostService_interface.go -destination=../../mocks/projectPostService_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	forms "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	models "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	gomock "go.uber.org/mock/gomock"
)

// MockProjectPostService is a mock of ProjectPostService interface.
type MockProjectPostService struct {
	ctrl     *gomock.Controller
	recorder *MockProjectPostServiceMockRecorder
}

// MockProjectPostServiceMockRecorder is the mock recorder for MockProjectPostService.
type MockProjectPostServiceMockRecorder struct {
	mock *MockProjectPostService
}

// NewMockProjectPostService creates a new mock instance.
func NewMockProjectPostService(ctrl *gomock.Controller) *MockProjectPostService {
	mock := &MockProjectPostService{ctrl: ctrl}
	mock.recorder = &MockProjectPostServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectPostService) EXPECT() *MockProjectPostServiceMockRecorder {
	return m.recorder
}

// CreateProjectPost mocks base method.
func (m *MockProjectPostService) CreateProjectPost(form *forms.ProjectPostCreationForm) (*models.ProjectPost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProjectPost", form)
	ret0, _ := ret[0].(*models.ProjectPost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProjectPost indicates an expected call of CreateProjectPost.
func (mr *MockProjectPostServiceMockRecorder) CreateProjectPost(form any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProjectPost", reflect.TypeOf((*MockProjectPostService)(nil).CreateProjectPost), form)
}

// GetProjectPost mocks base method.
func (m *MockProjectPostService) GetProjectPost(postID uint) (*models.ProjectPost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProjectPost", postID)
	ret0, _ := ret[0].(*models.ProjectPost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProjectPost indicates an expected call of GetProjectPost.
func (mr *MockProjectPostServiceMockRecorder) GetProjectPost(postID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProjectPost", reflect.TypeOf((*MockProjectPostService)(nil).GetProjectPost), postID)
}

// UpdateProjectPost mocks base method.
func (m *MockProjectPostService) UpdateProjectPost(updatedPost *models.ProjectPost) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProjectPost", updatedPost)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProjectPost indicates an expected call of UpdateProjectPost.
func (mr *MockProjectPostServiceMockRecorder) UpdateProjectPost(updatedPost any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProjectPost", reflect.TypeOf((*MockProjectPostService)(nil).UpdateProjectPost), updatedPost)
}
