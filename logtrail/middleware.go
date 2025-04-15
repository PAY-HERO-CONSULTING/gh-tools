// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logtrail

import (
	"regexp"
	"time"

	"github.com/PAY-HERO-CONSULTING/gh-tools/dtos"
	"github.com/gin-gonic/gin"
)

type Driver interface {
	// Store saves an audit log entry.
	Store(entry Entry) error
}

// Option configures an audit log middleware.
type Option interface {
	Apply(o *middlewareOptions)
}
type middlewareOptions struct {
	clock          dtos.Clock
	sensitivePaths []*regexp.Regexp
}

func NewMiddlewareOptions(
	driver Driver,
	opts ...Option,
) middlewareOptions {
	return middlewareOptions{
		clock: dtos.NewClock(),
		// userIDExtractor: func(req *http.Request) uint { return 0 },
	}
}

func (m *middlewareOptions) Now() time.Time {
	return m.clock.Now()
}

func (m *middlewareOptions) Since(startTime time.Time) time.Duration {
	return m.clock.Since(startTime)
}

func (m *middlewareOptions) SinceInMilliSeconds(startTime time.Time) int64 {
	return m.clock.Since(startTime).Microseconds()
}

func (m *middlewareOptions) SetupSentivePaths(
	c *gin.Context,
) bool {
	for _, r := range m.sensitivePaths {
		if r.MatchString(c.Request.URL.Path) {
			return true
		}
	}

	return false
}

type optionFunc func(o *middlewareOptions)

func (fn optionFunc) Apply(o *middlewareOptions) {
	fn(o)
}

// WithClock sets the clock in an audit log middleware.
func WithClock(clock dtos.Clock) Option {
	return optionFunc(func(o *middlewareOptions) {
		o.clock = clock
	})
}

// WithSensitivePaths marks API call paths as sensitive, causing the log entry to omit the request body.
func WithSensitivePaths(sensitivePaths []*regexp.Regexp) Option {
	return optionFunc(func(o *middlewareOptions) {
		o.sensitivePaths = sensitivePaths
	})
}
