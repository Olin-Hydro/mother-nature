// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/storage.go

// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
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

// CreateGardenReq mocks base method.
func (m *MockStorage) CreateGardenReq(gardenId string) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGardenReq", gardenId)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGardenReq indicates an expected call of CreateGardenReq.
func (mr *MockStorageMockRecorder) CreateGardenReq(gardenId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGardenReq", reflect.TypeOf((*MockStorage)(nil).CreateGardenReq), gardenId)
}

// CreateRALogsReq mocks base method.
func (m *MockStorage) CreateRALogsReq(RAId, limit string) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRALogsReq", RAId, limit)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRALogsReq indicates an expected call of CreateRALogsReq.
func (mr *MockStorageMockRecorder) CreateRALogsReq(RAId, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRALogsReq", reflect.TypeOf((*MockStorage)(nil).CreateRALogsReq), RAId, limit)
}

// CreateSensorLogsReq mocks base method.
func (m *MockStorage) CreateSensorLogsReq(SensorId, limit string) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSensorLogsReq", SensorId, limit)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSensorLogsReq indicates an expected call of CreateSensorLogsReq.
func (mr *MockStorageMockRecorder) CreateSensorLogsReq(SensorId, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSensorLogsReq", reflect.TypeOf((*MockStorage)(nil).CreateSensorLogsReq), SensorId, limit)
}
