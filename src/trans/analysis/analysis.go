package analysis

import (
	"bytes"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"
	"trans/analysis/lua"
	"trans/analysis/prefab"
	"trans/analysis/tabfile"
	"trans/dic"
	"trans/filetool"
	"trans/gpool"
	"trans/log"
)

type delegate interface {
	GetString(text []byte) ([][]byte, error)
	ReplaceOnce(context *[]byte, sText []byte, trans []byte) error
}

type analysis struct {
	rulesMap  map[string]string
	filterMap map[string]bool
}

const (
	const_rule_lua       = "lua_rules"
	const_rule_prefab    = "prefab_rules"
	const_rule_tablefile = "table_rules"
)

var instance *analysis
var once sync.Once

func GetInstance() *analysis {
	once.Do(func() {
		instance = &analysis{
			rulesMap:  make(map[string]string),
			filterMap: make(map[string]bool),
		}
	})
	return instance
}

func (a *analysis) SetRulesMap(k, v string) {
	a.rulesMap[path.Ext(k)] = v
}

func (a *analysis) SetFilterMap(key string) {
	a.filterMap[key] = true
}

func (a *analysis) getPool(file string) (delegate, error) {
	file_ex := path.Ext(file)
	rule, ok := a.rulesMap[file_ex]
	if !ok {
		return nil, errors.New(fmt.Sprintf("[not extract rule] %s", file))
	}
	switch rule {
	case const_rule_lua:
		return lua.GetInstance(), nil
	case const_rule_prefab:
		return prefab.GetInstance(), nil
	case const_rule_tablefile:
		return tabfile.GetInstance(), nil
	default:
		return nil, errors.New(fmt.Sprintf("[not extract rule] %s", file))
	}
}

func (a *analysis) filter(name string) error {
	namev := strings.Split(name, "/")
	for _, filename := range namev {
		if _, ok := a.filterMap[filename]; ok {
			return errors.New(fmt.Sprintf("[ingnore file] %s", name))
		}
	}
	return nil
}

func (a *analysis) GetString(dbname, root string) {
	root = strings.TrimRight(strings.Replace(root, "\\", "/", -1), "/")
	log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, fmt.Sprintf("extract chinese from %s", root))
	ft := filetool.GetInstance()
	fmap, err := ft.GetFilesMap(root)
	if err != nil {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
		return
	}
	var relativepath_total []string
	var entry_total [][][]byte
	for i := 0; i < len(fmap); i++ {
		if err := a.filter(fmap[i]); err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
			continue
		}
		ins, err := a.getPool(fmap[i])
		if err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
			continue
		}
		context, err := ft.ReadAll(fmap[i])
		if err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
			continue
		}
		entry, err := ins.GetString(context)
		if err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
		}
		var entry_norepeat [][]byte
		for i := 0; i < len(entry); i++ {
			bNoRepeat := true
			for j := 0; j < i; j++ {
				if bytes.Compare(entry[i], entry[j]) == 0 {
					bNoRepeat = false
					break
				}
			}
			if bNoRepeat {
				entry_norepeat = append(entry_norepeat, entry[i])
			}
		}
		if len(entry_norepeat) > 0 {
			relaticepath := strings.TrimLeft(strings.Split(fmap[i], root)[1], "/")
			if len(relaticepath) == 0 {
				relaticepath = path.Base(fmap[i])
			}
			relativepath_total = append(relativepath_total, relaticepath)
			entry_total = append(entry_total, entry_norepeat)
		}
	}
	newcount := 0
	db := dic.New(dbname)
	for i := 0; i < len(relativepath_total); i++ {
		path := relativepath_total[i]
		entry := entry_total[i]
		for j := 0; j < len(entry); j++ {
			if _, ok := db.Query(path, entry[j]); !ok {
				db.Append(path, entry[j])
				newcount += 1
			}
		}
	}
	if newcount > 0 {
		db.Save()
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO,
			fmt.Sprintf("generate %s, new line number: %d. finished!", dbname, newcount))
	} else {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO,
			fmt.Sprintf("nothing to do. finished!"))
	}
}

func (a *analysis) Translate(dbname, root, output string, queue int) {
	root = strings.TrimRight(strings.Replace(root, "\\", "/", -1), "/")
	output = strings.TrimRight(strings.Replace(output, "\\", "/", -1), "/")
	log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, fmt.Sprintf("translate %s to %s", root, output))
	ft := filetool.GetInstance()
	fmap, err := ft.GetFilesMap(root)
	if err != nil {
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
		return
	}
	db := dic.New(dbname)
	tatal, transcount, newcount := 0, 0, 0
	pool := gpool.New(queue)
	mutex := &sync.Mutex{}
	fwork := func(oldfile, newfile, relative string) {
		defer pool.Done()
		var entry [][]byte
		bv, err := ft.ReadAll(oldfile)
		if err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
			return
		}
		ins, err := a.getPool(oldfile)
		if err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
			goto Point
		}
		if err = a.filter(oldfile); err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO, err)
			goto Point
		}
		entry, err = ins.GetString(bv)
		if err != nil {
			log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
			goto Point
		}
		for _, v := range entry {
			trans, ok := db.Query(relative, v)
			if !ok {
				mutex.Lock()
				db.Append(relative, v)
				newcount += 1
				mutex.Unlock()
				continue
			}
			if len(trans) > 0 {
				if err := ins.ReplaceOnce(&bv, v, trans); err != nil {
					log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_ERROR, err)
				}
			}
		}
		transcount += 1
	Point:
		tatal += 1
		ft.WriteAll(newfile, bv)
	}
	for i := 0; i < len(fmap); i++ {
		pool.Add(1)
		fpath := strings.Replace(fmap[i], root, output, 1)
		frelative := strings.TrimLeft(strings.Split(fmap[i], root)[1], "/")
		if len(frelative) == 0 {
			frelative = path.Base(fmap[i])
		}
		go fwork(fmap[i], fpath, frelative)
	}
	pool.Wait()
	if newcount > 0 {
		db.Save()
		log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO,
			fmt.Sprintf("generate %s, new line number: %d.", dbname, newcount))
	}
	log.WriteLog(log.LOG_FILE|log.LOG_PRINT, log.LOG_INFO,
		fmt.Sprintf("translate file %d, copy file %d. finished!", transcount, tatal-transcount))
	return
}
