package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/abdollahpour/almaniha-draft/internal/db"
	"github.com/abdollahpour/almaniha-draft/internal/model"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Api interface {
	CreateDraft() http.Handler
	DraftsPaginatedHandler() http.Handler
	LivenessHandler() http.Handler
}

type DbApi struct {
	db db.Db
}

func NewDbApi(db db.Db) *DbApi {
	return &DbApi{db: db}
}

func getNumericParam(r *http.Request, name string, def int, min int, max int) int {
	value := r.URL.Query().Get(name)
	p, e := strconv.Atoi(value)
	if e != nil {
		return def
	} else if p < min {
		return min
	} else if p > max {
		return max
	}
	return p
}

func getParam(r *http.Request, name string, def string) string {
	param := r.URL.Query().Get(name)
	if len(param) > 0 {
		return param
	}
	return def
}

func getDateParam(r *http.Request, name string) *time.Time {
	param := r.URL.Query().Get(name)
	if len(param) > 0 {
		t, err := time.Parse(time.RFC3339, param)
		fmt.Println(t)
		if err != nil {
			log.WithField("param", param).WithError(err).Info("Time format is not currect")
		}
		if !t.IsZero() {
			return &t
		}
	}
	return nil
}

func getDraftSort(r *http.Request) model.DraftSort {
	return model.DraftSort{
		Field: model.SortField(getParam(r, "field", "createAt")),
		Order: model.SortOrder(getParam(r, "order", string(model.DESC))),
	}
}

func getCorrelationId(r *http.Request) string {
	correlationId := r.Header.Get("X-Correlation-ID")
	if correlationId == "" {
		correlationId = uuid.NewString()
		log.WithField("correlationId", correlationId).Debug("Correlation ID not found. New correlation ID key generated")
	}
	return correlationId
}

func getDraftFilter(r *http.Request) model.DraftFilter {
	return model.DraftFilter{
		Type: model.DraftType(r.URL.Query().Get("type")),
	}
}

func (api *DbApi) CreateDraft() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationId := getCorrelationId(r)

		var draft model.Draft
		err := json.NewDecoder(r.Body).Decode(&draft)
		if err != nil {
			log.WithError(err).WithField("correlationId", correlationId).Warn("Fail to parse object")
			w.WriteHeader(http.StatusBadRequest)
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  "Fail to parse object",
				Status: strconv.FormatInt(http.StatusBadRequest, 10),
				Code:   "REQ-102",
			}})
			return
		}

		created, err := api.db.CreateOrUpdate(draft)
		if err != nil {
			log.WithError(err).WithField("correlationId", correlationId).Warn("Fail to store")
			w.WriteHeader(http.StatusBadRequest)
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  "Fail to store",
				Status: strconv.FormatInt(http.StatusBadRequest, 10),
				Code:   "REQ-103",
			}})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(created)
		if err != nil {
			log.WithError(err).WithField("correlationId", correlationId).Error("Fail to serialize")
			w.WriteHeader(http.StatusBadRequest)
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  "Fail to serialize",
				Status: strconv.FormatInt(http.StatusInternalServerError, 10),
				Code:   "REQ-104",
			}})
			return
		}

		log.WithFields(log.Fields{
			"correlationId": correlationId,
			"created":       created,
		}).Info("Draft created")
	})
}

func (api *DbApi) ReadDraft() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		draftId := vars["draftId"]
		correlationId := getCorrelationId(r)

		draft, err := api.db.Read(draftId)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"correlationId": correlationId,
				"draftId":       draftId,
			}).Error("Fail to read the draft")
			w.WriteHeader(http.StatusBadRequest)
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  "Fail to read the draft",
				Detail: err.Error(),
				Status: strconv.FormatInt(http.StatusInternalServerError, 10),
				Code:   "REQ-105",
				Meta: &map[string]interface{}{
					"correlationId": correlationId,
					"draftId":       draftId,
				},
			}})
			return
		}
		if draft == nil {
			log.WithError(err).WithFields(log.Fields{
				"correlationId": correlationId,
				"draftId":       draftId,
			}).Warn("Not found")
			w.WriteHeader(http.StatusBadRequest)
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  "Not found",
				Status: strconv.FormatInt(http.StatusNotFound, 10),
				Code:   "REQ-106",
				Meta: &map[string]interface{}{
					"correlationId": correlationId,
					"draftId":       draftId,
				},
			}})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(draft)
		if err != nil {
			log.WithError(err).WithFields(log.Fields{
				"correlationId": correlationId,
				"draftId":       draftId,
			}).Error("Fail to serialize")
			w.WriteHeader(http.StatusBadRequest)
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  "Fail to serialize",
				Detail: err.Error(),
				Status: strconv.FormatInt(http.StatusInternalServerError, 10),
				Code:   "REQ-107",
				Meta: &map[string]interface{}{
					"correlationId": correlationId,
					"draftId":       draftId,
				},
			}})
			return
		}

		log.WithFields(log.Fields{
			"correlationId": correlationId,
			"draftId":       draftId,
			"draft":         draft,
		}).Info("Read draft")
	})
}

func (api *DbApi) DraftsPaginatedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := getNumericParam(r, "page", 1, math.MaxInt32, 100000)
		size := getNumericParam(r, "size", 20, 1, 100)
		filter := getDraftFilter(r)
		sort := getDraftSort(r)

		posts, err := api.db.Paginated(int64(page), int64(size), filter, sort)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(posts)
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (api *DbApi) LivenessHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
