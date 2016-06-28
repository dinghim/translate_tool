package prefab

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

var const_clean_flag1 [][]byte = [][]byte{
	{13, 10, 32, 32, 32, 32, 92, 32},
	{13, 32, 32, 32, 32, 92, 32},
	{10, 32, 32, 32, 32, 92, 32},
}

var const_clean_flag2 [][]byte = [][]byte{
	{13, 10, 32, 32, 32, 32},
	{13, 32, 32, 32, 32},
	{10, 32, 32, 32, 32},
}

type prefab struct{}

func New() *prefab {
	return &prefab{}
}

func (p *prefab) cleanWrap(text []byte) []byte {
	for _, v := range const_clean_flag1 {
		text = bytes.Replace(text, v, []byte{32}, -1)
	}
	for _, v := range const_clean_flag2 {
		text = bytes.Replace(text, v, []byte{}, -1)
	}
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

func (p *prefab) GetString(context []byte) ([][]byte, []int, []int, error) {
	var entryStart []int
	var entryEnd []int
	var entryTotal [][]byte
	tag := fmt.Sprintf("%c%c", sl, uu)
	frecord := func(start, end int) {
		unicode := string(p.cleanWrap(context[start:end]))
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
			entryStart = append(entryStart, start)
			entryEnd = append(entryEnd, end)
			entryTotal = append(entryTotal, []byte(unicode))
		}
	}
	nState := state_normal
	nStateStart := 0
	nSize := len(context)
	for i := 0; i < nSize; i++ {
		if context[i] == dq && i >= 7 && bytes.Compare(context[i-7:i-1], const_string_flag) == 0 {
			nStateStart = i + 1
			nState = state_double_quotes
			continue
		}
		switch nState {
		case state_double_quotes:
			if context[i] == dq {
				frecord(nStateStart, i)
				nState = state_normal
			}
		}
	}
	if nState != state_normal {
		return entryTotal, entryStart, entryEnd, errors.New(fmt.Sprintf("[file syntax error] state: %d", nState))
	}
	return entryTotal, entryStart, entryEnd, nil
}

func (p *prefab) Pretreat(trans []byte) []byte {
	sText := strconv.QuoteToASCII(string(trans))
	sText = sText[1 : len(sText)-1]
	for i := 0; i+5 < len(sText); i++ {
		if sText[i] == sl && sText[i+1] == uu {
			upper := strings.ToUpper(sText[i+2 : i+6])
			sText = strings.Replace(sText, sText[i+2:i+6], upper, 1)
		}
	}
	sText = strings.Replace(sText, "\\\\", "\\", -1)
	return []byte(sText)
}
