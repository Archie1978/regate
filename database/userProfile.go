package database

import (
	"bytes"
	"image/png"

	"github.com/o1egl/govatar"

	"gorm.io/gorm"
)

type UserProfile struct {
	gorm.Model
	Referance string

	Name     string
	LastName string
	Email    string
	Avatar   []byte
}

// Save
func SaveUser(user *UserProfile) error {
	if len(user.Avatar) == 0 {
		// Generate avatar if not exist
		imageAvatar, err := govatar.Generate(govatar.MALE)
		if err == nil {
			var buffer bytes.Buffer
			err := png.Encode(&buffer, imageAvatar)
			if err == nil {
				user.Avatar = buffer.Bytes()
			}
		}
	}
	ret := DB.Model(&UserProfile{}).Save(&user)
	return ret.Error
}

// Load User
func LoadUser(ref string) (*UserProfile, error) {
	var user UserProfile
	ret := DB.Debug().Model(&UserProfile{}).Where(UserProfile{Referance: ref}).Find(&user)
	return &user, ret.Error
}

// NewUser: Create new user with avatar
func NewUser() (*UserProfile, error) {
	var userProfil UserProfile

	// Generate avatar if not exist
	imageAvatar, err := govatar.Generate(govatar.MALE)
	if err == nil {
		var buffer bytes.Buffer
		err := png.Encode(&buffer, imageAvatar)
		if err == nil {
			userProfil.Avatar = buffer.Bytes()
		}
	}
	return &userProfil, nil
}
