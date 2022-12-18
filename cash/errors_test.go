package cash

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiError(t *testing.T) {
	var err MultiError
	assert.NoError(t, err.ErrorOrNil())
	err.Append(errors.New("a"))
	assert.Error(t, err.ErrorOrNil())
	err.Append(errors.New("b"))
	assert.Error(t, err.ErrorOrNil())
	assert.Equal(t, "a\nb", err.ErrorOrNil().Error())
}
