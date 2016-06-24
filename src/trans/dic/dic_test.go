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
	des := []byte("呵呵")
	d.Append(path, src, des)
	d.Save()
	if trans, ok := d.Query(src); !ok {
		t.Log("no tanslate")
	} else {
		fmt.Printf("%s\n", trans)
	}
}
