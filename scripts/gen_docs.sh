#!/bin/bash
#
# Copyright ADH Partnership
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

wget https://github.com/swaggo/swag/releases/download/v1.8.3/swag_1.8.3_Linux_x86_64.tar.gz -O /tmp/swag.tar.gz
mkdir /tmp/swag || true
tar -xzf /tmp/swag.tar.gz -C /tmp/swag
/tmp/swag/swag init --parseDependency --parseInternal --parseDepth 3 --generalInfo cmd/api/main.go
if [[ $? -ne 0 ]]; then
  rm -rf /tmp/swag
  echo "Failed to generate swagger documentation"
  exit 1
fi
rm -rf /tmp/swag