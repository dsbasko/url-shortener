// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dsbasko/yandex-go-shortener/internal/repository/storage (interfaces: Storage)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/storage.go -package=mock_storage github.com/dsbasko/yandex-go-shortener/internal/repository/storage Storage
//

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"

	entities "github.com/dsbasko/yandex-go-shortener/internal/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// CreateURL mocks base method.
func (m *MockStorage) CreateURL(arg0 context.Context, arg1 entities.URL) (entities.URL, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateURL", arg0, arg1)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateURL indicates an expected call of CreateURL.
func (mr *MockStorageMockRecorder) CreateURL(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateURL", reflect.TypeOf((*MockStorage)(nil).CreateURL), arg0, arg1)
}

// CreateURLs mocks base method.
func (m *MockStorage) CreateURLs(arg0 context.Context, arg1 []entities.URL) ([]entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateURLs", arg0, arg1)
	ret0, _ := ret[0].([]entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateURLs indicates an expected call of CreateURLs.
func (mr *MockStorageMockRecorder) CreateURLs(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateURLs", reflect.TypeOf((*MockStorage)(nil).CreateURLs), arg0, arg1)
}

// DeleteURLs mocks base method.
func (m *MockStorage) DeleteURLs(arg0 context.Context, arg1 []entities.URL) ([]entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteURLs", arg0, arg1)
	ret0, _ := ret[0].([]entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteURLs indicates an expected call of DeleteURLs.
func (mr *MockStorageMockRecorder) DeleteURLs(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteURLs", reflect.TypeOf((*MockStorage)(nil).DeleteURLs), arg0, arg1)
}

// GetURLByOriginalURL mocks base method.
func (m *MockStorage) GetURLByOriginalURL(arg0 context.Context, arg1 string) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLByOriginalURL", arg0, arg1)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLByOriginalURL indicates an expected call of GetURLByOriginalURL.
func (mr *MockStorageMockRecorder) GetURLByOriginalURL(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLByOriginalURL", reflect.TypeOf((*MockStorage)(nil).GetURLByOriginalURL), arg0, arg1)
}

// GetURLByShortURL mocks base method.
func (m *MockStorage) GetURLByShortURL(arg0 context.Context, arg1 string) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLByShortURL", arg0, arg1)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLByShortURL indicates an expected call of GetURLByShortURL.
func (mr *MockStorageMockRecorder) GetURLByShortURL(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLByShortURL", reflect.TypeOf((*MockStorage)(nil).GetURLByShortURL), arg0, arg1)
}

// GetURLsByUserID mocks base method.
func (m *MockStorage) GetURLsByUserID(arg0 context.Context, arg1 string) ([]entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURLsByUserID", arg0, arg1)
	ret0, _ := ret[0].([]entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURLsByUserID indicates an expected call of GetURLsByUserID.
func (mr *MockStorageMockRecorder) GetURLsByUserID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURLsByUserID", reflect.TypeOf((*MockStorage)(nil).GetURLsByUserID), arg0, arg1)
}

// Ping mocks base method.
func (m *MockStorage) Ping(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStorageMockRecorder) Ping(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStorage)(nil).Ping), arg0)
}
