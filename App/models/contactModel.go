package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	Email       string
	FirstName   string
	LastName    string
	PhoneNumber string
	FullName    string

	MailingLists []MailingList `gorm:"many2many:mailing_list_contacts"`
	Tags         []Tag         `gorm:"many2many:contact_tags"`
}
