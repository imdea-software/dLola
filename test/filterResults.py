import sys
import time
import json
import functools
import operator
import csv

class Result:
    def __init__(self, row):
        self.row = row
   
    def better(self, r):
        return self.row['totalMsgs'] < r['totalMsgs'] or self.row['totalPayload'] < r['totalPayload'] or self.row['memory'] < r['memory'] or self.row['avgDelay'] < r['avgDelay'] or self.row['avgSimplRounds'] < r['avgSimplRounds'] or self.row['totalRedirects'] < r['totalRedirects'] or self.row['maxDelay'] < r['maxDelay'] or self.row['maxSimplRounds'] < r['maxSimplRounds'] or self.row['minDelay'] < r['minDelay'] or self.row['minSimplRounds'] < r['minSimplRounds']


def filterCSV(inputt, outputt):
    with open(inputt) as csv_file:
        csv_reader = csv.DictReader(csv_file)
        with open(outputt, mode='w') as out:
            csvwr = csv.DictWriter(out, delimiter=',', quotechar='"', quoting=csv.QUOTE_MINIMAL, fieldnames=csv_reader.fieldnames)
            csvwr.writeheader()
            line_c = 0
            best={}
            for row in csv_reader:
                treatRow(row, csvwr, best)
                line_c += 1
            print line_c
            #print best
            for key in best: #clique, ring ...
                result = best[key]
                csvwr.writerow(result.row)
                
def treatRow(row, csvwr, best):
    if row["cent"] == "decent":
        csvwr.writerow(row)
    else:
        #print best
        topo = row["topo"]
        lazy = row["lazy"]
        nmons = row["nmons"]
        spec = row["spec"]
        tlen = row["tracelen"]
        b = getBest(best, topo, lazy, nmons,spec,tlen)
        if b == None or b.better(row): #if first occurrence or better store it
            #print 'storing better'
            insertRow(best, topo, lazy, nmons, spec, tlen, row)
        #else:
            #print 'not better'

def insertRow(best, topo, lazy, nmons, spec, tlen, row):
    best[topo+lazy+str(nmons)+spec+str(tlen)] = Result(row)

def getBest(best, topo, lazy, nmons, spec, tlen):
    #print topo + ',' + lazy + ',' + nmons + ',' + spec + ',' + tlen
    if topo+lazy+str(nmons)+spec+str(tlen)  in best:
        return best[topo+lazy+str(nmons)+spec+str(tlen)]
    else:
        return None
                        
#MAIN
if len(sys.argv) < 2:
    print "Need input and output file"
    exit(0)

print "Starting filtering"
filterCSV(sys.argv[1], sys.argv[2])
print "Filtering done"
"""
python -m py_compile *.py
chmod 770 *.pyc

python ./filterResults.pyc generated/results2019-04-24_f.txt generated/results2019-04-24_filtered.txt
"""
