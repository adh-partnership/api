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

set -e

exit_code=0
for file in "$@"; do
  if ! grep -Lqe "Apache License, Version 2" "${file}"; then
    echo "File is missing license: ${file}"
    exit_code=1
  fi

  if ! grep -Lqe "Copyright" "${file}"; then
    echo "File is missing copyright: ${file}"
    exit_code=1
  fi
done

exit ${exit_code}