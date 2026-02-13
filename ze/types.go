package ze

import "github.com/roidaradal/fn/clock"

type (
	ID       = uint
	Date     = string
	DateTime = string
)

// Embeddable ID property
type UniqueItem struct {
	ID ID `json:"-"`
}

func (x UniqueItem) GetID() ID {
	return x.ID
}

func (x *UniqueItem) SetID(id ID) {
	x.ID = id
}

// Embeddable Code property
type CodedItem struct {
	Code string `fx:"upper"`
}

func (x CodedItem) GetCode() string {
	return x.Code
}

// Embeddable Code property, with validate: required
type RequiredCode struct {
	Code string `fx:"upper" validate:"required"`
}

func (x RequiredCode) GetCode() string {
	return x.Code
}

// ID and Code
type Identity struct {
	UniqueItem
	CodedItem
}

type CodedIdentity struct {
	UniqueItem
	RequiredCode
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

func (x *ActiveItem) SetIsActive(isActive bool) {
	x.IsActive = isActive
}

// ID, Code, IsActive, CreatedAt
type Item struct {
	AutoItem
	CodedItem
}

// ID, RequiredCode, IsActive, CreatedAt
type ItemCoded struct {
	AutoItem
	RequiredCode
}

// ID, IsActive, CreatedAt
type AutoItem struct {
	UniqueItem
	CreatedItem
	ActiveItem
}

// Initialize the ID, CreatedAt, IsActive to default values
func (x *Item) Initialize() {
	x.AutoItem.Initialize()
}

// Initialize the ID, CreatedAt, IsActive to default values
func (x *ItemCoded) Initialize() {
	x.AutoItem.Initialize()
}

// Initialize the ID, CratedAt, IsActive to default values
func (x *AutoItem) Initialize() {
	x.ID = 0 // for auto-increment
	x.CreatedAt = clock.DateTimeNow()
	x.IsActive = true
}
