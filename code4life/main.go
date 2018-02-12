package main

import (
	"fmt"
	"os"
)

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

func Wait() {
	cmd("WAIT")
}

func ConnectSample(id Sid) {
	cmd("CONNECT %d", id)
}

func ConnectRank(id Rank) {
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

type Rank int

type Sid int

type Sample struct {
	ID            Sid
	CarriedBy     int
	Rank          int
	ExpertiseGain string // indicates the molecule for which expertise is gain
	Health        int
	MoleculeCost  [MoleculeCount]int
}

func (s Sample) Diagnosed() bool {
	return s.MoleculeCost[A] != -1
}

const (
	SAMP = "SAMPLES"
	DIAG = "DIAGNOSIS"
	MOLE = "MOLECULES"
	LABO = "LABORATORY"
)

const NoSample = -1
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

		/*
					TODO:
					 - movement may take more than one turn:
					 Robot Movement Matrix

			  					SAMPLES 	DIAGNOSIS 	MOLECULES 	LABORATORY
					Start area 		2 			2 			2 			2
					SAMPLES 		0 			3 			3 			3
					DIAGNOSIS 		3 			0 			3 			4
					MOLECULES 		3 			3 			0 			3
					LABORATORY 		3 			4 			3 			0

					Get samples:
					  - multiple rank : more rank -> more molecules and more points
					Diagnozed them
					  - determine the needed modulecules
					  - if one molecule > 5 -> Put in the cloud (MOLE only give 5 max)
					  - if TotalMolecule > 10 (amount a robot can hold) -> Put in the cloud (connect a diagnosed sample to DIAG)
					Collect molecule
					Send to labo
		*/

		var carriedSample int
		if !sampleCarried(samples, &carriedSample) {
			debug("no samples carried")

			id := findBestFreeSample(samples)

			if id == NoSample {
				debug("no samples available")
				if p[0].Target != SAMP {
					debug("not on samp target")
					Goto(SAMP)
					continue
				} else {
					debug("ask undiagnosed samples target")
					ConnectRank(Rank(1))
					continue
				}
			} else {
				if p[0].Target != DIAG {
					debug("want sample %s, but not on diag target", samples[id].ID)
					Goto(DIAG)
					continue
				}
				debug("get sample %d", id)
				ConnectSample(samples[id].ID)
				continue
			}
		}

		debug("carrying some samples %d", carriedSample)
		debug("sample %d is diagnozed: %t", carriedSample, samples[carriedSample].Diagnosed())

		if !samples[carriedSample].Diagnosed() {
			if p[0].Target != DIAG {
				debug("not on diag target")
				Goto(DIAG)
				continue
			} else {
				debug("diagonzed sample %d", samples[carriedSample].ID)
				ConnectSample(samples[carriedSample].ID)
				continue
			}
		} else {
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
