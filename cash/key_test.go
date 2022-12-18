package cash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapCashKey(t *testing.T) {
	k := MapCashKey{
		"spaceId":   "sid",
		"networkId": "nid",
	}
	assert.Equal(t, "networkId=nid;spaceId=sid", k.CashKey())
}
