package dtest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
)

func NormalisePointers(s string) string {
	re := regexp.MustCompile(`(0x[a-f0-9]+)`)
	counter := 1
	for {
		result := re.FindString(s)
		if result == "" {
			break
		}

		s = strings.Replace(s, result, fmt.Sprintf("0p%d", counter), -1)
		counter++
	}
	return s
}

func normalisePointers(s string) string {
	re := regexp.MustCompile(`(\(0x[a-f0-9]+\))`)
	counter := 1
	for {
		result := re.FindString(s)
		if result == "" {
			break
		}

		s = strings.Replace(s, result, fmt.Sprintf("(0p%d)", counter), -1)
		counter++
	}
	return s
}

type FailedCompare struct {
	Expected string
	Reported string
	Diff     string
}

func (err *FailedCompare) Error() string {
	return err.Diff
}

func CompareObjects(t *testing.T, expected, reported interface{}, message ...interface{}) error {
	spew := spew.ConfigState{SortKeys: true}
	expectedStr := normalisePointers(spew.Sdump(expected))
	reportedStr := normalisePointers(spew.Sdump(reported))

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(expectedStr),
		B:        difflib.SplitLines(reportedStr),
		FromFile: "Expected",
		ToFile:   "Reported",
		Context:  3,
	}

	text, _ := difflib.GetUnifiedDiffString(diff)
	if len(text) != 0 {
		if len(message) > 0 {
			t.Logf(message[0].(string), message[1:]...)
		}
		t.Log("\n" + text)
		t.Fail()
		return &FailedCompare{expectedStr, reportedStr, text}
	}
	return nil
}
func CompareSnapshot(t *testing.T, expected string, reported interface{}) {
	panic("NYI")
}

func CompareStrings(t *testing.T, expected string, reported string) {
	panic("NYI")
}

func CompareObjectExhibit(t *testing.T, exhibit string, reported interface{}) error {
	reportedStr := normalisePointers(spew.Sdump(reported))

	if os.Getenv("TEST_SNAPSHOT") == "TRUE" {
		err := ioutil.WriteFile(exhibit, []byte(reportedStr), 0644)
		if err != nil {
			t.Logf("File Error: %s", err.Error())
			t.FailNow()
		}
		return err
	}

	buf, err := ioutil.ReadFile(exhibit)
	if err != nil {
		t.Logf("File Error: %s", err.Error())
		t.FailNow()
		return err
	}

	expectedStr := string(buf[:])

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(expectedStr),
		B:        difflib.SplitLines(reportedStr),
		FromFile: "Expected",
		ToFile:   "Reported",
		Context:  3,
	}

	text, _ := difflib.GetUnifiedDiffString(diff)
	if len(text) != 0 {
		t.Log("Exhibit:", exhibit)
		t.Log("\n" + text)
		t.Fail()
		return &FailedCompare{expectedStr, reportedStr, text}
	}
	return nil
}

func CompareJSONExhibit(t *testing.T, exhibit string, reported interface{}) error {
	reportedJson, err := json.MarshalIndent(reported, "", "  ")
	if err != nil {
		t.Logf("Couldn't marshal: %s", err.Error())
		t.FailNow()
	}

	if os.Getenv("TEST_SNAPSHOT") == "TRUE" {
		err := ioutil.WriteFile(exhibit, reportedJson, 0644)
		if err != nil {
			t.Logf("File Error: %s", err.Error())
			t.FailNow()
		}
		return err
	}

	buf, err := ioutil.ReadFile(exhibit)
	if err != nil {
		t.Logf("File Error: %s", err.Error())
		t.Fail()
		return err
	}

	expectedStr := string(buf[:])
	reportedStr := string(reportedJson)

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(expectedStr),
		B:        difflib.SplitLines(reportedStr),
		FromFile: "Expected",
		ToFile:   "Reported",
		Context:  3,
	}

	text, _ := difflib.GetUnifiedDiffString(diff)
	if len(text) != 0 {
		t.Log("Exhibit:", exhibit)
		t.Log("\n" + text)
		t.Fail()
		return &FailedCompare{expectedStr, reportedStr, text}
	}
	return nil
}

func CompareExhibit(t *testing.T, exhibit string, reported string) error {
	if os.Getenv("TEST_SNAPSHOT") == "TRUE" && len(reported) != 0 {
		err := ioutil.WriteFile(exhibit, []byte(reported), 0644)
		if err != nil {
			t.Logf("File Error: %s", err.Error())
			t.Fail()
		}
		return err
	}

	buf, err := ioutil.ReadFile(exhibit)
	if err != nil && len(reported) != 0 {
		t.Logf("File Error: %s", err.Error())
		t.Fail()
		return err
	} else if len(reported) == 0 {
		//No exhibit, no reported.
		return nil
	}

	expectedStr := string(buf)

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(expectedStr),
		B:        difflib.SplitLines(reported),
		FromFile: "Expected",
		ToFile:   "Reported",
		Context:  3,
	}

	text, _ := difflib.GetUnifiedDiffString(diff)
	if len(text) != 0 {
		t.Log("Exhibit:", exhibit)
		t.Log("\n" + text)
		t.Fail()
		return &FailedCompare{expectedStr, reported, text}
	}
	return nil
}
