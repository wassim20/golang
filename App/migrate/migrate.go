package main

import (
	"github.com/wassim_p/App/initializers"
	"github.com/wassim_p/App/models"
)

func init() {
	initializers.LoadiEnv()
	initializers.Connect()
}
func main() {

	initializers.DB.AutoMigrate(&models.Company{},
		&models.Role{},
		&models.MailingList{},
		&models.Contact{},
		&models.Tag{},
		&models.User{},
		&models.Notification{},
	)
}
