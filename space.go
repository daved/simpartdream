package main

import (
	"io"
	"os"

	"github.com/daved/simpartsim"
	"github.com/tgreiser/etherdream"
)

type space struct {
	simpartsim.Space
}

func (s *space) run(ps simpartsim.Particles, frames int) <-chan []simpartsim.Coords {
	csc := make(chan []simpartsim.Coords)

	go func() {
		ps.Reset()
		s.Space.Run(ps, frames, csc)

		defer close(csc)
	}()

	return csc
}

func (s *space) toStdout(ps simpartsim.Particles, frames int) error {
	csc := s.run(ps, frames)

	for cs := range csc {
		if err := dumpToStdout(os.Stdout, cs); err != nil {
			return err
		}
	}

	return nil
}

func (s *space) pointStream(ps simpartsim.Particles, frames int) etherdream.PointStream {
	psFn := func(w io.WriteCloser) {
		defer func() {
			_ = w.Close()
		}()

		for {
			csc := s.run(ps, frames)

			for cs := range csc {
				if err := dumpInPointStream(w, cs); err != nil {
					return
				}
			}
		}
	}

	return psFn
}
