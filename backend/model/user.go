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
