package main

import (
	"fmt"
	"os"
	"time"
)

/**
 * Bring data on patient samples from the diagnosis machine to the laboratory with enough molecules to produce medicine!
 **/

func sum(s []int) int {
	acc := 0
	for _, v := range s {
		acc += v
	}
	return acc
}

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

func ConnectRank(id RankID) {
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

func (p Player) Cost(mol int, cost int) int {
	return cost - p.Expertise[mol]
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

type RankID int

type Rank struct {
	CostMin, CostMax int
}

var Ranks = []Rank{
	Rank{},
	Rank{CostMin: 3, CostMax: 5},
	Rank{CostMin: 4, CostMax: 8},
	Rank{CostMin: 7, CostMax: 14},
}

type Sid int

type Sample struct {
	ID            Sid
	CarriedBy     int
	Rank          int
	ExpertiseGain string // indicates the molecule for which expertise is gain
	Health        int
	MoleculeCost  Molecules
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
		start := time.Now()

		var p [2]Player

		for i := 0; i < 2; i++ {
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

					Start area 		2 			2 			2 			2
			  					SAMPLES 	DIAGNOSIS 	MOLECULES 	LABORATORY
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
					  - if one molecule > available molecule -> Put in the cloud (???)
					Collect molecule
					Send to labo
		*/
		debug("target is %s", p[0].Target)
		state, ok := states[p[0].Target]
		if !ok {
			StartGame(p[0], samples, available)
		} else {
			state(p[0], samples, available)
		}
		end := time.Now()
		debug("Turn completed in %fms", end.Sub(start).Nanoseconds()/1000000.0)
	}
}

var states = map[string]func(p Player, samples []Sample, available Molecules){
	SAMP: SamplesState,
	DIAG: DiagnosisState,
	MOLE: MoleculesState,
	LABO: LaboratoryState,
}

func StartGame(p Player, samples []Sample, available Molecules) {
	if p.Eta != 0 {
		Wait()
		return

	}
	Goto(SAMP)
}

func SamplesState(p Player, samples []Sample, available Molecules) {
	if p.Eta != 0 {
		Wait()
		return
	}

	carried := sampleCarried(samples)
	debug("%d samples carried", len(carried))
	if len(carried) < 3 {
		rank := 3
		totalExpertise := sum(p.Expertise[:])
		debug("expertise is %d (total: %d)", p.Expertise, totalExpertise)

		for ; rank > 1; rank-- {
			if Ranks[rank].CostMax-totalExpertise < 5 {
				break
			}
		}
		debug("ask undiagnosed samples target (rk %d)", rank)
		ConnectRank(RankID(rank))
	} else {
		Goto(DIAG)
	}
}

func DiagnosisState(p Player, samples []Sample, available Molecules) {
	if p.Eta != 0 {
		Wait()
		return
	}

	carried := sampleCarried(samples)
	undiagnosed := sampleUndiagnosed(carried, samples)

	if len(undiagnosed) > 0 {
		ConnectSample(samples[undiagnosed[0]].ID)
	} else {
		Goto(MOLE)
	}
}

func MoleculesState(p Player, samples []Sample, available Molecules) {
	if p.Eta != 0 {
		Wait()
		return
	}

	carried := sampleCarried(samples)

	debug("Storage is %d", p.Storage)
	if sum(p.Storage[:]) < 10 {
		uncompleted := sampleUncompleted(p, carried, samples)
		debug("%d uncompleted samples", len(uncompleted))

		for _, id := range uncompleted {
			s := samples[id]
			debug("sample %d cost: %d", id, s.MoleculeCost)
			debug("sample %d gain: %d", id, s.ExpertiseGain)
			for mol, cost := range s.MoleculeCost {
				if p.Cost(mol, cost)-p.Storage[mol] <= 0 {
					continue
				}
				if available[mol] <= 0 {
					continue
				}
				ConnectMol(MolName[mol])
				return
			}
		}
	}

	completed := sampleCompleted(p, carried, samples)
	if len(completed) <= 0 {
		Wait()
		return
	}
	Goto(LABO)
	return
}

func LaboratoryState(p Player, samples []Sample, available Molecules) {
	if p.Eta != 0 {
		Wait()
		return
	}

	carried := sampleCarried(samples)
	completed := sampleCompleted(p, carried, samples)

	if len(carried) == 0 {
		Goto(SAMP)
		return
	}
	if len(completed) == 0 {
		Goto(MOLE)
		return
	}

	health, bestId := -1, -1
	for _, id := range completed {
		s := samples[id]
		if s.Health > health {
			health = s.Health
			bestId = id
		}
	}
	ConnectSample(samples[bestId].ID)
}

func sampleCompleted(p Player, carried []int, samples []Sample) []int {
	result := make([]int, 0, len(carried))
	for _, id := range carried {
		s := samples[id]

		completed := true
		for mol, cost := range s.MoleculeCost {
			if p.Cost(mol, cost) > p.Storage[mol] {
				completed = false
				break
			}
		}
		if completed {
			result = append(result, id)
		}
	}
	return result
}

func sampleUncompleted(p Player, carried []int, samples []Sample) []int {
	result := make([]int, 0, len(carried))
	for _, id := range carried {
		s := samples[id]

		completed := true
		for mol, cost := range s.MoleculeCost {
			if p.Cost(mol, cost) > p.Storage[mol] {
				completed = false
				break
			}
		}
		if !completed {
			result = append(result, id)
		}
	}
	return result
}

func sampleUndiagnosed(ids []int, samples []Sample) []int {
	undiag := make([]int, 0, len(samples))
	for _, id := range ids {
		s := samples[id]
		if !s.Diagnosed() {
			undiag = append(undiag, id)
		}
	}
	return undiag
}

func sampleCarried(samples []Sample) []int {
	carried := make([]int, 0, 3)
	for id, s := range samples {
		if s.CarriedBy == ME {
			carried = append(carried, id)
		}
	}
	return carried
}
