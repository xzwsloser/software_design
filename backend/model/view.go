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


