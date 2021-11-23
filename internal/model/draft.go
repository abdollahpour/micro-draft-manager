package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Draft struct {
	ID             primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Type           DraftType              `bson:"type,omitempty" json:"type,omitempty"`
	Data           map[string]interface{} `bson:"data,omitempty" json:"data,omitempty"`
	CreatedBy      string                 `bson:"createdBy,omitempty" json:"createdBy,omitempty"`
	CreatedAt      time.Time              `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	ModifiedAt     time.Time              `bson:"modifiedAt,omitempty" json:"modifiedAt,omitempty"`
	IdempotencyKey string                 `bson:"idempotencyKey,omitempty" json:"idempotencyKey,omitempty"`
	CorrelationId  string                 `bson:"correlationId,omitempty" json:"correlationId,omitempty"`
}

type SortField string

const (
	CREATE_AT SortField = "createdAt"
)

type SortOrder string

const (
	DESC SortOrder = "DESC"
	ASC  SortOrder = "ASC"
)

type DraftType string

const (
	Business DraftType = "Business"
)

type DraftSort struct {
	Field SortField `json:"field,omitempty"`
	Order SortOrder `json:"order,omitempty"`
}

type DraftFilter struct {
	Type DraftType `json:"type,omitempty"`
}

type PaginatedDraft struct {
	Paginated
	Data   []Draft     `json:"data,omitempty"`
	Sort   DraftSort   `json:"sort,omitempty"`
	Filter DraftFilter `json:"filter,omitempty"`
}
