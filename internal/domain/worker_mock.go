// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/worker.go

// Package domain is a generated GoMock package.
package domain

import (
	context "context"
	reflect "reflect"

	civil "cloud.google.com/go/civil"
	gomock "github.com/golang/mock/gomock"
	cursor "github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/cursor"
)

// MockWorkerRepository is a mock of WorkerRepository interface.
type MockWorkerRepository struct {
	ctrl     *gomock.Controller
	recorder *MockWorkerRepositoryMockRecorder
}

// MockWorkerRepositoryMockRecorder is the mock recorder for MockWorkerRepository.
type MockWorkerRepositoryMockRecorder struct {
	mock *MockWorkerRepository
}

// NewMockWorkerRepository creates a new mock instance.
func NewMockWorkerRepository(ctrl *gomock.Controller) *MockWorkerRepository {
	mock := &MockWorkerRepository{ctrl: ctrl}
	mock.recorder = &MockWorkerRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorkerRepository) EXPECT() *MockWorkerRepositoryMockRecorder {
	return m.recorder
}

// GetAvailableShifts mocks base method.
func (m *MockWorkerRepository) GetAvailableShifts(ctx context.Context, queryCursor *cursor.Cursor, limit int, workerID int64, start, end civil.Date) (Shifts, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailableShifts", ctx, queryCursor, limit, workerID, start, end)
	ret0, _ := ret[0].(Shifts)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvailableShifts indicates an expected call of GetAvailableShifts.
func (mr *MockWorkerRepositoryMockRecorder) GetAvailableShifts(ctx, queryCursor, limit, workerID, start, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailableShifts", reflect.TypeOf((*MockWorkerRepository)(nil).GetAvailableShifts), ctx, queryCursor, limit, workerID, start, end)
}
