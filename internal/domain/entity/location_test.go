package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCep(t *testing.T) {
	location := NewLocation()

	err := location.AddCep("88010290")
	assert.Nil(t, err)

	err = location.AddCep("0212121212")
	assert.Error(t, err, InvalidCep)

	err = location.AddCep("toptopto")
	assert.Error(t, err, InvalidCep)
}

func TestFillOtherTempsFromCelsius(t *testing.T) {
	location := NewLocation()
	err := location.FillOtherTempsFromCelsius()
	assert.Error(t, err, ErrorCelsiusEmpty)

	location.TempCelsius = 20
	err = location.FillOtherTempsFromCelsius()
	assert.Nil(t, err)
}