package dLola

import (
	//	"errors"
	"encoding/json"
	"fmt"
	"math"
	//	"strconv"
)

type Time = int
type Id = int
type MsgType int

const (
	Res     MsgType = iota
	Req     MsgType = iota
	Trigger MsgType = iota
)

//instance Show MsgType where
//  show Res = "Res"
//  show Req = "Req"
//  show Trigger = "Trigger"

func (msg MsgType) String() string {
	// declare an array of strings operator counts how many items in the array (7)
	names := [...]string{
		"Res",
		"Req",
		"Trigger"}
	if msg < Res || msg > Trigger {
		return "Unknown"
	}
	return names[msg]
}

type Msg struct {
	Kind        MsgType
	Stream      InstStreamExpr
	Value       *InstExpr //instead of LolaType
	ResTime     *Time
	SimplRounds *SimplRounds
	//Resp *Resp
	Src Id
	Dst Id
}

func (msg Msg) String() string {
	value := ""
	if msg.Value != nil {
		val := *msg.Value
		value = val.Sprint()
	}
	resTime := ""
	if msg.ResTime != nil {
		resTime = fmt.Sprintf("%d", msg.ResTime)
	}
	simpl := ""
	if msg.SimplRounds != nil {
		simpl = fmt.Sprintf("%d", msg.SimplRounds)
	}
	return fmt.Sprintf("Msg{ kind = %s\nstream = %s\nvalue = %s\nresTime = %s\nsimplRounds = %s\nsrc = %d\ndst = %d\n}", msg.Kind.String(), msg.Stream.Sprint(), value, resTime, simpl, msg.Src, msg.Dst)
	//return fmt.Sprintf("Msg{ kind = %s\nstream = %s\nResp = %s\nsrc = %d\ndst = %d\n}", msg.Kind.String(), msg.Stream.Sprint(), Resp, msg.Src, msg.Dst)
}

func Equal(msg Msg, msg2 Msg) bool {
	return msg == msg2
}

/*payload of the msg in bits!!*/
func payload(msg Msg) int {
	payload := commonPayLoad(msg.Stream.GetName().Sprint())
	if msg.Value != nil {
		val := *msg.Value
		switch v := val.(type) {
		case InstTruePredicate:
			payload += 1
		case InstFalsePredicate:
			payload += 1
		case InstIntLiteralExpr:
			payload += 32
		case InstFloatLiteralExpr:
			payload += 32
		case InstStringLiteralExpr:
			payload += len(v.S) * 8
		default:

		}
	}
	return payload
}

/*
1 3-valued : kind
3 int : time of stream, src and dst ### resTime and simplRounds will not be counted because they are here fro profiling purposes only
1 string of characters of 8 bits, name of stream
*/
func commonPayLoad(s string) int {
	return 2 + 32*3 + 8*len(s)
}

type Received = []Msg                        //[Msg]
type Pending = []Msg                         //[Msg]
type Output = []Msg                          //[Msg]
type Requested = map[InstStreamExpr]struct{} //S.Set Stream

//type Type = Int // 0 Bool, 1 Float
type SimplRounds = int
type Eval = bool //won't need a Req to be sent to its destiny, will be flipped after creating the Msg to avoid duplicates
type Resp struct {
	Value       InstExpr //LolaType
	Eval        Eval
	ResTime     Time
	SimplRounds SimplRounds
	Ttl         Time
} //result, Eval|Lazy, time at which the result was obtained, #of calls to simplExp

func (r *Resp) Sprint() string {
	//return fmt.Sprintf("Resp{Value = %s, eval = %t, ResTime = %d, SimplRounds = %d, Ttl = %d}", r.Value.Sprint(), r.eval, r.ResTime, r.SimplRounds, r.Ttl)
	json, _ := json.Marshal(r)
	return string(json)
}

type Resolved struct {
	Stream InstStreamExpr
	Resp   Resp
}

func (r *Resolved) String() string {
	/*s := ""
	if r != nil {
		s = fmt.Sprintf("Resolved: {%s, %s}", r.stream.Sprint(), r.Resp.Sprint())
	}
	return s*/
	json, _ := json.Marshal(r)
	return string(json)
}

type Und struct {
	exp          InstExpr
	Eval         Eval
	SimplRounds  SimplRounds
	simplifiable bool //will be set to true at initialization and whenever something gets substituted, othw it will be false to avoid trying to simplify over and over the same expression without changes
}
type Unresolved struct {
	stream InstStreamExpr
	und    Und
}
type ExpEval struct {
	exp  InstExpr
	eval Eval
}

type RSet = map[InstStreamExpr]Resp //M.Map Stream Resp
type USet = map[InstStreamExpr]Und  //M.Map Stream Und
//type ExpSet = Spec                  //M.Map Stream (Exp, Eval)

type Metrics struct {
	NumMsgs        int
	SumPayload     int
	RedirectedMsgs int // part of numMsgs
	MaxDelay       *Resolved
	AvgDelay       float64
	MinDelay       *Resolved
	MaxSimplRounds *Resolved
	AvgSimplRounds float64
	MinSimplRounds *Resolved
	Memory         []int
}

func (m Metrics) String() string {
	return fmt.Sprintf("Metrics: {numMsgs: %d, sumPayload: %d, redirectedMsgs: %d\n MaxDelay: %s\n, avgDelay: %f\n, MinDelay: %s\n MaxSimplRounds: %s\n avgSimplRounds: %f\n MinSimplRounds: %s\n memory: %v}", m.NumMsgs, m.SumPayload, m.RedirectedMsgs, m.MaxDelay.String(), m.AvgDelay, m.MinDelay.String(), m.MaxSimplRounds.String(), m.AvgSimplRounds, m.MinSimplRounds.String(), m.Memory)
}

func (m Metrics) Short() string {
	return fmt.Sprintf("Metrics: {numMsgs: %d, sumPayload: %d, redirectedMsgs: %d,\n MaxDelay: %d, AvgDelay: %f, MinDelay: %d,\n MaxSimplRounds: %d, AvgSimplRounds: %f, MinSimplRounds: %d,\n memory: %v}", m.NumMsgs, m.SumPayload, m.RedirectedMsgs, m.MaxDelay.Resp.ResTime, m.AvgDelay, m.MinDelay.Resp.ResTime, m.MaxSimplRounds.Resp.SimplRounds, m.AvgSimplRounds, m.MinSimplRounds.Resp.SimplRounds, m.Memory)
}

type Monitor struct {
	nid           Id
	q             Received                            //input msgs queue
	i             []chan Resolved                     //input read at monitor(local)
	u             USet                                //unresolved expressions
	r             RSet                                //resolved expressions
	expr          Spec                                //specification, contains eval Streams
	pen           Pending                             //pending msgs to respond
	out           Output                              //outgoing msgs
	req           Requested                           //requested streams
	t             Time                                //actual tick
	routes        map[Id]Id                           //routes as map destiny nexthop
	delta         map[StreamName]Id                   //in which node the stream will be computed
	tracelen      int                                 //length of the input trace
	depGraph      DepGraphAdj                         //dependencies among streams
	dep           map[StreamName][]Id                 //Monitors that need the value of the stream (should be coherent to delta)
	trigger       []Resolved                          //resolved streams that were marked as triggers and will halt the execution of the system
	ttlMap        map[StreamName]Time                 //for Pruning R: ttl of each resolved stream (in R), will be decremented in each tick, when 0 it will be removed AT START >=1
	instStreamDep map[InstStreamExpr][]InstStreamExpr //list of all the other InstStreamExprs that an instanced stream need to be computed(without simplifying)
	metrics       Metrics                             //Metrics to measure performance
}

func (n Monitor) String() string {
	s := fmt.Sprintf("\n###############\n Node { nid = %d\n q = %s\n i = %v\n u = %v\n r = %s\n expr = %v\n pen = %v\n out = %v\n req = %v\n t = %d\n routes = %v\n delta = %v\n tracelen = %v\n"+
		" dep = %v\n trigger = %v\n ttlMap = %v\n instStreamDep = %v\n metrics = %v\n} ################\n",
		n.nid, n.q, n.i, printU(n.u), printR(n.r), PrettyPrintSpec(&(n.expr), ""), n.pen, n.out, n.req, n.t, n.routes, n.delta, n.tracelen, n.dep, n.trigger, n.ttlMap, n.instStreamDep, n.metrics)
	return s
}
func (m Monitor) triggered() bool {
	return len(m.trigger) > 0
}
func (m Monitor) finished() bool {
	return m.tracelen <= m.t && len(m.q) == 0 /*&& len(m.i) == 0*/ && len(m.pen) == 0 && len(m.out) == 0 //input will be consumed waiting until m.t == m.tracelen
}
func (m Monitor) computes(s StreamName) bool {
	return m.delta[s] == m.nid
}
func (m Monitor) isEval(s StreamName) bool {
	return m.expr.Output[s].Eval
}
func (m Monitor) isLazy(s StreamName) bool {
	return !m.isEval(s)
}
func printU(u USet) string {
	s := "map:["
	for istream, und := range u {
		s += fmt.Sprintf("%s : {%s, %t, %d, %t}; ", istream.Sprint(), und.exp.Sprint(), und.Eval, und.SimplRounds, und.simplifiable)
	}
	s += "]"
	return s
}
func printR(r RSet) string {
	s := "map:["
	for stream, resp := range r {
		s += fmt.Sprintf("%s : {%s, %t, %d, %d}; ", stream.Sprint(), resp.Value.Sprint(), resp.Eval, resp.ResTime, resp.SimplRounds)
	}
	s += "]"
	return s
}

func PrintMons(ms map[Id]*Monitor) string {
	s := ""
	for _, m := range ms {
		s += m.String()
	}
	return s
}

func NewMonitor(id, tracelen int, s Spec, received Received, routes map[Id]Id, delta map[StreamName]Id, depGraph DepGraphAdj, dep map[StreamName][]Id, inputChannels []chan Resolved, ttlMap map[StreamName]Time) Monitor {
	return Monitor{id, received, inputChannels, make(USet), make(RSet), s, Pending{}, Output{}, Requested{}, 0, routes, delta, tracelen, depGraph, dep, make([]Resolved, 0), ttlMap, make(map[InstStreamExpr][]InstStreamExpr), Metrics{0, 0, 0, nil, 0.0, nil, nil, 0.0, nil, make([]int, 0)}}
}

type Verdict struct {
	mons map[Id]*Monitor
	//totalMsgs, totalPayload, totalRedirects, maxDelay, maxSimplRounds int
	//maxDelayStream, maxSimplRoundsStream *Resolved
	Metrics  Metrics
	Triggers []Resolved
}

func (v Verdict) String() string {
	//return fmt.Sprintf("Verdict{mons = %s,\n metrics: %v,\n triggers: %v}", PrintMons(v.mons), v.metrics.String(), v.triggers)
	json, _ := json.Marshal(v)
	return string(json)
}
func (v Verdict) Short() string {
	//return fmt.Sprintf("Verdict{metrics: %v\ntriggers: %v}", v.metrics.Short(), v.triggers)
	//return fmt.Sprintf("Verdict: {%s, \ntriggers: %v}", v.metrics.Short(), v.triggers)
	json, _ := json.Marshal(v)
	return string(json)
}

func ConvergeCountTrigger(mons map[Id]*Monitor) Verdict {
	Converge(mons)
	totalMsgs := 0
	totalPayload := 0
	totalRedirects := 0
	/*var Maxdelay *Resolved
	var MaxSimplRounds *Resolved*/
	var Maxdelay *Resolved //Resolved{InstStreamFetchExpr{"s", -1}, Resp{InstIntLiteralExpr{0}, false, 0, 0, 0}}
	Avgdelay := 0.0
	var Mindelay *Resolved
	var Maxsimplrounds *Resolved //:= Resolved{InstStreamFetchExpr{"s", -1}, Resp{InstIntLiteralExpr{0}, false, 0, 0, 0}}
	Avgsimplrounds := 0.0
	var Minsimplrounds *Resolved
	var totalmemory []int
	triggers := make([]Resolved, 0)
	for _, m := range mons {
		//fmt.Printf("%s\n", m.String())
		totalMsgs += m.metrics.NumMsgs
		totalPayload += m.metrics.SumPayload
		totalRedirects += m.metrics.RedirectedMsgs
		//r := Resolved{InstStreamFetchExpr{"s", -1}, Resp{InstIntLiteralExpr{0}, false, 0, 0, 0}}
		//for _, m := range mons {
		if m.metrics.MaxDelay != nil && (Maxdelay == nil || Maxdelay.Resp.ResTime < m.metrics.MaxDelay.Resp.ResTime) {
			Maxdelay = m.metrics.MaxDelay
			//Maxdelay.Resp = Resp
		}
		Avgdelay += m.metrics.AvgDelay
		if m.metrics.MinDelay != nil && (Mindelay == nil || Mindelay.Resp.ResTime < m.metrics.MinDelay.Resp.ResTime) {
			Mindelay = m.metrics.MinDelay
			//mindelay.Resp = Resp
		}
		if m.metrics.MaxSimplRounds != nil && (Maxsimplrounds == nil || Maxsimplrounds.Resp.SimplRounds < m.metrics.MaxSimplRounds.Resp.SimplRounds) {
			Maxsimplrounds = m.metrics.MaxSimplRounds
			//Maxsimplrounds.Resp = Resp
		}
		Avgsimplrounds += m.metrics.AvgDelay
		if m.metrics.MinSimplRounds != nil && (Minsimplrounds == nil || Minsimplrounds.Resp.SimplRounds < m.metrics.MinSimplRounds.Resp.SimplRounds) {
			Minsimplrounds = m.metrics.MinSimplRounds
			//minsimplrounds.Resp = Resp
		}
		//TODO:addition of memories
		//fmt.Printf("Memory of mon: %d, %v", m.nid, m.metrics.memory)
		if len(totalmemory) == 0 {
			totalmemory = m.metrics.Memory
		} else {
			addMemories(totalmemory, m.metrics.Memory)
			m.metrics.Memory = make([]int, 0) //reset them to measure in the next tick
		}
		//}
		Avgdelay /= float64(len(mons))
		Avgsimplrounds /= float64(len(mons))
		triggers = append(triggers, m.trigger...)
	}
	return Verdict{mons, Metrics{totalMsgs, totalPayload, totalRedirects, Maxdelay, Avgdelay, Mindelay, Maxsimplrounds, Avgsimplrounds, Minsimplrounds, totalmemory}, triggers}
}

//both slices have the same length
func addMemories(totalmemory, memory []int) {
	//fmt.Printf("Adding memories: %v, %v\n", totalmemory, memory)
	for i, mt := range totalmemory {
		totalmemory[i] = mt + memory[i]
	}
}

func Converge(mons map[Id]*Monitor) {
	allfinished := false
	anytriggered := false
	cMons, cTicks := prepareMonitors(mons)
	//fmt.Printf("Converge: %s\n", PrintMons(mons))
	for !(anytriggered || allfinished) {
		//fmt.Printf("Should continue converging\n")
		Tick(mons, cMons, cTicks)
		anytriggered, allfinished = ShouldContinue(mons)
		//fmt.Printf("Should continue converging because %t, %t\n", anytriggered, allfinished)
	}
	//close(cTicks) //so go routines can shutdown properly
	//fmt.Printf("Finished\n")
	return
}

func ShouldContinue(mons map[Id]*Monitor) (bool, bool) {
	allfinished := true
	anytriggered := false
	for _, m := range mons {
		if m.triggered() {
			anytriggered = true
			break
		}
		allfinished = allfinished && m.finished()
	}
	return anytriggered, allfinished
}

func Tickn(mons map[Id]*Monitor, nticks int) {
	cMons, cTicks := prepareMonitors(mons)
	for i := 0; i < nticks; i++ {
		Tick(mons, cMons, cTicks)
	}
}

func prepareMonitors(mons map[Id]*Monitor) (chan *Monitor, []chan struct{}) {
	cMons := make(chan *Monitor, len(mons))
	cTicks := make([]chan struct{}, len(mons))
	for _, m := range mons {
		//m.process()       //sequential
		ci := make(chan struct{})
		cTicks[m.nid] = ci
		go processMon(m, cMons, ci) //process thread-safe
	}
	return cMons, cTicks
}

func Tick(mons map[Id]*Monitor, cMons chan *Monitor, cTicks []chan struct{}) {
	//fmt.Printf("Tick mons:%s\n", PrintMons(mons))
	nmons := len(mons)
	for i := 0; i < nmons; i++ {
		//fmt.Printf("Creating tick :%d\n", i)
		cTicks[i] <- struct{}{} //BLOCKING: false alternatives: buffered channel(some monitors run so much that stole the ticks of the others), buffered channel + id(each monitor will need to search for their tick: livelock!! REMEMBER in go context swap only occurs when a go routine gets blocked!!)
	}
	incomingQs := make(map[Id][]Msg)
	for i := 0; i < nmons; i++ { //retrieve the processed monitors
		newMon := <-cMons //BLOCKING if empty, non-blocking while it has elements
		classifyOut(newMon, incomingQs)
		mons[newMon.nid] = newMon //write back to the map
	}
	//fmt.Printf("incoming messages of each node:%v\n", incomingQs)
	for nid, m := range mons {
		//fmt.Printf("Before dispatch of mon %d:%s\n", m.nid, PrintMons(mons))
		//m.dispatch(mons, incomingQs)
		m.q = incomingQs[nid]
		//fmt.Printf("After dispatching of mon %d:%s\n", m.nid, PrintMons(mons))
	}
	//fmt.Printf("TICKED mons:%s\n", PrintMons(mons))
}

func processMon(m *Monitor, cMons chan *Monitor, cTicks chan struct{}) {
	open := true
	for open {
		_, open = <-cTicks //receive new tick BLOCKING, the channel need not be closed since it is not a buffered channel, the go routine will end when the main go routine returns
		if open {          //tick was received and channel still open
			//fmt.Printf("Tick received, processing...:%d\n", m.nid)
			m.process() //process
			cMons <- m  //send result should be NON-BLOCKING since it is a buffered channel
		}
	}
}

//prepare the msgs classifying them by their nextHop, clear m.out
func classifyOut(m *Monitor, incomingQs map[Id][]Msg) {
	//fmt.Printf("mon %d sends msgs:%v\n", m.nid, m.out)
	for _, msg := range m.out {
		nextHopMon := m.routes[msg.Dst]                              //we look for the nextHop in the route from m to msg.Dst
		incomingQs[nextHopMon] = append(incomingQs[nextHopMon], msg) //append msg to the incoming msgs of the destination
	}
	m.out = []Msg{} //clear out
}

/*func (m *Monitor) dispatch(mons map[Id]*Monitor) {
	for _, msg := range m.out {
		nextHopMon := m.routes[msg.Dst]                      //we look for the nextHop in the route from m to msg.Dst
		mons[nextHopMon].q = append(mons[nextHopMon].q, msg) //append msg to the incoming msgs of the destination
	}
	m.out = []Msg{}
}*/

func (m *Monitor) sendMsg(msg *Msg) {
	//fmt.Printf("Sending msg: %s\n", msg.String())
	m.out = append(m.out, *msg)
	m.metrics.NumMsgs++
	m.metrics.SumPayload += payload(*msg)
}

func (m *Monitor) process() {
	m.processQ()
	m.readInput()
	m.generateEquations()
	m.simplify()
	m.addRes()
	m.addReq()
	m.measureBeforePruning()
	m.pruneR()
	m.measureAfterPruning()
	m.t++
}

func (m *Monitor) processQ() {
	//fmt.Printf("[%d]:PROCESSQ: %s\n", m.nid, m.String())
	for _, msg := range m.q {
		if msg.Dst != m.nid { //redirect msgs whose dst is not this node
			m.sendMsg(&msg)
			m.metrics.RedirectedMsgs++
		} else {
			switch msg.Kind {
			case Res:
				m.r[msg.Stream] = msgToResp(&msg, m.ttlMap, &m.expr)
			case Req:
				m.pen = append(m.pen, msg)
			case Trigger:
				m.pen = append(m.pen, msg)
			}
		}
	}
	m.q = []Msg{}
}

func msgToResp(msg *Msg, ttlMap map[StreamName]Time, spec *Spec) Resp { //TODO: revise using ttl value as is, decrement by initialization time?
	ttl := 0
	if spec.isEval(msg.Stream.GetName()) { //if it is Eval the monitor will keep it, otw is lazy and the Unresolved eq is in U and the value need not be kept
		ttl = ttlMap[msg.Stream.GetName()]
	}
	/*r := *msg.Resp
	r.Eval = false //received responses are marked as LAZY, so as not to send them again and flood the net!!
	r.ttl = ttl*/
	return Resp{*msg.Value, false, *msg.ResTime, *msg.SimplRounds, ttl} //received responses are marked as LAZY, so as not to send them again and flood the net!!
}

func (m *Monitor) readInput() {
	//fmt.Printf("[%d]:READINPUT: %s\n", m.nid, m.String())
	//those events may be fed by some generator, a file, a socket... running in a go routine SEE inputReader.go
	if m.t < m.tracelen {
		for _, c := range m.i {
			r := <-c
			//fmt.Printf("Getting Input %v", r)
			m.r[r.Stream] = r.Resp
		}
	}

}
func (m *Monitor) generateEquations() {
	//fmt.Printf("[%d]:GENERATEEQ: %s\n", m.nid, m.String())
	if m.t < m.tracelen { //expressions will be intantiated from tick 0 to tracelen - 1
		for _, o := range m.expr.Output {
			if m.delta[o.Name] == m.nid {
				e := o.Expr
				i := e.InstantiateExpr(m.t, m.tracelen)
				//fmt.Printf("Instanced expr: %s\n", i.Sprint())
				u := Und{SimplifyExpr(i), o.Eval, 0, true} //Und gets generated with the eval specified and be able to simplify
				stream := InstStreamFetchExpr{o.Name, m.t}
				m.u[stream] = u
			}
		}
	}
}
func (m *Monitor) simplify() {
	//fmt.Printf("[%d]:SIMPLIFY: %s\n", m.nid, m.String())
	someSimpl := true //will control when to stop searching for substitutions and
	for someSimpl {
		//fmt.Printf("Again\n")
		someSimpl = false
		for stream, und := range m.u { //for each unresolved stream in U
			dep := m.depGraph[stream.GetName()]
			for _, d := range dep { //for each of its dependencies
				depStream := InstStreamFetchExpr{StreamName(d.Dest), stream.GetTick() + d.Weight} //build the depStream taking into account the tick of the stream and the weight of the dependency
				//fmt.Printf("[%d]need the stream %s to simplify\n", m.nid, depStream.Sprint())
				resp, ok := m.r[depStream]
				if ok { //we found the value of the dependency stream in R
					//fmt.Printf("[%d]and we have its value\n %s\n", m.nid, und.exp.Sprint())
					und.exp = und.exp.Substitute(depStream, resp.Value)
					und.simplifiable = true //set it to true so the expression will be simplified (if possible)
					//fmt.Printf("[%d]after subs %s\n", m.nid, und.exp.Sprint())
				}
			}
			if und.simplifiable {
				und.exp = SimplifyExpr(und.exp)
				und.SimplRounds++
				und.simplifiable = false
				m.u[stream] = und //store the substituted and simplified expression
				newresp, isResolved := toResp(stream, m.t, und, m.ttlMap)
				if isResolved {
					//fmt.Printf("[%d]New Resp: %v",m.nid, newresp)
					m.r[stream] = newresp //add it to R
					delete(m.u, stream)   //remove it from U
				}
				someSimpl = someSimpl || isResolved
			}
		}
	}
}

func toResp(stream InstStreamExpr, t int, und Und, ttlMap map[StreamName]Time) (Resp, bool) {
	//fmt.Printf("To Resp: %v\n", und.exp)
	var r Resp
	ttl := int(math.Max(0, float64(ttlMap[stream.GetName()]-(t-stream.GetTick())))) //time remaining to remove from R: ttl - max(0, now-instantiation)
	switch und.exp.(type) {
	case InstTruePredicate:
		return Resp{und.exp, und.Eval, t, und.SimplRounds, ttl}, true
	case InstFalsePredicate:
		return Resp{und.exp, und.Eval, t, und.SimplRounds, ttl}, true
	case InstIntLiteralExpr:
		return Resp{und.exp, und.Eval, t, und.SimplRounds, ttl}, true
	case InstFloatLiteralExpr:
		return Resp{und.exp, und.Eval, t, und.SimplRounds, ttl}, true
	case InstStringLiteralExpr:
		return Resp{und.exp, und.Eval, t, und.SimplRounds, ttl}, true
	default:
		//fmt.Printf("Not resolved")
		return r, false
	}
}

func (m *Monitor) addRes() {
	//fmt.Printf("[%d]:ADDRES: %s\n", m.nid, m.String())
	for stream, resp := range m.r {
		if resp.Eval && stream.GetTick() <= m.t && stream.GetTick() < m.tracelen { //EVAL streams
			if destinies, ok := m.dep[stream.GetName()]; ok {
				for _, d := range destinies {
					if d != m.nid {
						msg := createMsg(stream, &resp, m.nid, d) //i think the ref is producing the incorrect values as the resp for both msgs is the last one!!!
						//fmt.Printf("Creating Res msg of eval stream %s\n", msg.String())
						m.sendMsg(&msg)
					}
				}
				resp.Eval = false
				m.r[stream] = resp //we mark the resp as already sent to interested monitors, TODO: search for a way to do resp.Eval = false and make it persistent, so we avoid this unnecessary alloc
			}
		}
	}
	newPen := make([]Msg, 0)
	for _, penMsg := range m.pen { //LAZY streams need to be requested in order to send responses
		if resp, ok := m.r[penMsg.Stream]; ok { //note this msg will no longer be in pen
			newMsg := createMsg(penMsg.Stream, &resp, m.nid, penMsg.Src)
			//fmt.Printf("Resolved LAZY %s\n", newMsg.String())
			if newMsg.Dst != m.nid { //newMsg will be sent only if destiny is another monitor
				m.sendMsg(&newMsg)
			}
			if penMsg.Kind == Trigger {
				//fmt.Printf("Resolved Trigger %s\n", msg.Stream.Sprint())
				m.trigger = append(m.trigger, Resolved{penMsg.Stream, resp})
			}

		} else {
			newPen = append(newPen, penMsg) //if we do not have the resp, we will keep it in pen
		}
	}
	m.pen = newPen
}

func createMsg(stream InstStreamExpr, resp *Resp, id, dst Id) Msg {
	if resp == nil {
		return Msg{Req, stream, nil, nil, nil, id, dst}
	}
	val := resp.Value
	time := resp.ResTime
	simpl := resp.SimplRounds
	return Msg{Res, stream, &val, &time, &simpl, id, dst}
}
func (m *Monitor) addReq() { //TODO: think of extracting part to the offline
	//fmt.Printf("[%d]:ADDREQ: %s\n", m.nid, m.String())
	for _, p := range m.pen {
		if !m.expr.Output[p.Stream.GetName()].Eval { //we only need to analyze dependencies to create Reqs for LAZY streams
			createReqMsgsPen(p.Stream, m)
		}
	}
}

func createReqMsgsPen(stream InstStreamExpr, m *Monitor) {
	dependencies, found := m.instStreamDep[stream]
	if found {
		for _, d := range dependencies {
			createReqStream(d, m)
		}
	} else {
		//fmt.Printf("Creating REQS\n")
		dependencies := obtainDependencies(stream, m)
		for i := 0; i < len(dependencies); i++ {
			depStream := dependencies[i]
			//fmt.Printf("Adjacency %s\n", adj.Sprint())
			_, resolved := m.r[depStream]
			if !resolved { //if not resolved we need to Request it and continue analyzing sub-dependencies
				createReqStream(depStream, m)
				dependencies = addNextLevelDependencies(dependencies, obtainDependencies(depStream, m), m.r)
			}
			//fmt.Printf("next level dependencies after: %s\n", SprintStreams(dependencies))
		}
		m.instStreamDep[stream] = dependencies //save the dependencies of the stream
	}
}

func obtainDependencies(stream InstStreamExpr, m *Monitor) []InstStreamExpr {
	var dependencies []InstStreamExpr
	adjacencies := m.depGraph[stream.GetName()]
	uExpr, unresolved := m.u[stream]
	if unresolved { //get actual needed dependencies taking into account what was simplified
		dependencies = getUdependencies(stream, adjacencies, uExpr.exp)
	} else { //get the dependencies from the spec and analyze next level dependencies; case of not instantiated yet!
		dependencies = convertToStreams(stream, adjacencies, m.tracelen)
	} //othw should be resolved and the msg should have been responded in the addRes phase
	return dependencies
}

//will create a req msg for depStream and add it to out iff the stream is not in R, assigned to other monitor(delta), not previously requested, lazy and have been instanced
func createReqStream(depStream InstStreamExpr, m *Monitor) {
	_, resolved := m.r[depStream]
	_, requested := m.req[depStream]
	//fmt.Printf("Dependency could be intantiated tlen: %d\n%s\n !resolved %t, !requested %t, lazy %t, assigned to other monitor: %t\n", m.tracelen, depStream.Sprint(), !resolved, !requested, m.isLazy(depStream.GetName()), !m.computes(depStream.GetName()))
	if !resolved && !m.computes(depStream.GetName()) && m.isLazy(depStream.GetName()) && !requested && depStream.GetTick() <= m.t { //not in R, not assigned to this monitor, not eval and not already requested, allow request of streams that have not yet been instanced?
		//fmt.Printf("Creating Request: %s\n", depStream.Sprint())
		msg := createMsg(depStream, nil, m.nid, m.delta[depStream.GetName()])
		m.sendMsg(&msg)
		m.req[depStream] = struct{}{} //mark it as requested
	}
}
func convertToStreams(stream InstStreamExpr, adjacencies []Adj, tlen int) []InstStreamExpr {
	r := make([]InstStreamExpr, 0)
	for _, adj := range adjacencies {
		depTick := stream.GetTick() + adj.Weight
		if depTick >= 0 && depTick < tlen { //the dependency could be instantiated
			depStream := InstStreamFetchExpr{adj.Dest, depTick}
			r = append(r, depStream)
		}
	}
	return r
}

//will change dependencies
func addNextLevelDependencies(dependencies, candidates []InstStreamExpr, r RSet) []InstStreamExpr {
	//fmt.Printf("next level dependencies before: %s\ncandidates: %s\n", SprintStreams(dependencies), SprintStreams(candidates))
	for _, c := range candidates {
		_, resolved := r[c]
		if !elemStream(dependencies, c, EqInstStreamExpr) && !resolved { //add if not already present and not resolved(will also make its dependencies not be checked, since they're not useful)
			dependencies = append(dependencies, c)
		}
	}
	//fmt.Printf("next level dependencies after: %s\n", SprintStreams(dependencies))
	return dependencies
}

func (m *Monitor) measureBeforePruning() {
	//fmt.Printf("Measures: %v\n", m.metrics)
	Maxdelay := m.metrics.MaxDelay
	Avgdelay := 0.0
	Mindelay := m.metrics.MinDelay
	Maxsimplrounds := m.metrics.MaxSimplRounds
	Avgsimplrounds := 0.0
	Minsimplrounds := m.metrics.MinSimplRounds

	for stream, resp := range m.r {
		//fmt.Printf("R elem: %s, %v\n", stream.Sprint(), resp)
		if Maxdelay == nil || resp.ResTime > Maxdelay.Resp.ResTime {
			Maxdelay = &Resolved{stream, resp}
		}
		Avgdelay += float64(resp.ResTime)
		if Mindelay == nil || resp.ResTime < Mindelay.Resp.ResTime {
			Mindelay = &Resolved{stream, resp}
		}
		if Maxsimplrounds == nil || resp.SimplRounds > Maxsimplrounds.Resp.SimplRounds {
			Maxsimplrounds = &Resolved{stream, resp}
		}
		Avgsimplrounds += float64(resp.SimplRounds)
		if Minsimplrounds == nil || resp.SimplRounds < Minsimplrounds.Resp.SimplRounds {
			Minsimplrounds = &Resolved{stream, resp}
		}
	}
	m.metrics.MaxDelay = Maxdelay
	m.metrics.AvgDelay = Avgdelay
	m.metrics.MinDelay = Mindelay
	m.metrics.MaxSimplRounds = Maxsimplrounds
	m.metrics.AvgSimplRounds = Avgsimplrounds
	m.metrics.MinSimplRounds = Minsimplrounds
	//fmt.Printf("Updated Measures: %v\n", m.metrics)
}

//R Pruning
func (m *Monitor) pruneR() {
	for stream, resp := range m.r {
		if resp.Ttl == 0 {
			delete(m.r, stream)
		} else {
			resp.Ttl--
			m.r[stream] = resp
		}
	}
}

func (m *Monitor) measureAfterPruning() {
	//fmt.Printf("memory u: %v, r: %v\n", m.u, m.r)
	m.metrics.Memory = append(m.metrics.Memory, len(m.u)+len(m.r)) //R will be already pruned, a measure will be taken at each tick
}
