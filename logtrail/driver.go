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
	"emperror.dev/errors"
	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
)

// Drivers combine multiple drivers into a single instance.
type Drivers []Driver

// should enqueue to Elastic search
func (d Drivers) Store(entry Entry) error {
	var errs []error

	for _, driver := range d {
		logger.Infof("sending log trail payload: [%+v]", entry)
		errs = append(errs, driver.Store(entry))
	}

	return errors.Combine(errs...)
}
