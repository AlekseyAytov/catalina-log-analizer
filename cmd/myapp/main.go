package main

import (
	"aytovav/logAnalizer/internal/interactionSaving"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	// dataChannel := make(chan caseCreating.ParsedData)
	// go caseCreating.CaseCreating("logs_analizer/catalina_2.out", dataChannel)

	// err := caseCreating.CaseCreatingInfoWrite("logs_analizer/result_case_creating.txt", dataChannel)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	dataChannel := make(chan interactionSaving.ParsedData)
	go interactionSaving.InterSaving("logs_analizer/catalina_2.out", dataChannel)

	err := interactionSaving.InterSavingInfoWrite("logs_analizer/result_ineraction_saving.txt", dataChannel)
	if err != nil {
		fmt.Println(err)
		return
	}
}
