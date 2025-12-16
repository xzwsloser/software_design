package model

type Site struct {
	Id        int32   `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	Name      string  `gorm:"column:name" json:"name"`
	Score     float32 `gorm:"column:score" json:"score"`
	Address   string  `gorm:"column:address" json:"address"`
	HotDegree float32 `gorm:"column:hot_degree" json:"hotDegree"`
	Introduce string  `gorm:"column:introduce" json:"introduce"`
	OpenTime  string  `gorm:"column:open_time" json:"openTime"`
	Phone     string  `gorm:"column:phone" json:"phone"`
	Images    string  `gorm:"column:images" json:"images"`
	SiteIndex int32   `gorm:"column:site_idx" json:"siteIndex"`
}

func (*Site) TableName() string {
	return "tb_sites"
}

func (s *Site) QueryByPage(offset int32, limit int32) ([]Site, error) {
	var sites []Site
	err := GetMySqlClient().Table(s.TableName()).Limit(limit).Offset(offset).Find(&sites).Error
	return sites, err
}

func (s *Site) QueryByIndex() (Site, error) {
	var result Site
	err := GetMySqlClient().Table(s.TableName()).Where("site_idx = ?", s.SiteIndex).First(&result).Error
	return result, err
}

func (s *Site) QueryBySiteIndexes(siteIndexList []int32) ([]Site, error) {
	var sites []Site
	err := GetMySqlClient().Table(s.TableName()).Where("`site_idx` IN (?)", siteIndexList).Find(&sites).Error
	return sites, err
}



