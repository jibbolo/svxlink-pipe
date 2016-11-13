package parser

import (
	"fmt"
	"testing"
)

const row = "1474475227.316 Voter:sql_state Serano:+025 Assisi*+032 Bettona:-001 Lacugnano_+000"

func TestParse(t *testing.T) {
	row := []byte(row)
	res, err := Parse(row)
	if err != nil {
		t.Errorf("Error parsing: %v", err)
		t.FailNow()
	}
	fmt.Printf("%+v\n", res)

}
