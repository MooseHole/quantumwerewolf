package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// handleErr outputs messages based on the input error
// returns true if there was an error
func handleErr(c *gin.Context, err error, message string) bool {
	if err != nil {
		log.Printf("%s: %v", message, err)
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("$s: %v", message, err))
		return true
	}

	return false
}

// dbExec executes a database statement
func dbExec(c *gin.Context, db *sql.DB, statement string) {
	_, err := db.Exec(statement)
	handleErr(c, err, fmt.Sprintf("Error executing statement [%s]", statement))
}

// dbExecReturn executes a database statement and returns an int
func dbExecReturn(c *gin.Context, db *sql.DB, statement string) (returnValue int) {
	err := db.QueryRow(statement).Scan(&returnValue)
	handleErr(c, err, fmt.Sprintf("Error executing statement with return [%s]", statement))
	return
}

// factorial finds the factorial of the input number.  n!
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

// kthperm finds a specified lexical permutation of the input slice
// S is an int slice that contains the values to be permutated
// k is the index of the permutation to get
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

// minUint64 finds the minimum of two uint64 values
func minUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}

	return b
}

// PermUint64Trunc returns, as a slice of n uint64s, a pseudo-random, possibly truncated, permutation of the integers [0,n).
// r is the random number generator.
// n is the number of values in the slice.
// maxValues is the maximum size of the slice.  Values that would be past this amount are truncated.
func PermUint64Trunc(r *rand.Rand, n uint64, maxValues uint64) []uint64 {
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
