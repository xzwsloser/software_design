package model

// type User struct {
// 	Id 			int32 		`gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
// 	Username 	string 		`gorm:"column:username" json:"username"`
// 	Password 	string 		`gorm:"column:password" json:"password"`
// 	Gender		int32 		`gorm:"column:gender" json:"gender"`
// 	City 		string		`gorm:"column:city" json:"city"`
// }

type User struct {
	Id				int32		`gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	Username		string      `gorm:"column:username" json:"username"`
	Password		string	    `gorm:"column:password" json:"password"`
	Gender			int		    `gorm:"column:gender" json:"gender"`
	AddressId		int		    `gorm:"column:address_id" json:"addressId"`
	TouristType		int		    `gorm:"column:tourist_type" json:"touristType"`
	LikeType		string	    `gorm:"column:like_type" json:"likeType"`
	Targets			string	    `gorm:"column:targets" json:"targets"`
	PriceSensitive	int		    `gorm:"column:price_sensitive" json:"priceSensitive"`
	Attention		string	    `gorm:"column:attention" json:"attention"`
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

// func (u *User) UpdateByUserName() (error) {
// 	err := GetMySqlClient().Table(u.TableName()).Where("username = ?", u.Username).Select("password", "gender", "city").Updates(u).Error
// 	return err
// }

func (u *User) UpdateUserInfo() error {
    err := GetMySqlClient().Table(u.TableName()).
        Where("username = ?", u.Username).
        Select("address_id", "tourist_type", "like_type", "targets", "price_sensitive", "attention").
        Updates(u).Error
    return err
}
