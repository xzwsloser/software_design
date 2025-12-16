package model

type User struct {
	Id 			int32 		`gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	Username 	string 		`gorm:"column:username" json:"username"`
	Password 	string 		`gorm:"column:password" json:"password"`
	Gender		int32 		`gorm:"column:gender" json:"gender"`
	City 		string		`gorm:"column:city" json:"city"`
}

func (*User) TableName() string {
	return "tb_user"
}

func (u *User) QueryByUsername() (User,error) {
	var result User
	err := GetMySqlClient().Table(u.TableName()).Where("username = ?", u.Username).First(&result).Error
	return result, err
}

func (u *User) InsertUser() (error) {
	err := GetMySqlClient().Table(u.TableName()).Create(u).Error
	return err
}

func (u *User) UpdateByUserName() (error) {
	err := GetMySqlClient().Table(u.TableName()).Where("username = ?", u.Username).Select("password", "gender", "city").Updates(u).Error
	return err
}

