package uuidgen

import (
	"fmt"
	"math/rand"
	"time"
)

type IDGenerator struct{}

func NewIDGenerator() *IDGenerator { return &IDGenerator{} }

func (g *IDGenerator) GenerateTxID() string {
	return fmt.Sprintf("%016x%016x", rand.Uint64(), rand.Uint64())
}

func (g *IDGenerator) GenerateE2EID(ispb string) string {
	now := time.Now().UTC()
	suffix := fmt.Sprintf("%011x", rand.Uint64())[:11]
	return fmt.Sprintf("E%s%s%s", ispb, now.Format("20060102150405"), suffix)
}

func (g *IDGenerator) GenerateRtrID(ispb string) string {
	now := time.Now().UTC()
	suffix := fmt.Sprintf("%011x", rand.Uint64())[:11]
	return fmt.Sprintf("D%s%s%s", ispb, now.Format("20060102150405"), suffix)
}
