package rdb

type ID = uint
type DateTime = string

type UniqueItem struct {
	ID ID `json:"-"`
}

func (x UniqueItem) GetID() ID {
	return x.ID
}

type CreatedItem struct {
	CreatedAt DateTime `json:"CreatedAt,omitempty"`
}

func (x CreatedItem) GetDateTime() DateTime {
	return x.CreatedAt
}

type ActiveItem struct {
	IsActive bool
}

func (x ActiveItem) CheckIfActive() bool {
	return x.IsActive
}

type CodedItem struct {
	Code string
}

func (x CodedItem) GetCode() string {
	return x.Code
}

// ID, CreatedAt, IsActive, Code
type Item struct {
	UniqueItem
	CreatedItem
	ActiveItem
	CodedItem
}

type KV struct {
	Key           string `col:"AppKey"`
	Value         string `col:"AppValue"`
	LastUpdatedAt DateTime
}

type Access struct {
	ActiveItem
	Action string
	Role   string
}

type ActionDetails struct {
	Action  string
	Details string
}

type ActionLog struct {
	CreatedItem
	ActionDetails
	ActorID    ID     `json:"-"`
	ActorCode_ string `col:"-" json:"ActorCode"`
}

type BatchLog struct {
	CreatedItem
	CodedItem
	ActionDetails
}

type BatchLogItem struct {
	CodedItem
	Details string
}
