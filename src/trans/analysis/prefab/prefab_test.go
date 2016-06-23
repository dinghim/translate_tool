package prefab_test

import (
	"fmt"
	"testing"
	"trans/analysis/prefab"
	"trans/filetool"
)

func Test_Example(t *testing.T) {
	ft := filetool.GetInstance()
	ft.SetEncoding("prefab", "utf8")
	text, err := ft.ReadAll("../../test/cn/Boss.prefab")
	if err != nil {
		t.Fatal(err)
	}
	ins := prefab.GetInstance()
	entry, err := ins.GetString(text)
	if err != nil {
		t.Fatal(err)
	}
	trans := []string{"test", "测试"}
	for i := 0; i < len(entry); i++ {
		fmt.Printf("%d %s\n", i+1, entry[i])
		if err := ins.ReplaceOnce(&text, entry[i], []byte(fmt.Sprintf("%s-%d", trans[i%2], i))); err != nil {
			t.Fatal(err)
		}
	}
	if err = ft.WriteAll("../../test/temp/Boss.prefab", text); err != nil {
		t.Fatal(err)
	}
}
