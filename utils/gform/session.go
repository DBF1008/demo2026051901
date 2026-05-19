package gform

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gohouse/t"
)

// Session ...
type Session struct {
	IEngin
	IBinder
	master       *sql.DB
	tx           *sql.Tx
	slave        *sql.DB
	lastInsertId int64
	sqlLogs      []string
	lastSql      string
	union        interface{}
	transaction  bool
	err          error
}

var _ ISession = (*Session)(nil)

func NewSession(e IEngin) *Session {
	var s = new(Session)
	s.IEngin = e
	s.SetIBinder(NewBinder())

	s.master = e.GetExecuteDB()
	s.slave = e.GetQueryDB()

	return s
}

func (s *Session) Close() {
	s.master.Close()
	s.slave.Close()
}

func (s *Session) GetIEngin() IEngin {
	return s.IEngin
}

func (s *Session) SetIEngin(ie IEngin) {
	s.IEngin = ie
}

func (s *Session) Bind(tab interface{}) ISession {
	//fmt.Println(tab, NewBinder(tab))
	//s.SetIBinder(NewBinder(tab))
	s.GetIBinder().SetBindOrigin(tab)
	s.err = s.IBinder.BindParse(s.GetIEngin().GetPrefix())
	return s
}

func (s *Session) GetErr() error {
	return s.err
}

func (s *Session) SetIBinder(ib IBinder) {
	s.IBinder = ib
}

func (s *Session) GetIBinder() IBinder {
	return s.IBinder
}

func (s *Session) ResetBinderResult() {
	_ = s.IBinder.BindParse(s.GetIEngin().GetPrefix())
}

func (s *Session) GetTableName() (string, error) {
	//fmt.Println(s.GetIBinder())
	return s.GetIBinder().GetBindName(), s.err
}

// Begin ...
func (s *Session) Begin() (err error) {
	s.tx, err = s.master.Begin()
	s.SetTransaction(true)
	return
}

// Rollback ...
func (s *Session) Rollback() (err error) {
	err = s.tx.Rollback()
	s.tx = nil
	s.SetTransaction(false)
	return
}

// Commit ...
func (s *Session) Commit() (err error) {
	err = s.tx.Commit()
	s.tx = nil
	s.SetTransaction(false)
	return
}

// Transaction ...
func (s *Session) Transaction(closers ...func(ses ISession) error) (err error) {
	err = s.Begin()
	if err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return err
	}

	for _, closer := range closers {
		err = closer(s)
		if err != nil {
			s.GetIEngin().GetLogger().Error(err.Error())
			_ = s.Rollback()
			return
		}
	}
	return s.Commit()
}

// Query ...
func (s *Session) Query(sqlstring string, args ...interface{}) (result []Data, err error) {
	start := time.Now()
	//withRunTimeContext(func() {
	if s.err != nil {
		err = s.err
		s.GetIEngin().GetLogger().Error(err.Error())
	}
	s.lastSql = fmt.Sprint(sqlstring, ", ", args)
	//if s.IfEnableSqlLog() {

	var stmt *sql.Stmt
	if s.tx == nil {
		stmt, err = s.slave.Prepare(sqlstring)
	} else {
		stmt, err = s.tx.Prepare(sqlstring)
	}

	if err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return
	}

	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return
	}

	// make sure we always close rows
	defer rows.Close()

	err = s.scan(rows)
	if err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return
	}

	timeduration := time.Since(start)
	//if timeduration.Seconds() > 1 {
	s.GetIEngin().GetLogger().Slow(s.LastSql(), timeduration)
	s.GetIEngin().GetLogger().Sql(s.LastSql(), timeduration)

	result = s.GetIBinder().GetBindAll()
	return
}

// Execute ...
func (s *Session) Execute(sqlstring string, args ...interface{}) (rowsAffected int64, err error) {
	start := time.Now()
	//withRunTimeContext(func() {
	if s.err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return
	}
	s.lastSql = fmt.Sprint(sqlstring, ", ", args)
	//if s.IfEnableSqlLog() {

	var operType = strings.ToLower(sqlstring[0:6])
	if operType == "select" {
		s.GetIEngin().GetLogger().Error(err.Error())
		err = errors.New("Execute does not allow select operations, please use Query")
		return
	}

	var stmt *sql.Stmt
	if s.tx == nil {
		stmt, err = s.master.Prepare(sqlstring)
	} else {
		stmt, err = s.tx.Prepare(sqlstring)
	}

	if err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return
	}

	//var err error
	defer stmt.Close()
	result, err := stmt.Exec(args...)
	if err != nil {
		s.GetIEngin().GetLogger().Error(err.Error())
		return
	}

	if operType == "insert" {
		// get last insert id
		lastInsertId, err := result.LastInsertId()
		if err == nil {
			s.lastInsertId = lastInsertId
		} else {
			s.GetIEngin().GetLogger().Error(err.Error())
		}
	}
	// get rows affected
	rowsAffected, err = result.RowsAffected()
	timeduration := time.Since(start)
	if timeduration.Seconds() > 1 {
		s.GetIEngin().GetLogger().Slow(s.LastSql(), timeduration)
	} else {
		s.GetIEngin().GetLogger().Sql(s.LastSql(), timeduration)
	}
	return
}

// LastInsertId ...
func (s *Session) LastInsertId() int64 {
	return s.lastInsertId
}

// LastSql ...
func (s *Session) LastSql() string {
	return s.lastSql
}

func (s *Session) scan(rows *sql.Rows) (err error) {
	if s.GetIBinder() == nil {
		s.SetIBinder(NewBinder())
	}
	switch s.GetBindType() {
	case OBJECT_STRING:
		err = s.scanAll(rows)
	case OBJECT_STRUCT, OBJECT_STRUCT_SLICE:
		err = s.scanStructAll(rows)
	//case OBJECT_MAP, OBJECT_MAP_T:
	case OBJECT_MAP, OBJECT_MAP_T, OBJECT_MAP_SLICE, OBJECT_MAP_SLICE_T:
		err = s.scanMapAll(rows)
	case OBJECT_NIL:
		err = s.scanAll(rows)
	default:
		err = errors.New("Bind value error")
	}
	return
}

//	return s.scanMapAll(rows, dst)

func (s *Session) scanMapAll(rows *sql.Rows) (err error) {
	var columns []string
	if columns, err = rows.Columns(); err != nil {
		return
	}
	count := len(columns)

	for rows.Next() {
		values := make([]interface{}, count)
		scanArgs := make([]interface{}, count)
		for i := 0; i < count; i++ {
			scanArgs[i] = &values[i]
		}
		_ = rows.Scan(scanArgs...)

		var bindResultTmp = reflect.MakeMap(reflect.Indirect(reflect.ValueOf(s.GetBindResult())).Type())
		//var unionTmp = map[string]interface{}{}
		for i, col := range columns {
			var v interface{}
			val := values[i]
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = val
			}
			if s.GetUnion() != nil {
				s.union = v
				return
				//unionTmp[col] = v
				//s.union = unionTmp
			} else {
				br := reflect.Indirect(reflect.ValueOf(s.GetBindResult()))
				switch s.GetBindType() {
				case OBJECT_MAP_T, OBJECT_MAP_SLICE_T:
					br.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(t.New(v)))
					if s.GetBindType() == OBJECT_MAP_SLICE || s.GetBindType() == OBJECT_MAP_SLICE_T {
						bindResultTmp.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(t.New(v)))
					}
				default:
					br.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(v))
					if s.GetBindType() == OBJECT_MAP_SLICE || s.GetBindType() == OBJECT_MAP_SLICE_T {
						bindResultTmp.SetMapIndex(reflect.ValueOf(col), reflect.ValueOf(v))
					}
				}
			}
		}
		if s.GetUnion() == nil {
			if s.GetBindType() == OBJECT_MAP_SLICE || s.GetBindType() == OBJECT_MAP_SLICE_T {
				s.GetBindResultSlice().Set(reflect.Append(s.GetBindResultSlice(), bindResultTmp))
			}
		}
	}
	return
}

func (s *Session) scanStructAll(rows *sql.Rows) error {
	// check if there is data waiting
	//if !rows.Next() {
	//	if err := rows.Err(); err != nil {
	//		return err
	//	return sql.ErrNoRows
	var sfs = structForScan(s.GetBindResult())
	for rows.Next() {
		if s.GetUnion() != nil {
			var union interface{}
			err := rows.Scan(&union)
			if err != nil {
				s.GetIEngin().GetLogger().Error(err.Error())
				return err
			}
			s.union = union
			return err
		}
		// scan it
		err := rows.Scan(sfs...)
		if err != nil {
			s.GetIEngin().GetLogger().Error(err.Error())
			return err
		}
		if s.GetUnion() == nil {
			if s.GetBindType() == OBJECT_STRUCT_SLICE {
				// add to the result slice
				s.GetBindResultSlice().Set(reflect.Append(s.GetBindResultSlice(),
					reflect.Indirect(reflect.ValueOf(s.GetBindResult()))))
			}
		}
	}

	return rows.Err()
}

func (s *Session) scanAll(rows *sql.Rows) (err error) {
	var columns []string
	if columns, err = rows.Columns(); err != nil {
		return
	}
	count := len(columns)

	var result = []Data{}
	for rows.Next() {
		values := make([]interface{}, count)
		scanArgs := make([]interface{}, count)
		for i := 0; i < count; i++ {
			scanArgs[i] = &values[i]
		}
		_ = rows.Scan(scanArgs...)

		var resultTmp = Data{}
		//var unionTmp = map[string]interface{}{}
		for i, col := range columns {
			var v interface{}
			val := values[i]
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = val
			}
			if s.GetUnion() != nil {
				s.union = v
				return
				//unionTmp[col] = v
				//s.union = unionTmp
			}
			resultTmp[col] = v
		}
		result = append(result, resultTmp)
	}
	s.IBinder.SetBindAll(result)
	return
}

// SetUnion ...
func (s *Session) SetUnion(u interface{}) {
	s.union = u
}

// GetUnion ...
func (s *Session) GetUnion() interface{} {
	return s.union
}

// SetTransaction ...
func (s *Session) SetTransaction(b bool) {
	s.transaction = b
}

func (s *Session) GetTransaction() bool {
	return s.transaction
}
