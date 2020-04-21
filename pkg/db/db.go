package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DBInfo struct {
	Dialect  string
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

func OpenDB(i DBInfo) (db *gorm.DB, err error) {
	if i.Dialect == "postgres" {
		db, err = gorm.Open(i.Dialect, "host="+i.Host+" port="+i.Port+" user="+i.User+" dbname="+i.DBName+" password="+i.Password)
	} else if i.Dialect == "sqlite3" {
		db, err = gorm.Open(i.Dialect, "test.db")
	}
	return db, err
}

//func InitDB(i DBInfo, m interface{}) error {
//	db, err := OpenDB(i)
//	if err != nil {
//		return err
//	}
//	defer db.Close()
//
//	// Migrate the schema
//	db.AutoMigrate(&m)
//	return nil
//}

//func CreateRecord(i DBInfo, m interface{}, t string) error {
//	db, err := OpenDB(i)
//	if err != nil {
//		return err
//	}
//	defer db.Close()
//
//	db.Table(t).Create(&m)
//	log.Printf("[INFO] Create Record: %v", m)
//	return nil
//}

//func GetRecord(i DBInfo, q []string) (interface{}, error) {
//	db, err := OpenDB(i)
//	if err != nil {
//		return nil, err
//	}
//	defer db.Close()
//
//	log.Printf("[INFO] Query: %s", q)
//	var m []interface{}
//	if q[1] == "" {
//		db.Find(&m)
//	} else {
//		db.Where(q).First(&m)
//	}
//
//	return m, nil
//}
