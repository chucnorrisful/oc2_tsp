#!/usr/bin/env bash

rm -f binaries/*

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <filename prefix>"
  exit 1
fi

platforms=("linux/amd64" "windows/amd64" "windows/386" "darwin/amd64")

for platform in "${platforms[@]}"; do
  platform_split=(${platform//\// })
  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}
  output_name="binaries/"$package'-'$GOOS'_'$GOARCH
  if [ $GOOS = "windows" ]; then
    output_name+='.exe'
  fi

  env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name oc2_tsp
  if [ $? -ne 0 ]; then
    echo 'An error has occurred! Aborting the scmkript execution...'
    exit 1
  fi
done
