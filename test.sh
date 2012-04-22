#!/bin/bash
export GOPATH=`pwd`
export GOBIN=`pwd`/bin

go install spxs

time ./bin/spxs --debug --procs=4 --conf=conf/spxs.json ref=data/dna.gen
