package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const maxLocationNumber = 15

var locationRgx = regexp.MustCompile(`([a-zA-Z]+)(:|\*|\_)(\+|\-)([0-9]{3})`)
var recordSQLState = regexp.MustCompile(`^\d+\.\d+ `)

var recordLog = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4} \d{2}\:\d{2}\:\d{2}\: `)
var logActionRgx = regexp.MustCompile(`Client (disconnected|connected)\: (\d{1,3}\.\d{1,3}\.\d{1,3}.\d{1,3}\:\d{4,5})`)

// Location represents node location
type Location struct {
	Name          string `json:"name"`
	Status        string `json:"status"`
	Positive      bool   `json:"positive"`
	PositiveValue string `json:"positive_value"`
	Value         string `json:"value"`
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

	// if recordLog.Match(input) {
	// 	return parseLog(input)
	// }

	return parseRaw(input)

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
	matches := locationRgx.FindAllSubmatch(input, -1)
	c := 0
	for _, loc := range matches {
		// Max 15 location!
		if c == maxLocationNumber {
			break
		}
		res.Locations = append(res.Locations, Location{
			Name:          string(loc[1]),
			Status:        string(loc[2]),
			Positive:      loc[3][0] == byte(43), // 43 is +
			PositiveValue: string(loc[3][0]),
			Value:         string(loc[4]),
		})
		c++
	}

	encodedRes, err := json.Marshal(res)
	if err != nil {
		return []byte{}, fmt.Errorf("can't marshal: %v", err)
	}
	return encodedRes, nil
}

// LogAction represents log entries
type LogAction struct {
	Action string `json:"action"`
	Client string `json:"client"`
}

// Parse takes raw bytes from /tmp/sql_state and convert them in json
func parseLog(input []byte) ([]byte, error) {

	match := logActionRgx.FindSubmatch(input)
	if len(match) > 0 {
		la := LogAction{
			Action: string(match[1]),
			Client: string(match[2]),
		}

		encodedRes, _ := json.Marshal(la)
		return encodedRes, nil
	}
	return []byte{}, nil
}

// RawAction represents log entries
type RawAction struct {
	Raw string `json:"raw"`
}

// Parse takes raw bytes from /tmp/sql_state and convert them in json
func parseRaw(input []byte) ([]byte, error) {

	la := RawAction{
		Raw: string(input),
	}

	return json.Marshal(la)

}
