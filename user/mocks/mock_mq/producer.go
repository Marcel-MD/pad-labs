// Code generated by MockGen. DO NOT EDIT.
// Source: api/mq/producer.go
//
// Generated by this command:
//
//	mockgen -source=api/mq/producer.go -destination=mocks/mock_mq/producer.go
//
// Package mock_mq is a generated GoMock package.
package mock_mq

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockProducer is a mock of Producer interface.
type MockProducer struct {
	ctrl     *gomock.Controller
	recorder *MockProducerMockRecorder
}

// MockProducerMockRecorder is the mock recorder for MockProducer.
type MockProducerMockRecorder struct {
	mock *MockProducer
}

// NewMockProducer creates a new mock instance.
func NewMockProducer(ctrl *gomock.Controller) *MockProducer {
	mock := &MockProducer{ctrl: ctrl}
	mock.recorder = &MockProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProducer) EXPECT() *MockProducerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockProducer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockProducerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockProducer)(nil).Close))
}

// SendMsg mocks base method.
func (m *MockProducer) SendMsg(msgType string, msg any, queues []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendMsg", msgType, msg, queues)
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockProducerMockRecorder) SendMsg(msgType, msg, queues any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockProducer)(nil).SendMsg), msgType, msg, queues)
}