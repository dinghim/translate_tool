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
	trans map[string]map[string]string
}

func New(file string) *dic {
	ins := &dic{
		name:  file,
		trans: make(map[string]map[string]string),
	}
	ft := filetool.GetInstance()
	oldEncode, _ := ft.SetEncoding(file, "utf8")
	var err error
	ins.line, err = ft.ReadFileLine(file)
	if err != nil {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
		return ins
	}
	for _, v := range ins.line {
		linev := bytes.Split(v, []byte{0x09})
		if len(linev) != 4 {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, fmt.Sprintf("[dic abnormal] %s", v))
			continue
		}
		path := string(linev[1])
		key := string(linev[2])
		value := string(linev[3])
		if _, ok := ins.trans[path]; !ok {
			ins.trans[path] = make(map[string]string)
		}
		ins.trans[path][key] = value
	}
	ft.SetEncoding(file, oldEncode)
	return ins
}

func (d *dic) Query(path string, text []byte) (trans []byte, ok bool) {
	var strans string
	stext := string(text)
	_, ok = d.trans[path]
	if !ok {
		return
	}
	strans, ok = d.trans[path][stext]
	trans = []byte(strans)
	return
}

func (d *dic) Append(path string, text []byte) {
	stext := string(text)
	if _, ok := d.trans[path]; !ok {
		d.trans[path] = make(map[string]string)
	}
	if _, ok := d.trans[path][stext]; !ok {
		d.trans[path][stext] = ""
		line := []byte(fmt.Sprintf("%d\t%s\t%s\t%s", len(d.line)+1, path, stext, ""))
		d.line = append(d.line, line)
	}
}

func (d *dic) Save() {
	ft := filetool.GetInstance()
	oldEncode, _ := ft.SetEncoding(d.name, "utf8")
	err := ft.SaveFileLine(d.name, d.line)
	if err != nil {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
	}
	ft.SetEncoding(d.name, oldEncode)
}
