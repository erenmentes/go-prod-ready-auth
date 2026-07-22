package models

import (
	"time"
)

type User struct {
	ID                               int       `gorm:"primaryKey;autoIncrement;column:id"`
	Username                         string    `gorm:"type:varchar(50);unique;not null;column:username"`
	Email                            string    `gorm:"type:varchar(150);unique;not null;column:email"`
	UserRole                         string    `gorm:"type:varchar(25);not null;default:'User';column:userrole"`
	PhoneNumber                      *string   `gorm:"type:varchar(30);unique;column:phonenumber"`
	Pass                             string    `gorm:"type:text;not null;column:pass"`
	RefreshToken                     string    `gorm:"type:varchar(255);unique;not null;column:refreshtoken"`
	IsEmailVerified                  *bool     `gorm:"type:bool;not null;default:false"`
	IsTwoFactorVerificationActivated bool      `gorm:"type:bool;not null;default:false;column:istwofactorverificationactivated"`
	EmailVerificationCode            string    `gorm:"type:varchar(255);unique;not null;column:emailverificationcode"`
	EmailVerificationCodeExpiryDate  time.Time `gorm:"type:timestamptz;not null;column:emailverificationcodeexpirydate"`
	TwoFactorVerificationCode        *string   `gorm:"type:varchar(255);unique;column:twofactorverificationcode"`
}

func (User) TableName() string {
	return "users"
}
