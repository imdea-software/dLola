import sys
import time
import json
import functools
import operator
#import logic_str

def extractInfo(line, to_find, offset):
    result = ''    
    i = line.find(to_find) + offset
    while i < len(line) and line[i].isdigit():
        result = result + line[i]
        i += 1
    return result

def processOutput(filename, output, tracelen):
    """totalMsgs = extractInfo(output, "totalMsgs: ",11)
    totalPayload = extractInfo(output, "totalPayload: ",14)
    totalRedirects = extractInfo(output, "totalRedirects: ",16)
    maxDelay = extractInfo(output, "maxDelay: ",10)
    maxSimplRounds = extractInfo(output, "maxSimplRounds: ",16)"""
    verdict = json.loads(output)
    totalMsgs = verdict["Metrics"]["NumMsgs"]
    totalPayload = verdict["Metrics"]["SumPayload"]
    totalRedirects = verdict["Metrics"]["RedirectedMsgs"]
    maxDelay = verdict["Metrics"]["MaxDelay"]["Resp"]["SimplRounds"]
    avgDelay = verdict["Metrics"]["AvgDelay"]
    minDelay = verdict["Metrics"]["MinDelay"]["Resp"]["SimplRounds"]
    maxSimplRounds = verdict["Metrics"]["MaxSimplRounds"]["Resp"]["SimplRounds"]
    avgSimplRounds = verdict["Metrics"]["AvgSimplRounds"]
    minSimplRounds = verdict["Metrics"]["MinSimplRounds"]["Resp"]["SimplRounds"]
    memory = verdict["Metrics"]["Memory"]
    memorystr = list(map(lambda x : str(x),memory))
    memorystring = ""
    for m in memorystr:
        memorystring += m + ","
    memorystring = memorystring[:len(memorystring)-1]
    #memorystring = functools.reduce(operator.add, memorystr) #foldl
    #print memorystring
    (topo, tipe, lazy, cent, nmons, spec) = processFilename(filename)
    print convertCSV(topo, tipe, lazy, cent, nmons, spec, tracelen, str(totalMsgs), str(totalPayload), str(totalRedirects), str(maxDelay), str(avgDelay), str(minDelay), str(maxSimplRounds), str(avgSimplRounds), str(minSimplRounds), memorystring)

def processFilename(filename):
    i = filename.find("generated/") + 10
    topo = ""
    while i < len(filename) and filename[i] != "/":
        topo = topo + filename[i]
        i += 1
    i += 1
    tipe = ""
    while i < len(filename) and filename[i] != "/":
        tipe = tipe + filename[i]
        i += 1
    i += 1
    lazy = ""
    while i < len(filename) and filename[i] != "/":
        lazy = lazy + filename[i]
        i += 1
    i += 1
    cent = ""
    while i < len(filename) and filename[i] != "/":
        cent = cent + filename[i]
        i += 1
    i += 1
    nmons = ""
    while i < len(filename) and filename[i] != "/":
        nmons = nmons + filename[i]
        i += 1
    i += 1
    spec = ""
    while i < len(filename) and filename[i] != ".":
        spec = spec + filename[i]
        i += 1
    return (topo, tipe, lazy, cent, nmons, spec)

def convertCSV(topo, tipe, lazy, cent, nmons, spec, tracelen, totalMsgs, totalPayload, totalRedirects, maxDelay, avgDelay, minDelay, maxSimplRounds, avgSimplRounds, minSimplRounds, memory):
    #print topo+","+ tipe+","+ lazy+","+ cent+","+ nmons+","+ tracelen+","+ totalMsgs+","+ totalPayload+","+ totalRedirects+","+ maxDelay+","+ maxSimplRounds + "\n"
    return topo+","+ tipe+","+ lazy+","+ cent+","+ nmons+","+ spec + "," + tracelen+","+ totalMsgs+","+ totalPayload+","+ totalRedirects+","+ maxDelay+","+ avgDelay+","+ minDelay+","+ maxSimplRounds+ avgSimplRounds+","+ minSimplRounds+","+ memory

processOutput(sys.argv[1], sys.argv[2], sys.argv[3])
