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

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type filetool struct{}

var instance *filetool
var once sync.Once

func GetInstance() *filetool {
	once.Do(func() {
		instance = &filetool{}
	})
	return instance
}

func (ft *filetool) ReadFileLine(name string) ([][]byte, error) {
	var context [][]byte
	f, err := os.Open(name)
	defer f.Close()
	if err != nil {
		return context, err
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
	r := bufio.NewReader(f)
	err = nil
	var line []byte
	for err == nil {
		line, err = readline(r)
		if len(line) > 0 {
			context = append(context, line)
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
	w := bufio.NewWriter(f)
	length := len(context)
	if length < 1 {
		return w.Flush()
	} else {
		for _, v := range context[:length] {
			fmt.Fprintln(w, string(v))
		}
		return w.Flush()
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
		return nil, errors.New("Walk path Failed!")
	}
	return filemap, nil
}

func (ft *filetool) ReadAll(name string, bDecoder bool) ([]byte, error) {
	context, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	if bDecoder {
		reader := transform.NewReader(bytes.NewReader(context), simplifiedchinese.GBK.NewDecoder())
		dcontext, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		return dcontext, nil
	}
	return context, nil
}

func (ft *filetool) WriteAll(name string, text []byte, bGbkEncoder bool) error {
	if index := strings.LastIndex(name, "/"); index != -1 {
		err := os.MkdirAll(name[:index], os.ModePerm)
		if err != nil {
			return err
		}
	}
	if bGbkEncoder {
		reader := transform.NewReader(bytes.NewReader(text), simplifiedchinese.GBK.NewEncoder())
		gbktext, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(name, gbktext, os.ModePerm)
	}
	return ioutil.WriteFile(name, text, os.ModePerm)
}
