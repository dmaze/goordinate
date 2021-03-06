package coordinatetest

import (
	"fmt"
	"github.com/dmaze/goordinate/coordinate"
	"gopkg.in/check.v1"
	"reflect"
	"time"
)

// ---------------------------------------------------------------------------
// HasKeys

type hasKeysChecker struct {
	*check.CheckerInfo
}

func (c hasKeysChecker) Info() *check.CheckerInfo {
	return c.CheckerInfo
}

func (c hasKeysChecker) Check(params []interface{}, names []string) (result bool, error string) {
	if len(params) != 2 {
		return false, "incorrect number of parameters to HasKeys check"
	}
	if len(names) != 2 {
		return false, "incorrect number of names to HasKeys check"
	}
	obtained := reflect.ValueOf(params[0])
	if obtained.Type().Kind() != reflect.Map {
		return false, fmt.Sprintf("%v value is not a map", names[0])
	}
	expected, ok := params[1].([]string)
	if !ok {
		return false, "expected keys for HasKeys check not a []string"
	}
	for _, key := range expected {
		value := obtained.MapIndex(reflect.ValueOf(key))
		if !value.IsValid() {
			return false, fmt.Sprintf("missing key %v", key)
		}
	}
	return true, ""
}

// The HasKeys checker verifies that a map has an expected set of keys.
// The values of the provided map are not checked, and extra keys are not
// checked for (try a HasLen check in addition to this).
//
//     var map[string]interface{} actual
//     actual = ...
//     c.Check(actual, HasKeys, []string{"foo", "bar"})
var HasKeys check.Checker = &hasKeysChecker{
	&check.CheckerInfo{
		Name:   "HasKeys",
		Params: []string{"obtained", "expected"},
	},
}

// ---------------------------------------------------------------------------
// AttemptMatches

type attemptMatchesChecker struct {
	*check.CheckerInfo
}

func (c attemptMatchesChecker) Info() *check.CheckerInfo {
	return c.CheckerInfo
}

func (c attemptMatchesChecker) Check(params []interface{}, names []string) (result bool, error string) {
	if len(params) != 2 {
		return false, "incorrect number of parameters to AttemptMatches check"
	}
	obtained, ok := params[0].(coordinate.Attempt)
	if !ok {
		return false, "non-Attempt obtained value"
	}
	expected, ok := params[1].(coordinate.Attempt)
	if !ok {
		return false, "non-Attempt expected value"
	}
	if obtained.Worker().Name() != expected.Worker().Name() {
		return false, "mismatched workers"
	}
	if obtained.WorkUnit().Name() != expected.WorkUnit().Name() {
		return false, "mismatched work units"
	}
	if obtained.WorkUnit().WorkSpec().Name() != expected.WorkUnit().WorkSpec().Name() {
		return false, "mismatched work specs"
	}
	return true, ""
}

// The AttemptMatches checker verifies that two attempts are compatible
// based on their observable data.
var AttemptMatches check.Checker = &attemptMatchesChecker{
	&check.CheckerInfo{
		Name:   "AttemptMatches",
		Params: []string{"obtained", "expected"},
	},
}

// ---------------------------------------------------------------------------
// SameTime

type sameTimeChecker struct {
	*check.CheckerInfo
}

func (c sameTimeChecker) Info() *check.CheckerInfo {
	return c.CheckerInfo
}

func (c sameTimeChecker) Check(params []interface{}, names []string) (result bool, error string) {
	if len(params) != 2 {
		return false, "incorrect number of parameters to SameTime check"
	}
	obtained, ok := params[0].(time.Time)
	if !ok {
		return false, "non-Time obtained value"
	}
	expected, ok := params[1].(time.Time)
	if !ok {
		return false, "non-Time expected value"
	}
	// NB: the postgres backend rounds to the microsecond
	maxDelta := time.Duration(1) * time.Microsecond
	delta := obtained.Sub(expected)
	return delta < maxDelta && delta > -maxDelta, ""
}

// The SameTime checker verifies that two times are identical, or at
// least, very very close.
var SameTime check.Checker = &sameTimeChecker{
	&check.CheckerInfo{
		Name:   "SameTime",
		Params: []string{"obtained", "expected"},
	},
}

// ---------------------------------------------------------------------------
// Like

type likeChecker struct {
	*check.CheckerInfo
}

func (c likeChecker) Info() *check.CheckerInfo {
	return c.CheckerInfo
}

func (c likeChecker) Check(params []interface{}, names []string) (result bool, error string) {
	if len(params) != 2 {
		return false, "incorrect number of parameters to Like check"
	}
	obtained := params[0]
	expected := params[1]
	obtainedV := reflect.ValueOf(obtained)
	expectedV := reflect.ValueOf(expected)
	obtainedT := obtainedV.Type()
	expectedT := expectedV.Type()

	if !obtainedT.ConvertibleTo(expectedT) {
		return false, "wrong type"
	}

	convertedV := obtainedV.Convert(expectedT)
	converted := convertedV.Interface()
	return converted == expected, ""
}

// The SameTime checker verifies that two objects are equal, once the
// obtained value has been cast to the type of the expected value.
var Like check.Checker = &likeChecker{
	&check.CheckerInfo{
		Name:   "Like",
		Params: []string{"obtained", "expected"},
	},
}
