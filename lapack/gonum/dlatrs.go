// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"github.com/jak9708/gonummat/blas"
	"github.com/jak9708/gonummat/blas/blas64"
	"github.com/jak9708/gonummat/lapack"
)

// Dlatrs solves a triangular system of equations scaled to prevent overflow. It
// solves
//
//	A * x = scale * b if trans == blas.NoTrans
//	Aᵀ * x = scale * b if trans == blas.Trans
//
// where the scale s is set for numeric stability.
//
// A is an n×n triangular matrix. On entry, the slice x contains the values of
// b, and on exit it contains the solution vector x.
//
// If normin == true, cnorm is an input and cnorm[j] contains the norm of the off-diagonal
// part of the j^th column of A. If trans == blas.NoTrans, cnorm[j] must be greater
// than or equal to the infinity norm, and greater than or equal to the one-norm
// otherwise. If normin == false, then cnorm is treated as an output, and is set
// to contain the 1-norm of the off-diagonal part of the j^th column of A.
//
// Dlatrs is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dlatrs(uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, normin bool, n int, a []float64, lda int, x []float64, cnorm []float64) (scale float64) {
	switch {
	case uplo != blas.Upper && uplo != blas.Lower:
		panic(badUplo)
	case trans != blas.NoTrans && trans != blas.Trans && trans != blas.ConjTrans:
		panic(badTrans)
	case diag != blas.Unit && diag != blas.NonUnit:
		panic(badDiag)
	case n < 0:
		panic(nLT0)
	case lda < max(1, n):
		panic(badLdA)
	}

	// Quick return if possible.
	if n == 0 {
		return 1
	}

	switch {
	case len(a) < (n-1)*lda+n:
		panic(shortA)
	case len(x) < n:
		panic(shortX)
	case len(cnorm) < n:
		panic(shortCNorm)
	}

	upper := uplo == blas.Upper
	nonUnit := diag == blas.NonUnit

	smlnum := dlamchS / dlamchP
	bignum := 1 / smlnum
	scale = 1

	bi := blas64.Implementation()

	if !normin {
		if upper {
			cnorm[0] = 0
			for j := 1; j < n; j++ {
				cnorm[j] = bi.Dasum(j, a[j:], lda)
			}
		} else {
			for j := 0; j < n-1; j++ {
				cnorm[j] = bi.Dasum(n-j-1, a[(j+1)*lda+j:], lda)
			}
			cnorm[n-1] = 0
		}
	}
	// Scale the column norms by tscal if the maximum element in cnorm is greater than bignum.
	imax := bi.Idamax(n, cnorm, 1)
	var tscal float64
	if cnorm[imax] <= bignum {
		tscal = 1
	} else {
		tmax := cnorm[imax]
		// Avoid NaN generation if entries in cnorm exceed the overflow
		// threshold.
		if tmax <= math.MaxFloat64 {
			// Case 1: All entries in cnorm are valid floating-point numbers.
			tscal = 1 / (smlnum * tmax)
			bi.Dscal(n, tscal, cnorm, 1)
		} else {
			// Case 2: At least one column norm of A cannot be represented as
			// floating-point number. Find the offdiagonal entry A[i,j] with the
			// largest absolute value. If this entry is not +/- Infinity, use
			// this value as tscal.
			tmax = 0
			if upper {
				// A is upper triangular.
				for j := 1; j < n; j++ {
					tmax = math.Max(impl.Dlange(lapack.MaxAbs, j, 1, a[j:], lda, nil), tmax)
				}
			} else {
				// A is lower triangular.
				for j := 0; j < n-1; j++ {
					tmax = math.Max(impl.Dlange(lapack.MaxAbs, n-j-1, 1, a[(j+1)*lda+j:], lda, nil), tmax)
				}
			}
			if tmax <= math.MaxFloat64 {
				tscal = 1 / (smlnum * tmax)
				for j := 0; j < n; j++ {
					if cnorm[j] <= math.MaxFloat64 {
						cnorm[j] *= tscal
					} else {
						// Recompute the 1-norm without introducing Infinity in
						// the summation.
						cnorm[j] = 0
						if upper {
							for i := 0; i < j; i++ {
								cnorm[j] += tscal * math.Abs(a[i*lda+j])
							}
						} else {
							for i := j + 1; i < n; i++ {
								cnorm[j] += tscal * math.Abs(a[i*lda+j])
							}
						}
					}
				}
			} else {
				// At least one entry of A is not a valid floating-point entry.
				// Rely on Dtrsv to propagate Inf and NaN.
				bi.Dtrsv(uplo, trans, diag, n, a, lda, x, 1)
				return
			}
		}
	}

	// Compute a bound on the computed solution vector to see if bi.Dtrsv can be used.
	j := bi.Idamax(n, x, 1)
	xmax := math.Abs(x[j])
	xbnd := xmax
	var grow float64
	var jfirst, jlast, jinc int
	if trans == blas.NoTrans {
		if upper {
			jfirst = n - 1
			jlast = -1
			jinc = -1
		} else {
			jfirst = 0
			jlast = n
			jinc = 1
		}
		// Compute the growth in A * x = b.
		if tscal != 1 {
			grow = 0
			goto Solve
		}
		if nonUnit {
			grow = 1 / math.Max(xbnd, smlnum)
			xbnd = grow
			for j := jfirst; j != jlast; j += jinc {
				if grow <= smlnum {
					goto Solve
				}
				tjj := math.Abs(a[j*lda+j])
				xbnd = math.Min(xbnd, math.Min(1, tjj)*grow)
				if tjj+cnorm[j] >= smlnum {
					grow *= tjj / (tjj + cnorm[j])
				} else {
					grow = 0
				}
			}
			grow = xbnd
		} else {
			grow = math.Min(1, 1/math.Max(xbnd, smlnum))
			for j := jfirst; j != jlast; j += jinc {
				if grow <= smlnum {
					goto Solve
				}
				grow *= 1 / (1 + cnorm[j])
			}
		}
	} else {
		if upper {
			jfirst = 0
			jlast = n
			jinc = 1
		} else {
			jfirst = n - 1
			jlast = -1
			jinc = -1
		}
		if tscal != 1 {
			grow = 0
			goto Solve
		}
		if nonUnit {
			grow = 1 / (math.Max(xbnd, smlnum))
			xbnd = grow
			for j := jfirst; j != jlast; j += jinc {
				if grow <= smlnum {
					goto Solve
				}
				xj := 1 + cnorm[j]
				grow = math.Min(grow, xbnd/xj)
				tjj := math.Abs(a[j*lda+j])
				if xj > tjj {
					xbnd *= tjj / xj
				}
			}
			grow = math.Min(grow, xbnd)
		} else {
			grow = math.Min(1, 1/math.Max(xbnd, smlnum))
			for j := jfirst; j != jlast; j += jinc {
				if grow <= smlnum {
					goto Solve
				}
				xj := 1 + cnorm[j]
				grow /= xj
			}
		}
	}

Solve:
	if grow*tscal > smlnum {
		// Use the Level 2 BLAS solve if the reciprocal of the bound on
		// elements of X is not too small.
		bi.Dtrsv(uplo, trans, diag, n, a, lda, x, 1)
		if tscal != 1 {
			bi.Dscal(n, 1/tscal, cnorm, 1)
		}
		return scale
	}

	// Use a Level 1 BLAS solve, scaling intermediate results.
	if xmax > bignum {
		scale = bignum / xmax
		bi.Dscal(n, scale, x, 1)
		xmax = bignum
	}
	if trans == blas.NoTrans {
		for j := jfirst; j != jlast; j += jinc {
			xj := math.Abs(x[j])
			var tjj, tjjs float64
			if nonUnit {
				tjjs = a[j*lda+j] * tscal
			} else {
				tjjs = tscal
				if tscal == 1 {
					goto Skip1
				}
			}
			tjj = math.Abs(tjjs)
			if tjj > smlnum {
				if tjj < 1 {
					if xj > tjj*bignum {
						rec := 1 / xj
						bi.Dscal(n, rec, x, 1)
						scale *= rec
						xmax *= rec
					}
				}
				x[j] /= tjjs
				xj = math.Abs(x[j])
			} else if tjj > 0 {
				if xj > tjj*bignum {
					rec := (tjj * bignum) / xj
					if cnorm[j] > 1 {
						rec /= cnorm[j]
					}
					bi.Dscal(n, rec, x, 1)
					scale *= rec
					xmax *= rec
				}
				x[j] /= tjjs
				xj = math.Abs(x[j])
			} else {
				for i := 0; i < n; i++ {
					x[i] = 0
				}
				x[j] = 1
				xj = 1
				scale = 0
				xmax = 0
			}
		Skip1:
			if xj > 1 {
				rec := 1 / xj
				if cnorm[j] > (bignum-xmax)*rec {
					rec *= 0.5
					bi.Dscal(n, rec, x, 1)
					scale *= rec
				}
			} else if xj*cnorm[j] > bignum-xmax {
				bi.Dscal(n, 0.5, x, 1)
				scale *= 0.5
			}
			if upper {
				if j > 0 {
					bi.Daxpy(j, -x[j]*tscal, a[j:], lda, x, 1)
					i := bi.Idamax(j, x, 1)
					xmax = math.Abs(x[i])
				}
			} else {
				if j < n-1 {
					bi.Daxpy(n-j-1, -x[j]*tscal, a[(j+1)*lda+j:], lda, x[j+1:], 1)
					i := j + bi.Idamax(n-j-1, x[j+1:], 1)
					xmax = math.Abs(x[i])
				}
			}
		}
	} else {
		for j := jfirst; j != jlast; j += jinc {
			xj := math.Abs(x[j])
			uscal := tscal
			rec := 1 / math.Max(xmax, 1)
			var tjjs float64
			if cnorm[j] > (bignum-xj)*rec {
				rec *= 0.5
				if nonUnit {
					tjjs = a[j*lda+j] * tscal
				} else {
					tjjs = tscal
				}
				tjj := math.Abs(tjjs)
				if tjj > 1 {
					rec = math.Min(1, rec*tjj)
					uscal /= tjjs
				}
				if rec < 1 {
					bi.Dscal(n, rec, x, 1)
					scale *= rec
					xmax *= rec
				}
			}
			var sumj float64
			if uscal == 1 {
				if upper {
					sumj = bi.Ddot(j, a[j:], lda, x, 1)
				} else if j < n-1 {
					sumj = bi.Ddot(n-j-1, a[(j+1)*lda+j:], lda, x[j+1:], 1)
				}
			} else {
				if upper {
					for i := 0; i < j; i++ {
						sumj += (a[i*lda+j] * uscal) * x[i]
					}
				} else if j < n {
					for i := j + 1; i < n; i++ {
						sumj += (a[i*lda+j] * uscal) * x[i]
					}
				}
			}
			if uscal == tscal {
				x[j] -= sumj
				xj := math.Abs(x[j])
				var tjjs float64
				if nonUnit {
					tjjs = a[j*lda+j] * tscal
				} else {
					tjjs = tscal
					if tscal == 1 {
						goto Skip2
					}
				}
				tjj := math.Abs(tjjs)
				if tjj > smlnum {
					if tjj < 1 {
						if xj > tjj*bignum {
							rec = 1 / xj
							bi.Dscal(n, rec, x, 1)
							scale *= rec
							xmax *= rec
						}
					}
					x[j] /= tjjs
				} else if tjj > 0 {
					if xj > tjj*bignum {
						rec = (tjj * bignum) / xj
						bi.Dscal(n, rec, x, 1)
						scale *= rec
						xmax *= rec
					}
					x[j] /= tjjs
				} else {
					for i := 0; i < n; i++ {
						x[i] = 0
					}
					x[j] = 1
					scale = 0
					xmax = 0
				}
			} else {
				x[j] = x[j]/tjjs - sumj
			}
		Skip2:
			xmax = math.Max(xmax, math.Abs(x[j]))
		}
	}
	scale /= tscal
	if tscal != 1 {
		bi.Dscal(n, 1/tscal, cnorm, 1)
	}
	return scale
}
