package util

import (
	"testing"

	"github.com/abdollahpour/almaniha-draft/internal/config"
)

func TestPaginatedWrongKeyword(t *testing.T) {
	ConnectMongoDB(config.NewEnvConfiguration())
}
