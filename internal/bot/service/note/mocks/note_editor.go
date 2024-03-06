// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go

// Package mock_note is a generated GoMock package.
package mock_note

import (
	context "context"
	reflect "reflect"

	model "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	gomock "github.com/golang/mock/gomock"
)

// MocknoteEditor is a mock of noteEditor interface.
type MocknoteEditor struct {
	ctrl     *gomock.Controller
	recorder *MocknoteEditorMockRecorder
}

// MocknoteEditorMockRecorder is the mock recorder for MocknoteEditor.
type MocknoteEditorMockRecorder struct {
	mock *MocknoteEditor
}

// NewMocknoteEditor creates a new mock instance.
func NewMocknoteEditor(ctrl *gomock.Controller) *MocknoteEditor {
	mock := &MocknoteEditor{ctrl: ctrl}
	mock.recorder = &MocknoteEditorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocknoteEditor) EXPECT() *MocknoteEditorMockRecorder {
	return m.recorder
}

// DeleteAllByUserID mocks base method.
func (m *MocknoteEditor) DeleteAllByUserID(ctx context.Context, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllByUserID", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllByUserID indicates an expected call of DeleteAllByUserID.
func (mr *MocknoteEditorMockRecorder) DeleteAllByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllByUserID", reflect.TypeOf((*MocknoteEditor)(nil).DeleteAllByUserID), ctx, userID)
}

// DeleteNoteByID mocks base method.
func (m *MocknoteEditor) DeleteNoteByID(ctx context.Context, userID int64, noteID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNoteByID", ctx, userID, noteID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNoteByID indicates an expected call of DeleteNoteByID.
func (mr *MocknoteEditorMockRecorder) DeleteNoteByID(ctx, userID, noteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNoteByID", reflect.TypeOf((*MocknoteEditor)(nil).DeleteNoteByID), ctx, userID, noteID)
}

// GetAllByUserID mocks base method.
func (m *MocknoteEditor) GetAllByUserID(ctx context.Context, userID int64) ([]model.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByUserID", ctx, userID)
	ret0, _ := ret[0].([]model.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByUserID indicates an expected call of GetAllByUserID.
func (mr *MocknoteEditorMockRecorder) GetAllByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByUserID", reflect.TypeOf((*MocknoteEditor)(nil).GetAllByUserID), ctx, userID)
}

// GetByID mocks base method.
func (m *MocknoteEditor) GetByID(ctx context.Context, userID int64, noteID int) (*model.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, userID, noteID)
	ret0, _ := ret[0].(*model.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MocknoteEditorMockRecorder) GetByID(ctx, userID, noteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MocknoteEditor)(nil).GetByID), ctx, userID, noteID)
}

// Save mocks base method.
func (m *MocknoteEditor) Save(ctx context.Context, note model.Note) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, note)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MocknoteEditorMockRecorder) Save(ctx, note interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MocknoteEditor)(nil).Save), ctx, note)
}

// SearchByOneDate mocks base method.
func (m *MocknoteEditor) SearchByOneDate(ctx context.Context, searchNote model.SearchByOneDate) ([]model.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchByOneDate", ctx, searchNote)
	ret0, _ := ret[0].([]model.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchByOneDate indicates an expected call of SearchByOneDate.
func (mr *MocknoteEditorMockRecorder) SearchByOneDate(ctx, searchNote interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchByOneDate", reflect.TypeOf((*MocknoteEditor)(nil).SearchByOneDate), ctx, searchNote)
}

// SearchByText mocks base method.
func (m *MocknoteEditor) SearchByText(ctx context.Context, searchNote model.SearchByText) ([]model.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchByText", ctx, searchNote)
	ret0, _ := ret[0].([]model.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchByText indicates an expected call of SearchByText.
func (mr *MocknoteEditorMockRecorder) SearchByText(ctx, searchNote interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchByText", reflect.TypeOf((*MocknoteEditor)(nil).SearchByText), ctx, searchNote)
}
