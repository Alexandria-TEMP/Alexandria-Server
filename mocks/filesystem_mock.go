// Code generated by MockGen. DO NOT EDIT.
// Source: ./filesystem_interface.go
//
// Generated by this command:
//
//	mockgen -package=mocks -source=./filesystem_interface.go -destination=../../mocks/filesystem_mock.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	multipart "mime/multipart"
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	git "github.com/go-git/go-git/v5"
	plumbing "github.com/go-git/go-git/v5/plumbing"
	gomock "go.uber.org/mock/gomock"
)

// MockFilesystem is a mock of Filesystem interface.
type MockFilesystem struct {
	ctrl     *gomock.Controller
	recorder *MockFilesystemMockRecorder
}

// MockFilesystemMockRecorder is the mock recorder for MockFilesystem.
type MockFilesystemMockRecorder struct {
	mock *MockFilesystem
}

// NewMockFilesystem creates a new mock instance.
func NewMockFilesystem(ctrl *gomock.Controller) *MockFilesystem {
	mock := &MockFilesystem{ctrl: ctrl}
	mock.recorder = &MockFilesystemMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFilesystem) EXPECT() *MockFilesystemMockRecorder {
	return m.recorder
}

// CheckoutBranch mocks base method.
func (m *MockFilesystem) CheckoutBranch(branchName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckoutBranch", branchName)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckoutBranch indicates an expected call of CheckoutBranch.
func (mr *MockFilesystemMockRecorder) CheckoutBranch(branchName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckoutBranch", reflect.TypeOf((*MockFilesystem)(nil).CheckoutBranch), branchName)
}

// CheckoutDirectory mocks base method.
func (m *MockFilesystem) CheckoutDirectory(postID uint) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CheckoutDirectory", postID)
}

// CheckoutDirectory indicates an expected call of CheckoutDirectory.
func (mr *MockFilesystemMockRecorder) CheckoutDirectory(postID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckoutDirectory", reflect.TypeOf((*MockFilesystem)(nil).CheckoutDirectory), postID)
}

// CheckoutRepository mocks base method.
func (m *MockFilesystem) CheckoutRepository() (*git.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckoutRepository")
	ret0, _ := ret[0].(*git.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckoutRepository indicates an expected call of CheckoutRepository.
func (mr *MockFilesystemMockRecorder) CheckoutRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckoutRepository", reflect.TypeOf((*MockFilesystem)(nil).CheckoutRepository))
}

// CleanDir mocks base method.
func (m *MockFilesystem) CleanDir() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanDir")
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanDir indicates an expected call of CleanDir.
func (mr *MockFilesystemMockRecorder) CleanDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanDir", reflect.TypeOf((*MockFilesystem)(nil).CleanDir))
}

// CreateBranch mocks base method.
func (m *MockFilesystem) CreateBranch(branchName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBranch", branchName)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBranch indicates an expected call of CreateBranch.
func (mr *MockFilesystemMockRecorder) CreateBranch(branchName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBranch", reflect.TypeOf((*MockFilesystem)(nil).CreateBranch), branchName)
}

// CreateCommit mocks base method.
func (m *MockFilesystem) CreateCommit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCommit")
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCommit indicates an expected call of CreateCommit.
func (mr *MockFilesystemMockRecorder) CreateCommit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCommit", reflect.TypeOf((*MockFilesystem)(nil).CreateCommit))
}

// CreateRepository mocks base method.
func (m *MockFilesystem) CreateRepository() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRepository")
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRepository indicates an expected call of CreateRepository.
func (mr *MockFilesystemMockRecorder) CreateRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRepository", reflect.TypeOf((*MockFilesystem)(nil).CreateRepository))
}

// DeleteBranch mocks base method.
func (m *MockFilesystem) DeleteBranch(branchName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBranch", branchName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBranch indicates an expected call of DeleteBranch.
func (mr *MockFilesystemMockRecorder) DeleteBranch(branchName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBranch", reflect.TypeOf((*MockFilesystem)(nil).DeleteBranch), branchName)
}

// DeleteRepository mocks base method.
func (m *MockFilesystem) DeleteRepository() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRepository")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRepository indicates an expected call of DeleteRepository.
func (mr *MockFilesystemMockRecorder) DeleteRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRepository", reflect.TypeOf((*MockFilesystem)(nil).DeleteRepository))
}

// GetCurrentDirPath mocks base method.
func (m *MockFilesystem) GetCurrentDirPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentDirPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetCurrentDirPath indicates an expected call of GetCurrentDirPath.
func (mr *MockFilesystemMockRecorder) GetCurrentDirPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentDirPath", reflect.TypeOf((*MockFilesystem)(nil).GetCurrentDirPath))
}

// GetCurrentQuartoDirPath mocks base method.
func (m *MockFilesystem) GetCurrentQuartoDirPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentQuartoDirPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetCurrentQuartoDirPath indicates an expected call of GetCurrentQuartoDirPath.
func (mr *MockFilesystemMockRecorder) GetCurrentQuartoDirPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentQuartoDirPath", reflect.TypeOf((*MockFilesystem)(nil).GetCurrentQuartoDirPath))
}

// GetCurrentRenderDirPath mocks base method.
func (m *MockFilesystem) GetCurrentRenderDirPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentRenderDirPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetCurrentRenderDirPath indicates an expected call of GetCurrentRenderDirPath.
func (mr *MockFilesystemMockRecorder) GetCurrentRenderDirPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentRenderDirPath", reflect.TypeOf((*MockFilesystem)(nil).GetCurrentRenderDirPath))
}

// GetCurrentZipFilePath mocks base method.
func (m *MockFilesystem) GetCurrentZipFilePath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentZipFilePath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetCurrentZipFilePath indicates an expected call of GetCurrentZipFilePath.
func (mr *MockFilesystemMockRecorder) GetCurrentZipFilePath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentZipFilePath", reflect.TypeOf((*MockFilesystem)(nil).GetCurrentZipFilePath))
}

// GetFileTree mocks base method.
func (m *MockFilesystem) GetFileTree() (map[string]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileTree")
	ret0, _ := ret[0].(map[string]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileTree indicates an expected call of GetFileTree.
func (mr *MockFilesystemMockRecorder) GetFileTree() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileTree", reflect.TypeOf((*MockFilesystem)(nil).GetFileTree))
}

// GetLastCommit mocks base method.
func (m *MockFilesystem) GetLastCommit(branchName string) (*plumbing.Reference, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastCommit", branchName)
	ret0, _ := ret[0].(*plumbing.Reference)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastCommit indicates an expected call of GetLastCommit.
func (mr *MockFilesystemMockRecorder) GetLastCommit(branchName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastCommit", reflect.TypeOf((*MockFilesystem)(nil).GetLastCommit), branchName)
}

// Merge mocks base method.
func (m *MockFilesystem) Merge(toMerge, mergeInto string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Merge", toMerge, mergeInto)
	ret0, _ := ret[0].(error)
	return ret0
}

// Merge indicates an expected call of Merge.
func (mr *MockFilesystemMockRecorder) Merge(toMerge, mergeInto any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Merge", reflect.TypeOf((*MockFilesystem)(nil).Merge), toMerge, mergeInto)
}

// RenderExists mocks base method.
func (m *MockFilesystem) RenderExists() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenderExists")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RenderExists indicates an expected call of RenderExists.
func (mr *MockFilesystemMockRecorder) RenderExists() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderExists", reflect.TypeOf((*MockFilesystem)(nil).RenderExists))
}

// Reset mocks base method.
func (m *MockFilesystem) Reset() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reset")
	ret0, _ := ret[0].(error)
	return ret0
}

// Reset indicates an expected call of Reset.
func (mr *MockFilesystemMockRecorder) Reset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockFilesystem)(nil).Reset))
}

// SaveZipFile mocks base method.
func (m *MockFilesystem) SaveZipFile(c *gin.Context, file *multipart.FileHeader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveZipFile", c, file)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveZipFile indicates an expected call of SaveZipFile.
func (mr *MockFilesystemMockRecorder) SaveZipFile(c, file any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveZipFile", reflect.TypeOf((*MockFilesystem)(nil).SaveZipFile), c, file)
}

// Unzip mocks base method.
func (m *MockFilesystem) Unzip() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unzip")
	ret0, _ := ret[0].(error)
	return ret0
}

// Unzip indicates an expected call of Unzip.
func (mr *MockFilesystemMockRecorder) Unzip() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unzip", reflect.TypeOf((*MockFilesystem)(nil).Unzip))
}
