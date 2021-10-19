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

	"github.com/caser789/go-xorm/caches"
	"github.com/caser789/go-xorm/core"
	_ "github.com/caser789/go-xorm/drivers"
)

const (
	Version string = "0.4"
)

func init() {
	provided_dialects := map[string]struct {
		dbType core.DbType
		get    func() core.Dialect
	}{
		"odbc":     {"mssql", func() core.Dialect { return &mssql{} }},
		"mysql":    {"mysql", func() core.Dialect { return &mysql{} }},
		"mymysql":  {"mysql", func() core.Dialect { return &mysql{} }},
		"oci8":     {"oracle", func() core.Dialect { return &oracle{} }},
		"postgres": {"postgres", func() core.Dialect { return &postgres{} }},
		"sqlite3":  {"sqlite3", func() core.Dialect { return &sqlite3{} }},
	}

	for k, v := range provided_dialects {
		_, err := sql.Open(string(k), "")
		if err == nil {
			core.RegisterDialect(v.dbType, v.get())
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

func NewLRUCacher(store core.CacheStore, max int) *caches.LRUCacher {
	return caches.NewLRUCacher(store, core.CacheExpired, core.CacheMaxMemory, max)
}

func NewLRUCacher2(store core.CacheStore, expired time.Duration, max int) *caches.LRUCacher {
	return caches.NewLRUCacher(store, expired, 0, max)
}

func NewMemoryStore() *caches.MemoryStore {
	return caches.NewMemoryStore()
}
