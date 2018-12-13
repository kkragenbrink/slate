// Code generated by MockGen. DO NOT EDIT.
// Source: domains/domain.go

// Package domains is a generated GoMock package.
package domains

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCharacterRepository is a mock of CharacterRepository interface
type MockCharacterRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCharacterRepositoryMockRecorder
}

// MockCharacterRepositoryMockRecorder is the mock recorder for MockCharacterRepository
type MockCharacterRepositoryMockRecorder struct {
	mock *MockCharacterRepository
}

// NewMockCharacterRepository creates a new mock instance
func NewMockCharacterRepository(ctrl *gomock.Controller) *MockCharacterRepository {
	mock := &MockCharacterRepository{ctrl: ctrl}
	mock.recorder = &MockCharacterRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCharacterRepository) EXPECT() *MockCharacterRepositoryMockRecorder {
	return m.recorder
}

// FindByID mocks base method
func (m *MockCharacterRepository) FindByID(ctx context.Context, id int64) (*Character, error) {
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*Character)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID
func (mr *MockCharacterRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockCharacterRepository)(nil).FindByID), ctx, id)
}

// Store mocks base method
func (m *MockCharacterRepository) Store(ctx context.Context, c *Character) error {
	ret := m.ctrl.Call(m, "Store", ctx, c)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store
func (mr *MockCharacterRepositoryMockRecorder) Store(ctx, c interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockCharacterRepository)(nil).Store), ctx, c)
}

// MockSheet is a mock of Sheet interface
type MockSheet struct {
	ctrl     *gomock.Controller
	recorder *MockSheetMockRecorder
}

// MockSheetMockRecorder is the mock recorder for MockSheet
type MockSheetMockRecorder struct {
	mock *MockSheet
}

// NewMockSheet creates a new mock instance
func NewMockSheet(ctrl *gomock.Controller) *MockSheet {
	mock := &MockSheet{ctrl: ctrl}
	mock.recorder = &MockSheetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSheet) EXPECT() *MockSheetMockRecorder {
	return m.recorder
}

// System mocks base method
func (m *MockSheet) System() string {
	ret := m.ctrl.Call(m, "System")
	ret0, _ := ret[0].(string)
	return ret0
}

// System indicates an expected call of System
func (mr *MockSheetMockRecorder) System() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "System", reflect.TypeOf((*MockSheet)(nil).System))
}