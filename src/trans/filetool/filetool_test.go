package filetool_test

import (
	"fmt"
	"testing"
	"trans/filetool"
)

func Test_GetInstance(t *testing.T) {
	ft1 := filetool.GetInstance()
	ft2 := filetool.GetInstance()
	if ft1 != ft2 {
		t.Fatal("GetInstance diffrent value")
	}
}

func Test_GetFileMap(t *testing.T) {
	ft := filetool.GetInstance()
	fm, err := ft.GetFilesMap("../")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(fm); i++ {
		fmt.Println(i, fm[i])
	}
}

func Test_ReadFileLine(t *testing.T) {
	ft := filetool.GetInstance()
	context, err := ft.ReadFileLine("../test/cn.txt")
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range context {
		fmt.Printf("%d %s\n", k, v)
	}
}

func Test_SaveFileLine(t *testing.T) {
	ft := filetool.GetInstance()
	context, err := ft.ReadFileLine("../test/cn.txt")
	if err != nil {
		t.Fatal(err)
	}
	ft.SetEncoding(".txt", "gbk")
	err = ft.SaveFileLine("../test/cn1.txt", context)
	if err != nil {
		fmt.Println(err)
	}
	ft.SetEncoding(".txt", "hz-gb2312")
	err = ft.SaveFileLine("../test/cn2.txt", context)
	if err != nil {
		fmt.Println(err)
	}
	ft.SetEncoding(".txt", "gb18030")
	err = ft.SaveFileLine("../test/cn3.txt", context)
	if err != nil {
		fmt.Println(err)
	}
	_, err = ft.SetEncoding(".txt", "tcvn3")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_ReadAll(t *testing.T) {
	ft := filetool.GetInstance()
	bv, err := ft.ReadAll("../test/cn/test.lua")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", bv)
}

func Test_WriteAll(t *testing.T) {
	ft := filetool.GetInstance()
	bv, err := ft.ReadAll("../test/cn/test.lua")
	if err != nil {
		t.Fatal(err)
	}
	err = ft.WriteAll("../test/test/test/test.lua", bv)
	if err != nil {
		t.Fatal(err)
	}
	err = ft.WriteAll("../test/test2.txt", bv)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GbkFile(t *testing.T) {
	ft := filetool.GetInstance()
	bv, err := ft.ReadAll("../test/cn/ScriptItem.tab")
	if err != nil {
		t.Fatal(err)
	}
	ft.SetEncoding(".tab", "utf8")
	err = ft.WriteAll("../test/ScriptItem2.tab", bv)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Big5File(t *testing.T) {
	ft := filetool.GetInstance()
	ft.SetEncoding(".txt", "utf8")
	bv, err := ft.ReadAll("../test/big5.txt")
	if err != nil {
		t.Fatal(err)
	}
	ft.SetEncoding(".txt", "big5")
	err = ft.WriteAll("../test/big5_1.txt", bv)
	if err != nil {
		t.Fatal(err)
	}
	ft.SetEncoding(".txt", "gbk")
	err = ft.WriteAll("../test/big5_2.txt", bv)
	if err != nil {
		t.Fatal(err)
	}
}
