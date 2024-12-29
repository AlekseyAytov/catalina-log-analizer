package internal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
)

// ParsedHeaders содержит данные блока headers
type ParsedHeaders struct {
	StatusCode string
	CaseID     string
}

// ParsedData содержит итоговые данные
type ParsedData struct {
	ParsedHeaders
	ProjectID string
	RequestID string
	Date      string
}

var pattern, _ = regexp.Compile(`^(?P<id>\d*) (?P<date>.*) \[.*POST.*\/projects\/(?P<project_id>.*)\/cases`)
var endPattern, _ = regexp.Compile(`^\d* `)
var statusPattern, _ = regexp.Compile(`RESPONSE STATUS: (?P<status>\d{3})`)
var caseIDPattern, _ = regexp.Compile(`.*\/projects\/.*\/cases\/(?P<case_id>.*)`)

// ParseHeaders парсит блок данных headers и возвращает заполненную структуру
func ParseHeaders(str []string) (headers ParsedHeaders) {
	for _, s := range str {
		s = strings.Trim(s, " ")
		status := statusPattern.FindAllStringSubmatch(s, -1)
		if status != nil {
			headers.StatusCode = status[0][1]
		}
		caseID := caseIDPattern.FindAllStringSubmatch(s, -1)
		if caseID != nil {
			headers.CaseID = caseID[0][1]
		}
	}
	return
}

// ParseData анализирует строку данных на паттерн запроса на создание кейса
// возвращает булевый флаг и структуру
func ParseData(row string) (bool, ParsedData) {
	parsed := ParsedData{}
	// поиск паттерна запроса на создание кейса
	caseCreatingReq := pattern.FindAllStringSubmatch(row, -1)
	if caseCreatingReq != nil {

		// когда паттерн найден, парсим возожные данные со строки
		for _, v := range caseCreatingReq {
			for kk, vv := range pattern.SubexpNames() {
				if vv == "id" {
					parsed.RequestID = v[kk]
				}
				if vv == "date" {
					parsed.Date = v[kk]
				}
				if vv == "project_id" {
					parsed.ProjectID = v[kk]
				}
			}
		}
		return true, parsed
	}
	return false, parsed
}

// CaseCreating производит построчный анализ файла и передает в канал найденные данные
func CaseCreating(fileName string, out chan<- ParsedData) {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, _ := f.Stat()
	fmt.Printf("Parsing file %s, it's size %d Mb\n", fileInfo.Name(), fileInfo.Size()/1024/1024)
	defer f.Close()

	// result := make([]ParsedData, 0, 100)
	var isBlockOfData bool
	headerText := make([]string, 0, 30)
	parsed := ParsedData{}

	scanner := bufio.NewScanner(f)
	// newBuf := make([]byte, 100*1024)
	// scanner.Buffer(newBuf, 0)
	printMemUsage()
	for scanner.Scan() {

		row := strings.Trim(scanner.Text(), " ")

		if isBlockOfData {
			// поиск паттерна конца блока данных для анализа
			end := endPattern.FindAllStringSubmatch(row, -1)
			if end != nil {
				parsed.ParsedHeaders = ParseHeaders(headerText)
				out <- parsed

				// обнуление переменных
				isBlockOfData = false
				headerText = make([]string, 0, 30)
				parsed = ParsedData{}

				// поиск паттерна запроса на создание кейса
				// так как строка конца блока данных может быть началом запроса на создание кейса
				if ok, rowData := ParseData(row); ok {
					isBlockOfData = true
					parsed = rowData
				}
			} else {
				headerText = append(headerText, row)
			}
		} else {
			// поиск паттерна запроса на создание кейса
			if ok, rowData := ParseData(row); ok {
				isBlockOfData = true
				parsed = rowData
			}
		}
	}
	printMemUsage()
	close(out)
}

// CaseCreatingInfoWrite получает из канала итоговую структуру и записыввает в файл
func CaseCreatingInfoWrite(fileName string, in <-chan ParsedData) error {
	outF, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outF.Close()

	couter := 0
	for p := range in {
		_, err := outF.WriteString(fmt.Sprintf("%s %s - %s : %s %s\n", p.RequestID, p.Date, p.ProjectID, p.StatusCode, p.CaseID))
		if err != nil {
			return err
		}
		couter++
	}
	fmt.Printf("Successfully written %d rows to file: %s\n", couter, outF.Name())
	return nil
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
