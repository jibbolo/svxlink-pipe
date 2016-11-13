package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const maxLocationNumber = 15

var locationRgx = regexp.MustCompile(`([a-zA-Z]+)(:|\*|\_)(\+|\-)([0-9]{3})`)

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

// Parse takes raw bytes from /tmp/sql_state and convert them in Results
func Parse(input []byte) (*Result, error) {
	sec, err := strconv.Atoi(string(input[0:10]))
	if err != nil {
		return nil, fmt.Errorf("can't find time `%s: %s", input, err)
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
	return res, nil
}
