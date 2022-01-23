package dao

import (
	"testing"

	"github.com/funnyang/jump/conf"
	"github.com/funnyang/jump/model"
)

func TestGetDB(t *testing.T) {
	getDB("")
}

func TestInsertHost(t *testing.T) {
	var host = model.Host{
		IP:       "127.0.0.1",
		Port:     22,
		User:     "zz",
		Password: "123456",
		Desc:     "test",
	}
	d := NewDao(&conf.Config{})
	t.Log(d.InsertHost(host))
}
