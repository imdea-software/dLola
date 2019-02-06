#! /bin/bash

if [ $# -ne 1 ]; then
   echo "[generateLot]: need parameters: dir: root directory to place tests"
   exit 0
fi

DIR=$1 #root directory path for the generated specs

python -m py_compile *.py
chmod 770 *.pyc

SPECS=(lotAcc)
TOPOS=(clique) #ring ringshort line star
for SPEC in "${SPECS[@]}"; do
    for TOPO in "${TOPOS[@]}" ; do
	for n in {5..5..5} ; do
	    #generate code for both EVAL and LAZY strategies and both centralised and decentralised
	    python ./generateLot.pyc $SPEC $DIR $TOPO 'lazy' 'cent' $n
	    python ./generateLot.pyc $SPEC $DIR $TOPO 'lazy' 'decent' $n
	    python ./generateLot.pyc $SPEC $DIR $TOPO 'eval' 'cent' $n
	    python ./generateLot.pyc $SPEC $DIR $TOPO 'eval' 'decent' $n
      	done #while of num monitors
    done #for of topos
done #for of specs
#local
#./generateAllTests ./generated
