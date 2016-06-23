package prefab

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var (
	dq byte = 0x22 //双引号"
	sl byte = 0x5c //转义斜杠\\
	uu byte = 0x75 //u字符
	cr byte = 0x0d //回车CR
	lf byte = 0x0a //换行LF
)

const (
	state_normal        = iota //正常状态
	state_double_quotes        //"双引号"字符串
)

var const_string_flag []byte = []byte{109, 84, 101, 120, 116, 58}
var const_cr_lf_flag []byte = []byte{13, 10, 32, 32, 32, 32, 92}
var const_cr_flag []byte = []byte{13, 32, 32, 32, 32, 92}
var const_lf_flag []byte = []byte{10, 32, 32, 32, 32, 92}

type prefab struct{}

var instance *prefab
var once sync.Once

func GetInstance() *prefab {
	once.Do(func() {
		instance = &prefab{}
	})
	return instance
}

func (p *prefab) cleanWrap(text []byte) []byte {
	text = bytes.Replace(text, const_cr_lf_flag, []byte{}, -1)
	text = bytes.Replace(text, const_cr_flag, []byte{}, -1)
	text = bytes.Replace(text, const_lf_flag, []byte{}, -1)
	return text
}

func (p *prefab) uc2hanzi(uc string) (string, error) {
	val2int, err := strconv.ParseInt(uc, 16, 32)
	if err != nil {
		return uc, err
	}
	return fmt.Sprintf("%c", val2int), nil
}

func (p *prefab) filter(text []byte) bool {
	for i := 0; i < len(text); i++ {
		if text[i]&0x80 != 0 {
			return false
		}
	}
	return true
}

func (p *prefab) GetString(text []byte) ([][]byte, error) {
	var cnEntry [][]byte
	tag := fmt.Sprintf("%c%c", sl, uu)
	frecord := func(start, end int) {
		unicode := string(p.cleanWrap(text[start : end+1]))
		index := strings.Index(unicode, tag)
		for ; index != -1; index = strings.Index(unicode, tag) {
			hanzi, err := p.uc2hanzi(unicode[index+2 : index+6])
			if err != nil {
				panic(err)
			}
			unicode = strings.Replace(unicode, unicode[index:index+6], hanzi, 1)
		}
		slice := []byte(unicode)
		if !p.filter(slice) {
			cnEntry = append(cnEntry, []byte(unicode))
		}
	}
	nState := state_normal
	nStateStart := 0
	nSize := len(text)
	for i := 0; i < nSize; i++ {
		if text[i] == dq && i >= 7 && bytes.Compare(text[i-7:i-1], const_string_flag) == 0 {
			nStateStart = i + 1
			nState = state_double_quotes
			continue
		}
		switch nState {
		case state_double_quotes:
			if text[i] == dq {
				frecord(nStateStart, i-1)
				nState = state_normal
			}
		}
	}
	if nState != state_normal {
		return cnEntry, errors.New(fmt.Sprintf("%s state:%d", "file syntax error", nState))
	}
	return cnEntry, nil
}

func (p *prefab) ReplaceOnce(context *[]byte, sText []byte, trans []byte) error {
	prefabformat := func(s string) string {
		length := len(s)
		for i := 0; i+5 < length; i++ {
			if s[i] == sl && s[i+1] == uu {
				upper := strings.ToUpper(s[i+2 : i+6])
				s = strings.Replace(s, s[i+2:i+6], upper, 1)
			}
		}
		return s
	}

	textQuoted := strconv.QuoteToASCII(string(sText))
	textUnquoted := prefabformat(textQuoted[1 : len(textQuoted)-1])
	textUnquoted = strings.Replace(textUnquoted, "\\\\", "\\", -1)

	nState := state_normal
	nStateStart := 0
	text := *context
	nSize := len(text)
	found := false
	var sTextReal []byte
	for i := 0; i < nSize && !found; i++ {
		if text[i] == dq && i >= 7 && bytes.Compare(text[i-7:i-1], const_string_flag) == 0 {
			nStateStart = i + 1
			nState = state_double_quotes
			continue
		}
		switch nState {
		case state_double_quotes:
			if text[i] == dq {
				unicode := p.cleanWrap(text[nStateStart:i])
				if bytes.EqualFold(unicode, []byte(textUnquoted)) {
					sTextReal = text[nStateStart:i]
					found = true
				}
				nState = state_normal
			}
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("[can not find %s]", sText))
	}
	transQuoted := strconv.QuoteToASCII(string(trans))
	transUnquoted := prefabformat(transQuoted[1 : len(transQuoted)-1])
	transUnquoted = strings.Replace(transUnquoted, "\\\\", "\\", -1)
	*context = bytes.Replace(*context, sTextReal, []byte(transUnquoted), 1)
	return nil
}
