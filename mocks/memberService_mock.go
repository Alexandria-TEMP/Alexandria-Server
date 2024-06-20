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
func (m *MockMemberService) CreateMember(memberForm *forms.MemberCreationForm, userFields *models.ScientificFieldTagContainer) (*models.LoggedInMemberDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMember", memberForm, userFields)
	ret0, _ := ret[0].(*models.LoggedInMemberDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMember indicates an expected call of CreateMember.
func (mr *MockMemberServiceMockRecorder) CreateMember(memberForm, userFields any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMember", reflect.TypeOf((*MockMemberService)(nil).CreateMember), memberForm, userFields)
}

// DeleteMember mocks base method.
func (m *MockMemberService) DeleteMember(memberID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMember", memberID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMember indicates an expected call of DeleteMember.
func (mr *MockMemberServiceMockRecorder) DeleteMember(memberID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMember", reflect.TypeOf((*MockMemberService)(nil).DeleteMember), memberID)
}

// GetAllMembers mocks base method.
func (m *MockMemberService) GetAllMembers() ([]*models.MemberShortFormDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMembers")
	ret0, _ := ret[0].([]*models.MemberShortFormDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMembers indicates an expected call of GetAllMembers.
func (mr *MockMemberServiceMockRecorder) GetAllMembers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMembers", reflect.TypeOf((*MockMemberService)(nil).GetAllMembers))
}

// GetMember mocks base method.
func (m *MockMemberService) GetMember(memberID uint) (*models.Member, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMember", memberID)
	ret0, _ := ret[0].(*models.Member)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMember indicates an expected call of GetMember.
func (mr *MockMemberServiceMockRecorder) GetMember(memberID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMember", reflect.TypeOf((*MockMemberService)(nil).GetMember), memberID)
}

// LogInMember mocks base method.
func (m *MockMemberService) LogInMember(memberAuthForm *forms.MemberAuthForm) (*models.LoggedInMemberDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogInMember", memberAuthForm)
	ret0, _ := ret[0].(*models.LoggedInMemberDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogInMember indicates an expected call of LogInMember.
func (mr *MockMemberServiceMockRecorder) LogInMember(memberAuthForm any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogInMember", reflect.TypeOf((*MockMemberService)(nil).LogInMember), memberAuthForm)
}

// RefreshToken mocks base method.
func (m *MockMemberService) RefreshToken(form *forms.TokenRefreshForm) (*models.TokenPairDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", form)
	ret0, _ := ret[0].(*models.TokenPairDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockMemberServiceMockRecorder) RefreshToken(form any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockMemberService)(nil).RefreshToken), form)
}
