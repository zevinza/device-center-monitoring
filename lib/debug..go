package lib

import (
	"api/constant"
	"encoding/json"
	"fmt"
	"runtime"
)

var debugger = constant.GetEnvironment() == constant.EnvironmentDevelopment || constant.GetEnvironment() == constant.EnvironmentLocal

func Debug(v ...any) {
	if debugger {
		if _, file, line, ok := runtime.Caller(1); ok {
			fmt.Printf("\033[32m %s line:%d \n \033[0m", file, line)
		}
		fmt.Println(v...)
	}
}

func DebugJson(js ...any) {
	if debugger {
		if _, file, line, ok := runtime.Caller(1); ok {
			fmt.Printf("\033[32m %s line:%d \n \033[0m", file, line)
		}
		by, _ := json.MarshalIndent(js, "", "   ")
		fmt.Println(string(by))
	}
}
