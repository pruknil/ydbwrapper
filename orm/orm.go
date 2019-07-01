package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sync"
)

// field info collection
type fields struct {
	pk            *fieldInfo
	columns       map[string]*fieldInfo
	fields        map[string]*fieldInfo
	fieldsLow     map[string]*fieldInfo
	fieldsByType  map[int][]*fieldInfo
	fieldsRel     []*fieldInfo
	fieldsReverse []*fieldInfo
	fieldsDB      []*fieldInfo
	rels          []*fieldInfo
	orders        []string
	dbcols        []string
}

// single field info
type fieldInfo struct {
	mi         *modelInfo
	fieldIndex []int
	fieldType  int
	dbcol      bool // table column fk and onetoone
	inModel    bool
	name       string
	fullName   string
	column     string
	addrValue  reflect.Value
	sf         reflect.StructField
	auto       bool
	pk         bool
	null       bool
	index      bool
	unique     bool
	colDefault bool // whether has default tag
	//initial             StrTo // store the default value
	size                int
	toText              bool
	autoNow             bool
	autoNowAdd          bool
	rel                 bool // if type equal to RelForeignKey, RelOneToOne, RelManyToMany then true
	reverse             bool
	reverseField        string
	reverseFieldInfo    *fieldInfo
	reverseFieldInfoTwo *fieldInfo
	reverseFieldInfoM2M *fieldInfo
	relTable            string
	relThrough          string
	relThroughModelInfo *modelInfo
	relModelInfo        *modelInfo
	digits              int
	decimals            int
	isFielder           bool // implement Fielder interface
	onDelete            string
	description         string
}

// single model info
type modelInfo struct {
	pkg       string
	name      string
	fullName  string
	table     string
	model     interface{}
	fields    *fields
	manual    bool
	addrField reflect.Value //store the original struct value
	uniques   []string
	isThrough bool
}

// model info collection
type _modelCache struct {
	sync.RWMutex    // only used outsite for bootStrap
	orders          []string
	cache           map[string]*modelInfo
	cacheByFullName map[string]*modelInfo
	done            bool
}

var (
	modelCache = &_modelCache{
		cache:           make(map[string]*modelInfo),
		cacheByFullName: make(map[string]*modelInfo),
	}
)

// RegisterModel register models
func RegisterModel(models ...interface{}) {
	if modelCache.done {
		panic(fmt.Errorf("RegisterModel must be run before BootStrap"))
	}
	RegisterModelWithPrefix("", models...)
}

// RegisterModelWithPrefix register models with a prefix
func RegisterModelWithPrefix(prefix string, models ...interface{}) {
	if modelCache.done {
		panic(fmt.Errorf("RegisterModelWithPrefix must be run before BootStrap"))
	}

	//for _, model := range models {
	//	registerModel(prefix, model, true)
	//}
}

type orm struct {
}

func (orm) Read(md interface{}, cols ...string) error {
	panic("implement me")
}

func (orm) ReadForUpdate(md interface{}, cols ...string) error {
	panic("implement me")
}

func (orm) ReadOrCreate(md interface{}, col1 string, cols ...string) (bool, int64, error) {
	panic("implement me")
}

func (orm) Insert(interface{}) (int64, error) {

	return 0, nil
//	panic("implement me")
}

func (orm) InsertOrUpdate(md interface{}, colConflitAndArgs ...string) (int64, error) {
	panic("implement me")
}

func (orm) InsertMulti(bulk int, mds interface{}) (int64, error) {
	panic("implement me")
}

func (orm) Update(md interface{}, cols ...string) (int64, error) {
	panic("implement me")
}

func (orm) Delete(md interface{}, cols ...string) (int64, error) {
	panic("implement me")
}

func (orm) LoadRelated(md interface{}, name string, args ...interface{}) (int64, error) {
	panic("implement me")
}

func (orm) Using(name string) error {
	//	panic("implement me")
	return nil
}

func (orm) Begin() error {
	panic("implement me")
}

func (orm) BeginTx(ctx context.Context, opts *sql.TxOptions) error {
	panic("implement me")
}

func (orm) Commit() error {
	panic("implement me")
}

func (orm) Rollback() error {
	panic("implement me")
}

func (orm) DBStats() *sql.DBStats {
	panic("implement me")
}

func BootStrap() {
	modelCache.Lock()
	defer modelCache.Unlock()
	if modelCache.done {
		return
	}
	//bootStrap()
	modelCache.done = true
}
func NewOrm() Ormer {
	BootStrap() // execute only once

	o := new(orm)

	return o
}

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
