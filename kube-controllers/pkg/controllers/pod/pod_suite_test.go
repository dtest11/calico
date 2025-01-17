// Copyright (c) 2021 Tigera, Inc. All rights reserved.
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

package pod_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/dtest11/calico/libcalico-go/lib/testutils"

	"github.com/onsi/ginkgo/reporters"
)

func init() {
	testutils.HookLogrusForGinkgo()
	logrus.SetLevel(logrus.DebugLevel)
}

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../report/pod_controller_suite.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Pod controller suite", []Reporter{junitReporter})
}
