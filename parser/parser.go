package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const maxLocationNumber = 15

var recordSQLState = regexp.MustCompile(`^\d+\.\d+ Voter:sql_state`)
var recordTXState = regexp.MustCompile(`^\d+\.\d+ Tx:state`)

// {
//     "active": true,
//     "enabled": true,
//     "id": "?",
//     "name": "Serano",
//     "siglev": 39,
//     "sql_open": true
// }

// Location represents node location
type Location struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
	SigLev int32  `json:"siglev"`
}

// Result represents a single log record from /tmp/sql_state
type Result struct {
	Time      time.Time  `json:"time"`
	Locations []Location `json:"locations"`
	Raw       string     `json:"raw"`
}

func Parse(input []byte) ([]byte, error) {

	if recordSQLState.Match(input) {
		return parseSQLState(input)
	}
	if recordTXState.Match(input) {
		return []byte{}, nil
	}
	return []byte{}, fmt.Errorf("unhandled input `%s", input)
}

// Parse takes raw bytes from /tmp/sql_state and convert them in json
func parseSQLState(input []byte) ([]byte, error) {

	sec, err := strconv.Atoi(string(input[0:10]))
	if err != nil {
		return []byte{}, fmt.Errorf("can't find time `%s: %s", input, err)
	}

	timestamp := time.Unix(int64(sec), 0)
	res := &Result{
		Time: timestamp,
		Raw:  string(input),
	}
	jsonPart := input[31:len(input)]
	if err := json.Unmarshal(jsonPart, &res.Locations); err != nil {
		return []byte{}, fmt.Errorf("can't parse json locations `%s: %s", jsonPart, err)
	}

	encodedRes, err := json.Marshal(res)
	if err != nil {
		return []byte{}, fmt.Errorf("can't marshal: %v", err)
	}
	return encodedRes, nil
}
