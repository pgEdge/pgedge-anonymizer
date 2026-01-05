/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025 - 2026, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

package generator

import (
	"fmt"
	"strings"
)

// AUPostcodeGenerator generates Australian postcodes.
type AUPostcodeGenerator struct {
	BaseGenerator
}

// NewAUPostcodeGenerator creates a new Australian postcode generator.
func NewAUPostcodeGenerator() *AUPostcodeGenerator {
	return &AUPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "AU_POSTCODE"},
	}
}

// Generate produces an Australian postcode (4 digits).
func (g *AUPostcodeGenerator) Generate(input string) string {
	// Australian postcodes are 4 digits, first digit indicates state
	// 2xxx NSW, 3xxx VIC, 4xxx QLD, 5xxx SA, 6xxx WA, 7xxx TAS, 08xx NT, 02xx ACT
	firstDigit := 2 + randomInt(6) // 2-7
	return fmt.Sprintf("%d%03d", firstDigit, randomInt(1000))
}

// DEPostcodeGenerator generates German postcodes (PLZ).
type DEPostcodeGenerator struct {
	BaseGenerator
}

// NewDEPostcodeGenerator creates a new German postcode generator.
func NewDEPostcodeGenerator() *DEPostcodeGenerator {
	return &DEPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "DE_POSTCODE"},
	}
}

// Generate produces a German postcode (5 digits).
func (g *DEPostcodeGenerator) Generate(input string) string {
	// German PLZ are 5 digits, 01xxx to 99xxx
	return fmt.Sprintf("%05d", 1000+randomInt(99000))
}

// ESPostcodeGenerator generates Spanish postcodes.
type ESPostcodeGenerator struct {
	BaseGenerator
}

// NewESPostcodeGenerator creates a new Spanish postcode generator.
func NewESPostcodeGenerator() *ESPostcodeGenerator {
	return &ESPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "ES_POSTCODE"},
	}
}

// Generate produces a Spanish postcode (5 digits).
func (g *ESPostcodeGenerator) Generate(input string) string {
	// Spanish postcodes: 01xxx to 52xxx (provinces)
	return fmt.Sprintf("%05d", 1000+randomInt(52000))
}

// FIPostcodeGenerator generates Finnish postcodes.
type FIPostcodeGenerator struct {
	BaseGenerator
}

// NewFIPostcodeGenerator creates a new Finnish postcode generator.
func NewFIPostcodeGenerator() *FIPostcodeGenerator {
	return &FIPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "FI_POSTCODE"},
	}
}

// Generate produces a Finnish postcode (5 digits).
func (g *FIPostcodeGenerator) Generate(input string) string {
	// Finnish postcodes: 00100 to 99999
	return fmt.Sprintf("%05d", 100+randomInt(99900))
}

// FRPostcodeGenerator generates French postcodes.
type FRPostcodeGenerator struct {
	BaseGenerator
}

// NewFRPostcodeGenerator creates a new French postcode generator.
func NewFRPostcodeGenerator() *FRPostcodeGenerator {
	return &FRPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "FR_POSTCODE"},
	}
}

// Generate produces a French postcode (5 digits).
func (g *FRPostcodeGenerator) Generate(input string) string {
	// French postcodes: first 2 digits are department (01-95, 2A, 2B for Corsica)
	dept := 1 + randomInt(95)
	return fmt.Sprintf("%02d%03d", dept, randomInt(1000))
}

// IEPostcodeGenerator generates Irish Eircodes.
type IEPostcodeGenerator struct {
	BaseGenerator
}

// NewIEPostcodeGenerator creates a new Irish Eircode generator.
func NewIEPostcodeGenerator() *IEPostcodeGenerator {
	return &IEPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "IE_POSTCODE"},
	}
}

// Generate produces an Irish Eircode (A9A A9A9 format).
func (g *IEPostcodeGenerator) Generate(input string) string {
	// Eircode format: A9A A9A9 (routing key + unique identifier)
	// Valid routing key letters
	letters := "ACDEFHKNPRTVWXY"
	hasSpace := strings.Contains(input, " ")

	routing := fmt.Sprintf("%c%d%c",
		letters[randomInt(len(letters))],
		randomInt(10),
		letters[randomInt(len(letters))])

	unique := fmt.Sprintf("%c%d%c%d",
		letters[randomInt(len(letters))],
		randomInt(10),
		letters[randomInt(len(letters))],
		randomInt(10))

	if hasSpace {
		return routing + " " + unique
	}
	return routing + unique
}

// INPostcodeGenerator generates Indian PIN codes.
type INPostcodeGenerator struct {
	BaseGenerator
}

// NewINPostcodeGenerator creates a new Indian PIN code generator.
func NewINPostcodeGenerator() *INPostcodeGenerator {
	return &INPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "IN_POSTCODE"},
	}
}

// Generate produces an Indian PIN code (6 digits).
func (g *INPostcodeGenerator) Generate(input string) string {
	// Indian PIN codes: first digit 1-8 (zone), never starts with 0 or 9
	firstDigit := 1 + randomInt(8) // 1-8
	return fmt.Sprintf("%d%05d", firstDigit, randomInt(100000))
}

// ITPostcodeGenerator generates Italian postcodes (CAP).
type ITPostcodeGenerator struct {
	BaseGenerator
}

// NewITPostcodeGenerator creates a new Italian postcode generator.
func NewITPostcodeGenerator() *ITPostcodeGenerator {
	return &ITPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "IT_POSTCODE"},
	}
}

// Generate produces an Italian postcode (5 digits).
func (g *ITPostcodeGenerator) Generate(input string) string {
	// Italian CAP: 00010 to 98168
	return fmt.Sprintf("%05d", 10+randomInt(98160))
}

// JPPostcodeGenerator generates Japanese postal codes.
type JPPostcodeGenerator struct {
	BaseGenerator
}

// NewJPPostcodeGenerator creates a new Japanese postal code generator.
func NewJPPostcodeGenerator() *JPPostcodeGenerator {
	return &JPPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "JP_POSTCODE"},
	}
}

// Generate produces a Japanese postal code (XXX-XXXX format).
func (g *JPPostcodeGenerator) Generate(input string) string {
	// Japanese postal codes: 3 digits, hyphen, 4 digits
	hasDash := strings.Contains(input, "-") || strings.Contains(input, "〒")

	first := fmt.Sprintf("%03d", randomInt(1000))
	second := fmt.Sprintf("%04d", randomInt(10000))

	if hasDash || len(strings.ReplaceAll(input, " ", "")) <= 7 {
		return first + "-" + second
	}
	return first + second
}

// KRPostcodeGenerator generates South Korean postal codes.
type KRPostcodeGenerator struct {
	BaseGenerator
}

// NewKRPostcodeGenerator creates a new South Korean postal code generator.
func NewKRPostcodeGenerator() *KRPostcodeGenerator {
	return &KRPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "KR_POSTCODE"},
	}
}

// Generate produces a South Korean postal code (5 digits).
func (g *KRPostcodeGenerator) Generate(input string) string {
	// Korean postal codes: 5 digits, 01000 to 63644
	return fmt.Sprintf("%05d", 1000+randomInt(63000))
}

// MXPostcodeGenerator generates Mexican postal codes.
type MXPostcodeGenerator struct {
	BaseGenerator
}

// NewMXPostcodeGenerator creates a new Mexican postal code generator.
func NewMXPostcodeGenerator() *MXPostcodeGenerator {
	return &MXPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "MX_POSTCODE"},
	}
}

// Generate produces a Mexican postal code (5 digits).
func (g *MXPostcodeGenerator) Generate(input string) string {
	// Mexican postal codes: 5 digits, 01000 to 99999
	return fmt.Sprintf("%05d", 1000+randomInt(99000))
}

// NOPostcodeGenerator generates Norwegian postal codes.
type NOPostcodeGenerator struct {
	BaseGenerator
}

// NewNOPostcodeGenerator creates a new Norwegian postal code generator.
func NewNOPostcodeGenerator() *NOPostcodeGenerator {
	return &NOPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "NO_POSTCODE"},
	}
}

// Generate produces a Norwegian postal code (4 digits).
func (g *NOPostcodeGenerator) Generate(input string) string {
	// Norwegian postal codes: 4 digits, 0001 to 9991
	return fmt.Sprintf("%04d", 1+randomInt(9990))
}

// NZPostcodeGenerator generates New Zealand postal codes.
type NZPostcodeGenerator struct {
	BaseGenerator
}

// NewNZPostcodeGenerator creates a new New Zealand postal code generator.
func NewNZPostcodeGenerator() *NZPostcodeGenerator {
	return &NZPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "NZ_POSTCODE"},
	}
}

// Generate produces a New Zealand postal code (4 digits).
func (g *NZPostcodeGenerator) Generate(input string) string {
	// New Zealand postal codes: 4 digits, 0110 to 9893
	return fmt.Sprintf("%04d", 110+randomInt(9784))
}

// PKPostcodeGenerator generates Pakistani postal codes.
type PKPostcodeGenerator struct {
	BaseGenerator
}

// NewPKPostcodeGenerator creates a new Pakistani postal code generator.
func NewPKPostcodeGenerator() *PKPostcodeGenerator {
	return &PKPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "PK_POSTCODE"},
	}
}

// Generate produces a Pakistani postal code (5 digits).
func (g *PKPostcodeGenerator) Generate(input string) string {
	// Pakistani postal codes: 5 digits, 10000 to 97000
	return fmt.Sprintf("%05d", 10000+randomInt(87000))
}

// SEPostcodeGenerator generates Swedish postal codes.
type SEPostcodeGenerator struct {
	BaseGenerator
}

// NewSEPostcodeGenerator creates a new Swedish postal code generator.
func NewSEPostcodeGenerator() *SEPostcodeGenerator {
	return &SEPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "SE_POSTCODE"},
	}
}

// Generate produces a Swedish postal code (XXX XX format).
func (g *SEPostcodeGenerator) Generate(input string) string {
	// Swedish postal codes: 5 digits, often formatted as XXX XX
	hasSpace := strings.Contains(input, " ")
	first := fmt.Sprintf("%03d", 100+randomInt(900))
	second := fmt.Sprintf("%02d", randomInt(100))

	if hasSpace {
		return first + " " + second
	}
	return first + second
}

// SGPostcodeGenerator generates Singaporean postal codes.
type SGPostcodeGenerator struct {
	BaseGenerator
}

// NewSGPostcodeGenerator creates a new Singaporean postal code generator.
func NewSGPostcodeGenerator() *SGPostcodeGenerator {
	return &SGPostcodeGenerator{
		BaseGenerator: BaseGenerator{name: "SG_POSTCODE"},
	}
}

// Generate produces a Singaporean postal code (6 digits).
func (g *SGPostcodeGenerator) Generate(input string) string {
	// Singapore postal codes: 6 digits, first 2 digits indicate district (01-82)
	district := 1 + randomInt(82)
	return fmt.Sprintf("%02d%04d", district, randomInt(10000))
}
