package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
)

type Venue struct {
	gorm.Model
	Name 			string	`gorm:"size:100;not null;unique"	json:"name"`
	Description 	string	`gorm:"not null"					json:"description"`
	Location 		string	`gorm:"size:100;not null"			json:"location"`
	Capacity 		int		`gorm:"not null"					json:"capacity"`
	Category 		string	`gorm:"size:100;not null"			json:"category"'`
	CreatedBy 		User	`gorm:"foreignKey:UserID;"			json:"-"`
	UserID 			uint	`gorm:"not null"					json:"user_id"`
}

func (venue *Venue) Prepare() {
	venue.Name = strings.TrimSpace(venue.Name)
	venue.Description = strings.TrimSpace(venue.Description)
	venue.Location = strings.TrimSpace(venue.Location)
	venue.Category = strings.TrimSpace(venue.Category)
	venue.CreatedBy = User{}
}

func (venue *Venue) Validate() error {
	if venue.Name == "" {
		return errors.New("Name is required")
	}
	if venue.Description == "" {
		return errors.New("Description of venue is required")
	}
	if venue.Location == "" {

		return errors.New("Location of venue is required")
	}
	if venue.Category == "" {
		return errors.New("Category of venue is required")
	}
	if venue.Capacity == 0 {
		return errors.New("Capacity of venue is invalid")
	}
	return nil
}

func (venue *Venue) Save(db *gorm.DB) (*Venue, error) {
	var err error

	// Debug a single operation, show detailed log for this operation
	err = db.Debug().Create(&venue).Error
	if err != nil {
		return &Venue{}, err
	}
	return venue, nil
}

func (venue *Venue) GetVenue(db *gorm.DB) (*Venue, error) {
	rtnVenue := &Venue{}
	if err := db.Debug().Table("venues").Where("name = ?",
		venue.Name).First(venue).Error; err != nil {
		return nil, err
	}
	return rtnVenue, nil
}

func GetVenues(db *gorm.DB) (*[]Venue, error) {
	venues := []Venue{}
	if err := db.Debug().Table("venues").Find(&venues).Error; err != nil {
		return &[]Venue{}, err
	}
	return &venues, nil
}

func GetVenueById(id int, db *gorm.DB) (*Venue, error) {
	venue := &Venue{}
	if err := db.Debug().Table("venues").Where("id = ?",
		id).First(venue).Error; err != nil {
		return nil, err
	}
	return venue, nil
}

func (venue *Venue) UpdateVenue(id int, db *gorm.DB) (*Venue, error) {
	if err := db.Debug().Table("venues").Where("id = ?", id).Updates(Venue{
		Name:           venue.Name,
		Description:    venue.Description,
		Location:       venue.Location,
		Category:       venue.Category,
		Capacity:       venue.Capacity,
	}).Error; err != nil {
		return &Venue{}, err
	}
	return venue, nil
}

func DeleteVenue(id int, db *gorm.DB) error {
	if err := db.Debug().Table("venues").Where("id = ?",
		id).Delete(&Venue{}).Error; err != nil {
		return err
	}
	return nil
}