// Package trex
// Test Requests => T REQS => trex

package trex

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/gofiber/fiber/v2"
	"github.com/steinfletcher/apitest-jsonpath/jsonpath"
	"github.com/stretchr/testify/require"
)

func New(s IFiberTestSuite) *TestResponse {
	return &TestResponse{
		response: &http.Response{},
		suite:    s,
	}
}

type TestResponse struct {
	request  *http.Request
	response *http.Response
	suite    IFiberTestSuite
}

// Get sends a test GET request and returns a TestResponse.
func (tr *TestResponse) Get(url string, headers *http.Header) *TestResponse {
	tr.TestRequest(http.MethodGet, url, nil, headers)
	return tr
}

// Delete sends a test DELETE request and returns a TestResponse.
func (tr *TestResponse) Delete(url string, postData *url.Values, headers *http.Header) *TestResponse {
	tr.TestRequest(http.MethodDelete, url, postData, headers)
	return tr
}

// Patch sends a test PATCH request and returns a TestResponse.
func (tr *TestResponse) Patch(url string, postData *url.Values, headers *http.Header) *TestResponse {
	tr.TestRequest(http.MethodPatch, url, postData, headers)
	return tr
}

// Json sends a test json header and returns a TestResponse.
func (tr *TestResponse) Json(method string, url string, postData *url.Values, headers ...*http.Header) *TestResponse {
	hdr := &http.Header{}
	if len(headers) > 0 {
		if headers[0] != nil {
			hdr = headers[0]
		}
	}
	hdr.Set(fiber.HeaderContentType, "application/json")

	tr.TestRequest(method, url, postData, hdr)
	return tr
}

// Post sends a test POST request and returns a TestResponse.
func (tr *TestResponse) Post(url string, postData *url.Values, headers *http.Header) *TestResponse {
	tr.TestRequest(http.MethodPost, url, postData, headers)
	return tr
}

// Put sends a test PUT request and returns a TestResponse.
func (tr *TestResponse) Put(url string, postData *url.Values, headers *http.Header) *TestResponse {
	tr.TestRequest(http.MethodPut, url, postData, headers)
	return tr
}

// TestRequest will execute a fiber.App().Test() on the given request data and return a TestResponse.
func (tr *TestResponse) TestRequest(method, url string, postData *url.Values, headers *http.Header) *TestResponse {
	req := createRequest(method, tr.makeUrl(url), postData, headers)
	resp, _ := tr.suite.App().Test(req, 15000)

	// this allows you to see the post data in Dump()
	if postData != nil {
		req.Form = *postData
	}

	tr.response = resp
	tr.request = req
	return tr
}

// TestReaderRequest will execute a fiber.App().Test() on the given request data and return a TestResponse.
func (tr *TestResponse) TestReaderRequest(method, url string, reader io.Reader, headers *http.Header) *TestResponse {
	req := createReaderRequest(method, tr.makeUrl(url), reader, headers)
	resp, _ := tr.suite.App().Test(req, 15000)

	tr.response = resp
	tr.request = req
	return tr
}

// BodyReader returns the response body as io.ReadCloser.
//
// See: https://stackoverflow.com/a/60098300 for information about this issue.
// Drains the http.Response body and creates new readers every call. Allows us to call this function multiple times.
// This is not the right way to approach it, but since this is only being used in testing, we should be OK.
// Would love to figure out the most efficient way to do this.
func (tr *TestResponse) BodyReader() io.ReadCloser {
	rawBody, err := io.ReadAll(tr.response.Body)
	if err != nil {
		tr.suite.T().Error(err)
	}
	nb := io.NopCloser(bytes.NewBuffer(rawBody))
	tr.response.Body = io.NopCloser(bytes.NewBuffer(rawBody))
	return nb
}

// BodyBytes returns the response body as []byte
func (tr *TestResponse) BodyBytes() []byte {
	var b []byte
	save, b, err := drainBody(tr.BodyReader())
	if err != nil {
		return b
	}

	// reset response.body after draining body to keep stream alive
	tr.response.Body = save
	return b
}

// BodyString returns the response body as string
func (tr *TestResponse) BodyString() string {
	return string(tr.BodyBytes())
}

// ToJson is a helper function convert the response body's bytes into a simplejson.Json
func (tr *TestResponse) ToJson() *simplejson.Json {
	sj, err := simplejson.NewJson(tr.BodyBytes())
	if err != nil {
		tr.suite.T().Fatal(err)
	}
	return sj
}

// DataToJson is a helper function to easily pull the data field from SuccessResponse and return a simplejson.Json
func (tr *TestResponse) DataToJson() *simplejson.Json {
	return tr.ToJson().Get("data")
}

// DataToString is a helper function to easily pull the data field from SuccessResponse and return a string
func (tr *TestResponse) DataToString() string {
	return tr.ToJson().Get("data").MustString()
}

// DataToSlice is a helper function to easily pull the data field from SuccessResponse and return a []interface{}
func (tr *TestResponse) DataToSlice() []interface{} {
	return tr.ToJson().Get("data").MustArray()
}

// DataToMap is a helper function to easily pull the data field from SuccessResponse and return a map[string]interface{}
func (tr *TestResponse) DataToMap() map[string]interface{} {
	return tr.ToJson().Get("data").MustMap()
}

// Message is a helper function to easily pull the message field from SuccessResponse
func (tr *TestResponse) Message() string {
	return tr.ToJson().Get("message").MustString()
}

// Status returns the responses's status code
func (tr *TestResponse) Status() int {
	if tr.response == nil {
		return 0
	}
	return tr.response.StatusCode
}

// StatusString returns the response's status code as a string
func (tr *TestResponse) StatusString() string {
	return fmt.Sprintf("%d %s", tr.Status(), http.StatusText(tr.Status()))
}

// Request returns the underlying http.Request
func (tr *TestResponse) Request() *http.Request {
	return tr.request
}

// Dump dumps the response HTML and Header to stdout
func (tr *TestResponse) Dump() *TestResponse {
	fmt.Println(formatRequest(tr.Request()) + "\n------------------------------------------------------")
	fmt.Printf("Status: %s\n%s\n\n", tr.StatusString(), tr.BodyString())
	return tr
}

// AssertOk will check for status code http.StatusOK and fail if not found.
func (tr *TestResponse) AssertOk() *TestResponse {
	tr.AssertStatus(http.StatusOK)
	return tr
}

// AssertUnauthorized will check for status code http.StatusUnauthorized and fail if not found.
func (tr *TestResponse) AssertUnauthorized() *TestResponse {
	tr.AssertStatus(http.StatusUnauthorized)
	return tr
}

// AssertUnprocessable will check for status code http.StatusUnprocessableEntity and fail if not found.
func (tr *TestResponse) AssertUnprocessable() *TestResponse {
	tr.AssertStatus(http.StatusUnprocessableEntity)
	return tr
}

// AssertStatus will check for the given status code and fail if not found.
func (tr *TestResponse) AssertStatus(code int) *TestResponse {
	require.Equalf(tr.suite.T(), code, tr.Status(), "wanted status code: %v, got: %v", code, tr.Status())
	return tr
}

// AssertValidationErrors takes a list of field names and exhaustively checks the FormErrorResponse for matches.
func (tr *TestResponse) AssertValidationErrors(fields ...string) *TestResponse {
	res, err := tr.ParseFormErrors()
	if err != nil {
		tr.suite.T().Fatal(err)
		return tr
	}

	require.Truef(tr.suite.T(), res.Contains(fields...), "wanted: %v, got: %v", fields, res.Fields())
	return tr
}

// AssertJsonContains is a convenience function to assert that a jsonpath expression extracts a value in an array
func (tr *TestResponse) AssertJsonContains(expression string, expected interface{}) *TestResponse {
	require.NoError(tr.suite.T(), jsonpath.Contains(expression, expected, tr.BodyReader()))
	return tr
}

// AssertJsonEqual is a convenience function to assert that a jsonpath expression extracts a value in an array
func (tr *TestResponse) AssertJsonEqual(expression string, expected interface{}) *TestResponse {
	require.NoError(tr.suite.T(), jsonpath.Equal(expression, expected, tr.BodyReader()))
	return tr
}

// AssertJsonNotEqual is a function to check json path expression value is not equal to given value
func (tr *TestResponse) AssertJsonNotEqual(expression string, expected interface{}) *TestResponse {
	require.NoError(tr.suite.T(), jsonpath.NotEqual(expression, expected, tr.BodyReader()))
	return tr
}

// AssertJsonLen asserts that value is the expected length, determined by reflect.Len
func (tr *TestResponse) AssertJsonLen(expression string, expectedLen int) *TestResponse {
	require.NoError(tr.suite.T(), jsonpath.Length(expression, expectedLen, tr.BodyReader()))
	return tr
}

// AssertJsonGreaterThan asserts that value is greater than the given length, determined by reflect.Len
func (tr *TestResponse) AssertJsonGreaterThan(expression string, minimumLength int) *TestResponse {
	require.NoError(tr.suite.T(), jsonpath.GreaterThan(expression, minimumLength, tr.BodyReader()))
	return tr
}

// AssertJsonLessThan asserts that value is greater than the given length, determined by reflect.Len
func (tr *TestResponse) AssertJsonLessThan(expression string, maximumLength int) *TestResponse {
	require.NoError(tr.suite.T(), jsonpath.LessThan(expression, maximumLength, tr.BodyReader()))
	return tr
}

// AssertJsonPresent asserts that value returned by the expression is present
func (tr *TestResponse) AssertJsonPresent(expression string) *TestResponse {
	require.NoError(tr.suite.T(), jsonpath.Present(expression, tr.BodyReader()))
	return tr
}

// AssertJsonNotPresent asserts that value returned by the expression is not present
func (tr *TestResponse) AssertJsonNotPresent(expression string) *TestResponse {
	require.NoError(tr.suite.T(), jsonpath.NotPresent(expression, tr.BodyReader()))
	return tr
}

// AssertJsonMatches asserts that the value matches the given regular expression
func (tr *TestResponse) AssertJsonMatches(expression string, pattern string) *TestResponse {
	rgx, err := regexp.Compile(pattern)
	if err != nil {
		require.Fail(tr.suite.T(), errors.New(fmt.Sprintf("invalid pattern: '%s'", rgx)).Error())
	}
	value, _ := jsonpath.JsonPath(tr.BodyReader(), expression)
	if value == nil {
		require.Fail(tr.suite.T(), errors.New(fmt.Sprintf("no match for pattern: '%s'", expression)).Error())
	}
	kind := reflect.ValueOf(value).Kind()
	switch kind {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		if !rgx.Match([]byte(fmt.Sprintf("%v", value))) {
			require.Fail(tr.suite.T(), errors.New(fmt.Sprintf("value '%v' does not match pattern '%v'", value, rgx)).Error())
		}
		return nil
	default:
		require.Fail(tr.suite.T(), errors.New(fmt.Sprintf("unable to match using type: %s", kind.String())).Error())
	}
	return tr
}

// AssertDataCount checks the success response data array for a specific result count
func (tr *TestResponse) AssertDataCount(count int) *TestResponse {
	res, err := tr.ParseSuccess()
	if err != nil {
		tr.suite.T().Fatal(err)
		return tr
	}

	data := res.Data.([]interface{})
	if len(data) != count {
		tr.suite.T().Fatal(fmt.Errorf("wanted: len(%d), got: len(%d)", count, len(data)))
	}

	return tr
}

// ParseSuccess will parse the response body and convert it into a SuccessResponse
func (tr *TestResponse) ParseSuccess() (*SuccessResponse, error) {
	res, err := tr.parseBody(&SuccessResponse{})
	if err != nil {
		return nil, err
	}
	return res.(*SuccessResponse), nil
}

// ParseFormErrors will parse the response body and convert it into a FormErrorResponse
func (tr *TestResponse) ParseFormErrors() (*FormErrorResponse, error) {
	res, err := tr.parseBody(&FormErrorResponse{})
	if err != nil {
		return nil, err
	}
	return res.(*FormErrorResponse), nil
}

// parseBody will take an interface and json unmarshal the response body into it
func (tr *TestResponse) parseBody(castTo interface{}) (interface{}, error) {
	b := tr.BodyBytes()
	err := json.Unmarshal(b, &castTo)
	if err != nil {
		return castTo, err
	}

	return castTo, nil
}

// makeUrl will create a url in the appropriate format
func (tr *TestResponse) makeUrl(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

// convertUrlValuesToJson will convert a url.Values object into a json string and remove the array brackets on single items
func convertUrlValuesToJson(values *url.Values) string {
	json, err := json.Marshal(values)
	if err != nil {
		return ""
	}

	// remove the array brackets *only* from single items
	jsonString := string(json)
	jsonString = strings.Replace(jsonString, `["`, `"`, -1)
	jsonString = strings.Replace(jsonString, `"]`, `"`, -1)

	return jsonString
}

// createReaderRequest helper function to create httptest requests
func createReaderRequest(method, url string, reader io.Reader, headers *http.Header) *http.Request {
	req := httptest.NewRequest(method, url, reader)

	if headers != nil {
		// add the parameter headers into the requests headers
		for key, value := range *headers {
			req.Header.Add(key, value[0])
		}
	}

	if req.Header.Get(fiber.HeaderContentType) == "" {
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationForm)
	}

	return req
}

// createRequest helper function to create httptest requests
func createRequest(method, url string, postData *url.Values, headers *http.Header) *http.Request {
	var req *http.Request
	if postData == nil {
		req = httptest.NewRequest(method, url, nil)
		if headers != nil {
			// add the parameter headers into the requests headers
			for key, value := range *headers {
				req.Header.Add(key, value[0])
			}
		}
	} else {
		if headers != nil && headers.Get(fiber.HeaderContentType) == "application/json" {
			req = httptest.NewRequest(method, url, strings.NewReader(convertUrlValuesToJson(postData)))
		} else {
			req = httptest.NewRequest(method, url, strings.NewReader(postData.Encode()))
		}

		if headers != nil {
			// add the parameter headers into the requests headers
			for key, value := range *headers {
				req.Header.Add(key, value[0])
			}
		}

		if req.Header.Get(fiber.HeaderContentType) == "" {
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationForm)
		}
	}
	return req
}

func JsonHeader() *http.Header {
	headers := http.Header{}
	headers.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return &headers
}

// formatRequest generates ascii representation of a request. [https://medium.com/doing-things-right/pretty-printing-http-requests-in-golang-a918d5aaa000]
func formatRequest(r *http.Request) string {
	// Store return string
	var request []string
	// Add the request string
	uri := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, uri)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			return ""
		}
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

// drainBody reads all of b to memory and then returns two equivalent
// ReadClosers yielding the same bytes.
//
// It returns an error if the initial slurp of all bytes fails. It does not attempt
// to make the returned ReadClosers have identical error-matching behavior.
func drainBody(b io.ReadCloser) (io.ReadCloser, []byte, error) {
	var err error
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return nil, nil, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return io.NopCloser(bytes.NewReader(buf.Bytes())), buf.Bytes(), nil
}
