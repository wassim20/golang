package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	Email              string `gorm:"unique"`
	Password           string
	AccountType        string
	SubscriptionOption string

	MailingLists []MailingList
	Tags         []Tag
	Users        []User
}
