package main

import (
	"aytovav/logAnalizer/internal"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	dataChannel := make(chan internal.ParsedData)
	go internal.CaseCreating("logs_analizer/catalina_2.out", dataChannel)

	err := internal.CaseCreatingInfoWrite("logs_analizer/result_case_creating.txt", dataChannel)
	if err != nil {
		fmt.Println(err)
		return
	}
}
