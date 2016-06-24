package lua_test

import (
	"fmt"
	"testing"
	"trans/analysis/lua"
	"trans/filetool"
)

func Test_Example(t *testing.T) {
	ft := filetool.GetInstance()
	ft.SetEncoding("lua", "utf8")
	text, err := ft.ReadAll("../../test/cn/test.lua")
	if err != nil {
		t.Fatal(err)
	}
	entry, start, end, err := lua.New().GetString(text)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(entry); i++ {
		fmt.Printf("%d [%d:%d] %s\n", i+1, start[i], end[i], entry[i])
	}
}
