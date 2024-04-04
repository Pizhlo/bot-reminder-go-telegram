// Code generated by MockGen. DO NOT EDIT.
// Source: ./controller.go

// Package mock_controller is a generated GoMock package.
package mock_controller

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	telebot "gopkg.in/telebot.v3"
)

// MockteleCtx is a mock of teleCtx interface.
type MockteleCtx struct {
	ctrl     *gomock.Controller
	recorder *MockteleCtxMockRecorder
}

// MockteleCtxMockRecorder is the mock recorder for MockteleCtx.
type MockteleCtxMockRecorder struct {
	mock *MockteleCtx
}

// NewMockteleCtx creates a new mock instance.
func NewMockteleCtx(ctrl *gomock.Controller) *MockteleCtx {
	mock := &MockteleCtx{ctrl: ctrl}
	mock.recorder = &MockteleCtxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockteleCtx) EXPECT() *MockteleCtxMockRecorder {
	return m.recorder
}

// Accept mocks base method.
func (m *MockteleCtx) Accept(errorMessage ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range errorMessage {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Accept", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Accept indicates an expected call of Accept.
func (mr *MockteleCtxMockRecorder) Accept(errorMessage ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*MockteleCtx)(nil).Accept), errorMessage...)
}

// Answer mocks base method.
func (m *MockteleCtx) Answer(resp *telebot.QueryResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Answer", resp)
	ret0, _ := ret[0].(error)
	return ret0
}

// Answer indicates an expected call of Answer.
func (mr *MockteleCtxMockRecorder) Answer(resp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Answer", reflect.TypeOf((*MockteleCtx)(nil).Answer), resp)
}

// Args mocks base method.
func (m *MockteleCtx) Args() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Args")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Args indicates an expected call of Args.
func (mr *MockteleCtxMockRecorder) Args() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Args", reflect.TypeOf((*MockteleCtx)(nil).Args))
}

// Bot mocks base method.
func (m *MockteleCtx) Bot() *telebot.Bot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bot")
	ret0, _ := ret[0].(*telebot.Bot)
	return ret0
}

// Bot indicates an expected call of Bot.
func (mr *MockteleCtxMockRecorder) Bot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bot", reflect.TypeOf((*MockteleCtx)(nil).Bot))
}

// Callback mocks base method.
func (m *MockteleCtx) Callback() *telebot.Callback {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Callback")
	ret0, _ := ret[0].(*telebot.Callback)
	return ret0
}

// Callback indicates an expected call of Callback.
func (mr *MockteleCtxMockRecorder) Callback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Callback", reflect.TypeOf((*MockteleCtx)(nil).Callback))
}

// Chat mocks base method.
func (m *MockteleCtx) Chat() *telebot.Chat {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chat")
	ret0, _ := ret[0].(*telebot.Chat)
	return ret0
}

// Chat indicates an expected call of Chat.
func (mr *MockteleCtxMockRecorder) Chat() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chat", reflect.TypeOf((*MockteleCtx)(nil).Chat))
}

// ChatJoinRequest mocks base method.
func (m *MockteleCtx) ChatJoinRequest() *telebot.ChatJoinRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChatJoinRequest")
	ret0, _ := ret[0].(*telebot.ChatJoinRequest)
	return ret0
}

// ChatJoinRequest indicates an expected call of ChatJoinRequest.
func (mr *MockteleCtxMockRecorder) ChatJoinRequest() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChatJoinRequest", reflect.TypeOf((*MockteleCtx)(nil).ChatJoinRequest))
}

// ChatMember mocks base method.
func (m *MockteleCtx) ChatMember() *telebot.ChatMemberUpdate {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChatMember")
	ret0, _ := ret[0].(*telebot.ChatMemberUpdate)
	return ret0
}

// ChatMember indicates an expected call of ChatMember.
func (mr *MockteleCtxMockRecorder) ChatMember() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChatMember", reflect.TypeOf((*MockteleCtx)(nil).ChatMember))
}

// Data mocks base method.
func (m *MockteleCtx) Data() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Data")
	ret0, _ := ret[0].(string)
	return ret0
}

// Data indicates an expected call of Data.
func (mr *MockteleCtxMockRecorder) Data() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Data", reflect.TypeOf((*MockteleCtx)(nil).Data))
}

// Delete mocks base method.
func (m *MockteleCtx) Delete() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete")
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockteleCtxMockRecorder) Delete() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockteleCtx)(nil).Delete))
}

// DeleteAfter mocks base method.
func (m *MockteleCtx) DeleteAfter(d time.Duration) *time.Timer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAfter", d)
	ret0, _ := ret[0].(*time.Timer)
	return ret0
}

// DeleteAfter indicates an expected call of DeleteAfter.
func (mr *MockteleCtxMockRecorder) DeleteAfter(d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAfter", reflect.TypeOf((*MockteleCtx)(nil).DeleteAfter), d)
}

// Edit mocks base method.
func (m *MockteleCtx) Edit(what interface{}, opts ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{what}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Edit", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Edit indicates an expected call of Edit.
func (mr *MockteleCtxMockRecorder) Edit(what interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{what}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockteleCtx)(nil).Edit), varargs...)
}

// EditCaption mocks base method.
func (m *MockteleCtx) EditCaption(caption string, opts ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{caption}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EditCaption", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditCaption indicates an expected call of EditCaption.
func (mr *MockteleCtxMockRecorder) EditCaption(caption interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{caption}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditCaption", reflect.TypeOf((*MockteleCtx)(nil).EditCaption), varargs...)
}

// EditOrReply mocks base method.
func (m *MockteleCtx) EditOrReply(what interface{}, opts ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{what}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EditOrReply", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditOrReply indicates an expected call of EditOrReply.
func (mr *MockteleCtxMockRecorder) EditOrReply(what interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{what}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditOrReply", reflect.TypeOf((*MockteleCtx)(nil).EditOrReply), varargs...)
}

// EditOrSend mocks base method.
func (m *MockteleCtx) EditOrSend(what interface{}, opts ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{what}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EditOrSend", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditOrSend indicates an expected call of EditOrSend.
func (mr *MockteleCtxMockRecorder) EditOrSend(what interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{what}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditOrSend", reflect.TypeOf((*MockteleCtx)(nil).EditOrSend), varargs...)
}

// Entities mocks base method.
func (m *MockteleCtx) Entities() telebot.Entities {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Entities")
	ret0, _ := ret[0].(telebot.Entities)
	return ret0
}

// Entities indicates an expected call of Entities.
func (mr *MockteleCtxMockRecorder) Entities() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Entities", reflect.TypeOf((*MockteleCtx)(nil).Entities))
}

// Forward mocks base method.
func (m *MockteleCtx) Forward(msg telebot.Editable, opts ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{msg}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Forward", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Forward indicates an expected call of Forward.
func (mr *MockteleCtxMockRecorder) Forward(msg interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{msg}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Forward", reflect.TypeOf((*MockteleCtx)(nil).Forward), varargs...)
}

// ForwardTo mocks base method.
func (m *MockteleCtx) ForwardTo(to telebot.Recipient, opts ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{to}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ForwardTo", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForwardTo indicates an expected call of ForwardTo.
func (mr *MockteleCtxMockRecorder) ForwardTo(to interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{to}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForwardTo", reflect.TypeOf((*MockteleCtx)(nil).ForwardTo), varargs...)
}

// Get mocks base method.
func (m *MockteleCtx) Get(key string) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockteleCtxMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockteleCtx)(nil).Get), key)
}

// InlineResult mocks base method.
func (m *MockteleCtx) InlineResult() *telebot.InlineResult {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InlineResult")
	ret0, _ := ret[0].(*telebot.InlineResult)
	return ret0
}

// InlineResult indicates an expected call of InlineResult.
func (mr *MockteleCtxMockRecorder) InlineResult() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InlineResult", reflect.TypeOf((*MockteleCtx)(nil).InlineResult))
}

// Message mocks base method.
func (m *MockteleCtx) Message() *telebot.Message {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Message")
	ret0, _ := ret[0].(*telebot.Message)
	return ret0
}

// Message indicates an expected call of Message.
func (mr *MockteleCtxMockRecorder) Message() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Message", reflect.TypeOf((*MockteleCtx)(nil).Message))
}

// Migration mocks base method.
func (m *MockteleCtx) Migration() (int64, int64) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Migration")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(int64)
	return ret0, ret1
}

// Migration indicates an expected call of Migration.
func (mr *MockteleCtxMockRecorder) Migration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Migration", reflect.TypeOf((*MockteleCtx)(nil).Migration))
}

// Notify mocks base method.
func (m *MockteleCtx) Notify(action telebot.ChatAction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Notify", action)
	ret0, _ := ret[0].(error)
	return ret0
}

// Notify indicates an expected call of Notify.
func (mr *MockteleCtxMockRecorder) Notify(action interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Notify", reflect.TypeOf((*MockteleCtx)(nil).Notify), action)
}

// Poll mocks base method.
func (m *MockteleCtx) Poll() *telebot.Poll {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Poll")
	ret0, _ := ret[0].(*telebot.Poll)
	return ret0
}

// Poll indicates an expected call of Poll.
func (mr *MockteleCtxMockRecorder) Poll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Poll", reflect.TypeOf((*MockteleCtx)(nil).Poll))
}

// PollAnswer mocks base method.
func (m *MockteleCtx) PollAnswer() *telebot.PollAnswer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PollAnswer")
	ret0, _ := ret[0].(*telebot.PollAnswer)
	return ret0
}

// PollAnswer indicates an expected call of PollAnswer.
func (mr *MockteleCtxMockRecorder) PollAnswer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PollAnswer", reflect.TypeOf((*MockteleCtx)(nil).PollAnswer))
}

// PreCheckoutQuery mocks base method.
func (m *MockteleCtx) PreCheckoutQuery() *telebot.PreCheckoutQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreCheckoutQuery")
	ret0, _ := ret[0].(*telebot.PreCheckoutQuery)
	return ret0
}

// PreCheckoutQuery indicates an expected call of PreCheckoutQuery.
func (mr *MockteleCtxMockRecorder) PreCheckoutQuery() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreCheckoutQuery", reflect.TypeOf((*MockteleCtx)(nil).PreCheckoutQuery))
}

// Query mocks base method.
func (m *MockteleCtx) Query() *telebot.Query {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query")
	ret0, _ := ret[0].(*telebot.Query)
	return ret0
}

// Query indicates an expected call of Query.
func (mr *MockteleCtxMockRecorder) Query() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockteleCtx)(nil).Query))
}

// Recipient mocks base method.
func (m *MockteleCtx) Recipient() telebot.Recipient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recipient")
	ret0, _ := ret[0].(telebot.Recipient)
	return ret0
}

// Recipient indicates an expected call of Recipient.
func (mr *MockteleCtxMockRecorder) Recipient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recipient", reflect.TypeOf((*MockteleCtx)(nil).Recipient))
}

// Reply mocks base method.
func (m *MockteleCtx) Reply(what interface{}, opts ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{what}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Reply", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reply indicates an expected call of Reply.
func (mr *MockteleCtxMockRecorder) Reply(what interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{what}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reply", reflect.TypeOf((*MockteleCtx)(nil).Reply), varargs...)
}

// Respond mocks base method.
func (m *MockteleCtx) Respond(resp ...*telebot.CallbackResponse) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range resp {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Respond", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Respond indicates an expected call of Respond.
func (mr *MockteleCtxMockRecorder) Respond(resp ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Respond", reflect.TypeOf((*MockteleCtx)(nil).Respond), resp...)
}

// Send mocks base method.
func (m *MockteleCtx) Send(what interface{}, opts ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{what}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Send", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockteleCtxMockRecorder) Send(what interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{what}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockteleCtx)(nil).Send), varargs...)
}

// SendAlbum mocks base method.
func (m *MockteleCtx) SendAlbum(a telebot.Album, opts ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{a}
	for _, a_2 := range opts {
		varargs = append(varargs, a_2)
	}
	ret := m.ctrl.Call(m, "SendAlbum", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendAlbum indicates an expected call of SendAlbum.
func (mr *MockteleCtxMockRecorder) SendAlbum(a interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{a}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAlbum", reflect.TypeOf((*MockteleCtx)(nil).SendAlbum), varargs...)
}

// Sender mocks base method.
func (m *MockteleCtx) Sender() *telebot.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sender")
	ret0, _ := ret[0].(*telebot.User)
	return ret0
}

// Sender indicates an expected call of Sender.
func (mr *MockteleCtxMockRecorder) Sender() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sender", reflect.TypeOf((*MockteleCtx)(nil).Sender))
}

// Set mocks base method.
func (m *MockteleCtx) Set(key string, val interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", key, val)
}

// Set indicates an expected call of Set.
func (mr *MockteleCtxMockRecorder) Set(key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockteleCtx)(nil).Set), key, val)
}

// Ship mocks base method.
func (m *MockteleCtx) Ship(what ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range what {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Ship", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ship indicates an expected call of Ship.
func (mr *MockteleCtxMockRecorder) Ship(what ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ship", reflect.TypeOf((*MockteleCtx)(nil).Ship), what...)
}

// ShippingQuery mocks base method.
func (m *MockteleCtx) ShippingQuery() *telebot.ShippingQuery {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShippingQuery")
	ret0, _ := ret[0].(*telebot.ShippingQuery)
	return ret0
}

// ShippingQuery indicates an expected call of ShippingQuery.
func (mr *MockteleCtxMockRecorder) ShippingQuery() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShippingQuery", reflect.TypeOf((*MockteleCtx)(nil).ShippingQuery))
}

// Text mocks base method.
func (m *MockteleCtx) Text() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Text")
	ret0, _ := ret[0].(string)
	return ret0
}

// Text indicates an expected call of Text.
func (mr *MockteleCtxMockRecorder) Text() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Text", reflect.TypeOf((*MockteleCtx)(nil).Text))
}

// Update mocks base method.
func (m *MockteleCtx) Update() telebot.Update {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update")
	ret0, _ := ret[0].(telebot.Update)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockteleCtxMockRecorder) Update() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockteleCtx)(nil).Update))
}
