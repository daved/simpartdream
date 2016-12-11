package main

import (
	"flag"
	"fmt"
	"image/color"
	"io"
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
		Size:     100.0,
		Gravity:  9.81,
		Drag:     9.0,
	}

	flag.BoolVar(&stdout, "stdout", stdout, "to stdout")
	flag.IntVar(&parts, "parts", parts, "particle count")
	flag.IntVar(&frames, "frames", frames, "frame count")
	flag.Parse()

	spc := simpartsim.NewSimpleSpace(opts)
	ps := simpartsim.NewSimpleParticles(parts, spc.Termination())

	csc := make(chan []simpartsim.Coords)
	go func() {
		spc.Run(ps, frames, csc)
		defer close(csc)
	}()

	if stdout {
		if err := toStdout(csc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return
	}

	stream, err := pointStream(csc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Printf("Listening...\n")
	addr, _, err := etherdream.FindFirstDAC()
	if err != nil {
		log.Fatalf("Network error: %v", err)
	}

	log.Printf("Found DAC at %v\n", addr)

	dac, err := etherdream.NewDAC(addr.IP.String())
	if err != nil {
		log.Fatal(err)
	}
	defer dac.Close()

	debug := false
	dac.Play(stream, debug)

	_ = stream
}

func dumpToStdout(cs []simpartsim.Coords) error {
	var w io.Writer = os.Stdout

	for k := range cs {
		x, y := int(cs[k].X), int(cs[k].Y)
		bs := []byte(fmt.Sprintf("%d,%d\n", x, y))

		if _, err := w.Write(bs); err != nil {
			return err
		}
	}

	return nil
}

func toStdout(csc <-chan []simpartsim.Coords) error {
	for cs := range csc {
		if err := dumpToStdout(cs); err != nil {
			return err
		}
	}

	return nil
}

func dumpInPointStream(w io.Writer, cs []simpartsim.Coords) {
	for k := range cs {
		x, y := int(cs[k].X), int(cs[k].Y)
		c := color.RGBA{0xff, 0x00, 0x00, 0xff}

		_, _ = w.Write(etherdream.NewPoint(x*100, y*100, c).Encode())
	}
}

func pointStream(csc <-chan []simpartsim.Coords) (etherdream.PointStream, error) {
	ps := func(w io.WriteCloser) {
		defer func() {
			_ = w.Close()
		}()

		for cs := range csc {
			dumpInPointStream(w, cs)
		}
		for {

		}
	}

	return ps, nil
}
