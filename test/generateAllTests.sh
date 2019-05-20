#! /bin/bash

if [ $# -ne 1 ]; then
   echo "[generateLot]: need parameters: dir: root directory to place tests"
   exit 0
fi

DIR=$1 #root directory path for the generated specs

python -m py_compile *.py
chmod 770 *.pyc
rm -r $1/clique/
rm -r $1/line/
rm -r $1/star/
rm -r $1/ring*
total=0
SPECS=(lotAcc lotAccReset lotAvg lotUntil lotAutosar lotBox lotEvent lotMode) #(lotAcc lotAccReset lotAvg lotUntil lotAutosar lotBox lotEvent lotMode)
TOPOS=(clique ring ringshort line star)  #clique ring ringshort line star
NMONS=(4)
for SPEC in "${SPECS[@]}"; do
    for TOPO in "${TOPOS[@]}" ; do
	for n in "${NMONS[@]}" ; do
	    #generate code for both EVAL and LAZY strategies and both centralised and decentralised
	    python ./generateLot.pyc $SPEC $DIR $TOPO 'lazy' 'decent' $n 0
	    python ./generateLot.pyc $SPEC $DIR $TOPO 'eval' 'decent' $n 0
	    total=$[$total+2]
	    if [ $TOPO == "ring" ] || [ $TOPO == "ringshort" ] ; then
	       for (( c=0; c<$NMONS; c++ )) do
		   python ./generateLot.pyc $SPEC $DIR $TOPO 'lazy' 'cent' $n $c
		   python ./generateLot.pyc $SPEC $DIR $TOPO 'eval' 'cent' $n $c
		   total=$[$total+2]
	       done
	    else
		python ./generateLot.pyc $SPEC $DIR $TOPO 'lazy' 'cent' $n 0
		python ./generateLot.pyc $SPEC $DIR $TOPO 'eval' 'cent' $n 0
		total=$[$total+2]
	    fi 
      	done #while of num monitors
    done #for of topos
done #for of specs
echo $total
#local
#./generateAllTests.sh ./generated

