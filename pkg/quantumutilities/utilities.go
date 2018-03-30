package quantumutilities

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleErr outputs messages based on the input error
// returns true if there was an error
func HandleErr(c *gin.Context, err error, message string) bool {
	if err != nil {
		log.Printf("%s: %v", message, err)
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("%s: %v", message, err))
		return true
	}

	return false
}

// DbExec executes a database statement
func DbExec(c *gin.Context, db *sql.DB, statement string) {
	_, err := db.Exec(statement)
	HandleErr(c, err, fmt.Sprintf("Error executing statement [%s]", statement))
}

// DbExecReturn executes a database statement and returns an int
func DbExecReturn(c *gin.Context, db *sql.DB, statement string) (returnValue int) {
	err := db.QueryRow(statement).Scan(&returnValue)
	HandleErr(c, err, fmt.Sprintf("Error executing statement with return [%s]", statement))
	return
}

// Factorial finds the factorial of the input number.  n!
func Factorial(n int) uint64 {
	factVal := uint64(1)
	if n < 0 {
		panic("Factorial of negative number doesn't exist.")
	} else {
		for i := 1; i <= n; i++ {
			factVal *= uint64(i) // mismatched types int64 and int
		}

	}
	return factVal
}

// Kthperm finds a specified lexical permutation of the input slice
// S is an int slice that contains the values to be permutated
// k is the index of the permutation to get
func Kthperm(S []int, k uint64) []int {
	var P []int
	for len(S) > 0 {
		f := Factorial(len(S) - 1)
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

// GetBytes converts an arbitrary interface to a byte array
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetInterface converts abyte array to an arbitrary interface
func GetInterface(bts []byte, data interface{}) error {
	buf := bytes.NewBuffer(bts)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	return nil
}
