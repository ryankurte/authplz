package datastore

import "fmt"

import "github.com/jinzhu/gorm"
import _ "github.com/jinzhu/gorm/dialects/postgres"

// Datastore instance storage
type DataStore struct {
	db *gorm.DB
}

// Query filter types
type QueryFilter struct {
	Limit  uint // Number of objects to return
	Offset uint // Offset of objects to return
}

// Create a datastore instance
func NewDataStore(dbString string) (*DataStore, error) {
	// Attempt database connection
	db, err := gorm.Open("postgres", dbString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %s", dbString)
	}

	//db = db.LogMode(true)

	return &DataStore{db}, nil
}

// Close an open datastore instance
func (dataStore *DataStore) Close() {
	dataStore.db.Close()
}

// Drop and create existing tables to match required schema
// WARNING: do not run this on a live database...
func (ds *DataStore) ForceSync() {
	db := ds.db

	db = db.Exec("DROP TABLE IF EXISTS fido_tokens CASCADE;")
	db = db.Exec("DROP TABLE IF EXISTS totp_tokens CASCADE;")
	db = db.Exec("DROP TABLE IF EXISTS audit_events CASCADE;")
	db = db.Exec("DROP TABLE IF EXISTS users CASCADE;")

	db = db.AutoMigrate(&User{})
	db = db.AutoMigrate(&FidoToken{})
	db = db.AutoMigrate(&TotpToken{})
	db = db.AutoMigrate(&AuditEvent{})

	db = db.Model(&User{}).AddUniqueIndex("idx_user_email", "email")
	db = db.Model(&User{}).AddUniqueIndex("idx_user_ext_id", "ext_id")
}

func (ds *DataStore) AddAuditEvent(u *User, auditEvent *AuditEvent) (user *User, err error) {
	u.AuditEvents = append(u.AuditEvents, *auditEvent)
	u, err = ds.UpdateUser(u)
	return u, err
}

func (dataStore *DataStore) GetAuditEvents(u *User) ([]AuditEvent, error) {
	var auditEvents []AuditEvent

	err := dataStore.db.Model(u).Related(&auditEvents).Error

	return auditEvents, err
}
