package internal_test

import (
	"aytovav/logAnalizer/internal"
	"strings"
	"testing"
)

var str = `CLIENT: 10.149.100.4
USER: wsrest
REQUEST HEADERS:
        host: naucc.abb-win.akbars.ru:8443
        authorization: Basic d3NyZXN0OkxpaXNlZXI3YWlrZWkvYWNvaHF1ZWV0aQ==
        x-request-id: b2671cbb084940a2be6fbbe5ba133740
        x-operation-id: cc0c2a9848914d3daa7a4fa48e851e31
        x-ext-request-id: 5d3260940e334ea9b19e37434caf1f42
        x-origin: Abdt.Aimee.NaumenAdapter
        traceparent: 00-89215379dd43a0c8ed07198b0b86aeb2-54c50e0995152f3a-00
        content-type: application/json; charset=utf-8
        content-length: 532
RESPONSE STATUS: 201
RESPONSE HEADERS:
        Vary: Origin
        Vary: Origin
        Vary: Origin
        Location: https://naucc.abb-win.akbars.ru:8443/api/v2/projects/corebo00000000000om63bmce5j5en44/cases/ocpcas00000000000p4jdccq3hqicf1c
        X-Content-Type-Options: nosniff
        X-XSS-Protection: 0
        Cache-Control: no-cache, no-store, max-age=0, must-revalidate
        Pragma: no-cache
        Expires: 0
        Strict-Transport-Security: max-age=31536000 ; includeSubDomains
        Content-Type: application/json
        Transfer-Encoding: chunked
        Date: Wed, 02 Oct 2024 14:39:56 GMT`

func TestParseHeaders(t *testing.T) {
	data := strings.Split(str, "\n")
	expected := internal.ParsedHeaders{
		StatusCode: "201",
		CaseID:     "ocpcas00000000000p4jdccq3hqicf1c",
	}

	res := internal.ParseHeaders(data)
	if res != expected {
		t.Error()
	}
}

func TestCaseCreating(t *testing.T) {
	dataFile := "sample_caseCreating.txt"
	expected := []internal.ParsedData{
		{
			RequestID: "5605257630",
			Date:      "2024-10-02 17:39:57,092",
			ProjectID: "corebo00000000000om63bmce5j5en44",
			ParsedHeaders: internal.ParsedHeaders{
				StatusCode: "201",
				CaseID:     "ocpcas00000000000p4jdccq3hqicf1c",
			},
		},
		{
			RequestID: "5605081619",
			Date:      "2024-10-02 17:37:01,081",
			ProjectID: "corebo00000000000om63bmce5j5en44",
			ParsedHeaders: internal.ParsedHeaders{
				StatusCode: "201",
				CaseID:     "ocpcas00000000000p4jd9mrs283s8m8",
			},
		},
		{
			RequestID: "5605087673",
			Date:      "2024-10-02 17:37:07,135",
			ProjectID: "corebo00000000000om63bmce5j5en44",
			ParsedHeaders: internal.ParsedHeaders{
				StatusCode: "201",
				CaseID:     "ocpcas00000000000p4jd9pqi22gb68k",
			},
		},
	}

	dataChannel := make(chan internal.ParsedData)
	go internal.CaseCreating(dataFile, dataChannel)

	result := []internal.ParsedData{}
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
