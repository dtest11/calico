// Copyright (c) 2022 Tigera, Inc. All rights reserved.
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

package daemon

import (
	"github.com/dtest11/calico/felix/config"
	"github.com/dtest11/calico/libcalico-go/lib/clientv3"
	"github.com/dtest11/calico/typha/pkg/discovery"
)

func bootstrapWireguardAndFilterTyphaAddresses(
	_ *config.Config, _ clientv3.NodesClient, typhas []discovery.Typha,
) ([]discovery.Typha, error) {
	return typhas, nil
} // no filtering

func bootstrapRemoveWireguard(_ *config.Config, _ clientv3.Interface) error {
	return nil
} // no-op
