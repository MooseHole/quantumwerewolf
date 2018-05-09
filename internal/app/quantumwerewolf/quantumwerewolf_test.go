package quantumwerewolf_test

import (
	"quantumwerewolf/internal/app/quantumwerewolf"
	"testing"
)

func TestCollapseForFixedRole(t *testing.T) {
	quantumwerewolf.ResetVars()
	quantumwerewolf.Game.Name = "test"
	quantumwerewolf.Game.Number = 0
	quantumwerewolf.Game.Seed = 64907478
	playerA := quantumwerewolf.Player{}
	playerA.Name = "A"
	playerA.Num = 0
	playerA.Role.Totals = make(map[int]int)
	playerB := quantumwerewolf.Player{}
	playerB.Name = "B"
	playerB.Num = 1
	playerB.Role.Totals = make(map[int]int)
	playerC := quantumwerewolf.Player{}
	playerC.Name = "C"
	playerC.Num = 2
	playerC.Role.Totals = make(map[int]int)
	quantumwerewolf.Players = append(quantumwerewolf.Players, playerA)
	quantumwerewolf.Players = append(quantumwerewolf.Players, playerB)
	quantumwerewolf.Players = append(quantumwerewolf.Players, playerC)
	quantumwerewolf.GameSetup.Universes = 6
	quantumwerewolf.GameSetup.Roles["Villager"] = 1
	quantumwerewolf.GameSetup.Roles["Seer"] = 1
	quantumwerewolf.GameSetup.Roles["Wolf"] = 1
	quantumwerewolf.CreateMultiverse()

	if len(quantumwerewolf.Multiverse.Universes) != 6 {
		t.Errorf("CreateMultiverse did not generate correct number of universes.  expected %d != actual %d", 6, len(quantumwerewolf.Multiverse.Universes))
	}

	quantumwerewolf.Players[0].Actions = "0%2^|0@1|"
	quantumwerewolf.Players[1].Actions = "0%2~|0@0|"
	quantumwerewolf.Players[2].Actions = "0%0~|0@1|"

	quantumwerewolf.CollapseForAttack(2)
	if len(quantumwerewolf.Multiverse.Universes) != 6 {
		t.Errorf("CollapseForAttack did not generate correct number of universes.  expected %d != actual %d", 6, len(quantumwerewolf.Multiverse.Universes))
	}

	quantumwerewolf.Players[1].Actions = "0%2~|0@0|0#0|"

	quantumwerewolf.CollapseForFixedRole(1, 0)
	if len(quantumwerewolf.Multiverse.Universes) != 2 {
		t.Errorf("collapseForFixedRole did not generate correct number of universes.  expected %d != actual %d", 2, len(quantumwerewolf.Multiverse.Universes))
	}
	quantumwerewolf.CollapseForAttack(2)
	if len(quantumwerewolf.Multiverse.Universes) != 2 {
		t.Errorf("CollapseForAttack did not generate correct number of universes.  expected %d != actual %d", 2, len(quantumwerewolf.Multiverse.Universes))
	}
}