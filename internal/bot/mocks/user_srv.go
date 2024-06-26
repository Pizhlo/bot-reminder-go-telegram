// Code generated by MockGen. DO NOT EDIT.
// Source: ./user.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	user "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	gomock "github.com/golang/mock/gomock"
)

// MockuserEditor is a mock of userEditor interface.
type MockuserEditor struct {
	ctrl     *gomock.Controller
	recorder *MockuserEditorMockRecorder
}

// MockuserEditorMockRecorder is the mock recorder for MockuserEditor.
type MockuserEditorMockRecorder struct {
	mock *MockuserEditor
}

// NewMockuserEditor creates a new mock instance.
func NewMockuserEditor(ctrl *gomock.Controller) *MockuserEditor {
	mock := &MockuserEditor{ctrl: ctrl}
	mock.recorder = &MockuserEditorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuserEditor) EXPECT() *MockuserEditorMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockuserEditor) Get(ctx context.Context, userID int64) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, userID)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockuserEditorMockRecorder) Get(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockuserEditor)(nil).Get), ctx, userID)
}

// GetAll mocks base method.
func (m *MockuserEditor) GetAll(ctx context.Context) ([]*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockuserEditorMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockuserEditor)(nil).GetAll), ctx)
}

// GetState mocks base method.
func (m *MockuserEditor) GetState(ctx context.Context, id int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState", ctx, id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetState indicates an expected call of GetState.
func (mr *MockuserEditorMockRecorder) GetState(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockuserEditor)(nil).GetState), ctx, id)
}

// Save mocks base method.
func (m *MockuserEditor) Save(ctx context.Context, id int64, u *user.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, id, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockuserEditorMockRecorder) Save(ctx, id, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockuserEditor)(nil).Save), ctx, id, u)
}

// SaveState mocks base method.
func (m *MockuserEditor) SaveState(ctx context.Context, id int64, state string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveState", ctx, id, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveState indicates an expected call of SaveState.
func (mr *MockuserEditorMockRecorder) SaveState(ctx, id, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveState", reflect.TypeOf((*MockuserEditor)(nil).SaveState), ctx, id, state)
}

// MocktimezoneEditor is a mock of timezoneEditor interface.
type MocktimezoneEditor struct {
	ctrl     *gomock.Controller
	recorder *MocktimezoneEditorMockRecorder
}

// MocktimezoneEditorMockRecorder is the mock recorder for MocktimezoneEditor.
type MocktimezoneEditorMockRecorder struct {
	mock *MocktimezoneEditor
}

// NewMocktimezoneEditor creates a new mock instance.
func NewMocktimezoneEditor(ctrl *gomock.Controller) *MocktimezoneEditor {
	mock := &MocktimezoneEditor{ctrl: ctrl}
	mock.recorder = &MocktimezoneEditorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocktimezoneEditor) EXPECT() *MocktimezoneEditorMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MocktimezoneEditor) Get(ctx context.Context, userID int64) (*user.Timezone, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, userID)
	ret0, _ := ret[0].(*user.Timezone)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MocktimezoneEditorMockRecorder) Get(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MocktimezoneEditor)(nil).Get), ctx, userID)
}

// GetAll mocks base method.
func (m *MocktimezoneEditor) GetAll(ctx context.Context) ([]*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MocktimezoneEditorMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MocktimezoneEditor)(nil).GetAll), ctx)
}

// Save mocks base method.
func (m *MocktimezoneEditor) Save(ctx context.Context, id int64, tz *user.Timezone) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, id, tz)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MocktimezoneEditorMockRecorder) Save(ctx, id, tz interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MocktimezoneEditor)(nil).Save), ctx, id, tz)
}
