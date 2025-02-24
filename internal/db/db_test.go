package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupDb(t *testing.T) {
	db := SetupDb()
	assert.NotNil(t, db.DB)

}
