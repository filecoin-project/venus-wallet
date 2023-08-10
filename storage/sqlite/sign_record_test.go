package sqlite

import (
	"fmt"
	"testing"
	"time"

	"github.com/filecoin-project/venus/venus-shared/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSingRecord(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate the schema
	s, err := NewSqliteRecorder(db, nil)
	assert.NoError(t, err)

	err = s.Record(&types.SignRecord{
		RawMsg:   []byte("hello"),
		Err:      fmt.Errorf("error"),
		Type:     types.MTVerifyAddress,
		CreateAt: time.Now(),
	})
	assert.NoError(t, err)
	res, err := s.QueryRecord(&types.QuerySignRecordParams{})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(res))

}
