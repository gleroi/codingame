package main

import (
	"fmt"
	"os"
)

//import "os"

/**
 * Bring data on patient samples from the diagnosis machine to the laboratory with enough molecules to produce medicine!
 **/

func debug(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
}

func output(format string, v ...interface{}) {
	fmt.Fprintf(os.Stdout, format, v...)
}

func cmd(s string, v ...interface{}) {
	output(s+"\n", v...)
}

func Goto(target string) {
	cmd("GOTO " + target)
}

func ConnectSample(id int) {
	cmd("CONNECT %d", id)
}

func ConnectMol(mol string) {
	cmd("CONNECT " + mol)
}

const MoleculeCount = 5

type Player struct {
	Target    string
	Eta       int
	Score     int
	Storage   [MoleculeCount]int
	Expertise [MoleculeCount]int
}

type Molecules [MoleculeCount]int

const (
	A = 0
	B = 1
	C = 2
	D = 3
	E = 4
)

var MolName = [5]string{"A", "B", "C", "D", "E"}

type Sample struct {
	ID            int
	CarriedBy     int
	Rank          int
	ExpertiseGain string
	Health        int
	MoleculeCost  [MoleculeCount]int
}

const (
	DIAG = "DIAGNOSIS"
	MOLE = "MOLECULES"
	LABO = "LABORATORY"
)

const ME = 0

func main() {
	var projectCount int
	fmt.Scan(&projectCount)

	for i := 0; i < projectCount; i++ {
		var a, b, c, d, e int
		fmt.Scan(&a, &b, &c, &d, &e)
	}

	for {
		var p [2]Player

		for i := 0; i < 2; i++ {
			debug("reading player %d", i)
			fmt.Scan(&p[i].Target, &p[i].Eta, &p[i].Score,
				&p[i].Storage[A], &p[i].Storage[B], &p[i].Storage[C], &p[i].Storage[D], &p[i].Storage[E],
				&p[i].Expertise[A], &p[i].Expertise[B], &p[i].Expertise[C], &p[i].Expertise[D], &p[i].Expertise[E])
		}

		var available Molecules
		fmt.Scan(&available[A], &available[B], &available[C], &available[D], &available[E])

		var sampleCount int
		fmt.Scan(&sampleCount)

		samples := make([]Sample, sampleCount)
		for i := 0; i < sampleCount; i++ {
			fmt.Scan(&samples[i].ID, &samples[i].CarriedBy, &samples[i].Rank, &samples[i].ExpertiseGain, &samples[i].Health,
				&samples[i].MoleculeCost[A], &samples[i].MoleculeCost[B], &samples[i].MoleculeCost[C], &samples[i].MoleculeCost[D], &samples[i].MoleculeCost[E])
		}

		debug("sample count: %d", sampleCount)

		var carriedSample int
		if !sampleCarried(samples, &carriedSample) {
			debug("no samples carried")
			if p[0].Target != DIAG {
				debug("not on diag target")
				Goto(DIAG)
				continue
			}

			debug("on diag target")
			id := findBestFreeSample(samples)
			debug("get sample %d", id)
			ConnectSample(samples[id].ID)
			continue
		}

		debug("carrying some samples %d", carriedSample)
		if mol, ok := enoughMolecules(p[0], samples[carriedSample]); !ok {
			debug("no enough molecules")
			if p[0].Target != MOLE {
				Goto(MOLE)
				continue
			}
			ConnectMol(mol)
			continue
		}

		debug("enough molecule")
		if p[0].Target != LABO {
			Goto(LABO)
			continue
		}
		debug("send to lab %d", carriedSample)
		ConnectSample(samples[carriedSample].ID)
	}
}

func enoughMolecules(p Player, sid Sample) (string, bool) {
	debug("health %d", sid.Health)
	debug("storage %d", p.Storage)
	debug("cost %d", sid.MoleculeCost)
	for mol, cost := range sid.MoleculeCost {
		if cost > p.Storage[mol] {
			return MolName[mol], false
		}
	}
	return "", true
}

func findBestFreeSample(samples []Sample) int {
	health, bestID := -1, -1
	for id, s := range samples {
		if s.CarriedBy != -1 {
			continue
		}
		if s.Health > health {
			health = s.Health
			bestID = id
		}
	}
	return bestID
}

func sampleCarried(samples []Sample, carried *int) bool {
	for id, s := range samples {
		if s.CarriedBy == ME {
			*carried = id
			return true
		}
	}
	return false
}
