package sweets

import (
	"net/url"
	"strings"

	"github.com/stretchr/testify/suite"
)

// Suite wrapper for suite.Suite
type Suite struct {
	suite.Suite
}

func (suite *Suite) MakeUrlValues(data string) url.Values {
	v := url.Values{}
	ds := strings.Split(data, "&")
	for _, d := range ds {
		ns := strings.Split(d, "=")
		if len(ns) == 2 {
			v.Add(ns[0], ns[1])
		}
	}
	return v
}
