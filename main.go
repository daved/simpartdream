package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/daved/simpartsim"
	"github.com/tgreiser/etherdream"
)

func main() {
	stdout := false
	parts := 20
	frames := 200
	opts := simpartsim.SimpleSpaceOptions{
		FrameLen: .1,
		Size:     1000.0,
		Gravity:  9.81,
		Drag:     9.0,
	}

	flag.BoolVar(&stdout, "stdout", stdout, "to stdout")
	flag.IntVar(&parts, "parts", parts, "particle count")
	flag.IntVar(&frames, "frames", frames, "frame count")
	flag.Parse()

	spc := space{simpartsim.NewSimpleSpace(opts)}
	ps := simpartsim.NewSimpleParticles(parts, spc.Termination())

	if stdout {
		if err := spc.toStdout(ps, frames); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	stream := spc.pointStream(ps, frames)

	//addr, _, err := etherdream.FindFirstDAC()
	//if err != nil {
	//	log.Fatalf("Network error: %v", err)
	//}
	//log.Printf("Found DAC at %v\n", addr)

	dac, err := etherdream.NewDAC("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	defer dac.Close()

	debug := false
	dac.Play(stream, debug)
}
