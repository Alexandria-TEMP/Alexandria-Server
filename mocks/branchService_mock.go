// Code generated by MockGen. DO NOT EDIT.
// Source: ./branchService_interface.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=./branchService_interface.go -destination=../../mocks/branchService_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	multipart "mime/multipart"
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	flock "github.com/gofrs/flock"
	forms "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	models "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	gomock "go.uber.org/mock/gomock"
)

// MockBranchService is a mock of BranchService interface.
type MockBranchService struct {
	ctrl     *gomock.Controller
	recorder *MockBranchServiceMockRecorder
}

// MockBranchServiceMockRecorder is the mock recorder for MockBranchService.
type MockBranchServiceMockRecorder struct {
	mock *MockBranchService
}

// NewMockBranchService creates a new mock instance.
func NewMockBranchService(ctrl *gomock.Controller) *MockBranchService {
	mock := &MockBranchService{ctrl: ctrl}
	mock.recorder = &MockBranchServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBranchService) EXPECT() *MockBranchServiceMockRecorder {
	return m.recorder
}

// CreateBranch mocks base method.
func (m *MockBranchService) CreateBranch(branchCreationForm *forms.BranchCreationForm) (*models.Branch, error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBranch", branchCreationForm)
	ret0, _ := ret[0].(*models.Branch)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateBranch indicates an expected call of CreateBranch.
func (mr *MockBranchServiceMockRecorder) CreateBranch(branchCreationForm any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBranch", reflect.TypeOf((*MockBranchService)(nil).CreateBranch), branchCreationForm)
}

// CreateReview mocks base method.
func (m *MockBranchService) CreateReview(reviewCreationForm forms.ReviewCreationForm) (*models.BranchReview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReview", reviewCreationForm)
	ret0, _ := ret[0].(*models.BranchReview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReview indicates an expected call of CreateReview.
func (mr *MockBranchServiceMockRecorder) CreateReview(reviewCreationForm any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReview", reflect.TypeOf((*MockBranchService)(nil).CreateReview), reviewCreationForm)
}

// DeleteBranch mocks base method.
func (m *MockBranchService) DeleteBranch(branchID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBranch", branchID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBranch indicates an expected call of DeleteBranch.
func (mr *MockBranchServiceMockRecorder) DeleteBranch(branchID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBranch", reflect.TypeOf((*MockBranchService)(nil).DeleteBranch), branchID)
}

// GetAllBranchReviewStatuses mocks base method.
func (m *MockBranchService) GetAllBranchReviewStatuses(branchID uint) ([]models.BranchReviewDecision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllBranchReviewStatuses", branchID)
	ret0, _ := ret[0].([]models.BranchReviewDecision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllBranchReviewStatuses indicates an expected call of GetAllBranchReviewStatuses.
func (mr *MockBranchServiceMockRecorder) GetAllBranchReviewStatuses(branchID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBranchReviewStatuses", reflect.TypeOf((*MockBranchService)(nil).GetAllBranchReviewStatuses), branchID)
}

// GetBranch mocks base method.
func (m *MockBranchService) GetBranch(branchID uint) (*models.Branch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBranch", branchID)
	ret0, _ := ret[0].(*models.Branch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBranch indicates an expected call of GetBranch.
func (mr *MockBranchServiceMockRecorder) GetBranch(branchID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBranch", reflect.TypeOf((*MockBranchService)(nil).GetBranch), branchID)
}

// GetBranchProjectPost mocks base method.
func (m *MockBranchService) GetBranchProjectPost(branch *models.Branch) (*models.ProjectPost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBranchProjectPost", branch)
	ret0, _ := ret[0].(*models.ProjectPost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBranchProjectPost indicates an expected call of GetBranchProjectPost.
func (mr *MockBranchServiceMockRecorder) GetBranchProjectPost(branch any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBranchProjectPost", reflect.TypeOf((*MockBranchService)(nil).GetBranchProjectPost), branch)
}

// GetClosedBranch mocks base method.
func (m *MockBranchService) GetClosedBranch(closedBranchID uint) (*models.ClosedBranch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClosedBranch", closedBranchID)
	ret0, _ := ret[0].(*models.ClosedBranch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClosedBranch indicates an expected call of GetClosedBranch.
func (mr *MockBranchServiceMockRecorder) GetClosedBranch(closedBranchID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClosedBranch", reflect.TypeOf((*MockBranchService)(nil).GetClosedBranch), closedBranchID)
}

// GetFileFromProject mocks base method.
func (m *MockBranchService) GetFileFromProject(branchID uint, relFilepath string) (string, *flock.Flock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileFromProject", branchID, relFilepath)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*flock.Flock)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFileFromProject indicates an expected call of GetFileFromProject.
func (mr *MockBranchServiceMockRecorder) GetFileFromProject(branchID, relFilepath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileFromProject", reflect.TypeOf((*MockBranchService)(nil).GetFileFromProject), branchID, relFilepath)
}

// GetFiletree mocks base method.
func (m *MockBranchService) GetFiletree(branchID uint) (map[string]int64, error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFiletree", branchID)
	ret0, _ := ret[0].(map[string]int64)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFiletree indicates an expected call of GetFiletree.
func (mr *MockBranchServiceMockRecorder) GetFiletree(branchID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFiletree", reflect.TypeOf((*MockBranchService)(nil).GetFiletree), branchID)
}

// GetProject mocks base method.
func (m *MockBranchService) GetProject(branchID uint) (string, *flock.Flock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProject", branchID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*flock.Flock)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetProject indicates an expected call of GetProject.
func (mr *MockBranchServiceMockRecorder) GetProject(branchID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProject", reflect.TypeOf((*MockBranchService)(nil).GetProject), branchID)
}

// GetReview mocks base method.
func (m *MockBranchService) GetReview(reviewID uint) (*models.BranchReview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReview", reviewID)
	ret0, _ := ret[0].(*models.BranchReview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReview indicates an expected call of GetReview.
func (mr *MockBranchServiceMockRecorder) GetReview(reviewID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReview", reflect.TypeOf((*MockBranchService)(nil).GetReview), reviewID)
}

// MemberCanReview mocks base method.
func (m *MockBranchService) MemberCanReview(branchID, memberID uint) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MemberCanReview", branchID, memberID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MemberCanReview indicates an expected call of MemberCanReview.
func (mr *MockBranchServiceMockRecorder) MemberCanReview(branchID, memberID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MemberCanReview", reflect.TypeOf((*MockBranchService)(nil).MemberCanReview), branchID, memberID)
}

// UploadProject mocks base method.
func (m *MockBranchService) UploadProject(c *gin.Context, file *multipart.FileHeader, branchID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadProject", c, file, branchID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadProject indicates an expected call of UploadProject.
func (mr *MockBranchServiceMockRecorder) UploadProject(c, file, branchID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadProject", reflect.TypeOf((*MockBranchService)(nil).UploadProject), c, file, branchID)
}
