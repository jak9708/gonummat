# Gonum BLAS

[![go.dev reference](https://pkg.go.dev/badge/github.com/jak9708/gonummat/blas)](https://pkg.go.dev/github.com/jak9708/gonummat/blas)
[![GoDoc](https://godocs.io/github.com/jak9708/gonummat/blas?status.svg)](https://godocs.io/github.com/jak9708/gonummat/blas)

A collection of packages to provide BLAS functionality for the [Go programming
language](http://golang.org)

## Installation
```sh
  go get github.com/jak9708/gonummat/blas/...
```

## Packages

### blas

Defines [BLAS API](http://www.netlib.org/blas/blast-forum/cinterface.pdf) split in several
interfaces.

### blas/gonum

Go implementation of the BLAS API (incomplete, implements the `float32` and `float64` API).

### blas/blas64 and blas/blas32

Wrappers for an implementation of the double (i.e., `float64`) and single (`float32`)
precision real parts of the BLAS API.

```Go
package main

import (
	"fmt"

	"github.com/jak9708/gonummat/blas/blas64"
)

func main() {
	v := blas64.Vector{Inc: 1, Data: []float64{1, 1, 1}}
	v.N = len(v.Data)
	fmt.Println("v has length:", blas64.Nrm2(v))
}
```

### blas/cblas128 and blas/cblas64

Wrappers for an implementation of the double (i.e., `complex128`) and single (`complex64`) 
precision complex parts of the blas API.

Currently blas/cblas64 and blas/cblas128 require gonum.org/v1/netlib/blas.
