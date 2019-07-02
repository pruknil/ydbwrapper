package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sync"
	"time"
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

	for _, model := range models {
		registerModel(prefix, model, true)
	}
}

func registerModel(s string, model interface{}, b bool) {

}

type orm struct {
	alias *alias
	db    dbQuerier
	isTx  bool
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

func (o *orm) Insert(md interface{}) (int64, error) {
	mi := &modelInfo{}
	val := reflect.ValueOf(md)
	ind := reflect.Indirect(val)
	//mi, ind := o.getMiInd(md, true)
	id, err := o.alias.DbBaser.Insert(o.db, mi, ind, o.alias.TZ)
	if err != nil {
		return id, err
	}

	//	o.setPk(mi, ind, id)

	return id, nil
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

// create new mysql dbBaser.
func newdbBaseYotta() dbBaser {
	b := &dbBase{}
	b.ins = b
	return b
}

func NewOrm() Ormer {
	BootStrap() // execute only once

	o := &orm{}

	o.alias = &alias{
		DbBaser: newdbBaseYotta(),
	}
	return o
}

type alias struct {
	Name string
	//Driver       DriverType
	DriverName   string
	DataSource   string
	MaxIdleConns int
	MaxOpenConns int
	DB           *sql.DB
	DbBaser      dbBaser
	TZ           *time.Location
	Engine       string
}

// an instance of dbBaser interface/
type dbBase struct {
	ins dbBaser
}

var _ dbBaser = new(dbBase)

// execute insert sql dbQuerier with given struct reflect.Value.
func (d *dbBase) Insert(q dbQuerier, mi *modelInfo, ind reflect.Value, tz *time.Location) (int64, error) {
	//names := make([]string, 0, len(mi.fields.dbcols))
	//values, _, err := d.collectValues(mi, ind, mi.fields.dbcols, false, true, &names, tz)
	//if err != nil {
	//	return 0, err
	//}

	id, err := d.InsertValue(q, mi, false, []string{"name"}, []interface{}{"value"})
	if err != nil {
		return 0, err
	}

	return id, err
}

// get struct columns values as interface slice.
func (d *dbBase) collectValues(mi *modelInfo, ind reflect.Value, cols []string, skipAuto bool, insert bool, names *[]string, tz *time.Location) (values []interface{}, autoFields []string, err error) {
	if names == nil {
		ns := make([]string, 0, len(cols))
		names = &ns
	}
	values = make([]interface{}, 0, len(cols))

	return
}

// get one field value in struct column as interface.
func (d *dbBase) collectFieldValue(mi *modelInfo, fi *fieldInfo, ind reflect.Value, insert bool, tz *time.Location) (interface{}, error) {
	var value interface{}

	return value, nil
}

// execute insert sql with given struct and given values.
// insert the given values, not the field values in struct.
func (d *dbBase) InsertValue(q dbQuerier, mi *modelInfo, isMulti bool, names []string, values []interface{}) (int64, error) {

	var id int64
	return id, nil
}
