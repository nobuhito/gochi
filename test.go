package gochi

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// https://github.com/benbjohnson/testing

func (g *Gochi) Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: "+msg+"\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func (g *Gochi) Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: unexpected error: %s\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

func (g *Gochi) Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\texp: %#v\n\tgot: %#v\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func (g *Gochi) EqualsWithNumber(tb testing.TB, i int, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n%v)\texp: %#v\n\tgot: %#v\n", filepath.Base(file), line, i, exp, act)
		tb.FailNow()
	}
}

func testJSON() (interface{}, string) {
	json := retJSON()
	return json, "{\"id\":\"abcdefg\",\"value\":\"hijklmn\"}"
}

func retJSON() interface{} {
	json := struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	}{
		ID:    "abcdefg",
		Value: "hijklmn",
	}
	return json
}
