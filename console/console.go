package console

import (
	"flag"
	"fmt"
	"strings"
)

func Handler(cmd string) {
	fields := strings.Fields(cmd)
	switch fields[0] {
	case "command":
		command(cmd)
	case "func":
		//function(cmd)
	case "cfg":
	case "":

	}
}

//func function(cmd string) {
//	fields := strings.Fields(cmd)
//
//	switch fields[1] {
//	case "reboot":
//		sendCommand("shutdown -s -t 0")
//	case "req_host_name":
//
//	}
//	fs := flag.NewFlagSet("func", flag.ContinueOnError)
//
//	r := fs.String("r", "World", "Name to greet")
//	t := fs.Int("t", 0, "Age of the person")
//
//	err := fs.Parse(strings.Split(flagStr, " ")[1:])
//	if err != nil {
//		fmt.Println("Error parsing flags:", err)
//		return
//	}
//}

func command(cmd string) {
	fs := flag.NewFlagSet("cmd", flag.ContinueOnError)

	name := fs.String("r", "World", "Name to greet")
	age := fs.Int("t", 0, "Age of the person")

	err := fs.Parse(strings.Split(cmd, " ")[1:])
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}

	// 输出解析结果
	fmt.Printf("Hello, %s!\n", *name)
	if *age > 0 {
		fmt.Printf("You are %d years old.\n", *age)
	} else {
		fmt.Println("Age not provided.")
	}
}
