package main
import (
	"testing"
	"time"
	"reflect"
	"encoding/json"
)

var isValidUserTest = []struct{
	name	string
	expected bool
}{
	{"joel", true},
	{"alice", true},
	{"charles2", false},
	{"89monica", false},
}

func Test_validateUser(t *testing.T) {
	for _, u := range isValidUserTest {
		t.Run(u.name, func(t *testing.T) {
			got := validateUser(u.name)
			if got != u.expected {
				t.Errorf("Expected: %v, got %v", u.expected, got)
			}
		})
	}
}

var isValidDateTest = []struct{
	date	string
	expected	string
}{
	{"2023-05-25", ""},
	{"2023-15-25", "parsing time \"2023-15-25\": month out of range"},
	{"2023-05-40", "parsing time \"2023-05-40\": day out of range"},
}

func Test_validateDate(t *testing.T) {
	for _, d := range isValidDateTest {
		t.Run(d.date, func(t *testing.T) {
			err := validateDate(d.date)
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if errMsg != d.expected {
				t.Errorf("Expected error message %s, got %s", d.expected, errMsg)
			}
		})
	}
}
var convertStringToDateTest = []struct{
	date string
	expected string	
}{
	{"1982-05-25", ""},
	{"Monday, 02-Jan-06 15:04:05 MST", "parsing time \"Monday, 02-Jan-06 15:04:05 MST\" as \"2006-01-02\": cannot parse \"Monday, 02-Jan-06 15:04:05 MST\" as \"2006\""},
}
func Test_convertStringToDate(t *testing.T) {
	for _, d := range convertStringToDateTest {
		t.Run(d.date, func(t *testing.T) {
			tt, err := convertStringToDate(d.date)
			if reflect.TypeOf(tt) != reflect.TypeOf(time.Time{}) {
				t.Errorf("Return value type expected: time.Time, got %T", reflect.TypeOf(tt))
			}
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if errMsg != d.expected {
				t.Errorf("Expected error message %s, got %s", d.expected, errMsg)
			}
		})
	}
}

var isValidJSONTest = []struct{
	message string
	jsonMessage string
}{
	{"Hello, John! Your birthday is in 7 day(s)","{\"message\":\"Hello, John! Your birthday is in 7 day(s)\"}"},
	{"Hello, John! Happy birthday!","{\"message\":\"Hello, John! Happy birthday!\"}"},
}

func Test_sendJSONResponse(t *testing.T) {
	for _, r := range isValidJSONTest {
		t.Run(r.message, func(t *testing.T) {
			resp := make(map[string]string)
			resp["message"] = r.message
			jsonResp, err := json.Marshal(resp)
			if err == nil {
				if (string(jsonResp) != string(r.jsonMessage)) {
					t.Errorf("Expected %s, got %s", string(r.jsonMessage), jsonResp)
				}
			} else {
				t.Errorf("Got %v", err)
			}
		})
	}
}
