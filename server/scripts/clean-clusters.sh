
#!/usr/bin/env bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script handles the creation of multiple clusters using kind and the
# ability to create and configure an insecure container registry.

set -o errexit
set -o nounset
set -o pipefail

NUM_CLUSTERS="${NUM_CLUSTERS:-2}"

function clean-clusters() {
    local num_clusters=${1}
    for i in $(seq "${num_clusters}"); do
        kind delete clusters  "cluster${i}"
        kubectl config unset "users.cluster${i}"
        kubectl config unset "contexts.cluster${i}"
        kubectl config unset "clusters.cluster${i}"
        echo

    done
    kubectl config view

}
echo "Deleting ${NUM_CLUSTERS} clusters"
clean-clusters "${NUM_CLUSTERS}"
rm -rf submariner-cheatsheet
echo "Completed"