package model

type CommentPositive struct {
	Id			int32		`gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	SiteIndex 	int32		`gorm:"column:site_idx" json:"siteIndex"`
	Content		string		`gorm:"column:content" json:"content"`
}

type CommentNegative struct {
	Id			int32		`gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	SiteIndex 	int32		`gorm:"column:site_idx" json:"siteIndex"`
	Content		string		`gorm:"column:content" json:"content"`
}

func (*CommentPositive) TableName() string {
	return "tb_comments"
}

func (*CommentNegative) TableName() string {
	return "tb_comments_negative"
}

func (c *CommentPositive) QueryCommentsByPage(offset int32, limit int32) ([]CommentPositive, error) {
	var comments []CommentPositive
	err := GetMySqlClient().Table(c.TableName()).Where("site_idx = ?", c.SiteIndex).Limit(limit).Offset(offset).Find(&comments).Error
	return comments, err
}

func (c *CommentNegative) QueryCommentsByPage(offset int32, limit int32) ([]CommentNegative, error) {
	var comments []CommentNegative
	err := GetMySqlClient().Table(c.TableName()).Where("site_idx = ?", c.SiteIndex).Limit(limit).Offset(offset).Find(&comments).Error
	return comments, err
}

func (c *CommentPositive) CountPositiveComment() (int64, error) {
	var count int64
	err := GetMySqlClient().Table(c.TableName()).Where("site_idx = ?", c.SiteIndex).Count(&count).Error
	return count, err
}

func (c *CommentNegative) CountNegativeComment() (int64, error) {
	var count int64
	err := GetMySqlClient().Table(c.TableName()).Where("site_idx = ?", c.SiteIndex).Count(&count).Error
	return count, err
}
