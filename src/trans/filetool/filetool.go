package filetool

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type filetool struct {
	encodingmap map[string]encoding.Encoding
	file2coding map[string]string //默认以UTF8编码方式读写
}

var instance *filetool
var once sync.Once

func GetInstance() *filetool {
	once.Do(func() {
		instance = &filetool{
			encodingmap: map[string]encoding.Encoding{
				"utf8":      unicode.UTF8,
				"gbk":       simplifiedchinese.GBK,
				"hz-gb2312": simplifiedchinese.HZGB2312,
				"gb18030":   simplifiedchinese.GB18030,
				"big5":      traditionalchinese.Big5,
			},
			file2coding: make(map[string]string),
		}
	})
	return instance
}

func (ft *filetool) GetEncoding(file string) encoding.Encoding {
	filev := strings.Split(file, ".")
	file_ex := filev[len(filev)-1]
	codingstring, ok := ft.file2coding[file_ex]
	if !ok {
		return nil
	}
	coding, ok := ft.encodingmap[codingstring]
	if !ok {
		return nil
	}
	return coding
}

func (ft *filetool) SetEncoding(file, codingstring string) error {
	filev := strings.Split(file, ".")
	file_ex := filev[len(filev)-1]
	if _, ok := ft.encodingmap[codingstring]; !ok {
		return errors.New(fmt.Sprintf("encoding not exsit [%s] %s", codingstring, file))
	}
	ft.file2coding[file_ex] = codingstring
	return nil
}

func (ft *filetool) GetFilesMap(path string) (map[int]string, error) {
	index := 0
	filemap := make(map[int]string)
	_, err := os.Stat(path)
	if err != nil {
		return filemap, err
	}
	f := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filemap[index] = strings.Replace(path, "\\", "/", -1)
			index++
			return err
		} else {
			return nil
		}
	}
	fpErr := filepath.Walk(path, f)
	if fpErr != nil {
		return nil, errors.New(fmt.Sprintf("filepath.Walk Failed! %s", path))
	}
	return filemap, nil
}

func (ft *filetool) ReadAll(name string) ([]byte, error) {
	context, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	coding := ft.GetEncoding(name)
	if coding != nil {
		reader := transform.NewReader(bytes.NewReader(context), coding.NewDecoder())
		dcontext, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("%s %s", err.Error(), name))
		}
		return dcontext, nil
	}
	return context, nil
}

func (ft *filetool) WriteAll(name string, context []byte) error {
	if index := strings.LastIndex(name, "/"); index != -1 {
		err := os.MkdirAll(name[:index], os.ModePerm)
		if err != nil {
			return err
		}
	}
	coding := ft.GetEncoding(name)
	if coding != nil {
		reader := transform.NewReader(bytes.NewReader(context), coding.NewEncoder())
		econtext, err := ioutil.ReadAll(reader)
		if err != nil {
			return errors.New(fmt.Sprintf("%s %s", err.Error(), name))
		}
		return ioutil.WriteFile(name, econtext, os.ModePerm)
	}
	return ioutil.WriteFile(name, context, os.ModePerm)
}

func (ft *filetool) ReadFileLine(name string) ([][]byte, error) {
	var context [][]byte
	f, err := os.Open(name)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	readline := func(r *bufio.Reader) ([]byte, error) {
		var (
			isPrefix        bool  = true
			err             error = nil
			line, realyline []byte
		)
		for isPrefix && err == nil {
			line, isPrefix, err = r.ReadLine()
			realyline = append(realyline, line...)
		}
		return realyline, err
	}
	coding := ft.GetEncoding(name)
	r := bufio.NewReader(f)
	err = nil
	var line []byte
	for err == nil {
		line, err = readline(r)
		line = bytes.Trim(line, " ")
		if len(line) > 0 {
			if coding != nil {
				reader := transform.NewReader(bytes.NewReader(line), coding.NewDecoder())
				dline, err := ioutil.ReadAll(reader)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("%s %s", err.Error(), name))
				}
				context = append(context, dline)
			} else {
				context = append(context, line)
			}
		}
	}
	return context, nil
}

func (ft *filetool) SaveFileLine(name string, context [][]byte) error {
	f, err := os.Create(name)
	defer f.Close()
	if err != nil {
		return err
	}
	coding := ft.GetEncoding(name)
	w := bufio.NewWriter(f)
	length := len(context)
	if length >= 1 {
		for _, v := range context[:length] {
			v = bytes.Trim(v, " ")
			if len(v) > 0 {
				if coding != nil {
					reader := transform.NewReader(bytes.NewReader(v), coding.NewEncoder())
					ev, err := ioutil.ReadAll(reader)
					if err != nil {
						return errors.New(fmt.Sprintf("%s %s", err.Error(), name))
					}
					fmt.Fprintln(w, string(ev))
				} else {
					fmt.Fprintln(w, string(v))
				}
			}
		}
	}
	return w.Flush()
}
