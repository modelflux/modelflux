OUTDIR ?= $(CURDIR)/_output
BINARYNAME ?= tbd
PACKAGE := github.com/orelbn/tbd

build:
	go build -o $(OUTDIR)/bin/$(BINARYNAME) $(PACKAGE)/cmd

.DEFAULT_GOAL := build 