package main

import (
	"context"
	"os"
	"runtime/trace"
)

func main() {
	traceFile, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer traceFile.Close()

	err = trace.Start(traceFile)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	ctx := context.Background()

	trace.WithRegion(ctx, "myFunction", func() {
		trace.Log(ctx, "category", "Starting myFunction")
		myFunction()
		trace.Log(ctx, "category", "Finished myFunction")
	})
}

func myFunction() {
	trace.WithRegion(context.Background(), "subFunction", func() {
		trace.Log(context.Background(), "category", "Inside myFunction")
	})
}
