package conveyer

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// ShouldHaveHeader tests whether a *http.Request has a header on it.
// An additional argument can be specified asserting equality for the header
func ShouldHaveHeader(actual interface{}, args ...interface{}) string {
	req := actual.(*http.Request)
	headerName := args[0].(string)
	expectedVal, hasExpectedVal := args[1].(string)

	headerStrs, hasHeader := req.Header[headerName]
	actualVal := strings.Join(headerStrs, ",")

	if !hasHeader && (!hasExpectedVal || (hasExpectedVal && expectedVal != "")) {
		return fmt.Sprintf(`Header "%s" was missing from request`, headerName)
	} else if actualVal != expectedVal {
		return Explain(`Header "%s" had incorrect value`, expectedVal, actualVal, headerName)
	}
	return ""
}

// ShouldHaveQueryParam tests whether a *http.Request has a specified query parameter.
// An additional argument can be specified asserting equality for the parameter value
func ShouldHaveQueryParam(actual interface{}, args ...interface{}) string {
	req := actual.(*http.Request)
	queryParam := args[0].(string)
	expectedVal, hasExpectedVal := args[1].(string)

	queryVals := req.URL.Query()
	queryStrs, hasQuery := queryVals[queryParam]
	actualVal := strings.Join(queryStrs, ",")

	if !hasQuery && (!hasExpectedVal || (hasExpectedVal && expectedVal != "")) {
		return fmt.Sprintf(`Query parameter "%s" was missing from request`, queryParam)
	} else if actualVal != expectedVal {
		return Explain(`Query parameter "%s" had incorrect value`, expectedVal, actualVal, queryParam)
	}

	return ""
}

// ShouldHaveFormValues tests whether a *http.Request has the provided form values.
func ShouldHaveFormValues(actual interface{}, args ...interface{}) string {
	req := actual.(*http.Request)
	vals := args[0].(url.Values)

	if err := req.ParseForm(); err != nil {
		return fmt.Sprintf("Error parsing request form: %s", err.Error())
	}

	form := req.Form

	for key := range vals {
		if !reflect.DeepEqual(vals[key], form[key]) {
			return Explain(`Form value "%s" had incorrect value`, vals[key], form[key], key)
		}
	}

	return ""
}
