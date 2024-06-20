// Code generated by MockGen. DO NOT EDIT.
// Source: ./tagService_interface.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=./tagService_interface.go -destination=../../mocks/tagService_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	gomock "go.uber.org/mock/gomock"
)

// MockTagService is a mock of TagService interface.
type MockTagService struct {
	ctrl     *gomock.Controller
	recorder *MockTagServiceMockRecorder
}

// MockTagServiceMockRecorder is the mock recorder for MockTagService.
type MockTagServiceMockRecorder struct {
	mock *MockTagService
}

// NewMockTagService creates a new mock instance.
func NewMockTagService(ctrl *gomock.Controller) *MockTagService {
	mock := &MockTagService{ctrl: ctrl}
	mock.recorder = &MockTagServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTagService) EXPECT() *MockTagServiceMockRecorder {
	return m.recorder
}

// GetAllScientificFieldTags mocks base method.
func (m *MockTagService) GetAllScientificFieldTags() ([]*models.ScientificFieldTag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllScientificFieldTags")
	ret0, _ := ret[0].([]*models.ScientificFieldTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllScientificFieldTags indicates an expected call of GetAllScientificFieldTags.
func (mr *MockTagServiceMockRecorder) GetAllScientificFieldTags() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllScientificFieldTags", reflect.TypeOf((*MockTagService)(nil).GetAllScientificFieldTags))
}

// GetTagByID mocks base method.
func (m *MockTagService) GetTagByID(id uint) (*models.ScientificFieldTag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTagByID", id)
	ret0, _ := ret[0].(*models.ScientificFieldTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTagByID indicates an expected call of GetTagByID.
func (mr *MockTagServiceMockRecorder) GetTagByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTagByID", reflect.TypeOf((*MockTagService)(nil).GetTagByID), id)
}

// GetTagsFromIDs mocks base method.
func (m *MockTagService) GetTagsFromIDs(arg0 []uint) ([]*models.ScientificFieldTag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTagsFromIDs", arg0)
	ret0, _ := ret[0].([]*models.ScientificFieldTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTagsFromIDs indicates an expected call of GetTagsFromIDs.
func (mr *MockTagServiceMockRecorder) GetTagsFromIDs(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTagsFromIDs", reflect.TypeOf((*MockTagService)(nil).GetTagsFromIDs), arg0)
}
