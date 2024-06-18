package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRowToJson(t *testing.T) {
	db, cleanup := setupDatabase(t)
	defer cleanup()

	result, err := db.Select("SELECT * FROM albums WHERE AlbumId = 1 LIMIT 1;")

	assert.Nil(t, err)
	assert.Equal(t, result.GetNumberOfRows(), uint64(1))

	row, err := result.GetRow(0)
	assert.Nil(t, err)

	json, err := row.ToJSON()
	assert.Nil(t, err)

	assert.JSONEq(t, `{"Index":0,"ColumnValues":[{"Type":58,"Buffer":"MQ=="},{"Type":43,"Buffer":"Rm9yIFRob3NlIEFib3V0IFRvIFJvY2sgV2UgU2FsdXRlIFlvdQ=="},{"Type":58,"Buffer":"MQ=="}]}`, string(json))
}
