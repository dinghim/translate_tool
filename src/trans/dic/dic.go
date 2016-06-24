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
	var err error
	ins.line, err = ft.ReadFileLine(file)
	if err != nil {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
		return ins
	}
	for i := 0; i < len(ins.line); i++ {
		v := ins.line[i]
		linev := bytes.Split(v, []byte{0x09})
		if len(linev) != 2 || len(linev[0]) == 0 || len(linev[1]) == 0 {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, fmt.Sprintf("[dic abnormal] file:%s, line:%d, data:%s", file, i+1, v))
			continue
		}
		key := string(linev[0])
		if _, ok := ins.trans[key]; ok {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, fmt.Sprintf("[dic repeat] file:%s, line:%d, data:%s", file, i+1, key))
			continue
		}
		value := string(linev[1])
		ins.trans[key] = value
	}
	return ins
}

func NewOnly(file string) *dic {
	return &dic{
		name:  file,
		trans: make(map[string]string),
	}
}

func (d *dic) Query(text []byte) ([]byte, bool) {
	stext := string(text)
	strans, ok := d.trans[stext]
	return []byte(strans), ok
}

func (d *dic) Append(text []byte, trans []byte) bool {
	stext := string(text)
	strans := string(trans)
	if _, ok := d.trans[stext]; ok {
		return false
	}
	d.trans[stext] = strans
	line := []byte(fmt.Sprintf("%s\t%s", stext, strans))
	d.line = append(d.line, line)
	return true
}

func (d *dic) Merge(target *dic) int {
	succ := 0
	for k, v := range d.trans {
		if len(k) > 0 && len(v) > 0 {
			if target.Append([]byte(k), []byte(v)) {
				succ += 1
			}
		}
	}
	return succ
}

func (d *dic) Save() {
	ft := filetool.GetInstance()
	oldEncode, _ := ft.SetEncoding(d.name, "utf8")
	defer ft.SetEncoding(d.name, oldEncode)
	err := ft.SaveFileLine(d.name, d.line)
	if err != nil {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
	}
}
