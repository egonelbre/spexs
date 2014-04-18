#!/bin/bash

SPEXS=../spexs2

time $SPEXS -stats \
	-verbose=true \
	-conf=events/conf.json \
	inp=events/errors.txt \
	ref=events/events.txt