// Copyright (c) 2016-2020 Tigera, Inc. All rights reserved.
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

package calc

import (
	"errors"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/dtest11/calico/libcalico-go/lib/backend/api"
	"github.com/dtest11/calico/libcalico-go/lib/backend/model"
	v1v "github.com/dtest11/calico/libcalico-go/lib/validator/v1"
	v3v "github.com/dtest11/calico/libcalico-go/lib/validator/v3"
)

func NewValidationFilter(sink api.SyncerCallbacks) *ValidationFilter {
	return &ValidationFilter{
		sink: sink,
	}
}

type ValidationFilter struct {
	sink api.SyncerCallbacks
}

func (v *ValidationFilter) OnStatusUpdated(status api.SyncStatus) {
	// Pass through.
	v.sink.OnStatusUpdated(status)
}

func (v *ValidationFilter) OnUpdates(updates []api.Update) {
	filteredUpdates := make([]api.Update, len(updates))
	for i, update := range updates {
		logCxt := logrus.WithFields(logrus.Fields{
			"key":   update.Key,
			"value": update.Value,
		})
		logCxt.Debug("Validating KV pair.")
		validatorFunc := v1v.Validate
		if _, isV3 := update.Key.(model.ResourceKey); isV3 {
			logCxt.Debug("Use v3 validator")
			validatorFunc = v3v.Validate
		} else {
			logCxt.Debug("Use v1 validator")
		}
		if update.Value != nil {
			val := reflect.ValueOf(update.Value)
			if val.Kind() == reflect.Ptr {
				elem := val.Elem()
				if elem.Kind() == reflect.Struct {
					if err := validatorFunc(elem.Interface()); err != nil {
						logCxt.WithError(err).Warn("Validation failed; treating as missing")
						update.Value = nil
					}
				}
			}

			switch k := update.Key.(type) {
			case model.NodeKey:
				// TODO: This should be in its own filter.
				// Special case: we can't serialize Node keys but Felix only cares
				// about the host metadata anyway.  Extract the Host IP.
				update.Key = model.HostIPKey(k)
				if update.Value != nil {
					node, ok := update.Value.(*model.Node)
					if ok {
						update.Value = node.FelixIPv4
					} else {
						logCxt.WithField("value", update.Value).Warn(
							"NodeKey value wasn't a *Node; treating as missing")
						update.Value = nil
					}
				}
			}

			switch v := update.Value.(type) {
			case *model.WorkloadEndpoint:
				if v.Name == "" {
					logCxt.WithError(errors.New("Missing name")).Warn("Validation failed; treating as missing")
					update.Value = nil
				}
			}
		}
		filteredUpdates[i] = update
	}
	v.sink.OnUpdates(filteredUpdates)
}
