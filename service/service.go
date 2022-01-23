package service

import (
	"github.com/funnyang/jump/conf"
	"github.com/funnyang/jump/dao"
)

type Service struct {
	d *dao.Dao
}

func NewService(c *conf.Config) *Service {
	return &Service{
		d: dao.NewDao(c),
	}
}

func (s *Service) Close() {
	s.d.Close()
}
