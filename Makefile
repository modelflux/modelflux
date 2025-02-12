OUTDIR ?= $(CURDIR)/_output
BINARYNAME ?= mf
PACKAGE := github.com/modelflux/modelflux

build:
	go build -o $(OUTDIR)/bin/$(BINARYNAME) $(PACKAGE)/cmd

.DEFAULT_GOAL := build 