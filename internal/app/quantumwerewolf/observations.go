package quantumwerewolf

import (
	"log"
	"strconv"
	"strings"
)

var dirtyObservations bool

// PeekObservation keeps track of who peeked at whom and the result
type PeekObservation struct {
	Subject int
	Round   int
	Target  int
	IsEvil  bool
	Pending bool
}

// AttackObservation keeps track of who attacked whom
type AttackObservation struct {
	Subject int
	Round   int
	Target  int
	Pending bool
}

// LynchObservation keeps track of who lynched whom
type LynchObservation struct {
	Subject int
	Round   int
	Target  int
	Pending bool
}

// KillObservation keeps track of who was killed
type KillObservation struct {
	Subject int
	Round   int
	Role    int
	Pending bool
}

var peekObservations []PeekObservation
var attackObservations []AttackObservation
var lynchObservations []LynchObservation
var killObservations []KillObservation

// ResetObservations destroys all saved observation instances
func ResetObservations() {
	dirtyObservations = true
	peekObservations = nil
	attackObservations = nil
	lynchObservations = nil
	killObservations = nil
}

// FillObservations fills all observations with current player actions
func FillObservations() {
	if !dirtyObservations {
		return
	}

	ResetObservations()
	for _, p := range players {
		actionStrings := strings.Split(p.Actions, tokenEndAction)
		for _, action := range actionStrings {
			fillPeekObservation(p.Num, action)
			fillAttackObservation(p.Num, action)
			fillLynchObservation(p.Num, action)
			fillKillObservation(p.Num, action)
		}
	}
	dirtyObservations = false
}

// CommitObservations removes pending from all observations
func CommitObservations() {
	for _, o := range peekObservations {
		o.Pending = false
		addPeekObservation(o)
	}
	for _, o := range attackObservations {
		o.Pending = false
		addAttackObservation(o)
	}
	for _, o := range lynchObservations {
		o.Pending = false
		addLynchObservation(o)
	}
	for _, o := range killObservations {
		o.Pending = false
		addKillObservation(o)
	}
}

// FillActionsWithObservations takes the existing observations and fills the player Actions with a representation of them
func FillActionsWithObservations() {
	for i := range players {
		players[i].Actions = ""
	}
	for _, o := range peekObservations {
		result := tokenGood
		if o.IsEvil {
			result = tokenEvil
		}
		pending := ""
		if o.Pending {
			pending = tokenPending
		}

		players[getPlayerIndex(getPlayerByNumber(o.Subject))].Actions += strconv.Itoa(o.Round) + tokenPeek + strconv.Itoa(o.Target) + result + pending + tokenEndAction
	}
	for _, o := range attackObservations {
		pending := ""
		if o.Pending {
			pending = tokenPending
		}
		players[getPlayerIndex(getPlayerByNumber(o.Subject))].Actions += strconv.Itoa(o.Round) + tokenAttack + strconv.Itoa(o.Target) + pending + tokenEndAction
	}
	for _, o := range lynchObservations {
		pending := ""
		if o.Pending {
			pending = tokenPending
		}
		players[getPlayerIndex(getPlayerByNumber(o.Subject))].Actions += strconv.Itoa(o.Round) + tokenLynch + strconv.Itoa(o.Target) + pending + tokenEndAction
	}
	for _, o := range killObservations {
		pending := ""
		if o.Pending {
			pending = tokenPending
		}
		players[getPlayerIndex(getPlayerByNumber(o.Subject))].Actions += strconv.Itoa(o.Round) + tokenKilled + strconv.Itoa(o.Role) + pending + tokenEndAction
	}
}

func addAttackObservation(newO AttackObservation) {
	temp := make([]AttackObservation, 0, len(attackObservations)+1)
	entryReplaced := false
	for _, o := range attackObservations {
		if o.Subject == newO.Subject && o.Round == newO.Round {
			entryReplaced = true
			temp = append(temp, newO)
		} else {
			temp = append(temp, o)
		}
	}
	if !entryReplaced {
		temp = append(temp, newO)
	}
	attackObservations = nil
	for _, o := range temp {
		attackObservations = append(attackObservations, o)
	}
}

func addPeekObservation(newO PeekObservation) {
	temp := make([]PeekObservation, 0, len(peekObservations)+1)
	entryReplaced := false
	for _, o := range peekObservations {
		if o.Subject == newO.Subject && o.Round == newO.Round {
			entryReplaced = true
			temp = append(temp, newO)
		} else {
			temp = append(temp, o)
		}
	}
	if !entryReplaced {
		temp = append(temp, newO)
	}
	peekObservations = nil
	for _, o := range temp {
		peekObservations = append(peekObservations, o)
	}
}

func addLynchObservation(newO LynchObservation) {
	temp := make([]LynchObservation, 0, len(lynchObservations)+1)
	entryReplaced := false
	for _, o := range lynchObservations {
		if o.Subject == newO.Subject && o.Round == newO.Round {
			entryReplaced = true
			temp = append(temp, newO)
		} else {
			temp = append(temp, o)
		}
	}
	if !entryReplaced {
		temp = append(temp, newO)
	}
	lynchObservations = nil
	for _, o := range temp {
		lynchObservations = append(lynchObservations, o)
	}
}

func addKillObservation(newO KillObservation) {
	temp := make([]KillObservation, 0, len(killObservations)+1)
	entryReplaced := false
	for _, o := range killObservations {
		if o.Subject == newO.Subject && o.Round == newO.Round {
			entryReplaced = true
			temp = append(temp, newO)
		} else {
			temp = append(temp, o)
		}
	}
	if !entryReplaced {
		temp = append(temp, newO)
	}
	killObservations = nil
	for _, o := range temp {
		killObservations = append(killObservations, o)
	}
}

func fillPeekObservation(subject int, action string) {
	indexOfActionToken := strings.Index(action, tokenPeek)

	// If not correct action type
	if indexOfActionToken < 0 {
		return
	}

	round, err := strconv.ParseInt(action[0:indexOfActionToken], 10, 64)
	if err != nil {
		log.Printf("Error converting round for peek observation: %v", err)
	}
	pending := strings.Contains(action, tokenPending)
	// Leave a space for the result token
	endOfTarget := len(action) - 1
	if pending {
		endOfTarget--
	}
	target, err := strconv.ParseInt(action[indexOfActionToken+1:endOfTarget], 10, 64)
	if err != nil {
		log.Printf("Error converting target for peek observation: %v", err)
	}

	observation := PeekObservation{}
	observation.Subject = subject
	observation.Round = int(round)
	observation.Target = int(target)
	observation.IsEvil = action[len(action)-1:len(action)] == tokenEvil
	observation.Pending = pending
	addPeekObservation(observation)
}

func fillAttackObservation(subject int, action string) {
	indexOfActionToken := strings.Index(action, tokenAttack)

	// If not correct action type
	if indexOfActionToken < 0 {
		return
	}

	round, err := strconv.ParseInt(action[0:indexOfActionToken], 10, 64)
	if err != nil {
		log.Printf("Error converting round for attack observation: %v", err)
	}
	pending := strings.Contains(action, tokenPending)
	endOfTarget := len(action)
	if pending {
		endOfTarget--
	}
	target, err := strconv.ParseInt(action[indexOfActionToken+1:endOfTarget], 10, 64)
	if err != nil {
		log.Printf("Error converting target for attack observation: %v", err)
	}

	observation := AttackObservation{}
	observation.Subject = subject
	observation.Round = int(round)
	observation.Target = int(target)
	observation.Pending = pending
	addAttackObservation(observation)
}

func fillLynchObservation(subject int, action string) {
	indexOfActionToken := strings.Index(action, tokenLynch)

	// If not correct action type
	if indexOfActionToken < 0 {
		return
	}

	round, err := strconv.ParseInt(action[0:indexOfActionToken], 10, 64)
	if err != nil {
		log.Printf("Error converting round for lynch observation: %v", err)
	}
	pending := strings.Contains(action, tokenPending)
	endOfTarget := len(action)
	if pending {
		endOfTarget--
	}
	target, err := strconv.ParseInt(action[indexOfActionToken+1:endOfTarget], 10, 64)
	if err != nil {
		log.Printf("Error converting target for lynch observation: %v", err)
	}

	observation := LynchObservation{}
	observation.Subject = subject
	observation.Round = int(round)
	observation.Target = int(target)
	observation.Pending = pending
	addLynchObservation(observation)
}

func fillKillObservation(subject int, action string) {
	indexOfActionToken := strings.Index(action, tokenKilled)

	// If not correct action type
	if indexOfActionToken < 0 {
		return
	}

	round, err := strconv.ParseInt(action[0:indexOfActionToken], 10, 64)
	if err != nil {
		log.Printf("Error converting round for kill observation: %v", err)
	}
	pending := strings.Contains(action, tokenPending)
	endOfRole := len(action)
	if pending {
		endOfRole--
	}
	role, err := strconv.ParseInt(action[indexOfActionToken+1:endOfRole], 10, 64)
	if err != nil {
		log.Printf("Error converting role for kill observation: %v", err)
	}

	observation := KillObservation{}
	observation.Subject = subject
	observation.Round = int(round)
	observation.Role = int(role)
	observation.Pending = pending
	addKillObservation(observation)
}
