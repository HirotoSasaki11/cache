package cash

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNoCache = errors.New("cash: no cache found")
)

var (
	NilCacher Cacher = new(nilCacher)
)

type CashKey interface {
	CashKey() string
}

type CashValue interface {
	CashEncode() ([]byte, error)

	CashDecode([]byte) error
}

type Codec interface {
	Encode([]byte) ([]byte, error)

	Decode([]byte) ([]byte, error)
}

type Cash struct {
	// mu sync.RWMutex

	codecs []Codec

	cachers []Cacher
}

type CashOptions struct {
	Codecs []Codec

	Cachers []Cacher
}

type CacheLoader func() (value interface{}, err error)

func NewCash(opts *CashOptions) *Cash {
	csh := &Cash{}
	csh.codecs = opts.Codecs
	for _, cacher := range opts.Cachers {
		cacher.Init(csh)
		csh.cachers = append(csh.cachers, cacher)
	}
	return csh
}

func (c *Cash) Load(key interface{}, value interface{}) error {
	_, err := c.load(key, value)
	return err
}

func (c *Cash) load(key interface{}, value interface{}) (bool, error) {
	out := outValueOf(value)
	keyStr := c.EncodeKey(key)
	var (
		loaded     []byte
		foundLocal bool
	)
	for _, cacher := range c.cachers {
		b, err := cacher.Load(keyStr)
		if err == nil {
			loaded = b
			foundLocal = cacher.LocalCache()
			break
		} else if err != ErrNoCache {
			return false, err
		}
	}
	if len(loaded) == 0 {
		return false, ErrNoCache
	}
	decoded := reflect.New(out.Type()).Elem()
	if err := c.DecodeValue(loaded, decoded); err != nil {
		return foundLocal, err
	}
	out.Set(decoded)
	return foundLocal, nil
}

func (c *Cash) LoadOrStore(key interface{}, value interface{}, loader CacheLoader) error {
	foundLocal, err := c.load(key, value)
	if err == ErrNoCache {
		// go to store
	} else if err != nil {
		return err
	} else {
		if !foundLocal {
			// local cache is missed and hit in remote cache, push it to local
			return c.StoreLocal(key, value)
		} else {
			return nil
		}
	}
	loaded, err := loader()
	if err != nil {
		return err
	}
	lv := reflect.ValueOf(loaded)
	for lv.Kind() == reflect.Ptr {
		lv = lv.Elem()
	}
	out := outValueOf(value)
	if out.Kind() == reflect.Ptr && out.IsNil() {
		out.Set(reflect.New(out.Type().Elem()))
		out = out.Elem()
	}
	out.Set(lv)
	return c.Store(key, loaded)
}

func (c *Cash) StoreLocal(key interface{}, value interface{}) error {
	keyStr := c.EncodeKey(key)
	return c.store(c.localCachers(), keyStr, value)
}

func (c *Cash) Store(key interface{}, value interface{}) error {
	keyStr := c.EncodeKey(key)
	return c.store(c.cachers, keyStr, value)
}

func (c *Cash) store(cachers []Cacher, key string, value interface{}) error {
	valueBytes, err := c.EncodeValue(value)
	if err != nil {
		return err
	}
	var merr MultiError
	for _, cacher := range cachers {
		if err := cacher.Store(key, valueBytes); err != nil {
			merr.Append(err)
		}
	}
	return merr.ErrorOrNil()
}

func (c *Cash) DeleteLocal(key interface{}) error {
	keyStr := c.EncodeKey(key)
	return c.delete_(c.localCachers(), keyStr)
}

func (c *Cash) Delete(key interface{}) error {
	keyStr := c.EncodeKey(key)
	return c.delete_(c.cachers, keyStr)
}

func (c *Cash) localCachers() []Cacher {
	cachers := make([]Cacher, 0, len(c.cachers))
	for _, cacher := range c.cachers {
		if cacher.LocalCache() {
			cachers = append(cachers, cacher)
		}
	}
	return cachers
}

func (c *Cash) delete_(cachers []Cacher, key string) error {
	var merr MultiError
	for _, cacher := range cachers {
		if err := cacher.Delete(key); err != nil {
			merr.Append(err)
		}
	}
	return merr.ErrorOrNil()
}

func (c *Cash) EncodeKey(v interface{}) string {
	if ck, ok := v.(CashKey); ok {
		return ck.CashKey()
	}

	return fmt.Sprint(v)
}

func (c *Cash) EncodeValue(v interface{}) ([]byte, error) {
	if cv, ok := v.(CashValue); ok {
		return cv.CashEncode()
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}

	b := buf.Bytes()
	for _, codec := range c.codecs {
		var err error
		if b, err = codec.Encode(b); err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (c *Cash) DecodeValue(b []byte, v reflect.Value) error {
	var err error
	for _, codec := range c.codecs {
		if b, err = codec.Decode(b); err != nil {
			return err
		}
	}

	if cv, ok := v.Interface().(CashValue); ok {
		return cv.CashDecode(b)
	}

	return gob.NewDecoder(bytes.NewReader(b)).DecodeValue(v)
}
