#!/bin/bash
#
# This script is meant to be part of the init-container to take the configuration from a config map,
# apply some env variables to it, and place it for the application to read.
#

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