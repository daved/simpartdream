package main

import (
	"fmt"
	"image/color"
	"io"

	"github.com/daved/simpartsim"
	"github.com/tgreiser/etherdream"
)

func dumpToStdout(w io.Writer, cs []simpartsim.Coords) error {
	for k := range cs {
		x, y := int(cs[k].X), int(cs[k].Y)
		bs := []byte(fmt.Sprintf("%d,%d\n", x, y))

		if _, err := w.Write(bs); err != nil {
			return err
		}
	}

	return nil
}

func dumpInPointStream(w io.Writer, cs []simpartsim.Coords) error {
	for k := range cs {
		x, y := int(cs[k].X), int(cs[k].Y)
		c := color.RGBA{0xff, 0xff, 0x00, 0xff}

		if _, err := w.Write(etherdream.NewPoint(x*100, y*100, c).Encode()); err != nil {
			return err
		}
	}

	return nil
}
