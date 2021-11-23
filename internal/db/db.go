package db

import (
	"github.com/abdollahpour/almaniha-draft/internal/model"
)

type Db interface {
	CreateOrUpdate(draft model.Draft) (*model.Draft, error)
	Read(id string) (*model.Draft, error)
	Delete(id string) (*model.Draft, error)
	Paginated(page int64, size int64, filter model.DraftFilter, sort model.DraftSort) (*model.PaginatedDraft, error)
}
