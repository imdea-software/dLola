import os
import sys
import subprocess

def printDir(DIR, PREFIX, TEXT):
    #make subdierectories if they dont already exist and put in there the file
    if not os.path.isdir(DIR):
        subprocess.call("mkdir -p " + DIR, shell=True) # will fail if they already exist

    destiny = DIR + "/"+ PREFIX + ".spec"
    print destiny
    with open(destiny, "w") as f:
        f.write(TEXT)

def lotDecent(TOPO, TYPE, STREAM, LAZY, MAX, FIN, FOUT):
    s = TOPO
    i = 0
    while i < MAX:
        s += "\n@" + str(i) + "{\n"
        s += FIN(TYPE,STREAM, LAZY, MAX, i)
        s += FOUT(TYPE,STREAM, LAZY, MAX, i)
        s += "}\n"
        i += 1
    return s

def lotCent(TOPO, TYPE, STREAM, LAZY, MAX, FIN, FOUT):
    s = TOPO
    i = 0
    s += "\n@" + str(i) + "{\n" #start central node
    s += FIN(TYPE,STREAM, LAZY, MAX, i) #input central node
    while i < MAX: #centralized computations
        s += FOUT(TYPE,STREAM, LAZY, MAX, i)
        i += 1
    s += "}\n" #central node
    i = 1 #other inputs in monitors 1..N
    while i < MAX: #decentralized observations
        s += "\n@" + str(i) + "{\n"
        s += FIN(TYPE,STREAM, LAZY, MAX, i)
        s += "}\n" #no central node
        i += 1        
    return s

def lotAcciInput(TYPE, STREAM, LAZY, MAX, i):
    return "input " + TYPE + " " + STREAM + str(i) +" "+LAZY +"\n"

def lotAcciOutput(TYPE, STREAM, LAZY, MAX, i):
    return "output " + TYPE + " " + STREAM + str(i) + "r "+LAZY +" = " + STREAM + str(i) + "r[-1| 0] + " + STREAM + str(i) + "\n"


#MAIN
if len(sys.argv) < 7 :
   print "[generateLot]: need s : function to create Spec, topo: topology to use, lazy: lazy/eval, n : number of specs as parameters"
   sys.exit(0)

DIR=sys.argv[2]
TOPO = sys.argv[3]
if  sys.argv[4] == "lazy":
    LAZY = "lazy"
else:
    LAZY = "eval"
if  sys.argv[5] == "decent":
    DECENT="decent"
else:
    DECENT="cent"
N = int(sys.argv[6])
SPEC=sys.argv[1]
if SPEC == "lotAcc":
    TYPE = "num"
    if DECENT=="decent":
        spec = lotDecent(TOPO, TYPE, "a",LAZY, N, lotAcciInput, lotAcciOutput)
    elif DECENT=="cent":
        spec = lotCent(TOPO, TYPE, "a",LAZY, N, lotAcciInput, lotAcciOutput)
printDir(DIR+"/"+TOPO+"/"+TYPE+"/"+LAZY+"/"+DECENT+"/"+str(N), SPEC, spec)


""""
python -m py_compile *.py
chmod 770 *.pyc
python ./generateLot.pyc lotAcc ./generated clique num eval decent 1
"""
"""
event_10(p)
event_10[0] = p[0] or ...p[10]
event_10[1] = p[1] or ...p[11]
"""
