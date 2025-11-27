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

// AUPhoneGenerator generates Australian phone numbers.
type AUPhoneGenerator struct {
	BaseGenerator
}

// NewAUPhoneGenerator creates a new Australian phone generator.
func NewAUPhoneGenerator() *AUPhoneGenerator {
	return &AUPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "AU_PHONE"},
	}
}

// Generate produces an Australian phone number.
// Format: 04XX XXX XXX (mobile) or 0X XXXX XXXX (landline)
func (g *AUPhoneGenerator) Generate(input string) string {
	hasSpace := strings.Contains(input, " ")
	hasDash := strings.Contains(input, "-")
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+61")

	// Generate mobile (04XX) or landline (02-08)
	var number string
	if randomInt(2) == 0 {
		// Mobile: 04XX XXX XXX
		prefix := fmt.Sprintf("04%d%d", randomInt(10), randomInt(10))
		suffix := fmt.Sprintf("%d%d%d%d%d%d", randomInt(10), randomInt(10), randomInt(10),
			randomInt(10), randomInt(10), randomInt(10))
		if hasSpace {
			number = prefix + " " + suffix[:3] + " " + suffix[3:]
		} else if hasDash {
			number = prefix + "-" + suffix[:3] + "-" + suffix[3:]
		} else {
			number = prefix + suffix
		}
	} else {
		// Landline: 0X XXXX XXXX
		areaCode := fmt.Sprintf("0%d", 2+randomInt(7)) // 02-08
		suffix := fmt.Sprintf("%d%d%d%d%d%d%d%d", randomInt(10), randomInt(10), randomInt(10),
			randomInt(10), randomInt(10), randomInt(10), randomInt(10), randomInt(10))
		if hasSpace {
			number = areaCode + " " + suffix[:4] + " " + suffix[4:]
		} else if hasDash {
			number = areaCode + "-" + suffix[:4] + "-" + suffix[4:]
		} else {
			number = areaCode + suffix
		}
	}

	if hasCountryCode {
		return "+61 " + number[1:] // Remove leading 0
	}
	return number
}

// CAPhoneGenerator generates Canadian phone numbers (same format as US).
type CAPhoneGenerator struct {
	BaseGenerator
}

// NewCAPhoneGenerator creates a new Canadian phone generator.
func NewCAPhoneGenerator() *CAPhoneGenerator {
	return &CAPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "CA_PHONE"},
	}
}

// Generate produces a Canadian phone number.
// Format: (XXX) XXX-XXXX or XXX-XXX-XXXX using 555-01XX exchange
func (g *CAPhoneGenerator) Generate(input string) string {
	// Use fictional 555-01XX range
	format := detectPhoneFormat(input)
	areaCode := fmt.Sprintf("%d%d%d", 2+randomInt(8), randomInt(10), randomInt(10))
	lastFour := fmt.Sprintf("01%d%d", randomInt(10), randomInt(10))

	if format.hasParens {
		if format.separator != 0 {
			return fmt.Sprintf("(%s) 555%c%s", areaCode, format.separator, lastFour)
		}
		return fmt.Sprintf("(%s) 555-%s", areaCode, lastFour)
	}
	if format.separator != 0 {
		return fmt.Sprintf("%s%c555%c%s", areaCode, format.separator, format.separator, lastFour)
	}
	return areaCode + "555" + lastFour
}

// DEPhoneGenerator generates German phone numbers.
type DEPhoneGenerator struct {
	BaseGenerator
}

// NewDEPhoneGenerator creates a new German phone generator.
func NewDEPhoneGenerator() *DEPhoneGenerator {
	return &DEPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "DE_PHONE"},
	}
}

// Generate produces a German phone number.
// Format: +49 XXX XXXXXXXX or 0XXX XXXXXXXX
func (g *DEPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+49")
	hasSpace := strings.Contains(input, " ")

	areaCode := fmt.Sprintf("%d%d%d", 1+randomInt(9), randomInt(10), randomInt(10))
	subscriber := generateDigits(7)

	if hasCountryCode {
		if hasSpace {
			return "+49 " + areaCode + " " + subscriber
		}
		return "+49" + areaCode + subscriber
	}
	if hasSpace {
		return "0" + areaCode + " " + subscriber
	}
	return "0" + areaCode + subscriber
}

// ESPhoneGenerator generates Spanish phone numbers.
type ESPhoneGenerator struct {
	BaseGenerator
}

// NewESPhoneGenerator creates a new Spanish phone generator.
func NewESPhoneGenerator() *ESPhoneGenerator {
	return &ESPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "ES_PHONE"},
	}
}

// Generate produces a Spanish phone number.
// Format: +34 XXX XXX XXX or 9XX XXX XXX (landline) or 6XX XXX XXX (mobile)
func (g *ESPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+34")
	hasSpace := strings.Contains(input, " ")

	// Mobile (6XX) or landline (9XX)
	var prefix string
	if randomInt(2) == 0 {
		prefix = fmt.Sprintf("6%d%d", randomInt(10), randomInt(10))
	} else {
		prefix = fmt.Sprintf("9%d%d", randomInt(10), randomInt(10))
	}
	middle := generateDigits(3)
	suffix := generateDigits(3)

	if hasCountryCode {
		if hasSpace {
			return "+34 " + prefix + " " + middle + " " + suffix
		}
		return "+34" + prefix + middle + suffix
	}
	if hasSpace {
		return prefix + " " + middle + " " + suffix
	}
	return prefix + middle + suffix
}

// FIPhoneGenerator generates Finnish phone numbers.
type FIPhoneGenerator struct {
	BaseGenerator
}

// NewFIPhoneGenerator creates a new Finnish phone generator.
func NewFIPhoneGenerator() *FIPhoneGenerator {
	return &FIPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "FI_PHONE"},
	}
}

// Generate produces a Finnish phone number.
// Format: +358 XX XXX XXXX or 0XX XXX XXXX
func (g *FIPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+358")
	hasSpace := strings.Contains(input, " ")

	areaCode := fmt.Sprintf("%d%d", 4+randomInt(6), randomInt(10)) // 40-99
	subscriber := generateDigits(7)

	if hasCountryCode {
		if hasSpace {
			return "+358 " + areaCode + " " + subscriber[:3] + " " + subscriber[3:]
		}
		return "+358" + areaCode + subscriber
	}
	if hasSpace {
		return "0" + areaCode + " " + subscriber[:3] + " " + subscriber[3:]
	}
	return "0" + areaCode + subscriber
}

// FRPhoneGenerator generates French phone numbers.
type FRPhoneGenerator struct {
	BaseGenerator
}

// NewFRPhoneGenerator creates a new French phone generator.
func NewFRPhoneGenerator() *FRPhoneGenerator {
	return &FRPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "FR_PHONE"},
	}
}

// Generate produces a French phone number.
// Format: +33 X XX XX XX XX or 0X XX XX XX XX
func (g *FRPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+33")
	hasSpace := strings.Contains(input, " ")
	hasDot := strings.Contains(input, ".")

	// 01-05 landline, 06-07 mobile
	prefix := 1 + randomInt(7) // 1-7
	subscriber := generateDigits(8)

	var sep string
	if hasDot {
		sep = "."
	} else if hasSpace {
		sep = " "
	}

	if sep != "" {
		formatted := fmt.Sprintf("%s%s%s%s%s%s%s",
			subscriber[:2], sep, subscriber[2:4], sep,
			subscriber[4:6], sep, subscriber[6:8])
		if hasCountryCode {
			return fmt.Sprintf("+33 %d %s", prefix, formatted)
		}
		return fmt.Sprintf("0%d %s", prefix, formatted)
	}

	if hasCountryCode {
		return fmt.Sprintf("+33%d%s", prefix, subscriber)
	}
	return fmt.Sprintf("0%d%s", prefix, subscriber)
}

// IEPhoneGenerator generates Irish phone numbers.
type IEPhoneGenerator struct {
	BaseGenerator
}

// NewIEPhoneGenerator creates a new Irish phone generator.
func NewIEPhoneGenerator() *IEPhoneGenerator {
	return &IEPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "IE_PHONE"},
	}
}

// Generate produces an Irish phone number.
// Format: +353 XX XXX XXXX or 0XX XXX XXXX
func (g *IEPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+353")
	hasSpace := strings.Contains(input, " ")

	// Common area codes: 1 (Dublin), 21 (Cork), 61 (Limerick), 91 (Galway)
	areaCodes := []string{"1", "21", "61", "91", "22", "23", "24", "25", "26"}
	areaCode := areaCodes[randomInt(len(areaCodes))]
	subscriber := generateDigits(7)

	if hasCountryCode {
		if hasSpace {
			return "+353 " + areaCode + " " + subscriber[:3] + " " + subscriber[3:]
		}
		return "+353" + areaCode + subscriber
	}
	if hasSpace {
		return "0" + areaCode + " " + subscriber[:3] + " " + subscriber[3:]
	}
	return "0" + areaCode + subscriber
}

// INPhoneGenerator generates Indian phone numbers.
type INPhoneGenerator struct {
	BaseGenerator
}

// NewINPhoneGenerator creates a new Indian phone generator.
func NewINPhoneGenerator() *INPhoneGenerator {
	return &INPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "IN_PHONE"},
	}
}

// Generate produces an Indian phone number.
// Format: +91 XXXXX XXXXX or 0XXXXX XXXXX
func (g *INPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+91")
	hasSpace := strings.Contains(input, " ")

	// Mobile numbers start with 6-9
	first := 6 + randomInt(4) // 6-9
	subscriber := generateDigits(9)

	if hasCountryCode {
		if hasSpace {
			return fmt.Sprintf("+91 %d%s %s", first, subscriber[:4], subscriber[4:])
		}
		return fmt.Sprintf("+91%d%s", first, subscriber)
	}
	if hasSpace {
		return fmt.Sprintf("%d%s %s", first, subscriber[:4], subscriber[4:])
	}
	return fmt.Sprintf("%d%s", first, subscriber)
}

// ITPhoneGenerator generates Italian phone numbers.
type ITPhoneGenerator struct {
	BaseGenerator
}

// NewITPhoneGenerator creates a new Italian phone generator.
func NewITPhoneGenerator() *ITPhoneGenerator {
	return &ITPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "IT_PHONE"},
	}
}

// Generate produces an Italian phone number.
// Format: +39 XXX XXX XXXX or 0XX XXX XXXX
func (g *ITPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+39")
	hasSpace := strings.Contains(input, " ")

	// Mobile (3XX) or landline (0XX)
	var prefix string
	if randomInt(2) == 0 {
		prefix = fmt.Sprintf("3%d%d", randomInt(10), randomInt(10))
	} else {
		prefix = fmt.Sprintf("0%d", 2+randomInt(7)) // 02-08
	}
	subscriber := generateDigits(7)

	if hasCountryCode {
		if hasSpace {
			return "+39 " + prefix + " " + subscriber[:3] + " " + subscriber[3:]
		}
		return "+39" + prefix + subscriber
	}
	if hasSpace {
		return prefix + " " + subscriber[:3] + " " + subscriber[3:]
	}
	return prefix + subscriber
}

// JPPhoneGenerator generates Japanese phone numbers.
type JPPhoneGenerator struct {
	BaseGenerator
}

// NewJPPhoneGenerator creates a new Japanese phone generator.
func NewJPPhoneGenerator() *JPPhoneGenerator {
	return &JPPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "JP_PHONE"},
	}
}

// Generate produces a Japanese phone number.
// Format: +81 X-XXXX-XXXX or 0X-XXXX-XXXX
func (g *JPPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+81")
	hasDash := strings.Contains(input, "-")

	// Area codes: 3 (Tokyo), 6 (Osaka), etc.
	areaCode := fmt.Sprintf("%d", 1+randomInt(9))
	middle := generateDigits(4)
	suffix := generateDigits(4)

	if hasCountryCode {
		if hasDash {
			return "+81 " + areaCode + "-" + middle + "-" + suffix
		}
		return "+81" + areaCode + middle + suffix
	}
	if hasDash {
		return "0" + areaCode + "-" + middle + "-" + suffix
	}
	return "0" + areaCode + middle + suffix
}

// KRPhoneGenerator generates South Korean phone numbers.
type KRPhoneGenerator struct {
	BaseGenerator
}

// NewKRPhoneGenerator creates a new South Korean phone generator.
func NewKRPhoneGenerator() *KRPhoneGenerator {
	return &KRPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "KR_PHONE"},
	}
}

// Generate produces a South Korean phone number.
// Format: +82 XX-XXXX-XXXX or 0XX-XXXX-XXXX
func (g *KRPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+82")
	hasDash := strings.Contains(input, "-")

	// Mobile (010) or landline (02, 031-064)
	var areaCode string
	if randomInt(2) == 0 {
		areaCode = "10" // Mobile
	} else {
		areaCode = fmt.Sprintf("%d", 2+randomInt(63)) // 2-64
	}
	middle := generateDigits(4)
	suffix := generateDigits(4)

	if hasCountryCode {
		if hasDash {
			return "+82 " + areaCode + "-" + middle + "-" + suffix
		}
		return "+82" + areaCode + middle + suffix
	}
	if hasDash {
		return "0" + areaCode + "-" + middle + "-" + suffix
	}
	return "0" + areaCode + middle + suffix
}

// MXPhoneGenerator generates Mexican phone numbers.
type MXPhoneGenerator struct {
	BaseGenerator
}

// NewMXPhoneGenerator creates a new Mexican phone generator.
func NewMXPhoneGenerator() *MXPhoneGenerator {
	return &MXPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "MX_PHONE"},
	}
}

// Generate produces a Mexican phone number.
// Format: +52 XXX XXX XXXX or (XXX) XXX-XXXX
func (g *MXPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+52")
	hasSpace := strings.Contains(input, " ")
	hasParens := strings.Contains(input, "(")

	areaCode := generateDigits(3)
	middle := generateDigits(3)
	suffix := generateDigits(4)

	if hasCountryCode {
		if hasSpace {
			return "+52 " + areaCode + " " + middle + " " + suffix
		}
		return "+52" + areaCode + middle + suffix
	}
	if hasParens {
		return "(" + areaCode + ") " + middle + "-" + suffix
	}
	if hasSpace {
		return areaCode + " " + middle + " " + suffix
	}
	return areaCode + middle + suffix
}

// NOPhoneGenerator generates Norwegian phone numbers.
type NOPhoneGenerator struct {
	BaseGenerator
}

// NewNOPhoneGenerator creates a new Norwegian phone generator.
func NewNOPhoneGenerator() *NOPhoneGenerator {
	return &NOPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "NO_PHONE"},
	}
}

// Generate produces a Norwegian phone number.
// Format: +47 XXX XX XXX or XXX XX XXX
func (g *NOPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+47")
	hasSpace := strings.Contains(input, " ")

	// Mobile (4XX, 9XX) or landline (2X, 3X, 5X, 6X, 7X)
	var prefix string
	if randomInt(2) == 0 {
		prefix = fmt.Sprintf("%d%d%d", 4+randomInt(6), randomInt(10), randomInt(10))
	} else {
		prefix = fmt.Sprintf("%d%d%d", 2+randomInt(6), randomInt(10), randomInt(10))
	}
	suffix := generateDigits(5)

	if hasCountryCode {
		if hasSpace {
			return "+47 " + prefix + " " + suffix[:2] + " " + suffix[2:]
		}
		return "+47" + prefix + suffix
	}
	if hasSpace {
		return prefix + " " + suffix[:2] + " " + suffix[2:]
	}
	return prefix + suffix
}

// NZPhoneGenerator generates New Zealand phone numbers.
type NZPhoneGenerator struct {
	BaseGenerator
}

// NewNZPhoneGenerator creates a new New Zealand phone generator.
func NewNZPhoneGenerator() *NZPhoneGenerator {
	return &NZPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "NZ_PHONE"},
	}
}

// Generate produces a New Zealand phone number.
// Format: +64 X XXX XXXX or 0X XXX XXXX
func (g *NZPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+64")
	hasSpace := strings.Contains(input, " ")

	// Mobile (02X) or landline (03, 04, 06, 07, 09)
	var areaCode string
	if randomInt(2) == 0 {
		areaCode = fmt.Sprintf("2%d", randomInt(10))
	} else {
		landlines := []string{"3", "4", "6", "7", "9"}
		areaCode = landlines[randomInt(len(landlines))]
	}
	subscriber := generateDigits(7)

	if hasCountryCode {
		if hasSpace {
			return "+64 " + areaCode + " " + subscriber[:3] + " " + subscriber[3:]
		}
		return "+64" + areaCode + subscriber
	}
	if hasSpace {
		return "0" + areaCode + " " + subscriber[:3] + " " + subscriber[3:]
	}
	return "0" + areaCode + subscriber
}

// PKPhoneGenerator generates Pakistani phone numbers.
type PKPhoneGenerator struct {
	BaseGenerator
}

// NewPKPhoneGenerator creates a new Pakistani phone generator.
func NewPKPhoneGenerator() *PKPhoneGenerator {
	return &PKPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "PK_PHONE"},
	}
}

// Generate produces a Pakistani phone number.
// Format: +92 XXX XXXXXXX or 0XXX-XXXXXXX
func (g *PKPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+92")
	hasSpace := strings.Contains(input, " ")
	hasDash := strings.Contains(input, "-")

	// Mobile (03XX) or landline (0XX)
	var areaCode string
	if randomInt(2) == 0 {
		areaCode = fmt.Sprintf("3%d%d", randomInt(10), randomInt(10))
	} else {
		areaCode = fmt.Sprintf("%d%d", 2+randomInt(7), randomInt(10))
	}
	subscriber := generateDigits(7)

	if hasCountryCode {
		if hasSpace {
			return "+92 " + areaCode + " " + subscriber
		}
		return "+92" + areaCode + subscriber
	}
	if hasDash {
		return "0" + areaCode + "-" + subscriber
	}
	if hasSpace {
		return "0" + areaCode + " " + subscriber
	}
	return "0" + areaCode + subscriber
}

// SEPhoneGenerator generates Swedish phone numbers.
type SEPhoneGenerator struct {
	BaseGenerator
}

// NewSEPhoneGenerator creates a new Swedish phone generator.
func NewSEPhoneGenerator() *SEPhoneGenerator {
	return &SEPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "SE_PHONE"},
	}
}

// Generate produces a Swedish phone number.
// Format: +46 XX XXX XX XX or 0XX-XXX XX XX
func (g *SEPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+46")
	hasSpace := strings.Contains(input, " ")
	hasDash := strings.Contains(input, "-")

	areaCode := fmt.Sprintf("%d%d", 1+randomInt(9), randomInt(10))
	subscriber := generateDigits(7)

	if hasCountryCode {
		if hasSpace {
			return "+46 " + areaCode + " " + subscriber[:3] + " " + subscriber[3:5] + " " + subscriber[5:]
		}
		return "+46" + areaCode + subscriber
	}
	if hasDash {
		return "0" + areaCode + "-" + subscriber[:3] + " " + subscriber[3:5] + " " + subscriber[5:]
	}
	if hasSpace {
		return "0" + areaCode + " " + subscriber[:3] + " " + subscriber[3:5] + " " + subscriber[5:]
	}
	return "0" + areaCode + subscriber
}

// SGPhoneGenerator generates Singaporean phone numbers.
type SGPhoneGenerator struct {
	BaseGenerator
}

// NewSGPhoneGenerator creates a new Singaporean phone generator.
func NewSGPhoneGenerator() *SGPhoneGenerator {
	return &SGPhoneGenerator{
		BaseGenerator: BaseGenerator{name: "SG_PHONE"},
	}
}

// Generate produces a Singaporean phone number.
// Format: +65 XXXX XXXX or XXXX XXXX
func (g *SGPhoneGenerator) Generate(input string) string {
	hasCountryCode := strings.HasPrefix(strings.TrimSpace(input), "+65")
	hasSpace := strings.Contains(input, " ")

	// Mobile (8XXX, 9XXX) or landline (6XXX)
	var first string
	if randomInt(2) == 0 {
		first = fmt.Sprintf("%d", 8+randomInt(2)) // 8 or 9
	} else {
		first = "6"
	}
	subscriber := first + generateDigits(7)

	if hasCountryCode {
		if hasSpace {
			return "+65 " + subscriber[:4] + " " + subscriber[4:]
		}
		return "+65" + subscriber
	}
	if hasSpace {
		return subscriber[:4] + " " + subscriber[4:]
	}
	return subscriber
}
