package quantumwerewolf_test

import (
	"quantumwerewolf/internal/app/quantumwerewolf"
	"testing"
)

func TestCollapseAll(t *testing.T) {
	quantumwerewolf.ResetVars()
	quantumwerewolf.Game.Name = "test"
	quantumwerewolf.Game.Number = 0
	quantumwerewolf.Game.Seed = 64907478
	quantumwerewolf.Game.RoundNum = 0
	quantumwerewolf.Game.RoundNight = true
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

	expectedUniverses := 6
	if len(quantumwerewolf.Multiverse.Universes) != expectedUniverses {
		t.Errorf("CreateMultiverse did not generate correct number of universes.  expected %d != actual %d", expectedUniverses, len(quantumwerewolf.Multiverse.Universes))
	}

	quantumwerewolf.Game.RoundNum = 1
	quantumwerewolf.Game.RoundNight = false
	quantumwerewolf.Players[0].Actions = "0%2^|0@1|"
	quantumwerewolf.Players[1].Actions = "0%2~|0@0|"
	quantumwerewolf.Players[2].Actions = "0%0~|0@1|"
	quantumwerewolf.ResetObservations()
	quantumwerewolf.FillObservations()
	quantumwerewolf.CollapseAll()
	expectedUniverses = 3
	if len(quantumwerewolf.Multiverse.Universes) != expectedUniverses {
		t.Errorf("CollapseAll did not generate correct number of universes.  expected %d != actual %d", expectedUniverses, len(quantumwerewolf.Multiverse.Universes))
	}

	quantumwerewolf.Game.RoundNum = 0
	quantumwerewolf.Game.RoundNight = true
	quantumwerewolf.Players[1].Actions = "0%2~|0@0|0#0|"
	quantumwerewolf.ResetObservations()
	quantumwerewolf.FillObservations()
	quantumwerewolf.CollapseAll()
	expectedUniverses = 1
	if len(quantumwerewolf.Multiverse.Universes) != expectedUniverses {
		t.Errorf("CollapseAll did not generate correct number of universes.  expected %d != actual %d", expectedUniverses, len(quantumwerewolf.Multiverse.Universes))
	}
}

func TestCollapseAllEarlyDeath(t *testing.T) {
	quantumwerewolf.ResetVars()
	quantumwerewolf.Game.Name = "zyc"
	quantumwerewolf.Game.Number = 4
	quantumwerewolf.Game.Seed = 1014938898
	playerA := quantumwerewolf.Player{}
	playerA.Name = "y"
	playerA.Num = 0
	playerA.Role.Totals = make(map[int]int)
	playerB := quantumwerewolf.Player{}
	playerB.Name = "c"
	playerB.Num = 1
	playerB.Role.Totals = make(map[int]int)
	playerC := quantumwerewolf.Player{}
	playerC.Name = "z"
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

	expectedUniverses := 6
	if len(quantumwerewolf.Multiverse.Universes) != expectedUniverses {
		t.Errorf("CreateMultiverse did not generate correct number of universes.  expected %d != actual %d", expectedUniverses, len(quantumwerewolf.Multiverse.Universes))
	}

	if quantumwerewolf.PlayerDeadPercent(playerB) != 0 {
		t.Errorf("Dead amount: %d instead of %d", quantumwerewolf.PlayerDeadPercent(playerB), 0)
	}

	quantumwerewolf.Game.RoundNum = 1
	quantumwerewolf.Game.RoundNight = false
	quantumwerewolf.Players[0].Actions = "0%1^|0@1|"
	quantumwerewolf.Players[1].Actions = "0%0~|0@2|"
	quantumwerewolf.Players[2].Actions = "0%0~|0@1|"
	quantumwerewolf.ResetObservations()
	quantumwerewolf.FillObservations()
	quantumwerewolf.CollapseAll()
	expectedUniverses = 3
	if len(quantumwerewolf.Multiverse.Universes) != expectedUniverses {
		t.Errorf("CollapseAll did not generate correct number of universes.  expected %d != actual %d", expectedUniverses, len(quantumwerewolf.Multiverse.Universes))
	}

	if quantumwerewolf.PlayerDeadPercent(playerB) != 100 {
		t.Errorf("Dead amount: %d instead of %d", quantumwerewolf.PlayerDeadPercent(playerB), 100)
	}
}
