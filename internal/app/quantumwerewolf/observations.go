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
}

// AttackObservation keeps track of who attacked whom
type AttackObservation struct {
	Subject int
	Round   int
	Target  int
}

// LynchObservation keeps track of who lynched whom
type LynchObservation struct {
	Subject int
	Round   int
	Target  int
}

// KillObservation keeps track of who was killed
type KillObservation struct {
	Subject int
	Round   int
	Role    int
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
	target, err := strconv.ParseInt(action[indexOfActionToken+1:len(action)-1], 10, 64)
	if err != nil {
		log.Printf("Error converting target for peek observation: %v", err)
	}

	observation := PeekObservation{}
	observation.Subject = subject
	observation.Round = int(round)
	observation.Target = int(target)
	observation.IsEvil = action[len(action)-1:len(action)] == tokenEvil
	peekObservations = append(peekObservations, observation)
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
	target, err := strconv.ParseInt(action[indexOfActionToken+1:len(action)], 10, 64)
	if err != nil {
		log.Printf("Error converting target for attack observation: %v", err)
	}

	observation := AttackObservation{}
	observation.Subject = subject
	observation.Round = int(round)
	observation.Target = int(target)
	attackObservations = append(attackObservations, observation)
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
	target, err := strconv.ParseInt(action[indexOfActionToken+1:len(action)], 10, 64)
	if err != nil {
		log.Printf("Error converting target for lynch observation: %v", err)
	}

	observation := LynchObservation{}
	observation.Subject = subject
	observation.Round = int(round)
	observation.Target = int(target)
	lynchObservations = append(lynchObservations, observation)
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
	role, err := strconv.ParseInt(action[indexOfActionToken+1:len(action)], 10, 64)
	if err != nil {
		log.Printf("Error converting role for lynch observation: %v", err)
	}

	observation := KillObservation{}
	observation.Subject = subject
	observation.Round = int(round)
	observation.Role = int(role)
	killObservations = append(killObservations, observation)
}
