package ze

import (
	"github.com/roidaradal/fn/clock"
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
	Code string `fx:"upper"`
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

// Set ID of item
func (x *Item) SetID(id ID) {
	x.ID = id
}

// Set IsActive of item
func (x *Item) SetIsActive(isActive bool) {
	x.IsActive = isActive
}
