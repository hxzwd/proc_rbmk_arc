GOBIN=/usr/local/go/bin
GOCMD=$(GOBIN)/go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build

CURRENT_DIR=$(shell pwd)
PROJECT_NAME=proc_rbmk_arc
PROJECT_SUBDIR=p_arc_2
SOURCE_DIR=$(GOPATH)/$(PROJECT_NAME)/$(PROJECT_SUBDIR)/src
SOURCES=$(shell find $(SOURCE_DIR)/*.go)
SOURCES_NOTDIR=$(notdir $(SOURCE_FILES))
#SOURCES=$(shell cd $(SOURCE_DIR); find *.go; cd $(CURRENT_DIR))
#SOURCES=$(shell find $(SOURCE_DIR)/*.go)
#SOURCE_FILES=$(shell ls $(SOURCE_DIR)/*.go)
#SOURCES=$(notdir $(SOURCE_FILES))
#SOURCES=$(foreach SOURCE,$(SOURCES),$(shell basename $(SOURCE)))
BINARY_NAME=proc_rbmk_data_bin

MAKE_ARGS=$(shell echo $(MAKEFLAGS) | sed -r 's/(.)/-\1 /g')
MAKE_GOALS=$(MAKECMDGOALS)
#GOAL=$(firstword $(shell echo $(MAKE_GOALS) | sed -r 's/print//g' | ))
TMP=$(shell echo $(MAKE_GOALS) | grep -Ev "print" | tr -s "\s*" "\s")
GOAL=$(firstword $(TMP))
TARGETS=run print




$(SOURCES): $(SOURCE_DIR)




.PHONY: run
run: $(SOURCES)
	@$(GORUN) $(SOURCES)

.PHONY: build
build: $(SOURCES)
	@$(GOBUILD) -o $(BINARY_NAME) $(shell echo '-gcflags="-N -l"') $(SOURCES)


