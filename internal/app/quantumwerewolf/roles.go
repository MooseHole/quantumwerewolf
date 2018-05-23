package quantumwerewolf

import (
	"image/color"
)

// Role holds attributes of each role type
type Role struct {
	Name          string
	ID            int
	CanPeek       bool
	CanAttack     bool
	Evil          bool
	EnemyMustKill bool
	DefaultAmount float32
	Color         color.RGBA
}

// RoleTypes is an interface
type RoleTypes struct {
	Names interface{}
}

var roleTypes map[int]Role
var villager Role
var seer Role
var wolf Role
var rolesAreSet bool

func setupRoles() {
	if rolesAreSet {
		return
	}

	villager.Name = "Villager"
	villager.ID = 0 // 0 is reserved for the default type, Villager
	villager.CanPeek = false
	villager.CanAttack = false
	villager.Evil = false
	villager.EnemyMustKill = false
	villager.DefaultAmount = 0
	villager.Color = color.RGBA{0, 255, 0, 255}

	seer.Name = "Seer"
	seer.ID = 1
	seer.CanPeek = true
	seer.CanAttack = false
	seer.Evil = false
	seer.EnemyMustKill = false
	seer.DefaultAmount = 1 // Defaults to 1 of these
	seer.Color = color.RGBA{0, 0, 255, 255}

	wolf.Name = "Wolf"
	wolf.ID = -1
	wolf.CanPeek = false
	wolf.CanAttack = true
	wolf.Evil = true
	wolf.EnemyMustKill = true
	wolf.DefaultAmount = 1.0 / 3.0 // Defaults to 1/3 of total players
	wolf.Color = color.RGBA{255, 0, 0, 255}

	roleTypes = make(map[int]Role)
	roleTypes[villager.ID] = villager
	roleTypes[seer.ID] = seer
	roleTypes[wolf.ID] = wolf
	rolesAreSet = true
}

func getRoleTypes() Role {
	return Role{"Test", -2, false, true, false, true, 1.0, color.RGBA{0, 0, 0, 0}}
}
