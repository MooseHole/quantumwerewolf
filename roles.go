package main

import (
	_ "github.com/lib/pq"
)

type Role struct {
	Name         string
	Id           int
	CanPeek      bool
	CanAttack    bool
	Evil         bool
	GoodMustKill bool
}

var roleTypes map[int]Role
var villager Role
var seer Role
var wolf Role

func setupRoles() {
	villager.Name = "Villager"
	villager.Id = 0
	villager.CanPeek = false
	villager.CanAttack = false
	villager.Evil = false
	villager.GoodMustKill = false

	seer.Name = "Seer"
	seer.Id = 1
	seer.CanPeek = true
	seer.CanAttack = false
	seer.Evil = false
	seer.GoodMustKill = false

	wolf.Name = "Wolf"
	wolf.Id = -1
	wolf.CanPeek = false
	wolf.CanAttack = true
	wolf.Evil = true
	wolf.GoodMustKill = true

	roleTypes = make(map[int]Role)
	roleTypes[villager.Id] = villager
	roleTypes[seer.Id] = seer
	roleTypes[wolf.Id] = wolf
}
