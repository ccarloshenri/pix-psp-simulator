package uuidgen

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type IDGenerator struct{}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

func (g *IDGenerator) GenerateTxID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func (g *IDGenerator) GenerateE2EID(ispb string) string {
	now := time.Now().UTC()
	suffix := strings.ReplaceAll(uuid.New().String(), "-", "")[:11]
	return fmt.Sprintf("E%s%s%s", ispb, now.Format("20060102150405"), suffix)
}

func (g *IDGenerator) GenerateRtrID(ispb string) string {
	now := time.Now().UTC()
	suffix := strings.ReplaceAll(uuid.New().String(), "-", "")[:11]
	return fmt.Sprintf("D%s%s%s", ispb, now.Format("20060102150405"), suffix)
}
