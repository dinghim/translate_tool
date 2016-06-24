package analysis_test

import (
	"testing"
	"trans/analysis"
	"trans/filetool"
	"trans/log"
)

func Test_GetString(t *testing.T) {
	ft := filetool.GetInstance()
	ft.SetEncoding(".lua", "utf8")
	ft.SetEncoding(".prefab", "utf8")
	ft.SetEncoding(".tab", "gbk")
	ft.SetEncoding(".txt", "gbk")
	anal := analysis.GetInstance()
	anal.SetRulesMap(".lua", "lua_rules")
	anal.SetRulesMap(".prefab", "prefab_rules")
	anal.SetRulesMap(".tab", "table_rules")
	anal.SetRulesMap(".txt", "table_rules")
	anal.SetFilterMap("filter")
	anal.SetFilterMap("filter.lua")
	anal.GetString("../test/temp/dictionary.txt", "../test/cn/")
}

func Test_Translate(t *testing.T) {
	defer log.CloseLog()
	ft := filetool.GetInstance()
	ft.SetEncoding(".lua", "utf8")
	ft.SetEncoding(".prefab", "utf8")
	ft.SetEncoding(".tab", "gbk")
	ft.SetEncoding(".txt", "gbk")
	anal := analysis.GetInstance()
	anal.SetRulesMap(".lua", "lua_rules")
	anal.SetRulesMap(".prefab", "prefab_rules")
	anal.SetRulesMap(".tab", "table_rules")
	anal.SetRulesMap(".txt", "table_rules")
	anal.SetFilterMap("filter")
	anal.SetFilterMap("filter.lua")
	anal.Translate("../test/temp/dictionary.txt", "../test/cn/", "../test/temp/trans/", 1)
}
