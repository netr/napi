package sweets

import (
	"net/url"
	"reflect"
	"testing"
)

func TestSuite_MakePostData(t *testing.T) {
	tests := []struct {
		name string
		data string
		want url.Values
	}{
		{"should work with multiple params", "testing=this&and=testing&this=here", url.Values{"testing": {"this"}, "and": {"testing"}, "this": {"here"}}},
		{"it should work with no &", "testing_________=this", url.Values{"testing_________": {"this"}}},
		{"it should not work without a =", "testing_________", url.Values{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Suite{}
			if got := s.MakeUrlValues(tt.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakePostData() = %v, want %v", got, tt.want)
			}
		})
	}
}
