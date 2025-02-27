#!/usr/bin/env bash

# Copyright 2021 Antrea Authors
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

set -eo pipefail

function echoerr {
    >&2 echo "$@"
}

_usage="Usage: $0 [--release|-r] [--global|-g] [--leader|-l <namespace>] [--member|-m] [--help|-h]
Generate a YAML manifest for Antrea MultiCluster using Kustomize and print it to stdout.
        --release | -r                      Enable release mode which will set correct release variant.
        --global  | -g                      Generate a global manifest for a Cluster as leader in a ClusterSet
        --leader  | -l                      Generate a per-namespace manifest for a Cluster as leader in a ClusterSet.
                                            All resources will be in the given namespace
        --member  | -m                      Generate a manifest for a Cluster as member in a ClusterSet
        --help    | -h                      Print this message and exit

Environment variables IMG_NAME and IMG_TAG must be set when release mode is enabled.
"

function print_usage {
    echoerr "$_usage"
}

function print_help {
    echoerr "Try '$0 --help' for more information."
}

OVERLAY=member
NAMESPACE=changeme
MODE=""

while [[ $# -gt 0 ]]
do
key="$1"
case $key in
    --release|-r)
    MODE="release"
    shift
    ;;
    --leader|-l)
    OVERLAY=leader-ns
    NAMESPACE="$2"
    shift 2
    ;;
    --global|-g)
    OVERLAY=leader-global
    shift
    ;;
    --member|-m)
    OVERLAY=member
    shift
    ;;
    -h|--help)
    print_usage
    exit 0
    ;;
    *)    # unknown option
    echoerr "Unknown option $1"
    exit 1
    ;;
esac
done

if [ "$MODE" == "release" ] && [ -z "$IMG_TAG" ]; then
    echoerr "In 'release' mode, environment variable IMG_TAG must be set"
    print_help
    exit 1
fi

if [ "$MODE" == "release" ] && [ -z "$IMG_NAME" ]; then
    echoerr "In 'release' mode, environment variable IMG_NAME must be set"
    print_help
    exit 1
fi

WORK_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

source $WORK_DIR/../../hack/verify-kustomize.sh

if [ -z "$KUSTOMIZE" ]; then
    KUSTOMIZE="$(verify_kustomize)"
elif ! $KUSTOMIZE version > /dev/null 2>&1; then
    echoerr "$KUSTOMIZE does not appear to be a valid kustomize binary"
    print_help
    exit 1
fi

KUSTOMIZATION_DIR=$WORK_DIR/../config

cd $KUSTOMIZATION_DIR

TMP_DIR=$(mktemp -d $KUSTOMIZATION_DIR/overlays.XXXXXXXX)
pushd $TMP_DIR > /dev/null

if [ "$OVERLAY" == "leader-ns" ] ;
then
    mkdir config && cd config
    cp $KUSTOMIZATION_DIR/overlays/leader-ns/prefix_transformer.yaml .
    sed -ie "s/changeme/$NAMESPACE/g" prefix_transformer.yaml

cat << EOF > kustomization.yaml
namespace: $NAMESPACE

bases:
  - ../../overlays/$OVERLAY

transformers:
  - prefix_transformer.yaml
EOF
else
cat << EOF > kustomization.yaml
bases:
  - ../overlays/$OVERLAY
EOF
fi

if [ "$MODE" == "release" ]; then
    $KUSTOMIZE edit set image antrea/antrea-mc-controller=$IMG_NAME:$IMG_TAG
fi
$KUSTOMIZE build
rm -rf $TMP_DIR
