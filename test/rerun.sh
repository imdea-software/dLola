#! /bin/bash

function processDir {
    maxproc=$2
    proc=$2
    files="0"
    for file in $1/*; do
	#printf "###################\ncompiling and executing all tests\n###################\n"
	if [[ -d $file ]]; then
            # $f is a directory
	    #printf "$file is a directory\n"
	    processDir $file $proc &
	else
	    #printf "$file is a normal file\n"
	    files=$[$files+1]
	    if [[ $file == *.${EXT} ]]; then
		
		if [[ $proc -gt 0 ]]; then
		    processFile $1 $file &
		    #printf "processHS background $files\n"
		    proc=$[$proc-1]
		else
		    processFile $1 $file
		    #printf "processHS foreground $files\n"
		    proc=$[$proc+$maxproc]
		fi
	    fi
	fi
    done
    #printf "$files files in directory $1\n"
}

function processFile {
    dir=$1
    file=$2
    #echo $file
    DATESEC=$(date --rfc-3339='seconds')
    #DATE=$(date --rfc-3339='date')
    TEXT=""
    TEXT=$TEXT"[rerun.sh]: ($DATESEC) executing file $dir/$OFILE \n"
    #echo $PROGRAM $file $OPTIONS
    TEXT=$($PROGRAM $file $OPTIONS)
    #echo $TEXT
    #echo "python ./processOutput.pyc ${file} ${TEXT} ${TLEN}"
    RESULT=$(python ./processOutput.pyc "${file}" "${TEXT}" "${TLEN}")
    echo $RESULT
    #get exclusive lock of file before appending to it, had race conditions when I did not use this
    flock -x $CODEDIR/results$DATE.txt printf "${RESULT}\n" >> $CODEDIR/results$DATE.txt
}


if [ $# -lt 5 ] ; then
   echo "[rerun]: need root directory from which to start the search for files to compile and execute, num of processes, extension of the files, options and tracelen"
   exit 0
fi

CODEDIR=$1
PROC=$2
EXT=$3
OPTIONS=$4
TLEN=$5
OPTIONS="${OPTIONS} ${TLEN}"
PROGRAM="./dLolaSys " #$6 #generated/clique/num/eval/decent/10/lotAcc.txt #./dLolaSys
#printf $CODEDIR

if [ $# -ge 6 ] && [ $6 == "build" ]; then
    python -m py_compile processOutput.py &&
    chmod 770 *.pyc &&
    cd .. &&
    make &&
    cd ./test &&
    go build -buildmode=exe -o dLolaSys
fi
DATE=$(date --rfc-3339='date')
flock -x $CODEDIR/results$DATE.txt printf "topo,type,lazy,cent,nmons,spec,cnode,tracelen,totalMsgs,totalPayload,totalRedirects,maxDelay,avgDelay,minDelay,maxSimplRounds,avgSimplRounds,minSimplRounds,memory\n" > $CODEDIR/results$DATE.txt
processDir $CODEDIR $PROC $OPTIONS
#echo DONE
#run them as requests, NOT TRIGGERS, othw it will terminate as soon as the first trigger gets resolved!!! (and the performance won't be realistic)
#./rerun.sh /home/luismigueldanielsson/go/src/gitlab.software.imdea.org/luismiguel.danielsson/dLola/test/generated 2 spec "past req" 10 build
#./rerun.sh /home/luismigueldanielsson/go/src/gitlab.software.imdea.org/luismiguel.danielsson/dLola/test/generated 2 spec "past trigger" 10 build

#server
#compiling for the server in dLola/test where main.go is located
#go build -buildmode=exe -o dLolaSys
#scp dLolaSys rerun.sh generateAllTests.sh generateLot.py processOutput.py luismiguel.danielsson@zeus.software.imdea.org:~/RV/go/
#./generateAll.sh ./generated
#use just 8 processes in order not to exhaust the server resources (and not get errors therefore)
#./rerun.sh ~/RV/go/generated 6 spec "past req" 1000000
