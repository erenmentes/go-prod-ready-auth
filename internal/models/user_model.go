package models

import (
	"time"
)

type User struct {
	ID                              int       `gorm:"primaryKey;column:id"`
	Username                        string    `gorm:"type:varchar(50);unique;not null;column:username"`
	Email                           string    `gorm:"type:varchar(150);unique;not null;column:email"`
	PhoneNumber                     *string   `gorm:"type:varchar(30);unique;column:phonenumber"` // Boş (NULL) kalabileceği için pointer (*) yapıldı
	Pass                            string    `gorm:"type:text;not null;column:pass"`
	EmailVerificationCode           string    `gorm:"type:varchar(255);unique;not null;column:emailverificationcode"`
	EmailVerificationCodeExpiryDate time.Time `gorm:"type:date;not null;column:emailverificationcodeexpirydate"`
	TwoFactorVerificationCode       *string   `gorm:"type:varchar(255);unique;column:twofactorverificationcode"` // Boş (NULL) kalabileceği için pointer (*) yapıldı
}

func (User) TableName() string {
	return "users"
}
