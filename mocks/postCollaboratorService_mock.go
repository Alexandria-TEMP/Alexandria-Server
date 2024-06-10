// Code generated by MockGen. DO NOT EDIT.
// Source: ./postCollaboratorService_interface.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=./postCollaboratorService_interface.go -destination=../../mocks/postCollaboratorService_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	gomock "go.uber.org/mock/gomock"
)

// MockPostCollaboratorService is a mock of PostCollaboratorService interface.
type MockPostCollaboratorService struct {
	ctrl     *gomock.Controller
	recorder *MockPostCollaboratorServiceMockRecorder
}

// MockPostCollaboratorServiceMockRecorder is the mock recorder for MockPostCollaboratorService.
type MockPostCollaboratorServiceMockRecorder struct {
	mock *MockPostCollaboratorService
}

// NewMockPostCollaboratorService creates a new mock instance.
func NewMockPostCollaboratorService(ctrl *gomock.Controller) *MockPostCollaboratorService {
	mock := &MockPostCollaboratorService{ctrl: ctrl}
	mock.recorder = &MockPostCollaboratorServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostCollaboratorService) EXPECT() *MockPostCollaboratorServiceMockRecorder {
	return m.recorder
}

// GetPostCollaborator mocks base method.
func (m *MockPostCollaboratorService) GetPostCollaborator(id uint) (*models.PostCollaborator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostCollaborator", id)
	ret0, _ := ret[0].(*models.PostCollaborator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostCollaborator indicates an expected call of GetPostCollaborator.
func (mr *MockPostCollaboratorServiceMockRecorder) GetPostCollaborator(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostCollaborator", reflect.TypeOf((*MockPostCollaboratorService)(nil).GetPostCollaborator), id)
}

// MembersToPostCollaborators mocks base method.
func (m *MockPostCollaboratorService) MembersToPostCollaborators(IDs []uint, anonymous bool, collaborationType models.CollaborationType) ([]*models.PostCollaborator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MembersToPostCollaborators", IDs, anonymous, collaborationType)
	ret0, _ := ret[0].([]*models.PostCollaborator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MembersToPostCollaborators indicates an expected call of MembersToPostCollaborators.
func (mr *MockPostCollaboratorServiceMockRecorder) MembersToPostCollaborators(IDs, anonymous, collaborationType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MembersToPostCollaborators", reflect.TypeOf((*MockPostCollaboratorService)(nil).MembersToPostCollaborators), IDs, anonymous, collaborationType)
}
