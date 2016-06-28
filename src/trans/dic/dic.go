package dic

import (
	"bytes"
	"fmt"
	"trans/filetool"
	"trans/log"
)

type dic struct {
	name  string
	line  [][]byte
	trans map[string]string
}

func New(file string) *dic {
	ins := &dic{
		name:  file,
		trans: make(map[string]string),
	}
	ft := filetool.GetInstance()
	oldEncode, _ := ft.SetEncoding(file, "utf8")
	defer ft.SetEncoding(file, oldEncode)
	all, err := ft.ReadFileLine(file)
	if err != nil {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
		return ins
	}
	for i := 1; i < len(all); i++ {
		v := all[i]
		linev := bytes.Split(v, []byte{0x09})
		if len(linev) != 4 {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, fmt.Sprintf("[dic abnormal] file:%s, line:%d, data:%s", file, i+1, v))
			continue
		}
		key := string(linev[2])
		if _, ok := ins.trans[key]; ok {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, fmt.Sprintf("[dic repeat] file:%s, line:%d, data:%s", file, i+1, key))
			continue
		}
		ins.trans[key] = string(linev[3])
		ins.line = append(ins.line, v)
	}
	return ins
}

func (d *dic) Query(text []byte) ([]byte, bool) {
	stext := string(text)
	strans, ok := d.trans[stext]
	return []byte(strans), ok
}

func (d *dic) Append(path string, text []byte, trans []byte) bool {
	stext := string(text)
	strans := string(trans)
	if _, ok := d.trans[stext]; ok {
		return false
	}
	d.trans[stext] = strans
	line := []byte(fmt.Sprintf("%d\t%s\t%s\t%s", len(d.line)+1, path, stext, strans))
	d.line = append(d.line, line)
	return true
}

func (d *dic) Save() {
	ft := filetool.GetInstance()
	oldEncode, _ := ft.SetEncoding(d.name, "utf8")
	defer ft.SetEncoding(d.name, oldEncode)
	var all [][]byte
	all = append(all, []byte("ID\tFile\tOriginal\tTranslation"))
	all = append(all, d.line...)
	err := ft.SaveFileLine(d.name, all)
	if err != nil {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
	}
}
