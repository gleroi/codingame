package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGotoConnectMol(t *testing.T) {
	examples := []struct {
		input string
		p0Mol Molecule
		err   error
	}{
		{
			input: `MOLECULES 0 184 2 1 0 0 1 5 1 2 3 3
		MOLECULES 0 227 2 0 0 0 0 2 4 2 4 4
		1 4 5 5 4`,
			p0Mol: B,
		},
		{
			input: `MOLECULES 0 184 2 2 0 0 1 5 1 2 3 3
		LABORATORY 2 227 2 0 0 0 0 2 4 2 4 4
		1 3 5 5 4`,
			p0Mol: C,
		},
		{
			input: `MOLECULES 0 184 2 2 1 0 1 5 1 2 3 3
		LABORATORY 1 227 2 0 0 0 0 2 4 2 4 4
		1 3 4 5 4`,
			p0Mol: C,
		},
		{
			input: `MOLECULES 0 184 2 2 2 0 1 5 1 2 3 3
		 LABORATORY 0 227 2 0 0 0 0 2 4 2 4 4
		 1 3 3 5 4`,
			p0Mol: A,
		},
	}

	for i := 0; i < len(examples)-1; i++ {
		r := strings.NewReader(examples[i].input)
		players := readPlayers(r)
		avail := readAvailableMols(r)
		current := Game{
			Players:   players,
			Available: avail,
		}

		r = strings.NewReader(examples[i+1].input)
		players = readPlayers(r)
		avail = readAvailableMols(r)
		expected := Game{
			Players:   players,
			Available: avail,
		}

		connectMol := ConnectMolAct{examples[i].p0Mol}
		next := Apply(current, connectMol, GotoAct{LaboratoryState})
		assert.Equal(t, expected, next)
	}
}

func TestConnectMolNotAvailable(t *testing.T) {
	examples := []struct {
		input string
		p0Mol Molecule
		err   error
	}{
		{
			input: `MOLECULES 0 184 2 2 2 0 1 5 1 2 3 3
		 LABORATORY 0 227 2 0 0 0 0 2 4 2 4 4
		 0 3 3 5 4`,
			p0Mol: A,
		},
	}

	for i := 0; i < len(examples)-1; i++ {
		r := strings.NewReader(examples[i].input)
		players := readPlayers(r)
		avail := readAvailableMols(r)
		current := Game{
			Players:   players,
			Available: avail,
		}

		connectMol := ConnectMolAct{examples[i].p0Mol}
		err := connectMol.Validate(current, 0)
		if err == nil {
			t.Errorf("expected an error, got nil")
		}
	}
}

func TestConnectMolOneAvailablePlayersTakeSame(t *testing.T) {
	examples := []struct {
		input string
		p0Mol Molecule
		err   error
	}{
		{
			input: `MOLECULES 0 184 2 2 2 0 1 5 1 2 3 3
		 LABORATORY 0 227 2 0 0 0 0 2 4 2 4 4
		 1 3 3 5 4`,
			p0Mol: A,
		},
	}

	for i := 0; i < len(examples)-1; i++ {
		r := strings.NewReader(examples[i].input)
		players := readPlayers(r)
		avail := readAvailableMols(r)
		current := Game{
			Players:   players,
			Available: avail,
		}

		connectMol := ConnectMolAct{examples[i].p0Mol}
		next := Apply(current, connectMol, connectMol)

		if next.Available[examples[i].p0Mol] != -1 {
			// when two players take the same unique remaining molecules both get one
			// But both will have to but it in laboratory to become available again
			t.Errorf("expected -1, got %d", next.Available[examples[i].p0Mol])
		}
	}
}
