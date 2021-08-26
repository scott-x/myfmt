package db

import (
	"github.com/scott-x/myfmt/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path"
	"time"
)

var (
	err       error
	DB        *gorm.DB
	color_map map[string]string
)

type GoFile struct {
	Pth string `gorm:"column:pth;unique"`
	Md5 string `gorm:"column:md5"`
}

var dbname = path.Join(os.Getenv("GOPATH"), "src/github.com/scott-x/myfmt/myfmt.db")

func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,  // Ignore ErrRecordNotFound error for logger
			Colorful:                  false, // Disable color
		},
	)
	DB, err = gorm.Open(sqlite.Open(dbname), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
}

func IsMd5Same(file string) bool {
	md5_db, err := getMd5FromFile(file) //md5 in db
	if err != nil || len(md5_db) == 0 {
		return false
	}
	md5 := util.MD5(file) //current file md5
	return md5_db == md5
}

func getMd5FromFile(file string) (string, error) {
	var goFile GoFile
	err := DB.Where("pth=?", file).First(&goFile).Error
	if err != nil {
		return "", err
	}
	return goFile.Md5, nil
}

func Record(pth string) error {
	md5 := util.MD5(pth)
	gofile := GoFile{
		Pth: pth,
		Md5: md5,
	}
	return DB.Create(&gofile).Error
}

func IsItemExist(pth string) bool {
	var gofile GoFile
	err := DB.Where("pth=?", pth).First(&gofile).Error
	if err != nil {
		return false
	}
	return true
}

func UpdateRecordViaPth(pth string) error {
	return DB.Model(&GoFile{}).Where("pth=?", pth).Update("md5", util.MD5(pth)).Error
}
