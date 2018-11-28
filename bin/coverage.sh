#!/bin/bash
set -e

echo 'mode: count' > profile.cov

for dir in $(find . -maxdepth 10 -not -path './vendor*' -not -path './.git*' -not -path '*/_*' -type d);
do
if ls $dir/*.go &> /dev/null; then
    go test -short -covermode=count -coverprofile=$dir/profile.tmp $dir
    if [ -f $dir/profile.tmp ]
    then
        cat $dir/profile.tmp | tail -n +2 >> profile.cov
        rm $dir/profile.tmp
    fi
fi
done

coverage=$(go tool cover -func profile.cov | tail -n 1 | awk '{print $3}' | sed -e 's/[%]//g')

rm profile.cov

echo "coverage: $coverage%"
if (( $(echo "$coverage < 75" | bc -l) )); then
    exit 1
fi
