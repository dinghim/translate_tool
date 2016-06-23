package dic_test

import (
	"fmt"
	"testing"
	"trans/dic"
)

func Test_example(t *testing.T) {
	file := "../test/temp/dictionary.txt"
	d := dic.New(file)
	path := "test"
	src := []byte("测试")
	d.Append(path, src)
	d.Save()
	if trans, ok := d.Query(path, src); !ok {
		t.Log("no tanslate")
	} else {
		fmt.Printf("%s\n", trans)
	}
}
