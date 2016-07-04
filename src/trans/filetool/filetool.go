package filetool

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

func (ft *filetool) GetEncodeString(file string) string {
	file_ex := path.Ext(file)
	codingstring, ok := ft.file2coding[file_ex]
	if ok {
		return codingstring
	} else {
		return "Nil"
	}
}

func (ft *filetool) GetEncoding(file string) encoding.Encoding {
	codingstring := ft.GetEncodeString(file)
	coding, ok := ft.encodingmap[codingstring]
	if !ok {
		return nil
	}
	return coding
}

func (ft *filetool) SetEncoding(file, codingstring string) (string, error) {
	var oldstring string
	file_ex := path.Ext(file)
	if len(codingstring) > 0 {
		if _, ok := ft.encodingmap[codingstring]; !ok {
			return oldstring, errors.New(fmt.Sprintf("encoding not exsit [%s] %s", codingstring, file))
		}
		oldstring, _ = ft.file2coding[file_ex]
		ft.file2coding[file_ex] = codingstring
		return oldstring, nil
	} else {
		oldstring, _ = ft.file2coding[file_ex]
		delete(ft.file2coding, file_ex)
		return oldstring, nil
	}
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
			return nil, errors.New(fmt.Sprintf("%s %s %s", err.Error(), name, ft.GetEncodeString(name)))
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
			return errors.New(fmt.Sprintf("%s %s %s", err.Error(), name, ft.GetEncodeString(name)))
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
		temp := bytes.TrimSpace(line)
		if len(temp) > 0 {
			if coding != nil {
				reader := transform.NewReader(bytes.NewReader(line), coding.NewDecoder())
				dline, err := ioutil.ReadAll(reader)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("%s %s %s", err.Error(), name, ft.GetEncodeString(name)))
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
			temp := bytes.TrimSpace(v)
			if len(temp) > 0 {
				if coding != nil {
					reader := transform.NewReader(bytes.NewReader(v), coding.NewEncoder())
					ev, err := ioutil.ReadAll(reader)
					if err != nil {
						return errors.New(fmt.Sprintf("%s %s %s", err.Error(), name, ft.GetEncodeString(name)))
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
