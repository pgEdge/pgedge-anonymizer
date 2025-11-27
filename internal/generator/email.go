/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

package generator

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/pgedge/pgedge-anonymizer/internal/generator/data"
)

// EmailGenerator generates email addresses.
type EmailGenerator struct {
	BaseGenerator
	data *data.DataSet
}

// NewEmailGenerator creates a new email generator.
func NewEmailGenerator(d *data.DataSet) *EmailGenerator {
	return &EmailGenerator{
		BaseGenerator: BaseGenerator{name: "EMAIL"},
		data:          d,
	}
}

// Generate produces an email address.
// Uses a hash of the input to generate a unique local part, ensuring
// the same input always produces the same output while avoiding collisions.
func (g *EmailGenerator) Generate(input string) string {
	firstName := strings.ToLower(randomString(g.data.FirstNames))
	lastName := strings.ToLower(randomString(g.data.LastNames))
	domain := randomString(g.data.Domains)

	// Generate a unique suffix from input hash to avoid collisions
	hash := sha256.Sum256([]byte(input))
	hashStr := hex.EncodeToString(hash[:])
	// Use first 6 hex characters as unique suffix
	uniqueSuffix := hashStr[:6]

	// Vary email format randomly
	format := randomInt(5)
	switch format {
	case 0:
		// first.last.abc123@domain
		return firstName + "." + lastName + "." + uniqueSuffix + "@" + domain
	case 1:
		// flast.abc123@domain
		return string(firstName[0]) + lastName + "." + uniqueSuffix + "@" + domain
	case 2:
		// firstl.abc123@domain
		return firstName + string(lastName[0]) + "." + uniqueSuffix + "@" + domain
	case 3:
		// first_last_abc123@domain
		return firstName + "_" + lastName + "_" + uniqueSuffix + "@" + domain
	default:
		// firstlast.abc123@domain
		return firstName + lastName + "." + uniqueSuffix + "@" + domain
	}
}
