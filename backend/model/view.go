package model

type View struct {
	Id 			int 	`gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	UserId  	int		`gorm:"column:user_id" json:"userId"`
	SiteIndex	int		`gorm:"column:site_idx" json:"siteIndex"`
}

func (*View) TableName() string {
	return "tb_view"
}

func (v *View) CreateViewRecord() error {
	err := GetMySqlClient().Table(v.TableName()).Create(v).Error
	return err
}

func (v *View) QueryViewConnection() (View, error) {
	var result View
	err := GetMySqlClient().Table(v.TableName()).Where("user_id = ? and site_idx = ?",  v.UserId, v.SiteIndex).First(&result).Error
	return result, err
}

func (v *View) QueryViewList() ([]View, error) {
	var views []View
	err := GetMySqlClient().Table(v.TableName()).Where("user_id = ?", v.UserId).Find(&views).Error
	return views, err
}

func (v *View) QueryViewUserList() ([]View, error) {
	var views []View
	err := GetMySqlClient().Table(v.TableName()).Where("site_idx = ?", v.SiteIndex).Find(&views).Error
	return views, err
}





