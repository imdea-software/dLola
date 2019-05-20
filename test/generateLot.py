import os
import sys
import subprocess

def printDir(DIR, PREFIX, CNODE, TEXT):
    #make subdierectories if they dont already exist and put in there the file
    if not os.path.isdir(DIR):
        subprocess.call("mkdir -p " + DIR, shell=True) # will fail if they already exist

    destiny = DIR + "/"+ PREFIX + str(CNODE) + ".spec"
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

def lotCent(TOPO, TYPE, STREAM, LAZY, N, FIN, FOUT, CNODE):
    s = TOPO
    i = 0
    #s += "\n@" + str(i) + "{\n" #start central node
    #s += FIN(TYPE,STREAM, LAZY, i) #input central node
    #s += "}\n" #central node
    #i = 1 #other inputs in monitors 1..N
    while i < N: #decentralized observations
        s += "\n@" + str(i) + "{\n"
        s += FIN(TYPE,STREAM, LAZY, i)
        if i == CNODE: #centralized computations
            j = 0
            while j < N: 
                s += FOUT(TYPE,STREAM, LAZY, j)
                j += 1
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

def lotCentAccReset(TOPO, STREAM, LAZY, N, CNODE):
    s = TOPO
    i = 0
    #s += "\n@" + str(i) + "{\n" #start central node
    #s += lotInput("num",STREAM, LAZY, i) #input central node
    #s += "}\n" #central node
    #i = 1 #other inputs in monitors 1..N
    while i < N: #decentralized observations
        s += "\n@" + str(i) + "{\n"
        s += lotInput("num",STREAM, LAZY, i)
        if i == CNODE:#centralized computations
            j = 0
            while j < N: 
                s += lotInput("bool","reset", LAZY, j) #reset signals will be in the central node (so Lazy decent produces a great result)
                s += lotAccResetOutput(STREAM, LAZY, j, j)
                j += 1
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

def lotCentAvg(TOPO, TYPE, STREAM, LAZY, N, CNODE):
    s = TOPO
    i = 0
    #s += "\n@" + str(i) + "{\n" #start central node
    #s += lotInput(TYPE,STREAM, LAZY, i) #input central node
    #s += "}\n" #central node
    #i = 1 #other inputs in monitors 1..N
    while i < N: #decentralized observations
        s += "\n@" + str(i) + "{\n"
        s += lotInput(TYPE,STREAM, LAZY, i)
        if i == CNODE:
            j = 0
            while j < N: #centralized computations
                s += lotAvgOutput(TYPE,STREAM, LAZY, j)
                j += 1
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

def lotCentUntil(TOPO, TYPE, STREAM, LAZY, N, CNODE):
    s = TOPO
    i = 0
    #s += "\n@" + str(i) + "{\n" #start central node
    #s += lotInput(TYPE, "a", LAZY, i) #input central node
    #s += lotInput(TYPE, "b", LAZY, i)
    #s += "}\n" #central node
    #i = 1 #other inputs in monitors 1..N
    while i < N: #decentralized observations
        s += "\n@" + str(i) + "{\n"
        s += lotInput(TYPE, "a", LAZY, i)
        s += lotInput(TYPE, "b", LAZY, i)
        if i == CNODE:
            j = 0
            while j < N: #centralized computations
                s += lotUntilOutput(TYPE,STREAM, LAZY, j)
                j += 1
        s += "}\n" #no central node
        i += 1
    return s

def lotDecentAutosar(TOPO, TYPE, STREAM, LAZY, N):
    s= TOPO + '''\n//Monitor chassis system
@0{
const num direction_tolerance = 0.2
input num yaw ''' +LAZY + ''' //real wheel direction 100, 100 in percentage
input num steering ''' +LAZY + ''' //driver desired direction 100, 100
input bool drive_wheel_slip ''' +LAZY + '''
input num b1 ''' +LAZY + ''' //brakes front left
input num b2 ''' +LAZY + ''' //brakes front right
input num b3 ''' +LAZY + ''' //brake rear left
input num b4 ''' +LAZY + ''' //brakes rear right
'''+autosarCentStreams()+'''
}

//Monitor Engine
@1{
input num actual_throttle ''' +LAZY + '''
'''+autosarEngineStreams()+'''
}

//Monitor Transmission
@2{
input num actual_torque_distr ''' +LAZY + '''
'''+autosarTransmissionStreams()+'''
}

//Monitor PowerTrain Coordinator central monitor
@3{
'''+autosarPowerTrainStreams()+'''
}\n'''
    i = 4
    while i < N:
        s += "@" + str(i) + "{\n}\n"
        i+=1
    return s


def autosarCentStreams():
    return '''const num persistent_threshold = 25
define num direction_deviation ''' +LAZY + ''' = steering - yaw //this way we know the direction of the deviation
output bool activate_ESP ''' +LAZY + ''' = (direction_deviation > 0 and direction_deviation > direction_tolerance) or (direction_deviation < 0 and direction_deviation < direction_tolerance) or drive_wheel_slip
define bool under_steering_left ''' +LAZY + ''' = steering > 0 and yaw < steering + direction_tolerance //steering left but car is not moving enough to the left, to avoid an obstacle
define num count_under_steering_left ''' +LAZY + ''' = if under_steering_left then count_under_steering_left[-1|0] + 1 else 0
define bool persistent_under_steering_left ''' +LAZY + ''' = count_under_steering_left > persistent_threshold
define num effective_brake_left_rear ''' +LAZY + ''' = if persistent_under_steering_left then count_under_steering_left else 0


define bool over_steering_right ''' +LAZY + ''' = steering < 0 and yaw < steering + direction_tolerance //steering right but moving too much to the right, as reaction to the evasive maneouvre
define num count_over_steering_right ''' +LAZY + ''' = if over_steering_right then count_over_steering_right[-1|0] + 1 else 0
define bool persistent_over_steering_right ''' +LAZY + ''' = count_over_steering_right > persistent_threshold
define num effective_brake_right_rear ''' +LAZY + ''' = if persistent_over_steering_right then count_over_steering_right else 0

output num requested_throttle ''' +LAZY + ''' = if activate_ESP then requested_throttle[-1|0] - direction_deviation else 0 //or whatever correction needs to be applied to the throttle of the engine
output num requested_torque_distr ''' +LAZY + ''' = if activate_ESP then requested_torque_distr[-1|0] - direction_deviation else 0 // or whatever correction needs to be applied to the torque provided by the transmision
'''

def autosarEngineStreams():
    return '''//Monitor Engine
const num throttle_tolerance = 0.1
output bool correct_throttle ''' +LAZY + ''' = requested_throttle[-1|0]/actual_throttle <= throttle_tolerance
'''
def autosarTransmissionStreams():
    return '''//Monitor Transmission
const num torque_distr_tolerance = 0.1
output bool correct_torque_distr ''' +LAZY + ''' = requested_torque_distr[-1|0]/actual_torque_distr <= torque_distr_tolerance
'''
def autosarPowerTrainStreams():
    return '''//Monitor PowerTrain Coordinator central monitor
output bool all_correct ''' +LAZY + ''' = correct_throttle and correct_torque_distr
output bool box_all_correct ''' +LAZY + ''' = all_correct and box_all_correct[-1|true]
'''

def lotCentAutosar(TOPO, TYPE, STREAM, LAZY, N, CNODE):
    s= TOPO + '''\n
@0{
//Monitor chassis system
const num direction_tolerance = 0.2
input num yaw ''' +LAZY + ''' //real wheel direction 100, 100 in percentage
input num steering ''' +LAZY + ''' //driver desired direction 100, 100
input bool drive_wheel_slip ''' +LAZY + '''
input num b1 ''' +LAZY + ''' //brakes front left
input num b2 ''' +LAZY + ''' //brakes front right
input num b3 ''' +LAZY + ''' //brake rear left
input num b4 ''' +LAZY + ''' //brakes rear right
'''
    if CNODE == 0: 
        s += autosarCentStreams()+'''

        '''+ autosarEngineStreams()+'''

        '''+ autosarTransmissionStreams()+'''

        '''+ autosarPowerTrainStreams()
        
    s +='''}

//Monitor Engine
@1{
input num actual_throttle ''' +LAZY + '''
'''
    if CNODE == 1:
        s += autosarCentStreams()+'''

        '''+ autosarEngineStreams()+'''

        '''+ autosarTransmissionStreams()+'''

        '''+ autosarPowerTrainStreams()

    s +='''}

//Monitor Transmission
@2{
input num actual_torque_distr ''' +LAZY + '''
'''
    if CNODE == 2: 
        s += autosarCentStreams()+'''

        '''+ autosarEngineStreams()+'''

        '''+ autosarTransmissionStreams()+'''

        '''+ autosarPowerTrainStreams()
    s +='''}

//Monitor PowerTrain Coordinator central monitor
@3{
'''
    if CNODE == 3: 
        s += autosarCentStreams()+'''

        '''+ autosarEngineStreams()+'''

        '''+ autosarTransmissionStreams()+'''

        '''+ autosarPowerTrainStreams()

    s +='''}\n'''

    i = 4
    while i < N:
        s += "@" + str(i) + "{\n}\n"
        i+=1
    return s

def lotBoxOutput(TYPE, STREAM, LAZY, i):
    return "output " + TYPE + " box" + str(i) + " "+LAZY +" = "+ STREAM + str(i) + " and box" + str(i) + "[-1|true]\n"

def lotDecentBox(TOPO, TYPE, STREAM, LAZY, N):
    s = TOPO
    i = 0
    while i < N:
        s += "\n@" + str(i) + "{\n"
        s += lotInput(TYPE, STREAM, LAZY, i)
        s += lotBoxOutput(TYPE, STREAM, LAZY, i)
        s += "}\n"
        i += 1
    return s    
    

def lotCentBox(TOPO, TYPE, STREAM, LAZY, N, CNODE):
    s = TOPO
    i = 0
    while i < N: #decentralized observations
        s += "\n@" + str(i) + "{\n"
        s += lotInput(TYPE,STREAM, LAZY, i)
        if i == CNODE:#centralized computations
            j = 0
            while j < N: 
                s += lotBoxOutput(TYPE,STREAM, LAZY, j)
                j += 1
        s += "}\n" #no central node
        i += 1
    return s

def lotEventOutput(TYPE, STREAM, LAZY, i):
    return "output " + TYPE + " event" + str(i) + " "+LAZY +" = "+ STREAM + str(i) + " or event" + str(i) + "[-1|false]\n"

def lotDecentEvent(TOPO, TYPE, STREAM, LAZY, N):
    s = TOPO
    i = 0
    while i < N:
        s += "\n@" + str(i) + "{\n"
        s += lotInput(TYPE, STREAM, LAZY, i)
        s += lotEventOutput(TYPE, STREAM, LAZY, i)
        s += "}\n"
        i += 1
    return s    
    

def lotCentEvent(TOPO, TYPE, STREAM, LAZY, N, CNODE):
    s = TOPO
    i = 0
    while i < N: #decentralized observations
        s += "\n@" + str(i) + "{\n"
        s += lotInput(TYPE,STREAM, LAZY, i)
        if i == CNODE:#centralized computations
            j = 0
            while j < N: 
                s += lotEventOutput(TYPE,STREAM, LAZY, j)
                j += 1
        s += "}\n" #no central node
        i += 1
    return s

def lotModeOutput(TYPE, STREAM, LAZY, i, j):
    return "output " + TYPE + " m" + str(i) + " "+LAZY +" = if "+ STREAM + str(i) + " then a"+ str(j) +" and b"+ str(j) +" and c"+ str(j) +" else d"+ str(j) +" and e"+ str(j) +" and f"+ str(j) +"\n"

def lotDecentMode(TOPO, TYPE, STREAM, LAZY, N):
    s = TOPO
    i = 0
    while i < N:
        s += "\n@" + str(i) + "{\n"
        if i%2 == 0:
            s += lotInput(TYPE, "a", LAZY, i)
            s += lotInput(TYPE, "b", LAZY, i)
            s += lotInput(TYPE, "c", LAZY, i)
            s += lotInput(TYPE, "d", LAZY, i)
            s += lotInput(TYPE, "e", LAZY, i)
            s += lotInput(TYPE, "f", LAZY, i)
        else:
            s += lotInput(TYPE, STREAM, LAZY, i)
            s += lotModeOutput(TYPE, STREAM, LAZY, i,i-1)
        s += "}\n"
        i += 1
    return s    
    

def lotCentMode(TOPO, TYPE, STREAM, LAZY, N, CNODE):
    s = TOPO
    i = 0
    while i < N: #decentralized observations
        s += "\n@" + str(i) + "{\n"
        s += lotInput(TYPE, "a", LAZY, i)
        s += lotInput(TYPE, "b", LAZY, i)
        s += lotInput(TYPE, "c", LAZY, i)
        s += lotInput(TYPE, "d", LAZY, i)
        s += lotInput(TYPE, "e", LAZY, i)
        s += lotInput(TYPE, "f", LAZY, i)
        if i == CNODE:#centralized computations
            j = 0
            while j < N:
                s += lotInput(TYPE, STREAM, LAZY, j)
                s += lotModeOutput(TYPE,STREAM, LAZY, j, j)
                j += 1
        s += "}\n" #no central node
        i += 1
    return s

#MAIN
if len(sys.argv) < 7 :
   print "[generateLot]: need s : function to create Spec, topo: topology to use, lazy: lazy/eval, decent/cent, n : number of specs as parameters, cnode: central node for centralized specs"
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
CNODE = int(sys.argv[7])
SPEC=sys.argv[1]
if SPEC == "lotAcc":
    TYPE = "num"
    if DECENT=="decent":
        spec = lotDecent(TOPO, TYPE, "a",LAZY, N, lotInput, lotAcciOutput)
    elif DECENT=="cent":
        spec = lotCent(TOPO, TYPE, "a",LAZY, N, lotInput, lotAcciOutput, CNODE)

if SPEC == "lotAccReset":
    TYPE="num"
    if DECENT=="decent":
        spec = lotDecentAccReset(TOPO, "a",LAZY, N)
    elif DECENT=="cent":
        spec = lotCentAccReset(TOPO, "a",LAZY, N, CNODE)

if SPEC == "lotAvg":
    TYPE="num"
    if DECENT=="decent":
        spec = lotDecentAvg(TOPO, TYPE, "a",LAZY, N)
    elif DECENT=="cent":
        spec = lotCentAvg(TOPO, TYPE,"a",LAZY, N, CNODE)

if SPEC == "lotUntil":
    TYPE="bool"
    if DECENT=="decent":
        spec = lotDecentUntil(TOPO, TYPE, "a",LAZY, N)
    elif DECENT=="cent":
        spec = lotCentUntil(TOPO, TYPE,"a",LAZY, N, CNODE)

if SPEC == "lotAutosar":
    TYPE="num"
    if N >= 4: #needs at least 4 monitors
        if DECENT=="decent":
            spec = lotDecentAutosar(TOPO, TYPE, "a",LAZY, N)
        elif DECENT=="cent":
            spec = lotCentAutosar(TOPO, TYPE,"a",LAZY, N, CNODE)
    else:
        exit(0)

if SPEC == "lotBox":
    TYPE="bool"
    if DECENT=="decent":
        spec = lotDecentBox(TOPO, TYPE, "a",LAZY, N)
    elif DECENT=="cent":
        spec = lotCentBox(TOPO, TYPE,"a",LAZY, N, CNODE)

if SPEC == "lotEvent":
    TYPE="bool"
    if DECENT=="decent":
        spec = lotDecentEvent(TOPO, TYPE, "a",LAZY, N)
    elif DECENT=="cent":
        spec = lotCentEvent(TOPO, TYPE,"a",LAZY, N, CNODE)

if SPEC == "lotMode":
    TYPE="bool"
    if DECENT=="decent":
        spec = lotDecentMode(TOPO, TYPE, "mode",LAZY, N)
    elif DECENT=="cent":
        spec = lotCentMode(TOPO, TYPE,"mode",LAZY, N, CNODE)


printDir(DIR+"/"+TOPO+"/"+TYPE+"/"+LAZY+"/"+DECENT+"/"+str(N), SPEC, CNODE, spec)


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
