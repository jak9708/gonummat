#!/bin/bash

# Reset the tree to the current commit to handle
# any writes during the build.
git reset --hard

go generate github.com/jak9708/gonummat/blas
go generate github.com/jak9708/gonummat/blas/gonum
go generate github.com/jak9708/gonummat/unit
go generate github.com/jak9708/gonummat/unit/constant
go generate github.com/jak9708/gonummat/graph/formats/dot
go generate github.com/jak9708/gonummat/graph/formats/rdf
go generate github.com/jak9708/gonummat/stat/card

if [ -n "$(git diff)" ]; then	
	git diff
	exit 1
fi

rm -rf *
git reset --hard
