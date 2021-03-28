package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
)

func ToRadians(deg float64) float64 {
	return float64(deg) * (math.Pi / 180.0)
}

type Vector struct {
	x float64
	y float64
}

func (v1 Vector) Add(v2 Vector) Vector {
	return Vector{x: v1.x + v2.x, y: v1.y + v2.y}
}

type Particle struct {
	loc            Vector
	heading        float64
	sensorAngle    float64
	sensorDistance float64
	sensorReadings [3]uint8
}

type Grid struct {
	rows int
	cols int
	data []uint8
}

func main() {
	grid := initializeGrid(256, 256)
	particles := initializeParticles(20, grid.rows, grid.cols)

	runSimulation(10, particles, grid)

}

func runSimulation(nIterations int, particles []Particle, grid Grid) {
	iteration := 0
	for {

		particles = readSensors(particles, grid)
		particles = rotateParticles(particles)
		particles = moveParticles(particles)
		grid = deposit(particles, grid)

		writePPM(grid, fmt.Sprintf("img/%03d.ppm", iteration))
		if iteration == nIterations {
			break
		}
		iteration++
	}
}

// readSensors of given particles
func readSensors(particles []Particle, grid Grid) []Particle {
	// TODO make sure to 360 + 1 = 1 (angle wrap)
	for _, particle := range particles {
		sensorAngles := []float64{particle.heading - particle.sensorAngle, particle.heading, particle.heading + particle.sensorAngle}
		for i, sensorAngle := range sensorAngles {
			translationVector := Vector{x: math.Cos(ToRadians(sensorAngle)), y: math.Sin(ToRadians(sensorAngle))}
			sensor := particle.loc.Add(translationVector)
			particle.sensorReadings[i] = grid.data[int(math.Round(sensor.x))+int(math.Round(sensor.y))*grid.cols]
		}
	}
	return particles
}

func rotateParticles(particles []Particle) []Particle {
	for _, particle := range particles {
		rotateParticle(&particle)
	}
	return particles
}

func rotateParticle(particle *Particle) {
	f := particle.sensorReadings[1]
	fl := particle.sensorReadings[0]
	fr := particle.sensorReadings[2]

	if f < fl && f > fr {
		return
	} else if f < fl && f < fr {
		particle.heading += (rand.Float64() - 0.5) * 10
	} else if fl < fr {
		particle.heading += particle.sensorAngle
	} else if fr < fl {
		particle.heading -= particle.sensorAngle
	} else {
		return
	}

}

func moveParticles(particles []Particle) []Particle {
	for _, particle := range particles {
		moveParticle(&particle)
	}
	return particles
}

func moveParticle(particle *Particle) {
	translationVector := Vector{x: math.Cos(ToRadians(particle.heading)), y: math.Sin(ToRadians(particle.heading))}
	particle.loc = particle.loc.Add(translationVector)
	// TODO is this normalized? move should be 1 pixel
}

func deposit(particles []Particle, grid Grid) Grid {
	// TODO deposit 5
	return grid
}

func diffuse(grid []uint8, factor int) {

}

func decay(grid []uint8, rate float64) {

}

func initializeGrid(rows int, cols int) Grid {
	gridData := make([]uint8, rows*cols)
	for i := 0; i < rows*cols; i++ {
		gridData[i] = 0
	}
	return Grid{rows: rows, cols: cols, data: gridData}
}

func initializeParticles(n int, width int, height int) []Particle {
	particles := []Particle{}
	for i := 0; i < n; i++ {
		particle := Particle{loc: Vector{x: rand.Float64() * float64(width), y: rand.Float64() * float64(height)}, heading: rand.Float64() * 360, sensorAngle: 45, sensorDistance: 9}
		particles = append(particles, particle)
	}
	return particles
}

func writePPM(grid Grid, fp string) {
	f, err := os.Create(fp)
	if err != nil {
		log.Fatal(err)
		f.Close()
	}
	fmt.Fprintln(f, "P2")
	fmt.Fprintf(f, "%d %d\n", grid.rows, grid.cols)
	fmt.Fprintln(f, "256")

	for r := 0; r < grid.rows; r++ {
		for c := 0; c < grid.cols; c++ {
			fmt.Fprintf(f, "%d ", grid.data[c+r*grid.cols])
		}
		fmt.Fprintf(f, "\n")
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}
