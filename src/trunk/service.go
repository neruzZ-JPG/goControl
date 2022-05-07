package trunk

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var processMap map[string]string

func initAllProcesses() error {
	var err error
	processMap, err = readFile("reel")
	return err
}

func StartService() {
	err := initAllProcesses()
	if err != nil && os.IsNotExist(err) {
		//if the reel does not exist
		if _, erro := os.Create("reel"); erro != nil {
			fmt.Println("cannot find the file")
			return
		} else {
			initAllProcesses()
		}
	}

}

func QueryProcesses() {
	for key := range processMap {
		fmt.Println(key)
	}
}

func RegisterProcess(name, path string) error {
	//first see if the process already registered
	if _, exist := processMap[name]; exist {
		return fmt.Errorf("[registerProcess] process already registered")
	}
	//see if the path is correct(we can create)
	//!!!dont use open ! because open is readonly!!!
	if _, err := os.Open(path); err != nil && os.IsExist(err) {
		return fmt.Errorf("[regsiterProcess] config file already exists for this process")
	}
	//create file
	if _, err := os.Create(path); err != nil {
		return fmt.Errorf("[registerProcess] error while creating file: %v", err)
	}
	processMap[name] = path
	if err := SavingChanges(); err != nil {
		return err
	}
	return nil
}

func SavingChanges() error {
	return saveFile(processMap, "reel")
}

func cancel(name string) error {
	//check if registered
	if _, exist := processMap[name]; !exist {
		return fmt.Errorf("[cancel]process not registered")
	}
	if err := os.Remove(processMap[name]); err != nil {
		return err
	}
	delete(processMap, name)
	fmt.Print(processMap)
	SavingChanges()
	return nil
}

func selectProcess(name string) (map[string]string, error) {
	//existence check
	path, exist := processMap[name]
	if !exist {
		return nil, fmt.Errorf("[selectProcess]not registered")
	}
	configMap, err := readFile(path)
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

func changeConfig(name string, configMap map[string]string, key, value string) error {
	configMap[key] = value
	path := processMap[name]
	return saveFile(configMap, path)
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

func saveFile(ssMap map[string]string, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	str := ""
	for k, v := range ssMap {
		str += fmt.Sprintf("%s:%s\n", k, v)
	}
	if str != "" {
		str = str[:len(str)-1]
	}
	if _, err = file.WriteString(str); err != nil {
		return err
	}
	return nil
}
