package tabfile

import (
	"bytes"
	"sync"
)

var (
	cr byte = 0x0d //回车CR
	lf byte = 0x0a //换行LF
	tb byte = 0x09 //tab制表符
)

const (
	state_normal        = iota //正常状态
	state_double_quotes        //"双引号"字符串
)

type tabfile struct{}

var instance *tabfile
var once sync.Once

func GetInstance() *tabfile {
	once.Do(func() {
		instance = &tabfile{}
	})
	return instance
}

func (t *tabfile) filter(text []byte) bool {
	for i := 0; i < len(text); i++ {
		if text[i]&0x80 != 0 {
			return false
		}
	}
	return true
}

func (t *tabfile) GetString(text []byte) ([][]byte, error) {
	var cnEntry [][]byte
	frecord := func(nStart, nEnd int) {
		textv := bytes.Split(text[nStart:nEnd], []byte{tb})
		for _, v := range textv {
			v = bytes.TrimSpace(v)
			if !t.filter(v) {
				cnEntry = append(cnEntry, v)
			}
		}
	}
	nStart := 0
	length := len(text)
	for i := 0; i < length; i++ {
		if i+1 < length && text[i] == cr && text[i] == lf {
			frecord(nStart, i)
			nStart = i + 2
		} else if text[i] == cr || text[i] == lf {
			frecord(nStart, i)
			nStart = i + 1
		}
	}
	return cnEntry, nil
}

func (t *tabfile) ReplaceOnce(context *[]byte, sText []byte, trans []byte) error {
	*context = bytes.Replace(*context, sText, trans, 1)
	return nil
}
