package tabfile_test

import (
	"fmt"
	"testing"
	"trans/analysis/tabfile"
	"trans/filetool"
)

func Test_Example(t *testing.T) {
	ft := filetool.GetInstance()
	ft.SetEncoding(".tab", "gbk")
	text, err := ft.ReadAll("../../test/cn/ScriptItem.tab")
	if err != nil {
		t.Fatal(err)
	}
	entry, start, end, err := tabfile.New().GetString(text)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(entry); i++ {
		fmt.Printf("%d [%d:%d] %s\n", i+1, start[i], end[i], entry[i])
	}
}
