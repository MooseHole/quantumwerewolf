package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func handleErr(c *gin.Context, err error, message string) bool {
	if err != nil {
		log.Printf("%s: %v", message, err)
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("$s: %v", message, err))
		return true
	}

	return false
}

func resetVars() {
	rand.Seed(time.Now().UTC().UnixNano())
	players = nil
	roles.Name = ""
	roles.Total = 0
	roles.Villagers = 0
	roles.Seers = 0
	roles.Wolves = 0
	roles.Keep = 100
	game.Name = ""
	game.Number = -1
	game.RoundNight = true
	game.RoundNum = 0
	game.Seed = rand.Int63()
	multiverse.universes = nil
	multiverse.originalAssignments = nil
}

func dbExec(c *gin.Context, statement string) {
	_, err := db.Exec(statement)
	handleErr(c, err, fmt.Sprintf("Error executing statement [%s]", statement))
}

func dbExecReturn(c *gin.Context, statement string) (returnValue int) {
	err := db.QueryRow(statement).Scan(&returnValue)
	handleErr(c, err, fmt.Sprintf("Error executing statement with return [%s]", statement))
	return
}

func factorial(n int) uint64 {
	factVal := uint64(1)
	if n < 0 {
		fmt.Print("Factorial of negative number doesn't exist.")
	} else {
		for i := 1; i <= n; i++ {
			factVal *= uint64(i) // mismatched types int64 and int
		}

	}
	return factVal
}

func factorialCap(n int, ceiling uint64) uint64 {
	factVal := uint64(1)
	if n < 0 {
		fmt.Print("Factorial of negative number doesn't exist.")
	} else {
		for i := 1; i <= n; i++ {
			if factVal >= ceiling {
				return ceiling
			}
			factVal *= uint64(i) // mismatched types int64 and int
		}

	}
	return factVal
}

func kthperm(S []int, k uint64) []int {
	var P []int
	for len(S) > 0 {
		f := factorial(len(S) - 1)
		i := int(math.Floor(float64(k) / float64(f)))
		x := S[i]
		k = k % f
		P = append(P, x)
		S = append(S[:i], S[i+1:]...)
	}

	return P
}

// Uint63n returns, as an uint64, a pseudo-random number in [0,n).
func Uint63n(r *rand.Rand, n uint64) uint64 {
	max := uint64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := r.Uint64()
	for v > max {
		v = r.Uint64()
	}
	return v % n
}

func minUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}

	return b
}

// PermUint64 returns, as a slice of n uint64s, a pseudo-random permutation of the integers [0,n).
func PermUint64(r *rand.Rand, n uint64, maxValues uint64) []uint64 {
	log.Printf("n %d  maxValues %d  maxUint64 %d", n, maxValues, minUint64(n, maxValues))
	size := minUint64(n, maxValues)
	m := make([]uint64, size)
	b := make(map[uint64]bool)
	for i := uint64(0); i < size; i++ {
		j := Uint63n(r, n)
		for b[j] {
			j = Uint63n(r, n)
		}
		b[j] = true
		m[i] = j
	}

	return m
}
