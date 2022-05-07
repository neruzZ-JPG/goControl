package trunk

import "fmt"

func Homepage() {
	var choice int
	showWelcome()
	fmt.Scan(&choice)
	switch choice {
	case 1:
		showProcesses()
	case 2:
		selectProcessView()
	case 3:
		registerProcessView()
	case 4:
		cancelView()
	case 0:
		exit()
	}
}

func showWelcome() {
	fmt.Println("**************")
	fmt.Println("* goControl! *")
	fmt.Println("**************")
	fmt.Println("[1] show all registered processes")
	fmt.Println("[2] select a process")
	fmt.Println("[3] register a process")
	fmt.Println("[4] cancel a process")
	fmt.Println("[0] exit")
	inputBox()
}

func showProcesses() {
	fmt.Println("showing processs~~~")
	QueryProcesses()
	fmt.Println("---------------------------------------")
	Homepage()
}

func registerProcessView() {
	fmt.Println("enter the name an path as 'name path'")
	var name, path string
	inputBox()
	fmt.Scanf("%s %s", &name, &path)
	if err := RegisterProcess(name, path); err != nil {
		fmt.Println(err)
	}
	Homepage()
}

func exit() {
	fmt.Println("saving changes~~~")
	if err := SavingChanges(); err != nil {
		fmt.Printf("[exit]err occurs: %v\n", err)
	}
	fmt.Println("下班！开摆咯")
}

func cancelView() {
	var name string
	fmt.Println("say the name")
	inputBox()
	fmt.Scan(&name)
	var confirm string
	fmt.Printf("you'll delete %s, confirm?(y/n)\n", name)
	inputBox()
	fmt.Scan(&confirm)
	if confirm == "y" {
		if err := cancel(name); err != nil {
			fmt.Println(err)
		}
	}
	Homepage()
}

func inputBox() {
	fmt.Print("> ")
}

func selectProcessView() {
	var name string
	fmt.Println("name?")
	inputBox()
	fmt.Scan(&name)
	configMap, err := selectProcess(name)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range configMap {
		fmt.Println(k + ":" + v)
	}
	var reply string
	fmt.Println("do you want to change config?(y/n)")
	inputBox()
	fmt.Scan(&reply)
	if reply == "y" {
		changeConfigView(name, configMap)
	}
	Homepage()
}

func changeConfigView(name string, configMap map[string]string) {
	var key, value string
	fmt.Println("make a wish and i will try my best")
	fmt.Println("hint : format as 'key value'")
	inputBox()
	fmt.Scanf("%s %s", &key, &value)
	changeConfig(name, configMap, key, value)
}
