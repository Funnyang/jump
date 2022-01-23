package service

import (
	"time"

	"github.com/funnyang/jump/model"
)

func (s *Service) ListHostByKeyword(keyword string) ([]model.Host, error) {
	if keyword == "" {
		return s.ListHost()
	}
	return s.d.MatchHost(keyword)
}

func (s *Service) ExactMatchHost(keyword string) ([]model.Host, error) {
	return s.d.ExactMatchHost(keyword)
}

func (s *Service) ListHost() ([]model.Host, error) {
	return s.d.ListHost()
}

func (s *Service) InsertHost(host model.Host) error {
	host.CreateTime = time.Now()
	host.UpdateTime = time.Now()
	return s.d.InsertHost(host)
}

func (s *Service) DeleteHost(id int) error {
	return s.d.DeleteHost(id)
}

func (s *Service) UpdateHost(host model.Host) error {
	host.UpdateTime = time.Now()
	return s.d.UpdateHost(host)
}

func (s *Service) GetHost(id int) (model.Host, error) {
	return s.d.GetHost(id)
}
