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
	"fmt"
	"strings"
)

// AUTFNGenerator generates Australian Tax File Numbers.
type AUTFNGenerator struct {
	BaseGenerator
}

// NewAUTFNGenerator creates a new Australian TFN generator.
func NewAUTFNGenerator() *AUTFNGenerator {
	return &AUTFNGenerator{
		BaseGenerator: BaseGenerator{name: "AU_TFN"},
	}
}

// Generate produces an Australian Tax File Number (9 digits).
func (g *AUTFNGenerator) Generate(input string) string {
	// Australian TFN is 9 digits with a check digit algorithm
	// For anonymization, we generate valid-looking 9-digit numbers
	hasSpaces := strings.Contains(input, " ")

	d1 := randomInt(10)
	d2 := randomInt(10)
	d3 := randomInt(10)
	d4 := randomInt(10)
	d5 := randomInt(10)
	d6 := randomInt(10)
	d7 := randomInt(10)
	d8 := randomInt(10)
	d9 := randomInt(10)

	if hasSpaces {
		return fmt.Sprintf("%d%d%d %d%d%d %d%d%d", d1, d2, d3, d4, d5, d6, d7, d8, d9)
	}
	return fmt.Sprintf("%d%d%d%d%d%d%d%d%d", d1, d2, d3, d4, d5, d6, d7, d8, d9)
}

// CASINGenerator generates Canadian Social Insurance Numbers.
type CASINGenerator struct {
	BaseGenerator
}

// NewCASINGenerator creates a new Canadian SIN generator.
func NewCASINGenerator() *CASINGenerator {
	return &CASINGenerator{
		BaseGenerator: BaseGenerator{name: "CA_SIN"},
	}
}

// Generate produces a Canadian Social Insurance Number (XXX-XXX-XXX).
func (g *CASINGenerator) Generate(input string) string {
	// Canadian SIN is 9 digits, often formatted XXX-XXX-XXX
	// First digit indicates province/territory of registration
	hasDash := strings.Contains(input, "-")
	hasSpace := strings.Contains(input, " ")

	first := fmt.Sprintf("%d%02d", 1+randomInt(9), randomInt(100))
	second := fmt.Sprintf("%03d", randomInt(1000))
	third := fmt.Sprintf("%03d", randomInt(1000))

	if hasDash {
		return first + "-" + second + "-" + third
	}
	if hasSpace {
		return first + " " + second + " " + third
	}
	return first + second + third
}

// DESteurIDGenerator generates German tax identification numbers.
type DESteurIDGenerator struct {
	BaseGenerator
}

// NewDESteurIDGenerator creates a new German Steueridentifikationsnummer generator.
func NewDESteurIDGenerator() *DESteurIDGenerator {
	return &DESteurIDGenerator{
		BaseGenerator: BaseGenerator{name: "DE_STEUERID"},
	}
}

// Generate produces a German Steueridentifikationsnummer (11 digits).
func (g *DESteurIDGenerator) Generate(input string) string {
	// German Steuer-ID is 11 digits, never starts with 0
	hasSpaces := strings.Contains(input, " ")

	first := 1 + randomInt(9)
	rest := make([]int, 10)
	for i := range rest {
		rest[i] = randomInt(10)
	}

	if hasSpaces {
		return fmt.Sprintf("%d%d %d%d%d %d%d%d %d%d%d",
			first, rest[0], rest[1], rest[2], rest[3],
			rest[4], rest[5], rest[6], rest[7], rest[8], rest[9])
	}
	return fmt.Sprintf("%d%d%d%d%d%d%d%d%d%d%d",
		first, rest[0], rest[1], rest[2], rest[3],
		rest[4], rest[5], rest[6], rest[7], rest[8], rest[9])
}

// ESNIFGenerator generates Spanish tax identification numbers.
type ESNIFGenerator struct {
	BaseGenerator
}

// NewESNIFGenerator creates a new Spanish NIF generator.
func NewESNIFGenerator() *ESNIFGenerator {
	return &ESNIFGenerator{
		BaseGenerator: BaseGenerator{name: "ES_NIF"},
	}
}

// Generate produces a Spanish NIF (8 digits + letter).
func (g *ESNIFGenerator) Generate(input string) string {
	// Spanish NIF/DNI is 8 digits followed by a check letter
	letters := "TRWAGMYFPDXBNJZSQVHLCKE"
	number := randomInt(100000000)
	letter := letters[number%23]
	return fmt.Sprintf("%08d%c", number, letter)
}

// FIHETUGenerator generates Finnish personal identity codes.
type FIHETUGenerator struct {
	BaseGenerator
}

// NewFIHETUGenerator creates a new Finnish HETU generator.
func NewFIHETUGenerator() *FIHETUGenerator {
	return &FIHETUGenerator{
		BaseGenerator: BaseGenerator{name: "FI_HETU"},
	}
}

// Generate produces a Finnish HETU (DDMMYY-XXXC format).
func (g *FIHETUGenerator) Generate(input string) string {
	// Finnish HETU: DDMMYY-XXXC where C is check character
	// Century marker: + (1800s), - (1900s), A (2000s)
	day := 1 + randomInt(28)
	month := 1 + randomInt(12)
	year := randomInt(100)
	individual := randomInt(1000)
	checkChars := "0123456789ABCDEFHJKLMNPRSTUVWXY"

	// Calculate check character
	fullNumber := day*10000000 + month*100000 + year*1000 + individual
	checkIdx := fullNumber % 31
	checkChar := checkChars[checkIdx]

	return fmt.Sprintf("%02d%02d%02d-%03d%c", day, month, year, individual, checkChar)
}

// FRNIRGenerator generates French social security numbers.
type FRNIRGenerator struct {
	BaseGenerator
}

// NewFRNIRGenerator creates a new French NIR generator.
func NewFRNIRGenerator() *FRNIRGenerator {
	return &FRNIRGenerator{
		BaseGenerator: BaseGenerator{name: "FR_NIR"},
	}
}

// Generate produces a French NIR (15 digits).
func (g *FRNIRGenerator) Generate(input string) string {
	// French NIR: Sex(1) + YY + MM + Dept(2) + Commune(3) + Order(3) + Key(2)
	hasSpaces := strings.Contains(input, " ")

	sex := 1 + randomInt(2)    // 1 or 2
	year := randomInt(100)     // 00-99
	month := 1 + randomInt(12) // 01-12
	dept := 1 + randomInt(95)  // 01-95
	commune := randomInt(1000) // 000-999
	order := randomInt(1000)   // 000-999
	key := randomInt(100)      // 00-99

	if hasSpaces {
		return fmt.Sprintf("%d %02d %02d %02d %03d %03d %02d",
			sex, year, month, dept, commune, order, key)
	}
	return fmt.Sprintf("%d%02d%02d%02d%03d%03d%02d",
		sex, year, month, dept, commune, order, key)
}

// IEPPSGenerator generates Irish PPS numbers.
type IEPPSGenerator struct {
	BaseGenerator
}

// NewIEPPSGenerator creates a new Irish PPS number generator.
func NewIEPPSGenerator() *IEPPSGenerator {
	return &IEPPSGenerator{
		BaseGenerator: BaseGenerator{name: "IE_PPS"},
	}
}

// Generate produces an Irish PPS Number (7 digits + 1-2 letters).
func (g *IEPPSGenerator) Generate(input string) string {
	// Irish PPSN: 7 digits + 1 letter (+ optional W for married women)
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	number := randomInt(10000000)
	letter := letters[randomInt(len(letters))]

	// Sometimes add second letter (W or A)
	if randomInt(4) == 0 {
		secondLetters := "WA"
		return fmt.Sprintf("%07d%c%c", number, letter, secondLetters[randomInt(2)])
	}
	return fmt.Sprintf("%07d%c", number, letter)
}

// INAadhaarGenerator generates Indian Aadhaar numbers.
type INAadhaarGenerator struct {
	BaseGenerator
}

// NewINAadhaarGenerator creates a new Indian Aadhaar generator.
func NewINAadhaarGenerator() *INAadhaarGenerator {
	return &INAadhaarGenerator{
		BaseGenerator: BaseGenerator{name: "IN_AADHAAR"},
	}
}

// Generate produces an Indian Aadhaar number (12 digits).
func (g *INAadhaarGenerator) Generate(input string) string {
	// Aadhaar: 12 digits, first digit is 2-9
	hasSpaces := strings.Contains(input, " ")

	first := 2 + randomInt(8) // 2-9
	rest := make([]int, 11)
	for i := range rest {
		rest[i] = randomInt(10)
	}

	if hasSpaces {
		return fmt.Sprintf("%d%d%d%d %d%d%d%d %d%d%d%d",
			first, rest[0], rest[1], rest[2],
			rest[3], rest[4], rest[5], rest[6],
			rest[7], rest[8], rest[9], rest[10])
	}
	return fmt.Sprintf("%d%d%d%d%d%d%d%d%d%d%d%d",
		first, rest[0], rest[1], rest[2],
		rest[3], rest[4], rest[5], rest[6],
		rest[7], rest[8], rest[9], rest[10])
}

// INPANGenerator generates Indian PAN numbers.
type INPANGenerator struct {
	BaseGenerator
}

// NewINPANGenerator creates a new Indian PAN generator.
func NewINPANGenerator() *INPANGenerator {
	return &INPANGenerator{
		BaseGenerator: BaseGenerator{name: "IN_PAN"},
	}
}

// Generate produces an Indian PAN (AAAAA9999A format).
func (g *INPANGenerator) Generate(input string) string {
	// PAN: 5 letters + 4 digits + 1 letter
	// 4th letter indicates holder type (P=Person, C=Company, etc.)
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	holderTypes := "PCFATBLJG"

	return fmt.Sprintf("%c%c%c%c%c%04d%c",
		letters[randomInt(26)],
		letters[randomInt(26)],
		letters[randomInt(26)],
		holderTypes[randomInt(len(holderTypes))],
		letters[randomInt(26)],
		randomInt(10000),
		letters[randomInt(26)])
}

// ITCFGenerator generates Italian Codice Fiscale.
type ITCFGenerator struct {
	BaseGenerator
}

// NewITCFGenerator creates a new Italian Codice Fiscale generator.
func NewITCFGenerator() *ITCFGenerator {
	return &ITCFGenerator{
		BaseGenerator: BaseGenerator{name: "IT_CF"},
	}
}

// Generate produces an Italian Codice Fiscale (16 alphanumeric).
func (g *ITCFGenerator) Generate(input string) string {
	// Italian CF: SSSNNN YYXDD CCCC C
	// SSS=surname, NNN=name, YY=year, X=month, DD=day, CCCC=municipality, C=check
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	monthCodes := "ABCDEHLMPRST"

	surname := fmt.Sprintf("%c%c%c",
		letters[randomInt(26)], letters[randomInt(26)], letters[randomInt(26)])
	name := fmt.Sprintf("%c%c%c",
		letters[randomInt(26)], letters[randomInt(26)], letters[randomInt(26)])
	year := fmt.Sprintf("%02d", randomInt(100))
	month := string(monthCodes[randomInt(12)])
	day := fmt.Sprintf("%02d", 1+randomInt(31))
	municipality := fmt.Sprintf("%c%03d", letters[randomInt(26)], randomInt(1000))
	check := string(letters[randomInt(26)])

	return surname + name + year + month + day + municipality + check
}

// JPMyNumberGenerator generates Japanese My Number.
type JPMyNumberGenerator struct {
	BaseGenerator
}

// NewJPMyNumberGenerator creates a new Japanese My Number generator.
func NewJPMyNumberGenerator() *JPMyNumberGenerator {
	return &JPMyNumberGenerator{
		BaseGenerator: BaseGenerator{name: "JP_MYNUMBER"},
	}
}

// Generate produces a Japanese My Number (12 digits).
func (g *JPMyNumberGenerator) Generate(input string) string {
	// Japanese My Number is 12 digits
	hasSpaces := strings.Contains(input, " ")
	hasDash := strings.Contains(input, "-")

	digits := make([]int, 12)
	for i := range digits {
		digits[i] = randomInt(10)
	}

	if hasSpaces {
		return fmt.Sprintf("%d%d%d%d %d%d%d%d %d%d%d%d",
			digits[0], digits[1], digits[2], digits[3],
			digits[4], digits[5], digits[6], digits[7],
			digits[8], digits[9], digits[10], digits[11])
	}
	if hasDash {
		return fmt.Sprintf("%d%d%d%d-%d%d%d%d-%d%d%d%d",
			digits[0], digits[1], digits[2], digits[3],
			digits[4], digits[5], digits[6], digits[7],
			digits[8], digits[9], digits[10], digits[11])
	}
	return fmt.Sprintf("%d%d%d%d%d%d%d%d%d%d%d%d",
		digits[0], digits[1], digits[2], digits[3],
		digits[4], digits[5], digits[6], digits[7],
		digits[8], digits[9], digits[10], digits[11])
}

// KRRRNGenerator generates South Korean Resident Registration Numbers.
type KRRRNGenerator struct {
	BaseGenerator
}

// NewKRRRNGenerator creates a new South Korean RRN generator.
func NewKRRRNGenerator() *KRRRNGenerator {
	return &KRRRNGenerator{
		BaseGenerator: BaseGenerator{name: "KR_RRN"},
	}
}

// Generate produces a South Korean RRN (YYMMDD-XXXXXXX format).
func (g *KRRRNGenerator) Generate(input string) string {
	// Korean RRN: 6 digits (birthdate) + 7 digits (gender + registration)
	hasDash := strings.Contains(input, "-")

	year := randomInt(100)
	month := 1 + randomInt(12)
	day := 1 + randomInt(28)
	// First digit of second part: 1-2 (1900s male/female), 3-4 (2000s male/female)
	genderCentury := 1 + randomInt(4)
	rest := randomInt(1000000)

	first := fmt.Sprintf("%02d%02d%02d", year, month, day)
	second := fmt.Sprintf("%d%06d", genderCentury, rest)

	if hasDash {
		return first + "-" + second
	}
	return first + second
}

// MXCURPGenerator generates Mexican CURP numbers.
type MXCURPGenerator struct {
	BaseGenerator
}

// NewMXCURPGenerator creates a new Mexican CURP generator.
func NewMXCURPGenerator() *MXCURPGenerator {
	return &MXCURPGenerator{
		BaseGenerator: BaseGenerator{name: "MX_CURP"},
	}
}

// Generate produces a Mexican CURP (18 alphanumeric characters).
func (g *MXCURPGenerator) Generate(input string) string {
	// CURP: AAAA YYMMDD S EE CCC NN
	// AAAA=name initials, YYMMDD=birthdate, S=sex, EE=state, CCC=consonants, NN=check
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	vowels := "AEIOUX"
	consonants := "BCDFGHJKLMNPQRSTVWXYZ"
	states := []string{"AS", "BC", "BS", "CC", "CS", "CH", "CL", "CM", "DF", "DG",
		"GT", "GR", "HG", "JC", "MC", "MN", "MS", "NT", "NL", "OC",
		"PL", "QT", "QR", "SP", "SL", "SR", "TC", "TS", "TL", "VZ", "YN", "ZS"}
	sex := "HM"

	// First 4 letters: surname initial + first vowel + maternal surname initial + first name initial
	first4 := fmt.Sprintf("%c%c%c%c",
		letters[randomInt(26)],
		vowels[randomInt(len(vowels))],
		letters[randomInt(26)],
		letters[randomInt(26)])

	// Birthdate
	year := randomInt(100)
	month := 1 + randomInt(12)
	day := 1 + randomInt(28)
	birthdate := fmt.Sprintf("%02d%02d%02d", year, month, day)

	// Sex
	sexCode := string(sex[randomInt(2)])

	// State
	state := states[randomInt(len(states))]

	// 3 consonants
	cons := fmt.Sprintf("%c%c%c",
		consonants[randomInt(len(consonants))],
		consonants[randomInt(len(consonants))],
		consonants[randomInt(len(consonants))])

	// Homoclave (2 characters)
	homoclave := fmt.Sprintf("%c%d", letters[randomInt(26)], randomInt(10))

	return first4 + birthdate + sexCode + state + cons + homoclave
}

// NOFNRGenerator generates Norwegian national identity numbers.
type NOFNRGenerator struct {
	BaseGenerator
}

// NewNOFNRGenerator creates a new Norwegian Fødselsnummer generator.
func NewNOFNRGenerator() *NOFNRGenerator {
	return &NOFNRGenerator{
		BaseGenerator: BaseGenerator{name: "NO_FNR"},
	}
}

// Generate produces a Norwegian Fødselsnummer (11 digits).
func (g *NOFNRGenerator) Generate(input string) string {
	// Norwegian FNR: DDMMYY + 5 digits (individual number + 2 check digits)
	hasSpace := strings.Contains(input, " ")

	day := 1 + randomInt(28)
	month := 1 + randomInt(12)
	year := randomInt(100)
	individual := randomInt(1000)
	check := randomInt(100)

	first := fmt.Sprintf("%02d%02d%02d", day, month, year)
	second := fmt.Sprintf("%03d%02d", individual, check)

	if hasSpace {
		return first + " " + second
	}
	return first + second
}

// NZIRDGenerator generates New Zealand IRD numbers.
type NZIRDGenerator struct {
	BaseGenerator
}

// NewNZIRDGenerator creates a new New Zealand IRD number generator.
func NewNZIRDGenerator() *NZIRDGenerator {
	return &NZIRDGenerator{
		BaseGenerator: BaseGenerator{name: "NZ_IRD"},
	}
}

// Generate produces a New Zealand IRD number (8-9 digits).
func (g *NZIRDGenerator) Generate(input string) string {
	// NZ IRD: 8-9 digits, often formatted XXX-XXX-XXX
	hasDash := strings.Contains(input, "-")
	hasSpace := strings.Contains(input, " ")

	// Generate 8 or 9 digit number
	isNineDigit := randomInt(2) == 0
	var number int
	if isNineDigit {
		number = 10000000 + randomInt(90000000)
	} else {
		number = 10000000 + randomInt(90000000)
	}

	numStr := fmt.Sprintf("%d", number)

	if hasDash {
		if len(numStr) == 9 {
			return numStr[:3] + "-" + numStr[3:6] + "-" + numStr[6:]
		}
		return numStr[:2] + "-" + numStr[2:5] + "-" + numStr[5:]
	}
	if hasSpace {
		if len(numStr) == 9 {
			return numStr[:3] + " " + numStr[3:6] + " " + numStr[6:]
		}
		return numStr[:2] + " " + numStr[2:5] + " " + numStr[5:]
	}
	return numStr
}

// PKCNICGenerator generates Pakistani CNIC numbers.
type PKCNICGenerator struct {
	BaseGenerator
}

// NewPKCNICGenerator creates a new Pakistani CNIC generator.
func NewPKCNICGenerator() *PKCNICGenerator {
	return &PKCNICGenerator{
		BaseGenerator: BaseGenerator{name: "PK_CNIC"},
	}
}

// Generate produces a Pakistani CNIC (13 digits, XXXXX-XXXXXXX-X).
func (g *PKCNICGenerator) Generate(input string) string {
	// Pakistani CNIC: 13 digits formatted as XXXXX-XXXXXXX-X
	hasDash := strings.Contains(input, "-")

	region := 10000 + randomInt(90000) // 5 digits
	serial := randomInt(10000000)      // 7 digits
	gender := randomInt(10)            // 1 digit (odd=male, even=female)

	if hasDash {
		return fmt.Sprintf("%05d-%07d-%d", region, serial, gender)
	}
	return fmt.Sprintf("%05d%07d%d", region, serial, gender)
}

// SEPNRGenerator generates Swedish personal identity numbers.
type SEPNRGenerator struct {
	BaseGenerator
}

// NewSEPNRGenerator creates a new Swedish personnummer generator.
func NewSEPNRGenerator() *SEPNRGenerator {
	return &SEPNRGenerator{
		BaseGenerator: BaseGenerator{name: "SE_PNR"},
	}
}

// Generate produces a Swedish personnummer (YYMMDD-XXXX format).
func (g *SEPNRGenerator) Generate(input string) string {
	// Swedish personnummer: YYMMDD-XXXX or YYYYMMDD-XXXX
	hasDash := strings.Contains(input, "-")
	hasPlus := strings.Contains(input, "+") // Used for people over 100

	year := randomInt(100)
	month := 1 + randomInt(12)
	day := 1 + randomInt(28)
	serial := randomInt(10000)

	separator := "-"
	if hasPlus {
		separator = "+"
	}

	if hasDash || hasPlus {
		return fmt.Sprintf("%02d%02d%02d%s%04d", year, month, day, separator, serial)
	}
	return fmt.Sprintf("%02d%02d%02d%04d", year, month, day, serial)
}

// SGNRICGenerator generates Singaporean NRIC numbers.
type SGNRICGenerator struct {
	BaseGenerator
}

// NewSGNRICGenerator creates a new Singaporean NRIC generator.
func NewSGNRICGenerator() *SGNRICGenerator {
	return &SGNRICGenerator{
		BaseGenerator: BaseGenerator{name: "SG_NRIC"},
	}
}

// Generate produces a Singaporean NRIC (letter + 7 digits + letter).
func (g *SGNRICGenerator) Generate(input string) string {
	// NRIC: S/T (citizens) or F/G (foreigners) + 7 digits + check letter
	prefixes := "STFG"
	checkLetters := "JZIHGFEDCBA"

	prefix := prefixes[randomInt(len(prefixes))]
	number := randomInt(10000000)
	check := checkLetters[randomInt(len(checkLetters))]

	return fmt.Sprintf("%c%07d%c", prefix, number, check)
}

// USSSNGenerator generates US Social Security Numbers.
type USSSNGenerator struct {
	BaseGenerator
}

// NewUSSSNGenerator creates a new US SSN generator.
func NewUSSSNGenerator() *USSSNGenerator {
	return &USSSNGenerator{
		BaseGenerator: BaseGenerator{name: "US_SSN"},
	}
}

// Generate produces a US Social Security Number (XXX-XX-XXXX).
func (g *USSSNGenerator) Generate(input string) string {
	// SSN format: AAA-GG-SSSS
	// Area (AAA): 001-899, excluding 666
	// Group (GG): 01-99
	// Serial (SSSS): 0001-9999
	hasDash := strings.Contains(input, "-")
	hasSpace := strings.Contains(input, " ")

	// Generate area number (001-899, not 666)
	area := 1 + randomInt(899)
	if area == 666 {
		area = 667
	}
	group := 1 + randomInt(99)
	serial := 1 + randomInt(9999)

	if hasDash {
		return fmt.Sprintf("%03d-%02d-%04d", area, group, serial)
	}
	if hasSpace {
		return fmt.Sprintf("%03d %02d %04d", area, group, serial)
	}
	return fmt.Sprintf("%03d%02d%04d", area, group, serial)
}
