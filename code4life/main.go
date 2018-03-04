package main

import (
	"fmt"
	"io"
	"os"
	"sort"
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

func debugf(format string, v ...interface{}) {
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

func ConnectRank(id RankID, cmt string) {
	cmd("CONNECT %d %s", id, cmt)
}

func ConnectMol(mol string) {
	cmd("CONNECT " + mol)
}

const MoleculeCount = 5
const ProjectHealth = 50

type Project [MoleculeCount]int

func healthForProject(pl Player, p Project, s Sample) float64 {
	//TODO: project is based on experience not on molecule used.
	//TODO: use player experience to weight health gain
	plExp := pl.Expertise
	for mol, name := range MolName {
		if s.ExpertiseGain == name {
			plExp[mol]++
		}
	}

	turn := 0
	for i := range plExp {
		turn += p[i] - plExp[i]
	}

	if turn > 0 {
		return ProjectHealth / float64(turn)
	}
	return ProjectHealth
}

func healthForProjects(pl Player, ps []Project, s Sample) float64 {
	health := 0.0
	for _, p := range ps {
		health += healthForProject(pl, p, s)
	}
	return health
}

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
const NoBody = -1

const printDebug = true

func debug(v ...interface{}) {
	if printDebug {
		fmt.Fprintln(os.Stderr, v...)
	}
}

func readProjects(r io.Reader) []Project {
	var projectCount int
	fmt.Fscan(r, &projectCount)
	debug(projectCount)

	projects := make([]Project, projectCount)
	for i := 0; i < projectCount; i++ {
		fmt.Fscan(r, &projects[i][A], &projects[i][B], &projects[i][C], &projects[i][D], &projects[i][E])
		debug(projects[i][A], projects[i][B], projects[i][C], projects[i][D], projects[i][E])
	}
	return projects
}

func readPlayers(r io.Reader) [2]Player {
	var p [2]Player

	for i := 0; i < 2; i++ {
		fmt.Fscan(r, &p[i].Target, &p[i].Eta, &p[i].Score,
			&p[i].Storage[A], &p[i].Storage[B], &p[i].Storage[C], &p[i].Storage[D], &p[i].Storage[E],
			&p[i].Expertise[A], &p[i].Expertise[B], &p[i].Expertise[C], &p[i].Expertise[D], &p[i].Expertise[E])
		debug(p[i].Target, p[i].Eta, p[i].Score,
			p[i].Storage[A], p[i].Storage[B], p[i].Storage[C], p[i].Storage[D], p[i].Storage[E],
			p[i].Expertise[A], p[i].Expertise[B], p[i].Expertise[C], p[i].Expertise[D], p[i].Expertise[E])
	}
	return p
}

func readAvailableMols(r io.Reader) Molecules {
	var available Molecules
	fmt.Fscan(r, &available[A], &available[B], &available[C], &available[D], &available[E])
	debug(available[A], available[B], available[C], available[D], available[E])
	return available
}

func readSamples(r io.Reader) []Sample {
	var sampleCount int
	fmt.Fscan(r, &sampleCount)
	debug(sampleCount)

	samples := make([]Sample, sampleCount)
	for i := 0; i < sampleCount; i++ {
		fmt.Fscan(r, &samples[i].ID, &samples[i].CarriedBy, &samples[i].Rank, &samples[i].ExpertiseGain, &samples[i].Health,
			&samples[i].MoleculeCost[A], &samples[i].MoleculeCost[B], &samples[i].MoleculeCost[C], &samples[i].MoleculeCost[D], &samples[i].MoleculeCost[E])
		debug(samples[i].ID, samples[i].CarriedBy, samples[i].Rank, samples[i].ExpertiseGain, samples[i].Health,
			samples[i].MoleculeCost[A], samples[i].MoleculeCost[B], samples[i].MoleculeCost[C], samples[i].MoleculeCost[D], samples[i].MoleculeCost[E])
	}
	return samples
}

func SampleHealth(p Player, projects []Project, s Sample) float64 {
	si := healthForProjects(p, projects, s)
	return float64(s.Health) + si
}

func main() {
	projects := readProjects(os.Stdin)
	for {
		start := time.Now()

		p := readPlayers(os.Stdin)

		available := readAvailableMols(os.Stdin)

		samples := readSamples(os.Stdin)
		sort.Slice(samples, func(i, j int) bool {
			return SampleHealth(p[0], projects, samples[i]) >= SampleHealth(p[0], projects, samples[j])
		})
		debugf("sample count: %d", len(samples))

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
		debugf("target is %s", p[0].Target)
		state, ok := states[p[0].Target]
		if !ok {
			StartGame(p[0], samples, available)
		} else {
			state(p[0], samples, available, projects)
		}
		end := time.Now()
		debugf("Turn completed in %dms", end.Sub(start).Nanoseconds()/1000000.0)
	}
}

var states = map[string]func(p Player, samples []Sample, available Molecules, projects []Project){
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

func (p Player) SelectRank() RankID {
	rank := 3
	totalExpertise := sum(p.Expertise[:])
	debugf("expertise is %d (total: %d)", p.Expertise, totalExpertise)

	for ; rank > 1; rank-- {
		if float64(Ranks[rank].CostMax)-float64(totalExpertise) < 5 {
			break
		}
	}
	return RankID(rank)
}

func SamplesState(p Player, samples []Sample, available Molecules, projects []Project) {
	if p.Eta != 0 {
		Wait()
		return
	}

	carried := sampleCarried(samples)
	debugf("%d samples carried", len(carried))
	if len(carried) < 3 {
		rank := p.SelectRank()
		debugf("ask undiagnosed samples target (rk %d)", rank)
		ConnectRank(rank, fmt.Sprintf("carrying %d", len(carried)))
	} else {
		Goto(DIAG)
	}
}

func DiagnosisState(p Player, samples []Sample, available Molecules, projects []Project) {
	if p.Eta != 0 {
		Wait()
		return
	}

	carried := sampleCarried(samples)
	undiagnosed := sampleUndiagnosed(carried, samples)

	if len(undiagnosed) > 0 {
		ConnectSample(samples[undiagnosed[0]].ID)
	} else {
		if len(carried) < 3 {
			// check cloud for completable sample
			uncarried := sampleUncarried(samples)
			possible := samplePossibleToComplete(p, uncarried, available, samples)
			if len(possible) > 0 {
				ConnectSample(samples[possible[0]].ID)
				return
			}
		}

		if len(carried) <= 0 {
			Goto(SAMP)
			return
		}

		// get uncomplete sample and put back those that cannot be completed
		uncompleted := sampleUncompleted(p, carried, samples)
		impossible := sampleImpossibleToComplete(p, uncompleted, available, samples)

		if len(impossible) > 0 {
			debugf("available: %d", available)
			debugf("impossible sample %d: %d", samples[impossible[0]].ID, samples[impossible[0]].MoleculeCost)
			ConnectSample(samples[impossible[0]].ID)
			return
		}

		// take better sample from cloud if possible
		uncarried := sampleUncarried(samples)
		possible := samplePossibleToComplete(p, uncarried, available, samples)
		if len(possible) > 0 {
			for i := 0; i < len(uncompleted); i++ {
				ps := SampleHealth(p, projects, samples[possible[0]])
				us := SampleHealth(p, projects, samples[uncompleted[i]])
				if ps > us {
					debugf("sample %d is better than %d (%f > %f)", samples[possible[0]].ID, samples[uncompleted[i]].ID, ps, us)
					ConnectSample(samples[uncompleted[i]].ID)
					return
				}
			}
		}
		Goto(MOLE)
	}
}

var waitInMolecules = 0

func MoleculesState(p Player, samples []Sample, available Molecules, projects []Project) {
	if p.Eta != 0 {
		Wait()
		return
	}

	carried := sampleCarried(samples)

	debugf("Storage is %d", p.Storage)
	if sum(p.Storage[:]) < 10 {
		// while there is less than 10 mol in store
		// if there is some incomplete samples, try to fullfill
		// as much sample as possible
		uncompleted := sampleUncompleted(p, carried, samples)
		debugf("%d uncompleted samples", len(uncompleted))
		for _, id := range uncompleted {
			s := samples[id]
			debugf("sample %d cost: %d", id, s.MoleculeCost)
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

		//TODO: take more molecule than needed to prevent opponent to fullfill
		// or to optimize for sample in cloud
	}

	completed := sampleCompleted(p, carried, samples)
	if len(completed) <= 0 {
		// check if i can switch with something from cloud
		// check cloud for completable sample
		uncarried := sampleUncarried(samples)
		possible := samplePossibleToComplete(p, uncarried, available, samples)
		if len(possible) > 0 {
			Goto(DIAG)
			return
		}

		if waitInMolecules > 3 {
			Goto(DIAG)
			waitInMolecules = 0
			return
		}

		Wait()
		waitInMolecules++
		return
	}
	Goto(LABO)
	return
}

func LaboratoryState(p Player, samples []Sample, available Molecules, projects []Project) {
	if p.Eta != 0 {
		Wait()
		return
	}

	carried := sampleCarried(samples)
	if len(carried) == 0 {
		// check cloud for completable sample
		uncarried := sampleUncarried(samples)
		possible := samplePossibleToComplete(p, uncarried, available, samples)
		interesting := sampleWithRank(p, possible, p.SelectRank(), samples)
		if len(interesting) > 1 {
			Goto(DIAG)
			return
		}

		Goto(SAMP)
		return
	}

	completed := sampleCompleted(p, carried, samples)
	if len(completed) == 0 {
		possible := samplePossibleToComplete(p, carried, available, samples)
		if len(possible) > 0 {
			Goto(MOLE)
			return
		}

		if len(possible) == 0 {
			uncarried := sampleUncarried(samples)
			possible := samplePossibleToComplete(p, uncarried, available, samples)
			interesting := sampleWithRank(p, possible, p.SelectRank(), samples)
			if len(interesting) > 1 {
				Goto(DIAG)
				return
			}
		}
		Goto(SAMP)
		return
	}

	ConnectSample(samples[completed[0]].ID)
}

func sampleImpossibleToComplete(p Player, carried []int, availables Molecules, samples []Sample) []int {
	result := make([]int, 0, len(carried))
	for _, id := range carried {
		s := samples[id]

		possible := true
		for mol, cost := range s.MoleculeCost {
			if !canComplete(p, mol, cost, availables) {
				possible = false
				break
			}
		}
		if !possible {
			result = append(result, id)
		}
	}
	return result
}

func canComplete(p Player, mol int, cost int, availables Molecules) bool {
	return p.Cost(mol, cost)-p.Storage[mol] <= availables[mol]
}

func samplePossibleToComplete(p Player, carried []int, availables Molecules, samples []Sample) []int {
	result := make([]int, 0, len(carried))
	for _, id := range carried {
		s := samples[id]

		possible := true
		for mol, cost := range s.MoleculeCost {
			if !canComplete(p, mol, cost, availables) {
				possible = false
				break
			}
		}
		if possible {
			result = append(result, id)
		}
	}
	return result
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

func sampleUncarried(samples []Sample) []int {
	uncarried := make([]int, 0, 3)
	for id, s := range samples {
		if s.CarriedBy == NoBody {
			uncarried = append(uncarried, id)
		}
	}
	return uncarried
}

func sampleWithRank(p Player, carried []int, rank RankID, samples []Sample) []int {
	uncarried := make([]int, 0, len(carried))
	for _, id := range carried {
		s := samples[id]
		if s.Rank == int(rank) {
			uncarried = append(uncarried, id)
		}
	}
	return uncarried
}

func mostAvailableMolecule(availables Molecules) int {
	max, imax := -1, -1
	for id, s := range availables {
		if s > max {
			max = s
			imax = id
		}
	}
	return imax
}
