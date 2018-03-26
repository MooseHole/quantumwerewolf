package quantumwerewolf

// Role holds attributes of each role type
type Role struct {
	Name          string
	ID            int
	CanPeek       bool
	CanAttack     bool
	Evil          bool
	GoodMustKill  bool
	DefaultAmount float32
}

// RoleTypes is an interface
type RoleTypes struct {
	Names interface{}
}

var roleTypes map[int]Role
var villager Role
var seer Role
var wolf Role

func setupRoles() {
	villager.Name = "Villager"
	villager.ID = 0 // 0 is reserved for the default type, Villager
	villager.CanPeek = false
	villager.CanAttack = false
	villager.Evil = false
	villager.GoodMustKill = false
	villager.DefaultAmount = 0

	seer.Name = "Seer"
	seer.ID = 1
	seer.CanPeek = true
	seer.CanAttack = false
	seer.Evil = false
	seer.GoodMustKill = false
	villager.DefaultAmount = 1 // Defaults to 1 of these

	wolf.Name = "Wolf"
	wolf.ID = -1
	wolf.CanPeek = false
	wolf.CanAttack = true
	wolf.Evil = true
	wolf.GoodMustKill = true
	villager.DefaultAmount = 0.33 // Defaults to 33% of total players

	roleTypes = make(map[int]Role)
	roleTypes[villager.ID] = villager
	roleTypes[seer.ID] = seer
	roleTypes[wolf.ID] = wolf
}

func getRoleTypes() Role {
	return Role{"Test", -2, false, true, false, true, 1.0}
}
