package unit_tests

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/mock"
)

type MockPgxConn struct {
	mock.Mock
}

func (m *MockPgxConn) QueryRow(ctx context.Context, sql string, args ...interface{}) *pgx.Row {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(*pgx.Row)
}

func (m *MockPgxConn) Query(ctx context.Context, sql string, args ...interface{}) (*pgx.Rows, error) {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(*pgx.Rows), argsMock.Error(1)
}

func (m *MockPgxConn) Exec(ctx context.Context, sql string, args ...interface{}) (pgx.CommandTag, error) {
	argsMock := m.Called(ctx, sql, args)
	return argsMock.Get(0).(pgx.CommandTag), argsMock.Error(1)
}

func (m *MockPgxConn) Close() error {
	args := m.Called()
	return args.Error(0)
}
