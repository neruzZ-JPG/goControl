package branch

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type Branch struct {
	configMap map[string]string
	mutex     sync.RWMutex
}

var branch Branch

func Start(path string) {
	configMap, err := readFile(path)
	if err != nil {
		fmt.Print("[goControl]error read file")
		return
	}
	branch = Branch{
		configMap: configMap,
		mutex:     sync.RWMutex{},
	}
	go func() {
		for {
			time.Sleep(60 * time.Second)
			configMap, err := readFile(path)
			if err != nil {
				fmt.Print("[goControl]error read file")
				return
			}
			branch.mutex.Lock()
			branch.configMap = configMap
			branch.mutex.Unlock()
		}

	}()
}

func GetConfig(key string) (value string, err error) {
	branch.mutex.RLock()
	value, exist := branch.configMap[key]
	branch.mutex.RUnlock()
	if !exist {
		err = fmt.Errorf("[getConfig] no such config")
		return
	}
	err = nil
	return
}

func readFile(path string) (map[string]string, error) {
	resMap := make(map[string]string, 0)
	file, err := os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if string(data) == "" {
		return resMap, nil
	}
	rows := strings.Split(string(data), "\n")
	for _, row := range rows {
		file := strings.Split(row, ":")
		resMap[file[0]] = file[1]
	}
	return resMap, nil
}
