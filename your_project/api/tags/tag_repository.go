package tags

import (
	"errors"
	"labs/constants"
	"labs/domains"
	"regexp"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

// NewTagRepository performs automatic migration of tag-related structures in the database.
func NewTagRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Tag{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the tag structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of tags based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Tag, modelID uuid.UUID, limit, offset int) ([]domains.Tag, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadByID retrieves a company by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Tag, id uuid.UUID) (domains.Tag, error) {
	err := db.First(&model, id).Error
	return model, err
}

// ReadAllPagination  list of tags based on company ID.
func ReadAllTags(db *gorm.DB, model []domains.Tag, modelID uuid.UUID) ([]domains.Tag, error) {
	err := db.Where("company_id = ? ", modelID).Find(&model).Error
	return model, err
}
func AssignToMailinglist(db *gorm.DB, modelID uuid.UUID, mailinglistID uuid.UUID) error {

	var tagCount int64
	err := db.Model(&domains.Mailinglist{}).Where("id = ?", mailinglistID).Where("tags @> ARRAY[?]::uuid[]", modelID).Count(&tagCount).Error
	if err != nil {
		return err
	}

	if tagCount > 0 {
		// Tag already exists, no need to append
		logrus.Error("Tag already exists in the mailinglist")
		return errors.New("tag already exists in this mailing list")
	}

	if err := db.Exec("UPDATE mailinglists SET tags = array_append(tags, ?) WHERE id = ?", modelID, mailinglistID); err != nil {
		logrus.Error("An error occurred during updating mailinglist. Error: ", err)
	}

	return nil

	// mailingList := domains.Mailinglist{}
	// result := db.First(&mailingList, mailinglistID)
	// if result.Error != nil {
	// 	return result.Error
	// }
	// fmt.Println("hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh", mailingList.Tags)

	// for _, v := range mailingList.Tags {
	// 	if v == modelID {
	// 		return errors.New("tag already exists in this mailing list")
	// 	}
	// }

	// mailingList.Tags = append(mailingList.Tags, modelID)
	// if err := db.Save(&mailingList).Error; err != nil {
	// 	return err
	// }

	// return nil

}

func AssignToContact(db *gorm.DB, modelID uuid.UUID, ContactID uuid.UUID) error {

	var tagCount int64
	err := db.Model(&domains.Contact{}).Where("id = ?", ContactID).Where("tags @> ARRAY[?]::uuid[]", modelID).Count(&tagCount).Error
	if err != nil {
		return err
	}

	if tagCount > 0 {
		// Tag already exists, no need to append
		logrus.Error("Tag already exists in the contact")
		return errors.New("tag already exists in this mailing list")
	}

	if err := db.Exec("UPDATE contacts SET tags = array_append(tags, ?) WHERE id = ?", modelID, ContactID); err != nil {
		logrus.Error("An error occurred during updating contact. Error: ", err)
	}

	return nil

}
func Validate_color(tag *TagIn) error {
	re, err := regexp.Compile(constants.COLOR_REGEX)
	if err != nil {
		logrus.Error("Invalide format of color # +6 . Error: ", err.Error())
		return err
	}
	if !re.MatchString(tag.Color) {
		return errors.New("invalid color format") // Specific error for validation failure
	}
	return nil

}
