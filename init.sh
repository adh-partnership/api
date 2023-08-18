#!/bin/bash
#
# This script is meant to be part of the init-container to take the configuration from a config map,
# apply some env variables to it, and place it for the application to read.
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

configmap_path=${CONFIGMAP_PATH:-"/tmp/config"}
default_configmap=${DEFAULT_CONFIGMAP:-"../config.yaml.example"}
final_config=${CONFIG:-"/app/config.yaml"}
envpath=${ENV_PATH:-"out"}

if [[ -f "$configmap_path" ]]; then
    configmap=$(cat $configmap_path)
else
    configmap=$(cat $default_configmap)
fi

echo "$configmap" | PATH=$envpath envsubst > $final_config
