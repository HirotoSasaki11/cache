package cash

import (
	"sort"
	"strings"
)

type MapCashKey map[string]string

func (k MapCashKey) CashKey() string {
	keys := make([]string, 0, len(k))
	for key := range k {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	params := make([]string, len(keys))
	for i := range keys {
		params[i] = keys[i] + "=" + k[keys[i]]
	}
	return strings.Join(params, ";")
}
