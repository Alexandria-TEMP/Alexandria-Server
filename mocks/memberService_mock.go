// Code generated by MockGen. DO NOT EDIT.
// Source: ./memberService_interface.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=./memberService_interface.go -destination=../../mocks/memberService_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	forms "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	models "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	gomock "go.uber.org/mock/gomock"
)

// MockMemberService is a mock of MemberService interface.
type MockMemberService struct {
	ctrl     *gomock.Controller
	recorder *MockMemberServiceMockRecorder
}

// MockMemberServiceMockRecorder is the mock recorder for MockMemberService.
type MockMemberServiceMockRecorder struct {
	mock *MockMemberService
}

// NewMockMemberService creates a new mock instance.
func NewMockMemberService(ctrl *gomock.Controller) *MockMemberService {
	mock := &MockMemberService{ctrl: ctrl}
	mock.recorder = &MockMemberServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMemberService) EXPECT() *MockMemberServiceMockRecorder {
	return m.recorder
}

// CreateMember mocks base method.
func (m *MockMemberService) CreateMember(memberForm *forms.MemberCreationForm) *models.Member {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMember", memberForm)
	ret0, _ := ret[0].(*models.Member)
	return ret0
}

// CreateMember indicates an expected call of CreateMember.
func (mr *MockMemberServiceMockRecorder) CreateMember(memberForm any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMember", reflect.TypeOf((*MockMemberService)(nil).CreateMember), memberForm)
}

// GetCollaborator mocks base method.
func (m *MockMemberService) GetCollaborator(collaboratorID uint) (*models.PostCollaborator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCollaborator", collaboratorID)
	ret0, _ := ret[0].(*models.PostCollaborator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCollaborator indicates an expected call of GetCollaborator.
func (mr *MockMemberServiceMockRecorder) GetCollaborator(collaboratorID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCollaborator", reflect.TypeOf((*MockMemberService)(nil).GetCollaborator), collaboratorID)
}

// GetMember mocks base method.
func (m *MockMemberService) GetMember(userID uint) (*models.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMember", userID)
	ret0, _ := ret[0].(*models.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMember indicates an expected call of GetMember.
func (mr *MockMemberServiceMockRecorder) GetMember(userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMember", reflect.TypeOf((*MockMemberService)(nil).GetMember), userID)
}

// UpdateMember mocks base method.
func (m *MockMemberService) UpdateMember(updatedMember *models.Member) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMember", updatedMember)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMember indicates an expected call of UpdateMember.
func (mr *MockMemberServiceMockRecorder) UpdateMember(updatedMember any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMember", reflect.TypeOf((*MockMemberService)(nil).UpdateMember), updatedMember)
}
