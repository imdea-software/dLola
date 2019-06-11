# dLola
Decentralized Stream Runtime Verification

## Abstract of the related paper

We study the problem of decentralized monitoring of stream
runtime verification specifications. Decentralized monitoring consists of
organizing a monitoring activity to be performed by distributed compo-
nents that communicate using a synchronous network, a communication
setting common in many cyber-physical systems like automotive CPSs.
Previous approaches to decentralized monitoring were restricted to LTL
and similar logics whose monitors compute Boolean verdicts. We present
here an algorithm that solves the decentralized monitoring problem for
the more general setting of stream runtime verification. Additionally,
our algorithm handles network topologies while previous work assumed
a network in which all nodes can communicate directly.
Our algorithm is able to reach verdicts efficiently by exploiting simpli-
fiers and advanced communication strategies. Finally, we present the
results of an empirical evaluation of an implementation and compare
the expressive power and efficiency against state-of-the-art decentralized
monitoring tools like Themis.

## Tool phases
This implementation in Go of Decentralized Lola consists logically in 5 phases:
- Parsing the language.
- Semantic checking (variables definitions and matching types).
- Well-formed specification (which imply well-defined) see paper for further explanations.
- Creation/definition of the topology and the placement of stream variables into monitors.
- Execution of the decentralized stream runtime verification online algorithm.