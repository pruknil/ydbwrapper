package orm

import (
	"context"
	"database/sql"
	"reflect"
	"time"
)

// Ormer define the orm interface
type Ormer interface {
	// read data to model
	// for example:
	//	this will find User by Id field
	// 	u = &User{Id: user.Id}
	// 	err = Ormer.Read(u)
	//	this will find User by UserName field
	// 	u = &User{UserName: "astaxie", Password: "pass"}
	//	err = Ormer.Read(u, "UserName")
	Read(md interface{}, cols ...string) error
	// Like Read(), but with "FOR UPDATE" clause, useful in transaction.
	// Some databases are not support this feature.
	ReadForUpdate(md interface{}, cols ...string) error
	// Try to read a row from the database, or insert one if it doesn't exist
	ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error)
	// insert model data to database
	// for example:
	//  user := new(User)
	//  id, err = Ormer.Insert(user)
	//  user must be a pointer and Insert will set user's pk field
	Insert(interface{}) (int64, error)
	// mysql:InsertOrUpdate(model) or InsertOrUpdate(model,"colu=colu+value")
	// if colu type is integer : can use(+-*/), string : convert(colu,"value")
	// postgres: InsertOrUpdate(model,"conflictColumnName") or InsertOrUpdate(model,"conflictColumnName","colu=colu+value")
	// if colu type is integer : can use(+-*/), string : colu || "value"
	InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error)
	// insert some models to database
	InsertMulti(bulk int, mds interface{}) (int64, error)
	// update model to database.
	// cols set the columns those want to update.
	// find model by Id(pk) field and update columns specified by fields, if cols is null then update all columns
	// for example:
	// user := User{Id: 2}
	//	user.Langs = append(user.Langs, "zh-CN", "en-US")
	//	user.Extra.Name = "beego"
	//	user.Extra.Data = "orm"
	//	num, err = Ormer.Update(&user, "Langs", "Extra")
	Update(md interface{}, cols ...string) (int64, error)
	// delete model in database
	Delete(md interface{}, cols ...string) (int64, error)
	// load related models to md model.
	// args are limit, offset int and order string.
	//
	// example:
	// 	Ormer.LoadRelated(post,"Tags")
	// 	for _,tag := range post.Tags{...}
	//args[0] bool true useDefaultRelsDepth ; false  depth 0
	//args[0] int  loadRelationDepth
	//args[1] int limit default limit 1000
	//args[2] int offset default offset 0
	//args[3] string order  for example : "-Id"
	// make sure the relation is defined in model struct tags.
	LoadRelated(md interface{}, name string, args ...interface{}) (int64, error)

	Using(name string) error
	// begin transaction
	// for example:
	// 	o := NewOrm()
	// 	err := o.Begin()
	// 	...
	// 	err = o.Rollback()
	Begin() error
	// begin transaction with provided context and option
	// the provided context is used until the transaction is committed or rolled back.
	// if the context is canceled, the transaction will be rolled back.
	// the provided TxOptions is optional and may be nil if defaults should be used.
	// if a non-default isolation level is used that the driver doesn't support, an error will be returned.
	// for example:
	//  o := NewOrm()
	// 	err := o.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	//  ...
	//  err = o.Rollback()
	BeginTx(ctx context.Context, opts *sql.TxOptions) error
	// commit transaction
	Commit() error
	// rollback transaction
	Rollback() error
	// return a raw query seter for raw sql string.
	// for example:
	//	 ormer.Raw("UPDATE `user` SET `user_name` = ? WHERE `user_name` = ?", "slene", "testing").Exec()
	//	// update user testing's name to slene

	DBStats() *sql.DBStats
}

// base database struct
type dbBaser interface {
	//Read(dbQuerier, *modelInfo, reflect.Value, *time.Location, []string, bool) error
	Insert(dbQuerier, *modelInfo, reflect.Value, *time.Location) (int64, error)
	//InsertOrUpdate(dbQuerier, *modelInfo, reflect.Value, *alias, ...string) (int64, error)
	//InsertMulti(dbQuerier, *modelInfo, reflect.Value, int, *time.Location) (int64, error)
	InsertValue(dbQuerier, *modelInfo, bool, []string, []interface{}) (int64, error)
	//InsertStmt(stmtQuerier, *modelInfo, reflect.Value, *time.Location) (int64, error)
	//Update(dbQuerier, *modelInfo, reflect.Value, *time.Location, []string) (int64, error)
	//Delete(dbQuerier, *modelInfo, reflect.Value, *time.Location, []string) (int64, error)
	//ReadBatch(dbQuerier, *querySet, *modelInfo, *Condition, interface{}, *time.Location, []string) (int64, error)
	//SupportUpdateJoin() bool
	//UpdateBatch(dbQuerier, *querySet, *modelInfo, *Condition, Params, *time.Location) (int64, error)
	//DeleteBatch(dbQuerier, *querySet, *modelInfo, *Condition, *time.Location) (int64, error)
	//Count(dbQuerier, *querySet, *modelInfo, *Condition, *time.Location) (int64, error)
	//OperatorSQL(string) string
	//GenerateOperatorSQL(*modelInfo, *fieldInfo, string, []interface{}, *time.Location) (string, []interface{})
	//GenerateOperatorLeftCol(*fieldInfo, string, *string)
	//PrepareInsert(dbQuerier, *modelInfo) (stmtQuerier, string, error)
	//ReadValues(dbQuerier, *querySet, *modelInfo, *Condition, []string, interface{}, *time.Location) (int64, error)
	//RowsTo(dbQuerier, *querySet, *modelInfo, *Condition, interface{}, string, string, *time.Location) (int64, error)
	//MaxLimit() uint64
	//TableQuote() string
	//ReplaceMarks(*string)
	//HasReturningID(*modelInfo, *string) bool
	//TimeFromDB(*time.Time, *time.Location)
	//TimeToDB(*time.Time, *time.Location)
	//DbTypes() map[string]string
	//GetTables(dbQuerier) (map[string]bool, error)
	//GetColumns(dbQuerier, string) (map[string][3]string, error)
	//ShowTablesQuery() string
	//ShowColumnsQuery(string) string
	//IndexExists(dbQuerier, string, string) bool
	//collectFieldValue(*modelInfo, *fieldInfo, reflect.Value, bool, *time.Location) (interface{}, error)
	//setval(dbQuerier, *modelInfo, []string) error
}

// db querier
type dbQuerier interface {
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
