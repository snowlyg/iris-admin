package admin

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testDriverSequence atomic.Uint64

type testSQLDriver struct {
	query func(string) (driver.Rows, error)
	exec  func(string) error
}

func (d *testSQLDriver) Open(string) (driver.Conn, error) {
	return &testSQLConn{driver: d}, nil
}

type testSQLConn struct {
	driver *testSQLDriver
}

func (c *testSQLConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare not supported")
}

func (c *testSQLConn) Close() error {
	return nil
}

func (c *testSQLConn) Begin() (driver.Tx, error) {
	return testSQLTx{}, nil
}

func (c *testSQLConn) Ping(context.Context) error {
	return nil
}

func (c *testSQLConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	if c.driver.query == nil {
		return &testSQLRows{}, nil
	}
	return c.driver.query(query)
}

func (c *testSQLConn) ExecContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Result, error) {
	if c.driver.exec != nil {
		if err := c.driver.exec(query); err != nil {
			return nil, err
		}
	}
	return testSQLResult{lastInsertID: 1, rowsAffected: 1}, nil
}

type testSQLRows struct {
	columns []string
	values  [][]driver.Value
	index   int
}

func (r *testSQLRows) Columns() []string {
	return r.columns
}

func (r *testSQLRows) Close() error {
	return nil
}

func (r *testSQLRows) Next(dest []driver.Value) error {
	if r.index >= len(r.values) {
		return io.EOF
	}
	copy(dest, r.values[r.index])
	r.index++
	return nil
}

type testSQLTx struct{}

func (testSQLTx) Commit() error {
	return nil
}

func (testSQLTx) Rollback() error {
	return nil
}

type testSQLResult struct {
	lastInsertID int64
	rowsAffected int64
}

func (r testSQLResult) LastInsertId() (int64, error) {
	return r.lastInsertID, nil
}

func (r testSQLResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

func openTestGorm(t *testing.T, scripted *testSQLDriver) *gorm.DB {
	t.Helper()
	name := fmt.Sprintf("iris-admin-test-%d", testDriverSequence.Add(1))
	sql.Register(name, scripted)
	sqlDB, err := sql.Open(name, "")
	if err != nil {
		t.Fatalf("open test sql database: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test gorm database: %v", err)
	}
	return db
}
