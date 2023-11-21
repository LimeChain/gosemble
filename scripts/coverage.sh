#!/bin/bash

pkgs=$(go list ./... | grep -v runtime)
deps=`echo ${pkgs} | tr ' ' ","`
echo "mode: atomic" > coverage.txt

for pkg in $pkgs; do
    go test --tags nonwasmenv -race -cover -coverpkg "$deps" -coverprofile profile.tmp -covermode atomic $pkg

    if [ -f profile.tmp ]; then
        tail -n +2 profile.tmp  >> coverage.txt # skip mode line from each coverage file and append it to coverage.txt
        rm profile.tmp
    fi
done