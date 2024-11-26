module github.com/jak9708/gonummat

go 1.22.5

require (
	github.com/goccmack/gocc v0.0.0-20230228185258-2292f9e40198
	github.com/google/go-cmp v0.6.0
	golang.org/x/exp v0.0.0-20241108190413-2d47ceb2692f
	golang.org/x/tools v0.27.0
	gonum.org/v1/gonum v0.15.1
)

require golang.org/x/mod v0.22.0 // indirect

replace gonum.org/v1/gonum => github.com/jak9708/gonummat v0.0.2
