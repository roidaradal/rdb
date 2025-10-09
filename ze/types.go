package ze

import (
	"database/sql"

	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/rdb"
)

type ID = uint
type DateTime = string
type Date = string

// Embeddable ID property
type UniqueItem struct {
	ID ID `json:"-"`
}

func (x UniqueItem) GetID() ID {
	return x.ID
}

// Embeddable Code property
type CodedItem struct {
	Code string
}

func (x CodedItem) GetCode() string {
	return x.Code
}

// ID and Code
type Identity struct {
	UniqueItem
	CodedItem
}

// Embeddable CreatedAt property
type CreatedItem struct {
	CreatedAt DateTime
}

func (x CreatedItem) GetDateTime() DateTime {
	return x.CreatedAt
}

// Embeddable IsActive property
type ActiveItem struct {
	IsActive bool
}

func (x ActiveItem) CheckIfActive() bool {
	return x.IsActive
}

// ID, Code, IsActive, CreatedAt
type Item struct {
	Identity
	ActiveItem
	CreatedItem
}

// Initialize the ID, CreatedAt, IsActive to default values
func (x *Item) Initialize() {
	x.ID = 0 // for auto-increment
	x.CreatedAt = clock.DateTimeNow()
	x.IsActive = true
}

// Trio: DB, DBTx, Checker
type DBTrio struct {
	DB      *sql.DB
	DBTx    *sql.Tx
	Checker rdb.QueryResultChecker
}
