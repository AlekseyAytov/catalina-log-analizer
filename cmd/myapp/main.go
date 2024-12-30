package main

import (
	"aytovav/logAnalizer/internal/interactionSaving"
	"flag"
	"fmt"
	"os"
	"time"
)

var inFile string
var outFile string

func init() {
	const (
		defaultInFile  = "catalina.out"
		defaultOutFile = "result.txt"
		usageInFile    = "название входящего файла для анализа"
		usageOutFile   = "название исходящего файла для анализа"
	)

	flag.StringVar(&inFile, "input_file", defaultInFile, usageInFile)
	flag.StringVar(&inFile, "in", defaultInFile, usageInFile+" (shorthand)")
	flag.StringVar(&outFile, "out_file", defaultOutFile, usageOutFile)
	flag.StringVar(&outFile, "out", defaultOutFile, usageOutFile+" (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Скрипт парсит входящий файл на запросы создания кейсов по REST API\n")
		// fmt.Fprintf(flag.CommandLine.Output(), "Скрипт парсит входящий файл на запросы сохранения результата коммуникации\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Результат сохраняется в исходящий файл\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	inF, err := os.Open(inFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inF.Close()

	outF, err := os.Create(outFile)
	if err != nil {
		return
	}
	defer outF.Close()

	fileInfo, _ := inF.Stat()
	fmt.Printf("Parsing file %s, it's size %d Mb\n", fileInfo.Name(), fileInfo.Size()/1024/1024)

	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	// action caseCreating
	// go spinner()
	// dataChannel := make(chan caseCreating.ParsedData)
	// go caseCreating.CaseCreating(inF, dataChannel)

	// n, err := caseCreating.CaseCreatingInfoWrite(outF, dataChannel)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// action interactionSaving
	go spinner()
	dataChannel := make(chan interactionSaving.ParsedData)
	go interactionSaving.InterSaving(inF, dataChannel)

	n, err := interactionSaving.InterSavingInfoWrite(outF, dataChannel)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\rSuccessfully written %d rows to file: %s\n", n, outF.Name())
}

func spinner() {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(200 * time.Millisecond)
		}
	}
}
