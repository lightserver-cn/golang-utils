package datasource

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type DB struct {
	db *gorm.DB
}

type ConnectionOptions struct {
	Adapter      string        `yaml:"adapter"`
	Username     string        `yaml:"username"`
	Password     string        `yaml:"password"`
	Host         string        `yaml:"host"`
	Port         int64         `yaml:"port"`
	Database     string        `yaml:"database"`
	MaxIdleConns int           `yaml:"max_idle_conns"`
	MaxOpenConns int           `yaml:"max_open_conns"`
	MaxLifetime  time.Duration `yaml:"max_lifetime"`
	Prefix       string        `yaml:"prefix"`
	EnableLog    bool          `yaml:"enable_log"`
	DropTable    bool          `yaml:"drop_table"`
}

func NewDB(opts *ConnectionOptions) *gorm.DB {

	var url string
	switch opts.Adapter {
	default:
		panic("not supported database adapter")
	case "mysql":
		url = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", opts.Username, opts.Password, opts.Host, opts.Port, opts.Database)
	case "postgres":
		url = fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", opts.Username, opts.Password, opts.Host, opts.Database)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return opts.Prefix + defaultTableName
	}

	db, err := gorm.Open(opts.Adapter, url)
	if err != nil {
		logrus.Errorf("Opens DB failed: %s", err.Error())
	}

	db.LogMode(opts.EnableLog)
	db.SingularTable(true)
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
