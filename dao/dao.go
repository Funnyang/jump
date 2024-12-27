package dao

import (
	"fmt"
	"github.com/funnyang/jump/pkg/fileutil"
	"os"
	"path"

	"github.com/funnyang/jump/conf"
	"github.com/funnyang/jump/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dao struct {
	db *gorm.DB // db
}

func NewDao(c *conf.Config) *Dao {
	return &Dao{
		db: getDB(c.Database),
	}
}

func (d *Dao) Close() {
	rawDb, _ := d.db.DB()
	rawDb.Close()
}

func getDB(dbPath string) *gorm.DB {
	// 路径为空，使用默认数据库
	if dbPath == "" {
		// "~/.jump/host.db"
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic("获取用户家目录失败")
		}
		dbDir := path.Join(homeDir, ".jump")
		os.Mkdir(dbDir, 0755)
		dbPath = path.Join(dbDir, "host.db")
	}

	// 数据库是否已存在
	exist := fileutil.ExistPath(dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("连接数据库失败: %s", dbPath))
	}

	if !exist {
		// 迁移 schema
		db.AutoMigrate(&model.Host{})
	}

	return db
}

// existPath 返回路径是否存在
func existPath(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}
