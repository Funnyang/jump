package dao

import "github.com/funnyang/jump/model"

// MatchHost 按关键词搜索host
func (d *Dao) MatchHost(keyword string) (hosts []model.Host, err error) {
	keyword = "%" + keyword + "%"
	err = d.db.Where("id like ?", keyword).
		Or("hostname like ?", keyword).
		Or("ip like ?", keyword).
		Or("Desc like ?", keyword).Find(&hosts).Error
	return
}

// ExactMatchHost 基于关键字精确匹配
func (d *Dao) ExactMatchHost(keyword string) (hosts []model.Host, err error) {
	err = d.db.Where("id = ?", keyword).
		Or("hostname = ?", keyword).
		Or("ip = ?", keyword).Find(&hosts).Error
	return
}

func (d *Dao) ListHost() (hosts []model.Host, err error) {
	err = d.db.Find(&hosts).Error
	return
}

// InsertHost 添加host
func (d *Dao) InsertHost(host model.Host) error {
	return d.db.Create(&host).Error
}

// DeleteHost 删除host
func (d *Dao) DeleteHost(id int) error {
	return d.db.Delete(&model.Host{}, id).Error
}

// UpdateHost 更新host
func (d *Dao) UpdateHost(host model.Host) error {
	return d.db.Model(&host).Updates(host).Error
}

// GetHost 查询host
func (d *Dao) GetHost(id int) (host model.Host, err error) {
	err = d.db.Where("id = ?", id).First(&host).Error
	return
}
