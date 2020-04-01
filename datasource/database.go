package datasource

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

const (
	mysqlUrl    = "%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local"
	postgresUrl = "postgres://%v:%v@%v/%v?sslmode=disable"
)

type DB struct {
	db *gorm.DB
}

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
	DropTable     bool          `yaml:"drop_table"`     // 初始化表结构
}

func NewDB(opts *ConnectionOptions) *gorm.DB {

	var url string
	switch opts.Adapter {
	default:
		panic("not supported database adapter")
	case "mysql":
		url = fmt.Sprintf(mysqlUrl, opts.Username, opts.Password, opts.Host, opts.Port, opts.Database)
	case "postgres":
		url = fmt.Sprintf(postgresUrl, opts.Username, opts.Password, opts.Host, opts.Database)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return opts.Prefix + defaultTableName
	}

	db, err := gorm.Open(opts.Adapter, url)
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

// drop tables
func (DB *DB) DropTable(Models []interface{}) {
	fmt.Println("Drop tables from DB.")
	DB.db.DropTableIfExists(Models...)
}

// create tables
func (DB *DB) Migrate(Models []interface{}) {
	// create tables
	if err := DB.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(Models...).Error; err != nil {
		logrus.Errorf("Auto migrate tables failed: %s", err.Error())
	}

	fmt.Println("Create tables from DB.")
}
