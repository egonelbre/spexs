#!/bin/bash

SPEXS=../bin/spexs2

time $SPEXS -stats \
	-conf=transcript/conf.json \
	inp=transcript/transcripts.inp
