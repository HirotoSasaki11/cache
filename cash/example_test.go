package cash

import (
	"fmt"
)

func ExampleCash() {
	c := NewCash(&CashOptions{
		Cachers: []Cacher{
			&LRUCache{
				Size: 3,
			},
			&MapCache{},
		},
	})

	// Output:
	// err = <nil>
	// v = "", err = cash: no cache found
	// v = "hi", err = <nil>
	// v = "hi", err = <nil>
	var (
		v   string
		err error
	)
	err = c.Delete("testing")
	fmt.Printf("err = %v\n", err)

	v = ""
	err = c.Load("testing", &v)
	fmt.Printf("v = %q, err = %v\n", v, err)

	v = ""
	err = c.LoadOrStore("testing", &v, func() (interface{}, error) { return "hi", nil })
	fmt.Printf("v = %q, err = %v\n", v, err)

	err = c.cachers[0].Delete("testing")
	if err != nil {
		panic(err)
	}
	v = ""
	err = c.Load("testing", &v)
	fmt.Printf("v = %q, err = %v\n", v, err)
}

func ExampleCash_withCompaction() {
	c := NewCash(&CashOptions{
		Codecs: []Codec{
			new(DeflateCodec),
		},
		Cachers: []Cacher{
			&LRUCache{
				Size: 3,
			},
		},
	})

	// Output:
	// err = <nil>
	// v = "", err = cash: no cache found
	// v = "hi", err = <nil>
	// v = "hi", err = <nil>
	var (
		v   string
		err error
	)
	err = c.Delete("testing")
	fmt.Printf("err = %v\n", err)

	v = ""
	err = c.Load("testing", &v)
	fmt.Printf("v = %q, err = %v\n", v, err)

	v = ""
	err = c.LoadOrStore("testing", &v, func() (interface{}, error) { return "hi", nil })
	fmt.Printf("v = %q, err = %v\n", v, err)

	v = ""
	err = c.Load("testing", &v)
	fmt.Printf("v = %q, err = %v\n", v, err)
}
