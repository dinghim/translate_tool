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
	ins := lua.GetInstance()
	entry, err := ins.GetString(text)
	if err != nil {
		t.Fatal(err)
	}
	trans := []string{"test", "测试"}
	for i := 0; i < len(entry); i++ {
		fmt.Printf("%d %s\n", i+1, entry[i])
		ins.ReplaceOnce(&text, entry[i], []byte(fmt.Sprintf("%s-%d", trans[i%2], i)))
	}
	if err = ft.WriteAll("../../test/temp/test.lua", text); err != nil {
		t.Fatal(err)
	}
}
