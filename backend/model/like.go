package model

type Like struct {
	Id			int 	`gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	UserId  	int		`gorm:"column:user_id" json:"userId"`
	SiteIndex	int		`gorm:"column:site_idx" json:"siteIndex"` 
	Valid       int8	`gorm:"column:valid" json:"valid"`
}

func (*Like) TableName() string {
	return "tb_like"
}

func (l *Like) CreateRecord() error {
	l.Valid = 1
	err := GetMySqlClient().Table(l.TableName()).Create(l).Error
	return err
}

func (l *Like) UpdateLikeStatus(isValid int8) error {
	err := GetMySqlClient().Table(l.TableName()).Where("user_id = ? and site_idx = ?", l.UserId, l.SiteIndex).Update("valid", isValid).Error
	return err
}

func (l *Like) QueryConnection() (Like, error) {
	var r Like
	err := GetMySqlClient().Table(l.TableName()).Where("user_id = ? and site_idx = ?", l.UserId, l.SiteIndex).First(&r).Error
	return r, err
}

func (l *Like) QueryLikeList() ([]Like, error) {
	var likes []Like
	err := GetMySqlClient().Table(l.TableName()).Where("user_id = ? and valid = 1", l.UserId).Find(&likes).Error
	return likes, err
}

func (l *Like) QueryLikeUserList() ([]Like, error) {
	var likes []Like
	err := GetMySqlClient().Table(l.TableName()).Where("site_idx = ? and valid = 1", l.SiteIndex).Find(&likes).Error
	return likes, err
}
