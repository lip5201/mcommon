package mcommon

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	// 导入mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DbExeAble 数据库接口
type DbExeAble interface {
	Rebind(string) string
	Get(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// isShowSQL 是否显示执行的sql语句
var isShowSQL bool

// debugSQLMap 记录
var debugSQLMap map[string]string

// debugSQLCountMap sql次数
var debugSQLCountMap map[string]int64

// DbCreate 创建数据库链接
func DbCreate(dataSourceName string, showSQL bool) *sqlx.DB {
	isShowSQL = showSQL
	debugSQLMap = make(map[string]string)
	debugSQLCountMap = make(map[string]int64)

	var err error
	var db *sqlx.DB

	db, err = sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		Log.Fatalf("db connect error: %s", err.Error())
		return nil
	}

	count := runtime.NumCPU()*20 + 1
	db.SetMaxOpenConns(count)
	db.SetMaxIdleConns(count)
	db.SetConnMaxLifetime(1 * time.Hour)

	err = db.Ping()
	if err != nil {
		Log.Fatalf("db ping error: %s", err.Error())
		return nil
	}
	return db
}

// DbExecuteCountManyContent 返回sql语句并返回执行行数
func DbExecuteCountManyContent(ctx context.Context, tx DbExeAble, query string, n int, args ...interface{}) (int64, error) {
	var err error
	insertArgs := strings.Repeat("(?),", n)
	insertArgs = strings.TrimSuffix(insertArgs, ",")
	query = fmt.Sprintf(query, insertArgs)
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = tx.Rebind(query)
	if isShowSQL {
		queryStr := query + ";"
		for _, arg := range args {
			_, ok := arg.(string)
			if ok {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`"%s"`, arg), 1)
			} else {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`%v`, arg), 1)
			}
		}
		Log.Debugf(queryStr)
		debugSQLMap[query] = queryStr
		debugSQLCountMap[query] += 1
	}
	ret, err := tx.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return 0, err
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// DbExecuteLastIDNamedContent 执行sql语句并返回lastID
func DbExecuteLastIDNamedContent(ctx context.Context, tx DbExeAble, query string, argMap map[string]interface{}) (int64, error) {
	query, args, err := sqlx.Named(query, argMap)
	if err != nil {
		return 0, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = tx.Rebind(query)
	if isShowSQL {
		queryStr := query + ";"
		for _, arg := range args {
			_, ok := arg.(string)
			if ok {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`"%s"`, arg), 1)
			} else {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`%v`, arg), 1)
			}
		}
		Log.Debugf(queryStr)
		debugSQLMap[query] = queryStr
		debugSQLCountMap[query] += 1
	}
	ret, err := tx.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return 0, err
	}
	lastID, err := ret.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastID, nil
}

// DbExecuteCountNamedContent 执行sql语句返回执行个数
func DbExecuteCountNamedContent(ctx context.Context, tx DbExeAble, query string, argMap map[string]interface{}) (int64, error) {
	query, args, err := sqlx.Named(query, argMap)
	if err != nil {
		return 0, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return 0, err
	}
	query = tx.Rebind(query)
	if isShowSQL {
		queryStr := query + ";"
		for _, arg := range args {
			_, ok := arg.(string)
			if ok {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`"%s"`, arg), 1)
			} else {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`%v`, arg), 1)
			}
		}
		Log.Debugf(queryStr)
		debugSQLMap[query] = queryStr
		debugSQLCountMap[query] += 1
	}
	ret, err := tx.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return 0, err
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// DbGetNamedContent 执行sql查询并返回当个元素
func DbGetNamedContent(ctx context.Context, tx DbExeAble, dest interface{}, query string, argMap map[string]interface{}) (bool, error) {
	query, args, err := sqlx.Named(query, argMap)
	if err != nil {
		return false, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return false, err
	}
	query = tx.Rebind(query)
	if isShowSQL {
		queryStr := query + ";"
		for _, arg := range args {
			_, ok := arg.(string)
			if ok {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`"%s"`, arg), 1)
			} else {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`%v`, arg), 1)
			}
		}
		Log.Debugf(queryStr)
		debugSQLMap[query] = queryStr
		debugSQLCountMap[query] += 1
	}
	err = tx.GetContext(
		ctx,
		dest,
		query,
		args...,
	)
	if err == sql.ErrNoRows {
		// 没有元素
		return false, nil
	}
	if err != nil {
		// 执行错误
		return false, err
	}
	return true, nil
}

// DbSelectNamedContent 执行sql查询并返回多行
func DbSelectNamedContent(ctx context.Context, tx DbExeAble, dest interface{}, query string, argMap map[string]interface{}) error {
	query, args, err := sqlx.Named(query, argMap)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	query = tx.Rebind(query)
	if isShowSQL {
		queryStr := query + ";"
		for _, arg := range args {
			_, ok := arg.(string)
			if ok {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`"%s"`, arg), 1)
			} else {
				queryStr = strings.Replace(queryStr, "?", fmt.Sprintf(`%v`, arg), 1)
			}
		}
		Log.Debugf(queryStr)
		debugSQLMap[query] = queryStr
		debugSQLCountMap[query] += 1
	}
	err = tx.SelectContext(
		ctx,
		dest,
		query,
		args...,
	)
	if err == sql.ErrNoRows {
		// 没有元素
		return nil
	}
	if err != nil {
		// 执行错误
		return err
	}
	return nil
}

// DbUpdateKV 更新
func DbUpdateKV(ctx context.Context, tx DbExeAble, table string, updateMap H, keys []string, values []interface{}) (int64, error) {
	keysLen := len(keys)
	if 0 == keysLen {
		return 0, fmt.Errorf("keys len error")
	}
	if keysLen != len(values) {
		return 0, fmt.Errorf("value len error")
	}
	updateLastIndex := len(updateMap) - 1

	argMap := H{}
	query := strings.Builder{}
	query.WriteString("UPDATE\n")
	query.WriteString(table)
	query.WriteString("\nSET\n")
	var updateIndex int
	for k, v := range updateMap {
		query.WriteString(k)
		query.WriteString("=:")
		query.WriteString(k)
		if updateIndex == updateLastIndex {
			query.WriteString("\n")
		} else {
			query.WriteString(",\n")
		}
		updateIndex++
		argMap[k] = v
	}
	query.WriteString("WHERE\n")
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return 0, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(")")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	count, err := DbExecuteCountNamedContent(
		ctx,
		tx,
		query.String(),
		argMap,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// DbDeleteKV 删除
func DbDeleteKV(ctx context.Context, tx DbExeAble, table string, keys []string, values []interface{}) (int64, error) {
	keysLen := len(keys)
	if 0 == keysLen {
		return 0, fmt.Errorf("keys len error")
	}
	if keysLen != len(values) {
		return 0, fmt.Errorf("value len error")
	}
	argMap := H{}

	query := strings.Builder{}
	query.WriteString("DELETE\nFROM\n")
	query.WriteString(table)
	query.WriteString("\nWHERE\n")
	for i, key := range keys {
		if i != 0 {
			query.WriteString("AND ")
		}
		value := values[i]
		query.WriteString(key)
		rt := reflect.TypeOf(value)
		switch rt.Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)
			if s.Len() == 0 {
				return 0, nil
			}
			query.WriteString(" IN (:")
			query.WriteString(key)
			query.WriteString(")")
		default:
			query.WriteString("=:")
			query.WriteString(key)
		}
		query.WriteString("\n")
		argMap[key] = value
	}

	count, err := DbExecuteCountNamedContent(
		ctx,
		tx,
		query.String(),
		argMap,
	)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// DbTransaction 执行事物
func DbTransaction(ctx context.Context, db *sqlx.DB, f func(dbTx DbExeAble) error) error {
	isComment := false
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if !isComment {
			_ = tx.Rollback()
		}
	}()
	err = f(tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	isComment = true
	return nil
}

// DbGetDebugMap 获取debug sql 记录
func DbGetDebugMap() map[string]string {
	return debugSQLMap
}

// DbGetDebugCountMap 获取debug sql 次数
func DbGetDebugCountMap() map[string]int64 {
	return debugSQLCountMap
}
