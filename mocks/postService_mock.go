// Code generated by MockGen. DO NOT EDIT.
// Source: ./postService_interface.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=./postService_interface.go -destination=../../mocks/postService_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	multipart "mime/multipart"
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	forms "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	models "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	gomock "go.uber.org/mock/gomock"
)

// MockPostService is a mock of PostService interface.
type MockPostService struct {
	ctrl     *gomock.Controller
	recorder *MockPostServiceMockRecorder
}

// MockPostServiceMockRecorder is the mock recorder for MockPostService.
type MockPostServiceMockRecorder struct {
	mock *MockPostService
}

// NewMockPostService creates a new mock instance.
func NewMockPostService(ctrl *gomock.Controller) *MockPostService {
	mock := &MockPostService{ctrl: ctrl}
	mock.recorder = &MockPostServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostService) EXPECT() *MockPostServiceMockRecorder {
	return m.recorder
}

// CreatePost mocks base method.
func (m *MockPostService) CreatePost(form *forms.PostCreationForm) (*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost", form)
	ret0, _ := ret[0].(*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePost indicates an expected call of CreatePost.
func (mr *MockPostServiceMockRecorder) CreatePost(form any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockPostService)(nil).CreatePost), form)
}

// Filter mocks base method.
func (m *MockPostService) Filter(page, size int, form forms.FilterForm) ([]uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Filter", page, size, form)
	ret0, _ := ret[0].([]uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Filter indicates an expected call of Filter.
func (mr *MockPostServiceMockRecorder) Filter(page, size, form any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Filter", reflect.TypeOf((*MockPostService)(nil).Filter), page, size, form)
}

// GetMainFileFromProject mocks base method.
func (m *MockPostService) GetMainFileFromProject(postID uint, relFilepath string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMainFileFromProject", postID, relFilepath)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMainFileFromProject indicates an expected call of GetMainFileFromProject.
func (mr *MockPostServiceMockRecorder) GetMainFileFromProject(postID, relFilepath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMainFileFromProject", reflect.TypeOf((*MockPostService)(nil).GetMainFileFromProject), postID, relFilepath)
}

// GetMainFiletree mocks base method.
func (m *MockPostService) GetMainFiletree(branchID uint) (map[string]int64, error, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMainFiletree", branchID)
	ret0, _ := ret[0].(map[string]int64)
	ret1, _ := ret[1].(error)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMainFiletree indicates an expected call of GetMainFiletree.
func (mr *MockPostServiceMockRecorder) GetMainFiletree(branchID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMainFiletree", reflect.TypeOf((*MockPostService)(nil).GetMainFiletree), branchID)
}

// GetMainProject mocks base method.
func (m *MockPostService) GetMainProject(postID uint) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMainProject", postID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMainProject indicates an expected call of GetMainProject.
func (mr *MockPostServiceMockRecorder) GetMainProject(postID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMainProject", reflect.TypeOf((*MockPostService)(nil).GetMainProject), postID)
}

// GetPost mocks base method.
func (m *MockPostService) GetPost(postID uint) (*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPost", postID)
	ret0, _ := ret[0].(*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPost indicates an expected call of GetPost.
func (mr *MockPostServiceMockRecorder) GetPost(postID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPost", reflect.TypeOf((*MockPostService)(nil).GetPost), postID)
}

// UpdatePost mocks base method.
func (m *MockPostService) UpdatePost(updatedPost *models.Post) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePost", updatedPost)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePost indicates an expected call of UpdatePost.
func (mr *MockPostServiceMockRecorder) UpdatePost(updatedPost any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePost", reflect.TypeOf((*MockPostService)(nil).UpdatePost), updatedPost)
}

// UploadPost mocks base method.
func (m *MockPostService) UploadPost(c *gin.Context, file *multipart.FileHeader, postID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadPost", c, file, postID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadPost indicates an expected call of UploadPost.
func (mr *MockPostServiceMockRecorder) UploadPost(c, file, postID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadPost", reflect.TypeOf((*MockPostService)(nil).UploadPost), c, file, postID)
}
