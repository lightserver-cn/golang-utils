package datasource

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

const (
	// mysql 连接
	mysqlUrl = "%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local"
	// PostGreSQL 连接
	postGreSqlUrl = "postgres://%v:%v@%v/%v?sslmode=disable"
)

var db *gorm.DB

// ConnectionOptions
type ConnectionOptions struct {
	Adapter       string        `yaml:"adapter"`        // 适配器类型 mysql/postgres
	Username      string        `yaml:"username"`       // 用户名
	Password      string        `yaml:"password"`       // 密码
	Host          string        `yaml:"host"`           // 地址
	Port          int64         `yaml:"port"`           // 端口
	Database      string        `yaml:"database"`       // 数据库
	MaxIdleConns  int           `yaml:"max_idle_conns"` // 设置连接池中的最大闲置连接数
	MaxOpenConns  int           `yaml:"max_open_conns"` // 设置数据库的最大连接数量
	MaxLifetime   time.Duration `yaml:"max_lifetime"`   // 设置连接的最大可复用时间 int64 ns * 1000 * 1000 = s
	SingularTable bool          `yaml:"singular_table"` // 表生成结尾不带s
	Prefix        string        `yaml:"prefix"`         // 表前缀
	EnableLog     bool          `yaml:"enable_log"`     // 启用Logger，显示详细日志
}

// NewDB
func NewDB(opts *ConnectionOptions) *gorm.DB {
	var err error
	var url string

	switch opts.Adapter {
	default:
		panic("not supported database adapter")
	case "mysql":
		url = fmt.Sprintf(mysqlUrl, opts.Username, opts.Password, opts.Host, opts.Port, opts.Database)
	case "postgres":
		url = fmt.Sprintf(postGreSqlUrl, opts.Username, opts.Password, opts.Host, opts.Database)
	}

	// 设置默认表名前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return opts.Prefix + defaultTableName
	}

	db, err = gorm.Open(opts.Adapter, url)
	if err != nil {
		logrus.Errorf("Opens DB failed: %s", err.Error())
	}

	db.LogMode(opts.EnableLog)
	db.SingularTable(opts.SingularTable)
	db.DB().SetMaxIdleConns(opts.MaxIdleConns)
	db.DB().SetMaxOpenConns(opts.MaxOpenConns)
	db.DB().SetConnMaxLifetime(opts.MaxLifetime)

	fmt.Println("Connected to DB!")

	return db
}

// CloseDB
func CloseDB() {
	if db == nil {
		return
	}
	if err := db.Close(); nil != err {
		logrus.Errorf("Disconnect from database failed: %s", err.Error())
	}
}
