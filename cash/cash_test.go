package cash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testingValue struct {
	ID string
}

func TestCash(t *testing.T) {
	tStrct := testingValue{
		ID: "testing",
	}
	tStrctPtr := &testingValue{
		ID: "testingPtr",
	}
	newCash := func(state map[string]interface{}) *Cash {
		out := NewCash(&CashOptions{
			Cachers: []Cacher{
				&MapCache{},
			},
		})
		for k, v := range state {
			if err := out.Store(k, v); err != nil {
				panic(err)
			}
		}
		return out
	}
	t.Run("Load", func(t *testing.T) {
		c := newCash(map[string]interface{}{
			"string":    "hi",
			"int":       1,
			"bool":      true,
			"float":     1.0,
			"struct":    tStrct,
			"structPtr": tStrctPtr,
		})
		t.Run("string", func(t *testing.T) {
			var v string
			err := c.Load("string", &v)
			if assert.NoError(t, err) {
				assert.Equal(t, "hi", v)
			}
		})
		t.Run("int", func(t *testing.T) {
			var v int
			err := c.Load("int", &v)
			if assert.NoError(t, err) {
				assert.Equal(t, 1, v)
			}
		})
		t.Run("bool", func(t *testing.T) {
			var v bool
			err := c.Load("bool", &v)
			if assert.NoError(t, err) {
				assert.Equal(t, true, v)
			}
		})
		t.Run("float", func(t *testing.T) {
			var v float64
			err := c.Load("float", &v)
			if assert.NoError(t, err) {
				assert.Equal(t, 1.0, v)
			}
		})
		t.Run("struct", func(t *testing.T) {
			var v testingValue
			err := c.Load("struct", &v)
			if assert.NoError(t, err) {
				assert.Equal(t, tStrct, v)
			}
		})
		t.Run("structPtr", func(t *testing.T) {
			var v testingValue
			err := c.Load("structPtr", &v)
			if assert.NoError(t, err) {
				assert.Equal(t, tStrctPtr, &v)
			}
		})
	})
	t.Run("LoadOrStore", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			c := newCash(nil)
			var v string
			err := c.LoadOrStore("string", &v, func() (interface{}, error) {
				return "hi", nil
			})
			if assert.NoError(t, err) {
				assert.Equal(t, "hi", v)
			}
			v = ""
			err = c.Load("string", &v)
			if assert.NoError(t, err) {
				assert.Equal(t, "hi", v)
			}
		})
		t.Run("int", func(t *testing.T) {
			c := newCash(nil)
			var v int
			err := c.LoadOrStore("int", &v, func() (interface{}, error) {
				return 1, nil
			})
			if assert.NoError(t, err) {
				assert.Equal(t, 1, v)
			}
			v = 0
			err = c.Load("int", &v)
			if assert.NoError(t, err) {
				assert.Equal(t, 1, v)
			}
		})
		t.Run("bool", func(t *testing.T) {
			c := newCash(nil)
			var v bool
			err := c.LoadOrStore("bool", &v, func() (interface{}, error) {
				return true, nil
			})
			if assert.NoError(t, err) {
				assert.Equal(t, true, v)
			}
			v = false
			err = c.Load("bool", &v)
			if assert.NoError(t, err) {
				assert.Equal(t, true, v)
			}
		})
		t.Run("float", func(t *testing.T) {
			c := newCash(nil)
			var v float64
			err := c.LoadOrStore("float", &v, func() (interface{}, error) {
				return 1.0, nil
			})
			if assert.NoError(t, err) {
				assert.Equal(t, 1.0, v)
			}
			v = 0
			err = c.Load("float", &v)
			if assert.NoError(t, err) {
				assert.Equal(t, 1.0, v)
			}
		})
		t.Run("struct", func(t *testing.T) {
			t.Run("store as value", func(t *testing.T) {
				c := newCash(nil)
				var v testingValue
				err := c.LoadOrStore("struct", &v, func() (interface{}, error) {
					return tStrct, nil
				})
				if assert.NoError(t, err) {
					assert.Equal(t, tStrct, v)
				}
				v = testingValue{}
				err = c.Load("struct", &v)
				if assert.NoError(t, err) {
					assert.Equal(t, tStrct, v)
				}
			})
			t.Run("store as pointer", func(t *testing.T) {
				c := newCash(nil)
				var v testingValue
				err := c.LoadOrStore("struct", &v, func() (interface{}, error) {
					return &tStrct, nil
				})
				if assert.NoError(t, err) {
					assert.Equal(t, tStrct, v)
				}
				v = testingValue{}
				err = c.Load("struct", &v)
				if assert.NoError(t, err) {
					assert.Equal(t, tStrct, v)
				}
			})
		})
		t.Run("structPtr", func(t *testing.T) {
			t.Run("store as value", func(t *testing.T) {
				c := newCash(nil)
				var v *testingValue
				err := c.LoadOrStore("structPtr", &v, func() (interface{}, error) {
					return *tStrctPtr, nil
				})
				if assert.NoError(t, err) {
					assert.Equal(t, tStrctPtr, v)
				}
				v = nil
				err = c.Load("structPtr", &v)
				if assert.NoError(t, err) {
					assert.Equal(t, tStrctPtr, v)
				}
			})
			t.Run("store as pointer", func(t *testing.T) {
				c := newCash(nil)
				var v *testingValue
				err := c.LoadOrStore("structPtr", &v, func() (interface{}, error) {
					return tStrctPtr, nil
				})
				if assert.NoError(t, err) {
					assert.Equal(t, tStrctPtr, v)
				}
				v = nil
				err = c.Load("structPtr", &v)
				if assert.NoError(t, err) {
					assert.Equal(t, tStrctPtr, v)
				}
			})
		})
	})
}
