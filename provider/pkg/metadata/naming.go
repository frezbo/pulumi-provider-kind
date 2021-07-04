// Copyright 2016-2021, Pulumi Corporation.
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

package metadata

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
)

var dns1123Alphabet = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

// AssignNameIfAutonamable generates a name for an object. Uses DNS-1123-compliant characters.
func AssignNameIfAutonamable(obj *v1alpha4.Cluster, propMap resource.PropertyMap, base tokens.QName) {
	contract.Assert(base != "")

	// Check if the name is set and is a computed value. If so, do not auto-name.
	if name, ok := propMap["name"]; ok && name.IsComputed() {
		return
	}

	if obj.Name == "" {
		obj.Name = fmt.Sprintf("%s-%s", base, randString(8))
	}
}

// AdoptOldAutonameIfUnnamed checks if `newObj` has a name, and if not, "adopts" the name of `oldObj`
// instead. If `oldObj` was autonamed, then we mark `newObj` as autonamed, too.
func AdoptOldAutonameIfUnnamed(newObj, oldObj *v1alpha4.Cluster) {
	contract.Assert(oldObj.Name != "")
	if newObj.Name == "" {
		newObj.Name = oldObj.Name
	}
}
func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		// nolint:gosec
		b[i] = dns1123Alphabet[rand.Intn(len(dns1123Alphabet))]
	}
	return string(b)
}

// Seed RNG to get different random names at each suffix.
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
