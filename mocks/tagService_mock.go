// Code generated by MockGen. DO NOT EDIT.
// Source: ./tagService_interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	tags "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
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

// GetTagsFromIDs mocks base method.
func (m *MockTagService) GetTagsFromIDs(arg0 []string) ([]*tags.ScientificFieldTag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTagsFromIDs", arg0)
	ret0, _ := ret[0].([]*tags.ScientificFieldTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTagsFromIDs indicates an expected call of GetTagsFromIDs.
func (mr *MockTagServiceMockRecorder) GetTagsFromIDs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTagsFromIDs", reflect.TypeOf((*MockTagService)(nil).GetTagsFromIDs), arg0)
}
