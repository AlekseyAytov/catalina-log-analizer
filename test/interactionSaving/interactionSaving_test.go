package interactionsaving_test

import (
	"aytovav/logAnalizer/internal/interactionSaving"
	"fmt"
	"os"
	"strings"
	"testing"
)

var str = `CLIENT: 10.149.100.48
USER: ZijatdinovaRR
REQUEST HEADERS:
        host: naucc.abb-win.akbars.ru:8443
        authorization: Basic d3NyZXN0OkxpaXNlZXI3YWlrZWkvYWNvaHF1ZWV0aQ==
        run_as_user: ZijatdinovaRR
        x-request-id: 334ce928e268491ebfd1c478c90e8853
        x-operation-id: 38e4b3fcc7ec471fb6508d37de95f07d
        x-ext-request-id: f77474352541453dbf3ce6cb388f2163
        x-origin: Abdt.Aimee.NaumenAdapter
        traceparent: 00-51e8f6222f413a739efa7ef4494d696d-abb0b4999fe3833a-00
        content-type: application/json; charset=utf-8
        content-length: 232
RESPONSE STATUS: 201
RESPONSE HEADERS:
        Vary: Origin
        Vary: Origin
        Vary: Origin
        X-Content-Type-Options: nosniff
        X-XSS-Protection: 0
        Cache-Control: no-cache, no-store, max-age=0, must-revalidate
        Pragma: no-cache
        Expires: 0
        Strict-Transport-Security: max-age=31536000 ; includeSubDomains
        Content-Type: application/json
        Content-Language: ru-RU
        Transfer-Encoding: chunked
        Date: Mon, 25 Nov 2024 00:19:51 GMT`

func TestParseHeaders(t *testing.T) {
	data := strings.Split(str, "\n")
	expected := interactionSaving.ParsedHeaders{
		StatusCode: "201",
		User:       "ZijatdinovaRR",
	}

	res := interactionSaving.ParseHeaders(data)
	if res != expected {
		t.Error()
	}
}

func TestInteractionSaving(t *testing.T) {
	dataFile := "sample_interactionSaving.txt"
	expected := []interactionSaving.ParsedData{
		{
			RequestID: "279021683",
			Date:      "2024-11-25 03:19:51,301",
			ProjectID: "corebo00000000000ng9hgl5vhtinhmo",
			SessionID: "node_1_domain_0_nauss_0_1732214988_47691",
			ParsedHeaders: interactionSaving.ParsedHeaders{
				StatusCode: "201",
				User:       "ZijatdinovaRR",
			},
		},
		{
			RequestID: "279021755",
			Date:      "2024-11-25 03:19:51,373",
			ProjectID: "corebo00000000000ogs9g8e2307ag4k",
			SessionID: "node_1_domain_0_nauss_0_1732214988_47691",
			ParsedHeaders: interactionSaving.ParsedHeaders{
				StatusCode: "405",
				User:       "ZijatdinovaRR",
			},
		},
		{
			RequestID: "281775493",
			Date:      "2024-11-25 04:05:45,111",
			ProjectID: "corebo00000000000ng9hgl5vhtinhmo",
			SessionID: "node_1_domain_2_nauss_0_1732214988_47694",
			ParsedHeaders: interactionSaving.ParsedHeaders{
				StatusCode: "201",
				User:       "MrX",
			},
		},
	}

	inF, err := os.Open(dataFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inF.Close()

	dataChannel := make(chan interactionSaving.ParsedData)
	go interactionSaving.InterSaving(inF, dataChannel)

	result := []interactionSaving.ParsedData{}
	for val := range dataChannel {
		result = append(result, val)
	}

	if len(result) != len(expected) {
		t.Error("Not equal lenth")
	}
	for i := range result {
		if expected[i] != result[i] {
			t.Errorf("Not equal %d element", i)
		}
	}
}
