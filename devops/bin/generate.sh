#!/bin/bash

set -o errexit
set -o pipefail
set -o errtrace

# -------------------------------
# Config & Colors
# -------------------------------
YQ_MIN_VERSION="4.25.1"
RED='\033[0;31m' # Red Color
NC='\033[0m'      # No Color

# -------------------------------
# Error handling
# -------------------------------
err_report() {
    echo "Error running '$1' [rc=$2] line $3"
}
trap 'err_report "$BASH_COMMAND" $? $LINENO' ERR

# -------------------------------
# Version check
# -------------------------------
function check_package_version() {
    CURR_VERSION=$1
    MIN_VERSION=$2
    PACKAGE=$3

    IFS="." read -r -a CURRPARTS <<< "$CURR_VERSION"
    IFS="." read -r -a MINPARTS <<< "$MIN_VERSION"

    if [[ ${#CURRPARTS[@]} -lt 3 ]] || [[ ${#MINPARTS[@]} -lt 3 ]]; then
        echo "Incorrect version format for $PACKAGE: MAJOR.MINOR.PATCH required."
        return 1
    fi

    if [[ ${CURRPARTS[0]} -lt ${MINPARTS[0]} ]] || \
       [[ ${CURRPARTS[0]} -eq ${MINPARTS[0]} && ${CURRPARTS[1]} -lt ${MINPARTS[1]} ]] || \
       [[ ${CURRPARTS[0]} -eq ${MINPARTS[0]} && ${CURRPARTS[1]} -eq ${MINPARTS[1]} && ${CURRPARTS[2]} -lt ${MINPARTS[2]} ]]; then
        echo "Minimum version required for $PACKAGE is $MIN_VERSION. Your current $PACKAGE version is $CURR_VERSION."
        return 1
    fi
}

# -------------------------------
# Pre-requirements check
# -------------------------------
check_Pre_Requirements() {
    if ! hash curl 2>/dev/null; then
        echo "curl command is required: https://curl.haxx.se/download.html"
        exit 1
    fi

    if ! hash swagger 2>/dev/null; then
        echo "swagger command is required: https://github.com/go-swagger/go-swagger/releases"
        exit 1
    fi

    # Determine AWK
    if awk --version 2>&1 | grep -q "GNU Awk"; then
        AWK=GAWK
    else
        AWK=NAWK
    fi

    # Check Swagger latest version
    if [ "$AWK" = "NAWK" ]; then
        latest=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | awk -F\" '/\"name\":/ {print $4;}' | grep -E 'v[[:digit:]\.]+' || true)
    else
        latest=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | awk -F\" '/"name":/ {print $4;}' | grep -E 'v[[:digit:]\.]+' || true)
    fi

    if [ -z "$latest" ]; then
        latest=$(curl -s -i https://github.com/go-swagger/go-swagger/releases/latest | awk -F/ 'tolower($0) ~ /location:/ {sub("\r", ""); print $NF}')
    fi

    current=$(swagger version | awk '/version/ {print $2}')
    if [ "$latest" != "$current" ]; then
        echo "Update recommended: Swagger version $current, latest is $latest."
        # exit 1
    fi

    # Check for yq
    if ! hash yq 2>/dev/null; then
        echo "yq command is missing: https://github.com/mikefarah/yq"
        echo "Skipping API Docs processing."
        exit 1
    fi

    YQ_CURR_VERSION=$(yq -V | awk '{print $4}' | sed 's/v//g')
    check_package_version "$YQ_CURR_VERSION" $YQ_MIN_VERSION "yq"
    if [ $? -eq 1 ]; then
        exit 1
    fi
}

# -------------------------------
# Main Script
# -------------------------------
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

# Check for pre-requirements
check_Pre_Requirements

# Initialize YAML files array
ALL_YAML_FILES=()

# Primary definitions YAML
PRIMARY_YAML="swagger/components/definitions.yaml"
if [ -f "$PRIMARY_YAML" ]; then
    ALL_YAML_FILES+=("$PRIMARY_YAML")
fi

# Additional YAMLs in components
while IFS= read -r -d $'\0' file; do
    ALL_YAML_FILES+=("$file")
done < <(find "swagger/components/" -name "*.yaml" ! -name "definitions.yaml" -print0)

# Move to project root
cd "$DIR/../.."

# Remove previously generated Go files
find . -type f -name '*.go' -not -path "./vendor/*" -exec sh -c 'grep -qE "^// Code generated .* DO NOT EDIT\.$" $0 && rm -f $0' {} \;

# -------------------------------
# Mix YAMLs into swagger.yaml
# -------------------------------
set +o errexit
echo -e "${RED}Executing Swagger mixin for swagger/swagger.yaml${NC}"

if [ ${#ALL_YAML_FILES[@]} -gt 0 ]; then
    cp "${ALL_YAML_FILES[0]}" swagger/swagger.yaml

    for yamlFile in "${ALL_YAML_FILES[@]:1}"; do
        echo "Mixing in $yamlFile"
        swagger mixin --format=yaml -o swagger/swagger.yaml swagger/swagger.yaml "$yamlFile"
    done

    # Ensure top-level 'swagger' field exists
    if ! grep -q '^swagger:' swagger/swagger.yaml; then
        echo "Adding top-level 'swagger: \"2.0\"' to swagger.yaml"
        yq eval '.swagger = "2.0"' -i swagger/swagger.yaml
    fi

    # Ensure .info section exists
    if ! grep -q '^info:' swagger/swagger.yaml; then
        echo "Adding .info section to swagger.yaml"
        yq eval '.info = {"title": "Adornme", "version": "1.0.0"}' -i swagger/swagger.yaml
    fi

    # Sort keys
    yq -i -P 'sort_keys(..)' swagger/swagger.yaml
else
    echo "Error: No swagger YAML files found to process."
    exit 1
fi

set -o errexit

# -------------------------------
# Generate Go server
# -------------------------------
echo -e "${RED}Executing Swagger generate SB Server${NC}"
swagger generate server -A AdronmeCode -P models.Principal -f swagger/swagger.yaml || {
    echo "Error generating server. Please check the swagger.yaml for issues."
    exit 1
}

# -------------------------------
# Generate swagger.json
# -------------------------------
echo -e "${RED}Generating swagger.json from swagger.yaml${NC}"
if hash yq 2>/dev/null; then
    yq eval -o=json swagger/swagger.yaml > swagger/swagger.json
    echo "swagger.json created at swagger/swagger.json"
else
    echo "yq command not found, cannot generate swagger.json"
fi

echo -e "${RED}Finish Execution${NC}"
