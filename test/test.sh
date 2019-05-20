#! /bin/bash
if [ $# -lt 2 ]; then
    echo "need #processes and tlen"
    exit 0
fi
DATE=$(date --rfc-3339='date')
PROC=$1
TLEN=$2
./generateAllTests.sh generated/ &&
    ./rerun.sh ./generated $PROC spec "past req" $TLEN &&
    sleep 300 &&
    python ./filterResults.pyc generated/results${DATE}.txt generated/results${DATE}_filtered.txt

#./test 8 100000
