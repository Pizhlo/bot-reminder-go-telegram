// Code generated by MockGen. DO NOT EDIT.
// Source: ./note.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	elastic "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockelasticClient is a mock of elasticClient interface.
type MockelasticClient struct {
	ctrl     *gomock.Controller
	recorder *MockelasticClientMockRecorder
}

// MockelasticClientMockRecorder is the mock recorder for MockelasticClient.
type MockelasticClientMockRecorder struct {
	mock *MockelasticClient
}

// NewMockelasticClient creates a new mock instance.
func NewMockelasticClient(ctrl *gomock.Controller) *MockelasticClient {
	mock := &MockelasticClient{ctrl: ctrl}
	mock.recorder = &MockelasticClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockelasticClient) EXPECT() *MockelasticClientMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockelasticClient) Delete(ctx context.Context, search elastic.Data) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, search)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockelasticClientMockRecorder) Delete(ctx, search interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockelasticClient)(nil).Delete), ctx, search)
}

// DeleteAllByUserID mocks base method.
func (m *MockelasticClient) DeleteAllByUserID(ctx context.Context, data elastic.Data) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllByUserID", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllByUserID indicates an expected call of DeleteAllByUserID.
func (mr *MockelasticClientMockRecorder) DeleteAllByUserID(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllByUserID", reflect.TypeOf((*MockelasticClient)(nil).DeleteAllByUserID), ctx, data)
}

// Save mocks base method.
func (m *MockelasticClient) Save(ctx context.Context, search elastic.Data) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, search)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockelasticClientMockRecorder) Save(ctx, search interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockelasticClient)(nil).Save), ctx, search)
}

// SearchByID mocks base method.
func (m *MockelasticClient) SearchByID(ctx context.Context, search elastic.Data) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchByID", ctx, search)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchByID indicates an expected call of SearchByID.
func (mr *MockelasticClientMockRecorder) SearchByID(ctx, search interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchByID", reflect.TypeOf((*MockelasticClient)(nil).SearchByID), ctx, search)
}

// SearchByText mocks base method.
func (m *MockelasticClient) SearchByText(ctx context.Context, search elastic.Data) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchByText", ctx, search)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchByText indicates an expected call of SearchByText.
func (mr *MockelasticClientMockRecorder) SearchByText(ctx, search interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchByText", reflect.TypeOf((*MockelasticClient)(nil).SearchByText), ctx, search)
}

// Update mocks base method.
func (m *MockelasticClient) Update(ctx context.Context, search elastic.Data) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, search)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockelasticClientMockRecorder) Update(ctx, search interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockelasticClient)(nil).Update), ctx, search)
}
