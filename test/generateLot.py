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

def lotDecent(TOPO, TYPE, STREAM, LAZY, N, FIN, FOUT):
    s = TOPO
    i = 0
    while i < N:
        s += "\n@" + str(i) + "{\n"
        s += FIN(TYPE,STREAM, LAZY, i)
        s += FOUT(TYPE,STREAM, LAZY, i)
        s += "}\n"
        i += 1
    return s

def lotCent(TOPO, TYPE, STREAM, LAZY, N, FIN, FOUT):
    s = TOPO
    i = 0
    s += "\n@" + str(i) + "{\n" #start central node
    s += FIN(TYPE,STREAM, LAZY, i) #input central node
    while i < N: #centralized computations
        s += FOUT(TYPE,STREAM, LAZY, i)
        i += 1
    s += "}\n" #central node
    i = 1 #other inputs in monitors 1..N
    while i < N: #decentralized observations
        s += "\n@" + str(i) + "{\n"
        s += FIN(TYPE,STREAM, LAZY, i)
        s += "}\n" #no central node
        i += 1
    return s

def lotInput(TYPE, STREAM, LAZY, i):
    return "input " + TYPE + " " + STREAM + str(i) +" "+LAZY +"\n"

def lotAcciOutput(TYPE, STREAM, LAZY, i):
    return "output " + TYPE + " " + STREAM + str(i) + "r "+LAZY +" = " + STREAM + str(i) + "r[-1| 0] + " + STREAM + str(i) + "\n"


#acc + reset
def lotAccResetOutput(STREAM, LAZY, i, iin):
    return "output num accReset"+ str(i)+" "+LAZY +" = if reset"+ str(i) +" then 0 else accReset"+ str(i) + "[-1|0] + " + STREAM + str(iin) + "\n"

def lotDecentAccReset(TOPO, STREAM, LAZY, N):
    s = TOPO
    i = 0
    while i < N:
        if i %2 == 0:
            s += "\n@" + str(i) + "{\n"
            s += lotInput("num",STREAM, LAZY, i)
            #s += lotAcciOutput("num",STREAM, LAZY, i) #will be 'STREAMr'
            s += "}\n"
        else :
            s += "\n@" + str(i) + "{\n"
            s += lotInput("bool","reset", LAZY, i)
            s += lotAccResetOutput(STREAM, LAZY, i, i-1)
            s += "}\n"
        i += 1
    return s    

def lotCentAccReset(TOPO, STREAM, LAZY, N):
    s = TOPO
    i = 0
    s += "\n@" + str(i) + "{\n" #start central node
    s += lotInput("num",STREAM, LAZY, i) #input central node
    s += lotInput("bool","reset", LAZY, i)
    while i < N: #centralized computations
        s += lotAccResetOutput(STREAM, LAZY, i, i)
        i += 1
    s += "}\n" #central node
    i = 1 #other inputs in monitors 1..N
    while i < N: #decentralized observations
        s += "\n@" + str(i) + "{\n"
        s += lotInput("num",STREAM, LAZY, i)
        s += lotInput("bool","reset", LAZY, i)
        s += "}\n" #no central node
        i += 1
    return s

#AVG
def lotAvgOutput(TYPE, STREAM, LAZY, i):
    s = "define " + TYPE + " acc"+ str(i) + " "+LAZY +" = acc"+ str(i) + "[-1| 0] + " + STREAM + str(i) + "\n" #acc
    s += "define " + TYPE + " counter" + str(i) + " "+LAZY +" = counter"+ str(i) + "[-1| 0] + 1\n" #counter
    s += "output " + TYPE + " avg" + str(i) + " "+LAZY +" = acc" + str(i) + "/ counter" + str(i) + "\n" #avg
    return s


def lotDecentAvg(TOPO, TYPE,STREAM, LAZY, N):
    s = TOPO
    i = 0
    while i < N:
        s += "\n@" + str(i) + "{\n"
        s += lotInput(TYPE, STREAM, LAZY, i)
        s += lotAvgOutput(TYPE, STREAM, LAZY, i)
        s += "}\n"
        i += 1
    return s    

def lotCentAvg(TOPO, TYPE, STREAM, LAZY, N):
    s = TOPO
    i = 0
    s += "\n@" + str(i) + "{\n" #start central node
    s += lotInput(TYPE,STREAM, LAZY, i) #input central node
    while i < N: #centralized computations
        s += lotAvgOutput(TYPE,STREAM, LAZY, i)
        i += 1
    s += "}\n" #central node
    i = 1 #other inputs in monitors 1..N
    while i < N: #decentralized observations
        s += "\n@" + str(i) + "{\n"
        s += lotInput(TYPE,STREAM, LAZY, i)
        s += "}\n" #no central node
        i += 1
    return s

#Until
def lotUntilOutput(TYPE,STREAM, LAZY, i):
    return "output " + TYPE + " until" + str(i) + " "+LAZY +" = b" + str(i) + " or (a" +str(i)+" and until" + str(i) + "[-1|false])\n"

def lotDecentUntil(TOPO, TYPE, STREAM, LAZY, N):
    s = TOPO
    i = 0
    while i < N:
        s += "\n@" + str(i) + "{\n"
        s += lotInput(TYPE, "a", LAZY, i)
        s += lotInput(TYPE, "b", LAZY, i)
        s += lotUntilOutput(TYPE, STREAM, LAZY, i)
        s += "}\n"
        i += 1
    return s        

def lotCentUntil(TOPO, TYPE, STREAM, LAZY, N):
    s = TOPO
    i = 0
    s += "\n@" + str(i) + "{\n" #start central node
    s += lotInput(TYPE, "a", LAZY, i) #input central node
    s += lotInput(TYPE, "b", LAZY, i)
    while i < N: #centralized computations
        s += lotUntilOutput(TYPE,STREAM, LAZY, i)
        i += 1
    s += "}\n" #central node
    i = 1 #other inputs in monitors 1..N
    while i < N: #decentralized observations
        s += "\n@" + str(i) + "{\n"
        s += lotInput(TYPE, "a", LAZY, i)
        s += lotInput(TYPE, "b", LAZY, i)
        s += "}\n" #no central node
        i += 1
    return s

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
        spec = lotDecent(TOPO, TYPE, "a",LAZY, N, lotInput, lotAcciOutput)
    elif DECENT=="cent":
        spec = lotCent(TOPO, TYPE, "a",LAZY, N, lotInput, lotAcciOutput)

if SPEC == "lotAccReset":
    TYPE="num"
    if DECENT=="decent":
        spec = lotDecentAccReset(TOPO, "a",LAZY, N)
    elif DECENT=="cent":
        spec = lotCentAccReset(TOPO, "a",LAZY, N)

if SPEC == "lotAvg":
    TYPE="num"
    if DECENT=="decent":
        spec = lotDecentAvg(TOPO, TYPE, "a",LAZY, N)
    elif DECENT=="cent":
        spec = lotCentAvg(TOPO, TYPE,"a",LAZY, N)

if SPEC == "lotUntil":
    TYPE="bool"
    if DECENT=="decent":
        spec = lotDecentUntil(TOPO, TYPE, "a",LAZY, N)
    elif DECENT=="cent":
        spec = lotCentUntil(TOPO, TYPE,"a",LAZY, N)

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
