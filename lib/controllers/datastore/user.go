/* AuthPlz Authentication and Authorization Microservice
 * Datastore - User objects
 *
 * Copyright 2018 Ryan Kurte
 */

package datastore

import (
	"errors"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"

	"github.com/authplz/authplz-core/lib/controllers/datastore/oauth2"
)

var ErrInvalidQuery = errors.New("Invalid DB Query argument")

// User represents the user for this application
type User struct {
	ID              uint      `gorm:"primary_key" description:"External user ID"`
	CreatedAt       time.Time `description:"Creation time"`
	UpdatedAt       time.Time `description:"Last update time"`
	DeletedAt       *time.Time
	ExtID           string `gorm:"not null;unique"`
	Email           string `gorm:"not null;unique"`
	Username        string `gorm:"not null;unique"`
	Password        string `gorm:"not null"`
	PasswordChanged time.Time
	Activated       bool `gorm:"not null; default:false"`
	Enabled         bool `gorm:"not null; default:false"`
	Locked          bool `gorm:"not null; default:false"`
	Admin           bool `gorm:"not null; default:false"`
	LoginRetries    uint `gorm:"not null; default:0"`
	LastLogin       time.Time

	ActionTokens []ActionToken
	FidoTokens   []FidoToken
	TotpTokens   []TotpToken
	BackupTokens []BackupToken
	AuditEvents  []AuditEvent

	OauthClients               []oauthstore.OauthClient
	OauthAccessTokenSessions   []oauthstore.OauthAccessToken
	OauthAuthorizeCodeSessions []oauthstore.OauthAuthorizeCode
	OauthRefreshTokenSessions  []oauthstore.OauthRefreshToken
}

// Getters and Setters

// GetIntID fetches a users internal ID
func (u *User) GetIntID() uint { return u.ID }

// GetExtID fetches a users ExtID
func (u *User) GetExtID() string { return u.ExtID }

// GetEmail fetches a users Email
func (u *User) GetEmail() string { return u.Email }

// GetUsername fetches a users Username
func (u *User) GetUsername() string { return u.Username }

// GetPassword fetches a users Password
func (u *User) GetPassword() string { return u.Password }

// GetPasswordChanged fetches a users PasswordChanged time
func (u *User) GetPasswordChanged() time.Time { return u.PasswordChanged }

// IsActivated checks if a user is activated
func (u *User) IsActivated() bool { return u.Activated }

// SetActivated sets a users activated status
func (u *User) SetActivated(activated bool) { u.Activated = activated }

// IsEnabled checks if a user is enabled
func (u *User) IsEnabled() bool { return u.Enabled }

// SetEnabled sets a users enabled status
func (u *User) SetEnabled(enabled bool) { u.Enabled = enabled }

// IsLocked checkes if a user account is locked
func (u *User) IsLocked() bool { return u.Locked }

// SetLocked sets a users locked status
func (u *User) SetLocked(locked bool) { u.Locked = locked }

// IsAdmin checks if a user is an admin
func (u *User) IsAdmin() bool { return u.Admin }

// SetAdmin sets a users admin status
func (u *User) SetAdmin(admin bool) { u.Admin = admin }

// GetLoginRetries fetches a users login retry count
func (u *User) GetLoginRetries() uint { return u.LoginRetries }

// SetLoginRetries sets a users login retry count
func (u *User) SetLoginRetries(retries uint) { u.LoginRetries = retries }

// ClearLoginRetries clears a users login retry count
func (u *User) ClearLoginRetries() { u.LoginRetries = 0 }

// GetLastLogin fetches a users LastLogin time
func (u *User) GetLastLogin() time.Time { return u.LastLogin }

// GetCreatedAt fetches a users account creation time
func (u *User) GetCreatedAt() time.Time { return u.CreatedAt }

// SetLastLogin sets a users LastLogin time
func (u *User) SetLastLogin(t time.Time) { u.LastLogin = t }

// SecondFactors Checks if a user has attached second factors
func (u *User) SecondFactors() bool {
	return (len(u.FidoTokens) > 0) || (len(u.TotpTokens) > 0)
}

// SetPassword sets a user password
func (u *User) SetPassword(pass string) {
	u.Password = pass
	u.PasswordChanged = time.Now()
}

// AddUser Adds a user to the datastore
func (dataStore *DataStore) AddUser(email, username, pass string) (interface{}, error) {

	if !govalidator.IsEmail(email) {
		return nil, fmt.Errorf("invalid email address %s", email)
	}

	user := &User{
		Email:     email,
		Username:  username,
		Password:  pass,
		ExtID:     uuid.NewV4().String(),
		Enabled:   true,
		Activated: false,
		Locked:    false,
		Admin:     false,
		CreatedAt: time.Now(),
	}

	dataStore.db = dataStore.db.Create(user)
	err := dataStore.db.Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail Fetches a user account by email
func (dataStore *DataStore) GetUserByEmail(email string) (interface{}, error) {
	if email == "" {
		return nil, ErrInvalidQuery
	}

	var user User
	err := dataStore.db.Where(&User{Email: email}).First(&user).Error
	if (err != nil) && (err != gorm.ErrRecordNotFound) {
		return nil, err
	} else if (err != nil) && (err == gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, nil
}

// GetUserByExtID Fetch a user account by external id
func (dataStore *DataStore) GetUserByExtID(extID string) (interface{}, error) {
	if extID == "" {
		return nil, ErrInvalidQuery
	}

	var user User
	err := dataStore.db.Where(&User{ExtID: extID}).First(&user).Error
	if (err != nil) && (err != gorm.ErrRecordNotFound) {
		return nil, err
	} else if (err != nil) && (err == gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, nil
}

// GetUserByUsername Fetches a user account by username
func (dataStore *DataStore) GetUserByUsername(username string) (interface{}, error) {
	if username == "" {
		return nil, ErrInvalidQuery
	}

	var user User
	err := dataStore.db.Where(&User{Username: username}).First(&user).Error
	if (err != nil) && (err != gorm.ErrRecordNotFound) {
		return nil, err
	} else if (err != nil) && (err == gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, nil
}

// UpdateUser Update a user object
func (dataStore *DataStore) UpdateUser(user interface{}) (interface{}, error) {
	u := user.(*User)

	dataStore.db = dataStore.db.Save(&u)
	if dataStore.db.Error != nil {
		return nil, dataStore.db.Error
	}

	return user, nil
}

// GetTokens Fetches tokens attached to a user account
func (dataStore *DataStore) GetTokens(user interface{}) (interface{}, error) {
	var err error

	u := user.(*User)

	err = dataStore.db.Model(user).Related(&u.FidoTokens).Error
	if err != nil {
		return nil, err
	}
	err = dataStore.db.Model(user).Related(&u.TotpTokens).Error
	if err != nil {
		return nil, err
	}
	err = dataStore.db.Model(user).Related(&u.BackupTokens).Error
	if err != nil {
		return nil, err
	}

	return u, nil
}
