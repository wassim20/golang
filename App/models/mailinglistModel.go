package models

import "gorm.io/gorm"

type MailingList struct {
	gorm.Model
	Name        string
	Description string
	CompanyID   uint

	Contacts []Contact `gorm:"many2many:mailing_list_contacts"`
	Tags     []Tag     `gorm:"many2many:mailing_list_tags"`
}
