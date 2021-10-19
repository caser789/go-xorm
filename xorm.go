package xorm

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/caser789/go-xorm/core"
)

const (
	Version string = "0.4"
)

func init() {

	providedDrvsNDialects := map[string]struct {
		dbType     core.DbType
		getDriver  func() core.Driver
		getDialect func() core.Dialect
	}{
		"odbc":     {"mssql", func() core.Driver { return &odbcDriver{} }, func() core.Dialect { return &mssql{} }}, // !nashtsai! TODO change this when supporting MS Access
		"mysql":    {"mysql", func() core.Driver { return &mysqlDriver{} }, func() core.Dialect { return &mysql{} }},
		"mymysql":  {"mysql", func() core.Driver { return &mymysqlDriver{} }, func() core.Dialect { return &mysql{} }},
		"postgres": {"postgres", func() core.Driver { return &pqDriver{} }, func() core.Dialect { return &postgres{} }},
		"sqlite3":  {"sqlite3", func() core.Driver { return &sqlite3Driver{} }, func() core.Dialect { return &sqlite3{} }},
		"oci8":     {"oracle", func() core.Driver { return &oci8Driver{} }, func() core.Dialect { return &oracle{} }},
		"goracle":  {"oracle", func() core.Driver { return &goracleDriver{} }, func() core.Dialect { return &oracle{} }},
	}

	for driverName, v := range providedDrvsNDialects {
		_, err := sql.Open(driverName, "")
		if err == nil {
			core.RegisterDriver(driverName, v.getDriver())
			core.RegisterDialect(v.dbType, v.getDialect())
		}
	}
}

func close(engine *Engine) {
	engine.Close()
}

// new a db manager according to the parameter. Currently support four
// drivers
func NewEngine(driverName string, dataSourceName string) (*Engine, error) {
	driver := core.QueryDriver(driverName)
	if driver == nil {
		return nil, errors.New(fmt.Sprintf("Unsupported driver name: %v", driverName))
	}

	uri, err := driver.Parse(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	dialect := core.QueryDialect(uri.DbType)
	if dialect == nil {
		return nil, errors.New(fmt.Sprintf("Unsupported dialect type: %v", uri.DbType))
	}

	err = dialect.Init(uri, driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	engine := &Engine{
		DriverName:     driverName,
		DataSourceName: dataSourceName,
		dialect:        dialect,
	}

	engine.SetMapper(core.NewCacheMapper(new(core.SnakeMapper)))

	engine.Filters = dialect.Filters()

	engine.Tables = make(map[reflect.Type]*core.Table)

	engine.mutex = &sync.RWMutex{}
	engine.TagIdentifier = "xorm"

	engine.Logger = NewSimpleLogger(os.Stdout)

	//engine.Pool = NewSimpleConnectPool()
	//engine.Pool = NewNoneConnectPool()
	//engine.Cacher = NewLRUCacher()
	err = engine.SetPool(NewSysConnectPool())
	runtime.SetFinalizer(engine, close)
	return engine, err
}

// func NewLRUCacher(store core.CacheStore, max int) *LRUCacher {
// 	return NewLRUCacher(store, core.CacheExpired, core.CacheMaxMemory, max)
// }

func NewLRUCacher2(store core.CacheStore, expired time.Duration, max int) *LRUCacher {
	return NewLRUCacher(store, expired, 0, max)
}

// func NewMemoryStore() *MemoryStore {
// 	return NewMemoryStore()
// }
