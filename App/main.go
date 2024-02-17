package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wassim_p/App/api/company"
	"github.com/wassim_p/App/api/contact"
	"github.com/wassim_p/App/api/mailinglist"
	"github.com/wassim_p/App/api/notification"
	"github.com/wassim_p/App/api/role"
	"github.com/wassim_p/App/api/tag"
	"github.com/wassim_p/App/api/user"
	"github.com/wassim_p/App/initializers"
)

func init() {
	initializers.LoadiEnv()
	initializers.Connect()

}

func main() {
	initializers.Connect()

	// Create Gin router
	router := gin.Default()
	company.RegisterRoutes(router)
	mailinglist.RegisterRoutes(router)
	contact.RegisterRoutes(router)
	tag.RegisterRoutes(router)
	notification.RegisterRoutes(router)
	user.RegisterRoutes(router)
	role.RegisterRoutes(router)

	// //COMPANY
	// r.POST("company", controllers.CreateCompany)
	// r.GET("company", controllers.AllCompany)
	// r.GET("company/:id", controllers.GetCompany)
	// r.PUT("company/:id", controllers.UpdateCompany)
	// //advanced commands
	// r.POST("company/:companyID/mailinglists/:mailingListID/assign", controllers.AssignMailingListToCompany)

	// //CONTACT
	// r.POST("contact", controllers.CreateContact)
	// r.GET("contact", controllers.AllContacts)
	// r.GET("contact/:id", controllers.GetContact)
	// r.PUT("contact/:id", controllers.UpdateContact)
	// // advanced commands
	// r.POST("contact/:contactID/tags/:tagID/assign", controllers.AssignTagToContact)

	// //TAGS
	// r.POST("tag", controllers.CreateTag)
	// r.GET("tag", controllers.AllTags)
	// r.GET("tag/:id", controllers.GetTag)
	// r.PUT("tag/:id", controllers.UpdateTag)

	// //MAILING_LIST
	// r.POST("mailing_list", controllers.CreateMailingList)
	// r.GET("mailing_list", controllers.AllMailingLists)
	// r.GET("mailing_list/:id", controllers.GetMailingList)
	// r.PUT("mailing_list/:id", controllers.UpdateMailingList)
	// //advanced commands
	// r.POST("mailing_list/:mailingListID/contacts/:contactID/assign", controllers.AssignContactToMailingList)
	// r.POST("mailing_list/:mailingListID/tags/:tagID/assign", controllers.AssignTagToMailingList)

	router.Run() // listen and serve on 0.0.0.0:8080
}
