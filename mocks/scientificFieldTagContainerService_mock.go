// Code generated by MockGen. DO NOT EDIT.
// Source: ./scientificFieldTagContainerService_interface.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=./scientificFieldTagContainerService_interface.go -destination=../../mocks/scientificFieldTagContainerService_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	gomock "go.uber.org/mock/gomock"
)

// MockScientificFieldTagContainerService is a mock of ScientificFieldTagContainerService interface.
type MockScientificFieldTagContainerService struct {
	ctrl     *gomock.Controller
	recorder *MockScientificFieldTagContainerServiceMockRecorder
}

// MockScientificFieldTagContainerServiceMockRecorder is the mock recorder for MockScientificFieldTagContainerService.
type MockScientificFieldTagContainerServiceMockRecorder struct {
	mock *MockScientificFieldTagContainerService
}

// NewMockScientificFieldTagContainerService creates a new mock instance.
func NewMockScientificFieldTagContainerService(ctrl *gomock.Controller) *MockScientificFieldTagContainerService {
	mock := &MockScientificFieldTagContainerService{ctrl: ctrl}
	mock.recorder = &MockScientificFieldTagContainerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScientificFieldTagContainerService) EXPECT() *MockScientificFieldTagContainerServiceMockRecorder {
	return m.recorder
}

// GetScientificFieldTagContainer mocks base method.
func (m *MockScientificFieldTagContainerService) GetScientificFieldTagContainer(containerID uint) (*models.ScientificFieldTagContainer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScientificFieldTagContainer", containerID)
	ret0, _ := ret[0].(*models.ScientificFieldTagContainer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScientificFieldTagContainer indicates an expected call of GetScientificFieldTagContainer.
func (mr *MockScientificFieldTagContainerServiceMockRecorder) GetScientificFieldTagContainer(containerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScientificFieldTagContainer", reflect.TypeOf((*MockScientificFieldTagContainerService)(nil).GetScientificFieldTagContainer), containerID)
}
