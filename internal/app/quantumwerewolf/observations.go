package quantumwerewolf

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var dirtyObservations bool

type observation interface {
	getSubject() int
	getPending() bool
	getRound() int
	getOrder() int
	getTarget() (int, error)
	getIsEvil() (bool, error)
	getRole() (int, error)
	actionMessage() string
	action() string
	getType() string
	add()
	commit() observation
}

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

// VoteObservation keeps track of who voted to lynch whom
type VoteObservation struct {
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

// LynchObservation keeps track of who was lynched
type LynchObservation struct {
	Subject int
	Round   int
	Pending bool
}

func (o PeekObservation) getSubject() int {
	return o.Subject
}
func (o PeekObservation) getPending() bool {
	return o.Pending
}
func (o PeekObservation) getRound() int {
	return o.Round
}
func (o PeekObservation) getOrder() int {
	return 0
}
func (o PeekObservation) getTarget() (int, error) {
	return o.Target, nil
}
func (o PeekObservation) getIsEvil() (bool, error) {
	return o.IsEvil, nil
}
func (o PeekObservation) getRole() (int, error) {
	return -1, errors.New("PeekObservation does not produce a Role")
}
func (o PeekObservation) actionMessage() string {
	resultString := "good"
	if o.IsEvil {
		resultString = "evil"
	}
	return fmt.Sprintf("%s peeked at %s on night %d and found them %s.", getPlayerByNumber(o.Subject).Name, getPlayerByNumber(o.Target).Name, o.Round, resultString)
}
func (o PeekObservation) getType() string {
	return "PeekObservation"
}
func (o PeekObservation) commit() observation {
	if o.Pending {
		o.IsEvil = Peek(o.Subject, o.Target)
	}
	o.Pending = false
	return o
}
func (o PeekObservation) action() string {
	result := tokenGood
	if o.IsEvil {
		result = tokenEvil
	}
	pending := ""
	if o.Pending {
		pending = tokenPending
	}

	return strconv.Itoa(o.Round) + tokenPeek + strconv.Itoa(o.Target) + result + pending + tokenEndAction
}

func (o AttackObservation) action() string {
	pending := ""
	if o.Pending {
		pending = tokenPending
	}
	return strconv.Itoa(o.Round) + tokenAttack + strconv.Itoa(o.Target) + pending + tokenEndAction
}
func (o VoteObservation) action() string {
	pending := ""
	if o.Pending {
		pending = tokenPending
	}
	return strconv.Itoa(o.Round) + tokenVote + strconv.Itoa(o.Target) + pending + tokenEndAction
}
func (o KillObservation) action() string {
	pending := ""
	if o.Pending {
		pending = tokenPending
	}
	return strconv.Itoa(o.Round) + tokenKilled + strconv.Itoa(o.Role) + pending + tokenEndAction
}
func (o LynchObservation) action() string {
	pending := ""
	if o.Pending {
		pending = tokenPending
	}
	return strconv.Itoa(o.Round) + tokenLynched + pending + tokenEndAction
}

func (o AttackObservation) getSubject() int {
	return o.Subject
}
func (o AttackObservation) getPending() bool {
	return o.Pending
}
func (o AttackObservation) getRound() int {
	return o.Round
}
func (o AttackObservation) getOrder() int {
	return 1
}
func (o AttackObservation) getTarget() (int, error) {
	return o.Target, nil
}
func (o AttackObservation) getIsEvil() (bool, error) {
	return false, errors.New("AttackObservation does not produce IsEvil")
}
func (o AttackObservation) getRole() (int, error) {
	return -1, errors.New("AttackObservation does not produce a Role")
}
func (o AttackObservation) actionMessage() string {
	return fmt.Sprintf("%s attacked %s on night %d.", getPlayerByNumber(o.Subject).Name, getPlayerByNumber(o.Target).Name, o.Round)
}
func (o AttackObservation) getType() string {
	return "AttackObservation"
}
func (o AttackObservation) commit() observation {
	o.Pending = false
	return o
}

func (o VoteObservation) getSubject() int {
	return o.Subject
}
func (o VoteObservation) getPending() bool {
	return o.Pending
}
func (o VoteObservation) getRound() int {
	return o.Round
}
func (o VoteObservation) getOrder() int {
	return 2
}
func (o VoteObservation) getTarget() (int, error) {
	return o.Target, nil
}
func (o VoteObservation) getIsEvil() (bool, error) {
	return false, errors.New("VoteObservation does not produce IsEvil")
}
func (o VoteObservation) getRole() (int, error) {
	return -1, errors.New("VoteObservation does not produce a Role")
}
func (o VoteObservation) actionMessage() string {
	return fmt.Sprintf("%s voted to lynch %s on day %d.", getPlayerByNumber(o.Subject).Name, getPlayerByNumber(o.Target).Name, o.Round)
}
func (o VoteObservation) getType() string {
	return "VoteObservation"
}
func (o VoteObservation) commit() observation {
	o.Pending = false
	return o
}

func (o LynchObservation) getSubject() int {
	return o.Subject
}
func (o LynchObservation) getPending() bool {
	return o.Pending
}
func (o LynchObservation) getRound() int {
	return o.Round
}
func (o LynchObservation) getOrder() int {
	return 3
}
func (o LynchObservation) getTarget() (int, error) {
	return -1, errors.New("LynchObservation does not produce a Target")
}
func (o LynchObservation) getIsEvil() (bool, error) {
	return false, errors.New("LynchObservation does not produce IsEvil")
}
func (o LynchObservation) getRole() (int, error) {
	return -1, errors.New("LynchObservation does not produce a Role")
}
func (o LynchObservation) actionMessage() string {
	return fmt.Sprintf("%s got lynched on day %d.", getPlayerByNumber(o.Subject).Name, o.Round)
}
func (o LynchObservation) getType() string {
	return "LynchObservation"
}
func (o LynchObservation) commit() observation {
	o.Pending = false
	return o
}

func (o KillObservation) getSubject() int {
	return o.Subject
}
func (o KillObservation) getPending() bool {
	return o.Pending
}
func (o KillObservation) getRound() int {
	return o.Round
}
func (o KillObservation) getOrder() int {
	return 5
}
func (o KillObservation) getTarget() (int, error) {
	return -1, errors.New("KillObservation does not produce a Target")
}
func (o KillObservation) getIsEvil() (bool, error) {
	return false, errors.New("KillObservation does not produce IsEvil")
}
func (o KillObservation) getRole() (int, error) {
	return o.Role, nil
}
func (o KillObservation) actionMessage() string {
	return fmt.Sprintf("%s died in round %d and was a %s.", getPlayerByNumber(o.Subject).Name, o.Round, roleTypes[o.Role].Name)
}
func (o KillObservation) getType() string {
	return "KillObservation"
}
func (o KillObservation) commit() observation {
	o.Pending = false
	return o
}

func (o AttackObservation) add() {
	addObservation(o, false)
}

func (o PeekObservation) add() {
	addObservation(o, false)
}

func (o VoteObservation) add() {
	addObservation(o, false)
}

func (o KillObservation) add() {
	addObservation(o, true)
}

func (o LynchObservation) add() {
	addObservation(o, true)
}

func addObservation(o observation, exclusive bool) {
	temp := make([]observation, 0, len(observations)+1)
	exists := false
	addedBefore := false
	for _, current := range observations {
		replace := false
		if current.getType() == o.getType() && current.getSubject() == o.getSubject() {
			addedBefore = true
			if current.getRound() == o.getRound() {
				exists = true
				replace = true
				temp = append(temp, o)
			}
		}
		if !replace {
			temp = append(temp, current)
		}
	}

	if !exists && (!exclusive || !addedBefore) {
		temp = append(temp, o)
	}
	observations = nil
	for _, current := range temp {
		observations = append(observations, current)
	}
}

var observations []observation

// ResetObservations destroys all saved observation instances
func ResetObservations() {
	dirtyObservations = true
	observations = nil
}

// FillObservations fills all observations with current player actions
func FillObservations() {
	if !dirtyObservations {
		return
	}

	ResetObservations()
	for _, p := range Players {
		actionStrings := strings.Split(p.Actions, tokenEndAction)
		for _, action := range actionStrings {
			fillPeekObservation(p.Num, action)
			fillAttackObservation(p.Num, action)
			fillVoteObservation(p.Num, action)
			fillKillObservation(p.Num, action)
			fillLynchObservation(p.Num, action)
		}
	}
	dirtyObservations = false
}

// CommitObservations removes pending from all observations
func CommitObservations() {
	for _, o := range observations {
		o.commit().add()
	}
}

// FillActionsWithObservations takes the existing observations and fills the player Actions with a representation of them
func FillActionsWithObservations() {
	for i := range Players {
		Players[i].Actions = ""
	}
	for _, o := range observations {
		Players[getPlayerIndex(getPlayerByNumber(o.getSubject()))].Actions += o.action()
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
	observation.add()
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
	observation.add()
}

func fillVoteObservation(subject int, action string) {
	indexOfActionToken := strings.Index(action, tokenVote)

	// If not correct action type
	if indexOfActionToken < 0 {
		return
	}

	round, err := strconv.ParseInt(action[0:indexOfActionToken], 10, 64)
	if err != nil {
		log.Printf("Error converting round for vote observation: %v", err)
	}
	pending := strings.Contains(action, tokenPending)
	endOfTarget := len(action)
	if pending {
		endOfTarget--
	}
	target, err := strconv.ParseInt(action[indexOfActionToken+1:endOfTarget], 10, 64)
	if err != nil {
		log.Printf("Error converting target for vote observation: %v", err)
	}

	observation := VoteObservation{}
	observation.Subject = subject
	observation.Round = int(round)
	observation.Target = int(target)
	observation.Pending = pending
	observation.add()
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
	observation.add()
}

func fillLynchObservation(subject int, action string) {
	indexOfActionToken := strings.Index(action, tokenLynched)

	// If not correct action type
	if indexOfActionToken < 0 {
		return
	}

	round, err := strconv.ParseInt(action[0:indexOfActionToken], 10, 64)
	if err != nil {
		log.Printf("Error converting round for lynch observation: %v", err)
	}
	pending := strings.Contains(action, tokenPending)

	observation := LynchObservation{}
	observation.Subject = subject
	observation.Round = int(round)
	observation.Pending = pending
	observation.add()
}
