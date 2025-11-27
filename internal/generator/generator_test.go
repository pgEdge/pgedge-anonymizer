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
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pgedge/pgedge-anonymizer/internal/generator/data"
	"github.com/pgedge/pgedge-anonymizer/internal/generator/data/countries"
)

// TestManager tests the generator manager and registry
func TestManager(t *testing.T) {
	m := NewManager()

	t.Run("core generators registered", func(t *testing.T) {
		// Core generators that must always be present
		coreGenerators := []string{
			// Phone generators
			"US_PHONE", "UK_PHONE", "INTERNATIONAL_PHONE", "WORLDWIDE_PHONE",
			// Country-specific phone generators
			"AU_PHONE", "CA_PHONE", "DE_PHONE", "ES_PHONE", "FI_PHONE",
			"FR_PHONE", "IE_PHONE", "IN_PHONE", "IT_PHONE", "JP_PHONE",
			"KR_PHONE", "MX_PHONE", "NO_PHONE", "NZ_PHONE", "PK_PHONE",
			"SE_PHONE", "SG_PHONE",
			// Person data generators
			"PERSON_NAME", "PERSON_FIRST_NAME", "PERSON_LAST_NAME",
			"EMAIL", "ADDRESS", "CITY",
			// Worldwide name generators
			"WORLDWIDE_FIRST_NAME", "WORLDWIDE_LAST_NAME",
			"WORLDWIDE_NAME", "WORLDWIDE_CITY",
			// Postal code generators
			"US_ZIP", "UK_POSTCODE", "CA_POSTCODE", "WORLDWIDE_POSTCODE",
			"AU_POSTCODE", "DE_POSTCODE", "ES_POSTCODE", "FI_POSTCODE",
			"FR_POSTCODE", "IE_POSTCODE", "IN_POSTCODE", "IT_POSTCODE",
			"JP_POSTCODE", "KR_POSTCODE", "MX_POSTCODE", "NO_POSTCODE",
			"NZ_POSTCODE", "PK_POSTCODE", "SE_POSTCODE", "SG_POSTCODE",
			// Address generators
			"US_ADDRESS", "UK_ADDRESS", "CA_ADDRESS", "AU_ADDRESS",
			"DE_ADDRESS", "ES_ADDRESS", "FI_ADDRESS", "FR_ADDRESS",
			"IE_ADDRESS", "IN_ADDRESS", "IT_ADDRESS", "JP_ADDRESS",
			"KR_ADDRESS", "MX_ADDRESS", "NO_ADDRESS", "NZ_ADDRESS",
			"PK_ADDRESS", "SE_ADDRESS", "SG_ADDRESS",
			// Financial generators
			"CREDIT_CARD", "CREDIT_CARD_EXPIRY", "CREDIT_CARD_CVV",
			// ID number generators
			"US_SSN", "UK_NI", "UK_NHS", "PASSPORT",
			"AU_TFN", "CA_SIN", "DE_STEUERID", "ES_NIF", "FI_HETU",
			"FR_NIR", "IE_PPS", "IN_AADHAAR", "IN_PAN", "IT_CF",
			"JP_MYNUMBER", "KR_RRN", "MX_CURP", "NO_FNR", "NZ_IRD",
			"PK_CNIC", "SE_PNR", "SG_NRIC", "US_SSN",
			// Date generators
			"DOB", "DOB_OVER_13", "DOB_OVER_16", "DOB_OVER_18", "DOB_OVER_21",
			// Text generators
			"LOREMIPSUM",
			// Network generators
			"IPV4_ADDRESS", "IPV6_ADDRESS", "HOSTNAME",
		}

		for _, name := range coreGenerators {
			if _, ok := m.Get(name); !ok {
				t.Errorf("core generator %s not registered", name)
			}
		}

		// Ensure we have a reasonable number of generators registered
		// (includes country-specific name/city generators for 19 countries)
		registered := m.List()
		if len(registered) < 100 {
			t.Errorf("expected at least 100 generators, got %d", len(registered))
		}
	})

	t.Run("unknown generator returns false", func(t *testing.T) {
		if _, ok := m.Get("UNKNOWN"); ok {
			t.Error("expected unknown generator to return false")
		}
	})
}

// TestUSPhoneGenerator tests US phone number generation
func TestUSPhoneGenerator(t *testing.T) {
	g := NewUSPhoneGenerator()

	if g.Name() != "US_PHONE" {
		t.Errorf("expected name US_PHONE, got %s", g.Name())
	}

	tests := []struct {
		name    string
		input   string
		pattern string
	}{
		// Pattern verifies 555 exchange is used (reserved for fictional use)
		{"dash format", "212-123-4567", `^\d{3}-555-01\d{2}$`},
		{"dot format", "212.123.4567", `^\d{3}\.555\.01\d{2}$`},
		{"parens format", "(212) 123-4567", `^\(\d{3}\) 555[ -]01\d{2}$`},
		{"no separator", "2121234567", `^\d{3}55501\d{2}$`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := g.Generate(tt.input)
			matched, _ := regexp.MatchString(tt.pattern, result)
			if !matched {
				t.Errorf("input %q: result %q doesn't match pattern %s",
					tt.input, result, tt.pattern)
			}
		})
	}
}

// TestUKPhoneGenerator tests UK phone number generation
func TestUKPhoneGenerator(t *testing.T) {
	g := NewUKPhoneGenerator()

	if g.Name() != "UK_PHONE" {
		t.Errorf("expected name UK_PHONE, got %s", g.Name())
	}

	// Ofcom-reserved fictional landline patterns (area code + exchange prefix)
	landlinePatterns := []string{
		`^\+44 20 7946 0\d{3}$`, // London
		`^\+44 117 496 0\d{3}$`, // Bristol
		`^\+44 131 496 0\d{3}$`, // Edinburgh
		`^\+44 161 496 0\d{3}$`, // Manchester
	}
	localLandlinePatterns := []string{
		`^020 7946 0\d{3}$`, // London
		`^0117 496 0\d{3}$`, // Bristol
		`^0131 496 0\d{3}$`, // Edinburgh
		`^0161 496 0\d{3}$`, // Manchester
	}

	t.Run("with country code landline uses Ofcom range", func(t *testing.T) {
		// Generate multiple times to test randomness
		for i := 0; i < 20; i++ {
			result := g.Generate("+44 20 7946 0958")
			if !strings.HasPrefix(result, "+44 ") {
				t.Errorf("expected +44 prefix, got %s", result)
			}
			// Check it matches one of the Ofcom-reserved patterns
			matched := false
			for _, pattern := range landlinePatterns {
				if m, _ := regexp.MatchString(pattern, result); m {
					matched = true
					break
				}
			}
			if !matched {
				t.Errorf("result %q doesn't match any Ofcom-reserved pattern", result)
			}
		}
	})

	t.Run("without country code landline uses Ofcom range", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			result := g.Generate("020 7946 0958")
			if !strings.HasPrefix(result, "0") {
				t.Errorf("expected 0 prefix, got %s", result)
			}
			// Check it matches one of the Ofcom-reserved patterns
			matched := false
			for _, pattern := range localLandlinePatterns {
				if m, _ := regexp.MatchString(pattern, result); m {
					matched = true
					break
				}
			}
			if !matched {
				t.Errorf("result %q doesn't match any Ofcom-reserved pattern", result)
			}
		}
	})

	t.Run("mobile with country code uses Ofcom range", func(t *testing.T) {
		result := g.Generate("+447700900123")
		// Ofcom mobile range: +44 7700 900xxx
		matched, _ := regexp.MatchString(`^\+44 7700 900\d{3}$`, result)
		if !matched {
			t.Errorf("expected +44 7700 900xxx format, got %s", result)
		}
	})

	t.Run("mobile without country code uses Ofcom range", func(t *testing.T) {
		result := g.Generate("07700900123")
		// Ofcom mobile range: 07700 900xxx
		matched, _ := regexp.MatchString(`^07700 900\d{3}$`, result)
		if !matched {
			t.Errorf("expected 07700 900xxx format, got %s", result)
		}
	})
}

// TestInternationalPhoneGenerator tests international phone generation
func TestInternationalPhoneGenerator(t *testing.T) {
	g := NewInternationalPhoneGenerator()

	if g.Name() != "INTERNATIONAL_PHONE" {
		t.Errorf("expected name INTERNATIONAL_PHONE, got %s", g.Name())
	}

	result := g.Generate("+1 555 123 4567")
	if !strings.HasPrefix(result, "+") {
		t.Errorf("expected + prefix, got %s", result)
	}
}

// TestWorldwidePhoneGenerator tests worldwide phone generation
func TestWorldwidePhoneGenerator(t *testing.T) {
	g := NewWorldwidePhoneGenerator()

	if g.Name() != "WORLDWIDE_PHONE" {
		t.Errorf("expected name WORLDWIDE_PHONE, got %s", g.Name())
	}

	t.Run("matches digit count", func(t *testing.T) {
		result := g.Generate("1234567890123")
		if len(result) != 13 {
			t.Errorf("expected 13 digits, got %d: %s", len(result), result)
		}
	})

	t.Run("minimum 10 digits", func(t *testing.T) {
		result := g.Generate("12345")
		if len(result) != 10 {
			t.Errorf("expected 10 digits minimum, got %d: %s",
				len(result), result)
		}
	})
}

// TestNameGenerator tests name generation
func TestNameGenerator(t *testing.T) {
	d := data.Load()
	g := NewNameGenerator(d)

	if g.Name() != "PERSON_NAME" {
		t.Errorf("expected name PERSON_NAME, got %s", g.Name())
	}

	t.Run("space separated", func(t *testing.T) {
		result := g.Generate("John Smith")
		if !strings.Contains(result, " ") {
			t.Errorf("expected space in name, got %s", result)
		}
	})

	t.Run("comma separated", func(t *testing.T) {
		result := g.Generate("Smith, John")
		if !strings.Contains(result, ",") {
			t.Errorf("expected comma in name, got %s", result)
		}
	})

	t.Run("uppercase preservation", func(t *testing.T) {
		result := g.Generate("JOHN SMITH")
		if result != strings.ToUpper(result) {
			t.Errorf("expected uppercase, got %s", result)
		}
	})

	t.Run("lowercase preservation", func(t *testing.T) {
		result := g.Generate("john smith")
		if result != strings.ToLower(result) {
			t.Errorf("expected lowercase, got %s", result)
		}
	})
}

// TestFirstNameGenerator tests first name generation
func TestFirstNameGenerator(t *testing.T) {
	d := data.Load()
	g := NewFirstNameGenerator(d)

	if g.Name() != "PERSON_FIRST_NAME" {
		t.Errorf("expected name PERSON_FIRST_NAME, got %s", g.Name())
	}

	result := g.Generate("John")
	if result == "" {
		t.Error("expected non-empty result")
	}
	if strings.Contains(result, " ") {
		t.Errorf("expected single name, got %s", result)
	}
}

// TestLastNameGenerator tests last name generation
func TestLastNameGenerator(t *testing.T) {
	d := data.Load()
	g := NewLastNameGenerator(d)

	if g.Name() != "PERSON_LAST_NAME" {
		t.Errorf("expected name PERSON_LAST_NAME, got %s", g.Name())
	}

	result := g.Generate("Smith")
	if result == "" {
		t.Error("expected non-empty result")
	}
	if strings.Contains(result, " ") {
		t.Errorf("expected single name, got %s", result)
	}
}

// TestEmailGenerator tests email generation
func TestEmailGenerator(t *testing.T) {
	d := data.Load()
	g := NewEmailGenerator(d)

	if g.Name() != "EMAIL" {
		t.Errorf("expected name EMAIL, got %s", g.Name())
	}

	result := g.Generate("test@example.com")
	if !strings.Contains(result, "@") {
		t.Errorf("expected @ in email, got %s", result)
	}
	if !strings.Contains(result, ".") {
		t.Errorf("expected . in email domain, got %s", result)
	}
}

// TestCreditCardGenerator tests credit card generation
func TestCreditCardGenerator(t *testing.T) {
	g := NewCreditCardGenerator()

	if g.Name() != "CREDIT_CARD" {
		t.Errorf("expected name CREDIT_CARD, got %s", g.Name())
	}

	t.Run("dash separator", func(t *testing.T) {
		result := g.Generate("4532-1234-5678-9012")
		matched, _ := regexp.MatchString(`^\d{4}-\d{4}-\d{4}-\d{4}$`, result)
		if !matched {
			t.Errorf("expected dash-separated format, got %s", result)
		}
	})

	t.Run("space separator", func(t *testing.T) {
		result := g.Generate("4532 1234 5678 9012")
		matched, _ := regexp.MatchString(`^\d{4} \d{4} \d{4} \d{4}$`, result)
		if !matched {
			t.Errorf("expected space-separated format, got %s", result)
		}
	})

	t.Run("no separator", func(t *testing.T) {
		result := g.Generate("4532123456789012")
		matched, _ := regexp.MatchString(`^\d{16}$`, result)
		if !matched {
			t.Errorf("expected 16 digits, got %s", result)
		}
	})

	t.Run("valid Luhn checksum", func(t *testing.T) {
		result := g.Generate("4532123456789012")
		if !isValidLuhn(result) {
			t.Errorf("invalid Luhn checksum for %s", result)
		}
	})
}

// isValidLuhn validates a credit card number using the Luhn algorithm
func isValidLuhn(number string) bool {
	// Remove non-digit characters
	digits := ""
	for _, c := range number {
		if c >= '0' && c <= '9' {
			digits += string(c)
		}
	}

	sum := 0
	nDigits := len(digits)
	parity := nDigits % 2

	for i := 0; i < nDigits; i++ {
		d := int(digits[i] - '0')
		if i%2 == parity {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}

	return sum%10 == 0
}

// TestCreditCardExpiryGenerator tests expiry date generation
func TestCreditCardExpiryGenerator(t *testing.T) {
	g := NewCreditCardExpiryGenerator()

	if g.Name() != "CREDIT_CARD_EXPIRY" {
		t.Errorf("expected name CREDIT_CARD_EXPIRY, got %s", g.Name())
	}

	t.Run("short format", func(t *testing.T) {
		result := g.Generate("12/25")
		matched, _ := regexp.MatchString(`^\d{2}/\d{2}$`, result)
		if !matched {
			t.Errorf("expected MM/YY format, got %s", result)
		}
	})

	t.Run("long format", func(t *testing.T) {
		result := g.Generate("12/2025")
		matched, _ := regexp.MatchString(`^\d{2}/20\d{2}$`, result)
		if !matched {
			t.Errorf("expected MM/YYYY format, got %s", result)
		}
	})
}

// TestCreditCardCVVGenerator tests CVV generation
func TestCreditCardCVVGenerator(t *testing.T) {
	g := NewCreditCardCVVGenerator()

	if g.Name() != "CREDIT_CARD_CVV" {
		t.Errorf("expected name CREDIT_CARD_CVV, got %s", g.Name())
	}

	t.Run("3 digit CVV", func(t *testing.T) {
		result := g.Generate("123")
		if len(result) != 3 {
			t.Errorf("expected 3 digits, got %s", result)
		}
	})

	t.Run("4 digit CVV (Amex)", func(t *testing.T) {
		result := g.Generate("1234")
		if len(result) != 4 {
			t.Errorf("expected 4 digits, got %s", result)
		}
	})
}

// TestSSNGenerator tests SSN generation
func TestSSNGenerator(t *testing.T) {
	g := NewSSNGenerator()

	if g.Name() != "US_SSN" {
		t.Errorf("expected name US_SSN, got %s", g.Name())
	}

	t.Run("dash format", func(t *testing.T) {
		result := g.Generate("123-45-6789")
		matched, _ := regexp.MatchString(`^\d{3}-\d{2}-\d{4}$`, result)
		if !matched {
			t.Errorf("expected XXX-XX-XXXX format, got %s", result)
		}
	})

	t.Run("no separator", func(t *testing.T) {
		result := g.Generate("123456789")
		matched, _ := regexp.MatchString(`^\d{9}$`, result)
		if !matched {
			t.Errorf("expected 9 digits, got %s", result)
		}
	})

	t.Run("valid area number", func(t *testing.T) {
		// Generate multiple and check none have invalid area numbers
		for i := 0; i < 100; i++ {
			result := g.Generate("123-45-6789")
			area := result[:3]
			if area == "000" || area == "666" || area >= "900" {
				t.Errorf("invalid area number: %s", area)
			}
		}
	})
}

// TestUKNIGenerator tests UK National Insurance number generation
func TestUKNIGenerator(t *testing.T) {
	g := NewUKNIGenerator()

	if g.Name() != "UK_NI" {
		t.Errorf("expected name UK_NI, got %s", g.Name())
	}

	result := g.Generate("AB123456C")
	matched, _ := regexp.MatchString(`^[A-Z]{2}\d{6}[A-D]$`, result)
	if !matched {
		t.Errorf("expected LLNNNNNNL format, got %s", result)
	}
}

// TestUKNHSGenerator tests UK NHS number generation
func TestUKNHSGenerator(t *testing.T) {
	g := NewUKNHSGenerator()

	if g.Name() != "UK_NHS" {
		t.Errorf("expected name UK_NHS, got %s", g.Name())
	}

	t.Run("spaced format", func(t *testing.T) {
		result := g.Generate("485 777 3456")
		matched, _ := regexp.MatchString(`^\d{3} \d{3} \d{4}$`, result)
		if !matched {
			t.Errorf("expected XXX XXX XXXX format, got %s", result)
		}
	})

	t.Run("no spaces", func(t *testing.T) {
		result := g.Generate("4857773456")
		matched, _ := regexp.MatchString(`^\d{10}$`, result)
		if !matched {
			t.Errorf("expected 10 digits, got %s", result)
		}
	})
}

// TestPassportGenerator tests passport number generation
func TestPassportGenerator(t *testing.T) {
	g := NewPassportGenerator()

	if g.Name() != "PASSPORT" {
		t.Errorf("expected name PASSPORT, got %s", g.Name())
	}

	result := g.Generate("A12345678")
	matched, _ := regexp.MatchString(`^[A-Z0-9]{9}$`, result)
	if !matched {
		t.Errorf("expected 9 alphanumeric chars, got %s", result)
	}
}

// TestDOBGenerator tests date of birth generation
func TestDOBGenerator(t *testing.T) {
	t.Run("DOB any age", func(t *testing.T) {
		g := NewDOBGenerator()
		if g.Name() != "DOB" {
			t.Errorf("expected name DOB, got %s", g.Name())
		}
		result := g.Generate("1985-03-15")
		matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, result)
		if !matched {
			t.Errorf("expected YYYY-MM-DD format, got %s", result)
		}
	})

	t.Run("DOB over 18", func(t *testing.T) {
		g := NewDOBOver18Generator()
		if g.Name() != "DOB_OVER_18" {
			t.Errorf("expected name DOB_OVER_18, got %s", g.Name())
		}

		result := g.Generate("1985-03-15")
		dob, err := time.Parse("2006-01-02", result)
		if err != nil {
			t.Errorf("failed to parse date: %s", result)
		}

		age := time.Since(dob).Hours() / 24 / 365
		if age < 18 {
			t.Errorf("expected age >= 18, got %.1f for %s", age, result)
		}
	})

	t.Run("DOB over 21", func(t *testing.T) {
		g := NewDOBOver21Generator()
		if g.Name() != "DOB_OVER_21" {
			t.Errorf("expected name DOB_OVER_21, got %s", g.Name())
		}

		result := g.Generate("1985-03-15")
		dob, err := time.Parse("2006-01-02", result)
		if err != nil {
			t.Errorf("failed to parse date: %s", result)
		}

		age := time.Since(dob).Hours() / 24 / 365
		if age < 21 {
			t.Errorf("expected age >= 21, got %.1f for %s", age, result)
		}
	})

	t.Run("US date format", func(t *testing.T) {
		g := NewDOBGenerator()
		result := g.Generate("03/15/1985")
		matched, _ := regexp.MatchString(`^\d{2}/\d{2}/\d{4}$`, result)
		if !matched {
			t.Errorf("expected MM/DD/YYYY format, got %s", result)
		}
	})

	t.Run("long date format", func(t *testing.T) {
		g := NewDOBGenerator()
		result := g.Generate("March 15, 1985")
		if !strings.Contains(result, ",") {
			t.Errorf("expected long format with comma, got %s", result)
		}
	})
}

// TestLoremGenerator tests lorem ipsum generation
func TestLoremGenerator(t *testing.T) {
	d := data.Load()
	g := NewLoremGenerator(d)

	if g.Name() != "LOREMIPSUM" {
		t.Errorf("expected name LOREMIPSUM, got %s", g.Name())
	}

	t.Run("short text", func(t *testing.T) {
		result := g.Generate("Short text")
		if result == "" {
			t.Error("expected non-empty result")
		}
	})

	t.Run("long text", func(t *testing.T) {
		input := "This is a much longer piece of text that should " +
			"generate more lorem ipsum words to match the length."
		result := g.Generate(input)
		if result == "" {
			t.Error("expected non-empty result")
		}
	})
}

// TestAddressGenerator tests address generation
func TestAddressGenerator(t *testing.T) {
	cd := countries.Load()
	g := NewAddressGenerator(cd)

	if g.Name() != "ADDRESS" {
		t.Errorf("expected name ADDRESS, got %s", g.Name())
	}

	t.Run("simple address", func(t *testing.T) {
		result := g.Generate("123 Main St")
		if result == "" {
			t.Error("expected non-empty result")
		}
	})

	t.Run("full address with city", func(t *testing.T) {
		// Needs newline or >30 chars with comma to include city
		input := "123 Main Street, Springfield, CA 90210"
		result := g.Generate(input)
		if result == "" {
			t.Error("expected non-empty result")
		}
		if !strings.Contains(result, ",") {
			t.Errorf("expected comma in full address, got %s", result)
		}
	})

	t.Run("multiline address input generates address", func(t *testing.T) {
		// Note: The worldwide address generator uses diverse international formats
		// and may not preserve multiline format from input
		input := "123 Main St\nSpringfield 90210"
		result := g.Generate(input)
		if result == "" {
			t.Error("expected non-empty address")
		}
	})
}

// TestLuhnCheckDigit tests the Luhn check digit calculation
func TestLuhnCheckDigit(t *testing.T) {
	tests := []struct {
		digits   string
		expected byte
	}{
		{"453201234567901", '2'},
		{"000000000000000", '0'},
		{"123456789012345", '0'},
	}

	for _, tt := range tests {
		t.Run(tt.digits, func(t *testing.T) {
			result := luhnCheckDigit(tt.digits)
			// Verify the resulting number is valid
			fullNumber := tt.digits + string(result)
			if !isValidLuhn(fullNumber) {
				t.Errorf("luhnCheckDigit(%s) = %c, but full number %s "+
					"fails Luhn validation", tt.digits, result, fullNumber)
			}
		})
	}
}

// TestRandomHelpers tests helper functions
func TestRandomHelpers(t *testing.T) {
	t.Run("randomInt bounds", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			r := randomInt(10)
			if r < 0 || r >= 10 {
				t.Errorf("randomInt(10) out of bounds: %d", r)
			}
		}
	})

	t.Run("randomInt zero", func(t *testing.T) {
		r := randomInt(0)
		if r != 0 {
			t.Errorf("randomInt(0) should return 0, got %d", r)
		}
	})

	t.Run("randomDigit", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			d := randomDigit()
			if d < '0' || d > '9' {
				t.Errorf("randomDigit out of bounds: %c", d)
			}
		}
	})

	t.Run("randomDigitNonZero", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			d := randomDigitNonZero()
			if d < '1' || d > '9' {
				t.Errorf("randomDigitNonZero out of bounds: %c", d)
			}
		}
	})

	t.Run("randomString empty", func(t *testing.T) {
		result := randomString([]string{})
		if result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
	})

	t.Run("randomString single", func(t *testing.T) {
		result := randomString([]string{"test"})
		if result != "test" {
			t.Errorf("expected 'test', got %s", result)
		}
	})

	t.Run("generateDigits", func(t *testing.T) {
		result := generateDigits(5)
		if len(result) != 5 {
			t.Errorf("expected 5 digits, got %d: %s", len(result), result)
		}
		for _, c := range result {
			if c < '0' || c > '9' {
				t.Errorf("non-digit in result: %c", c)
			}
		}
	})
}

// TestPhoneFormatDetection tests phone format detection
func TestPhoneFormatDetection(t *testing.T) {
	tests := []struct {
		input     string
		hasParens bool
		separator byte
	}{
		{"555-123-4567", false, '-'},
		{"555.123.4567", false, '.'},
		{"555 123 4567", false, ' '},
		// Note: space comes before dash in "(555) 123-4567", so space is detected
		{"(555) 123-4567", true, ' '},
		{"5551234567", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			format := detectPhoneFormat(tt.input)
			if format.hasParens != tt.hasParens {
				t.Errorf("hasParens: expected %v, got %v",
					tt.hasParens, format.hasParens)
			}
			if format.separator != tt.separator {
				t.Errorf("separator: expected %c, got %c",
					tt.separator, format.separator)
			}
		})
	}
}

// TestDataLoad tests that embedded data loads correctly
func TestDataLoad(t *testing.T) {
	d := data.Load()

	if len(d.FirstNames) == 0 {
		t.Error("FirstNames not loaded")
	}
	if len(d.LastNames) == 0 {
		t.Error("LastNames not loaded")
	}
	if len(d.StreetNames) == 0 {
		t.Error("StreetNames not loaded")
	}
	if len(d.Cities) == 0 {
		t.Error("Cities not loaded")
	}
	if len(d.Domains) == 0 {
		t.Error("Domains not loaded")
	}
	if len(d.LoremWords) == 0 {
		t.Error("LoremWords not loaded")
	}
}

// TestFormatGenerator tests format-based generation
func TestFormatGenerator(t *testing.T) {
	t.Run("date format strftime codes", func(t *testing.T) {
		g := NewFormatGenerator("TEST_DATE", FormatConfig{
			Format:  "%Y-%m-%d",
			Type:    FormatTypeDate,
			MinYear: 1990,
			MaxYear: 2020,
		})

		if g.Name() != "TEST_DATE" {
			t.Errorf("expected name TEST_DATE, got %s", g.Name())
		}

		result := g.Generate("")
		matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, result)
		if !matched {
			t.Errorf("expected YYYY-MM-DD format, got %s", result)
		}

		// Verify year is in range
		year := result[:4]
		if year < "1990" || year > "2020" {
			t.Errorf("year %s out of range 1990-2020", year)
		}
	})

	t.Run("date format US style", func(t *testing.T) {
		g := NewFormatGenerator("TEST_DATE_US", FormatConfig{
			Format: "%m/%d/%Y",
			Type:   FormatTypeDate,
		})

		result := g.Generate("")
		matched, _ := regexp.MatchString(`^\d{2}/\d{2}/\d{4}$`, result)
		if !matched {
			t.Errorf("expected MM/DD/YYYY format, got %s", result)
		}
	})

	t.Run("date format with time", func(t *testing.T) {
		g := NewFormatGenerator("TEST_DATETIME", FormatConfig{
			Format: "%Y-%m-%d %H:%M:%S",
			Type:   FormatTypeDate,
		})

		result := g.Generate("")
		matched, _ := regexp.MatchString(
			`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, result)
		if !matched {
			t.Errorf("expected datetime format, got %s", result)
		}
	})

	t.Run("date format with month names", func(t *testing.T) {
		g := NewFormatGenerator("TEST_DATE_LONG", FormatConfig{
			Format: "%B %d, %Y",
			Type:   FormatTypeDate,
		})

		result := g.Generate("")
		// Should contain a month name like "January", "February", etc.
		months := []string{"January", "February", "March", "April", "May",
			"June", "July", "August", "September", "October", "November",
			"December"}
		found := false
		for _, month := range months {
			if strings.Contains(result, month) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected month name in result, got %s", result)
		}
	})

	t.Run("mask format digits", func(t *testing.T) {
		g := NewFormatGenerator("TEST_MASK_DIGITS", FormatConfig{
			Format: "###-##-####",
			Type:   FormatTypeMask,
		})

		result := g.Generate("")
		matched, _ := regexp.MatchString(`^\d{3}-\d{2}-\d{4}$`, result)
		if !matched {
			t.Errorf("expected SSN-like format, got %s", result)
		}
	})

	t.Run("mask format letters", func(t *testing.T) {
		g := NewFormatGenerator("TEST_MASK_LETTERS", FormatConfig{
			Format: "AA-####",
			Type:   FormatTypeMask,
		})

		result := g.Generate("")
		matched, _ := regexp.MatchString(`^[A-Z]{2}-\d{4}$`, result)
		if !matched {
			t.Errorf("expected XX-NNNN format, got %s", result)
		}
	})

	t.Run("mask format mixed", func(t *testing.T) {
		g := NewFormatGenerator("TEST_MASK_MIXED", FormatConfig{
			Format: "XXX-aaa-###",
			Type:   FormatTypeMask,
		})

		result := g.Generate("")
		// XXX = alphanumeric upper, aaa = lowercase, ### = digits
		matched, _ := regexp.MatchString(
			`^[A-Z0-9]{3}-[a-z]{3}-\d{3}$`, result)
		if !matched {
			t.Errorf("expected mixed format, got %s", result)
		}
	})

	t.Run("mask format with escape", func(t *testing.T) {
		g := NewFormatGenerator("TEST_MASK_ESCAPE", FormatConfig{
			Format: `\A###`,
			Type:   FormatTypeMask,
		})

		result := g.Generate("")
		matched, _ := regexp.MatchString(`^A\d{3}$`, result)
		if !matched {
			t.Errorf("expected A followed by 3 digits, got %s", result)
		}
	})

	t.Run("number format printf", func(t *testing.T) {
		g := NewFormatGenerator("TEST_NUMBER", FormatConfig{
			Format: "%08d",
			Type:   FormatTypeNumber,
			Min:    1,
			Max:    99999999,
		})

		result := g.Generate("")
		if len(result) != 8 {
			t.Errorf("expected 8 digit padded number, got %s", result)
		}
		matched, _ := regexp.MatchString(`^\d{8}$`, result)
		if !matched {
			t.Errorf("expected 8 digits, got %s", result)
		}
	})

	t.Run("number format simple", func(t *testing.T) {
		g := NewFormatGenerator("TEST_NUMBER_SIMPLE", FormatConfig{
			Format: "%d",
			Type:   FormatTypeNumber,
			Min:    100,
			Max:    200,
		})

		result := g.Generate("")
		// Result should be a number between 100 and 200
		matched, _ := regexp.MatchString(`^\d+$`, result)
		if !matched {
			t.Errorf("expected number, got %s", result)
		}
	})

	t.Run("auto detect date type", func(t *testing.T) {
		detected := DetectFormatType("%Y-%m-%d")
		if detected != FormatTypeDate {
			t.Errorf("expected date type, got %s", detected)
		}
	})

	t.Run("auto detect number type", func(t *testing.T) {
		detected := DetectFormatType("%08d")
		if detected != FormatTypeNumber {
			t.Errorf("expected number type, got %s", detected)
		}
	})

	t.Run("auto detect mask type", func(t *testing.T) {
		detected := DetectFormatType("###-##-####")
		if detected != FormatTypeMask {
			t.Errorf("expected mask type, got %s", detected)
		}
	})
}

// TestManagerRegisterFormatPattern tests dynamic format pattern registration
func TestManagerRegisterFormatPattern(t *testing.T) {
	m := NewManager()

	// Register a custom format pattern
	err := m.RegisterFormatPattern(FormatPatternConfig{
		Name:   "CUSTOM_DATE",
		Format: "%Y-%m-%d",
		Type:   "date",
	})
	if err != nil {
		t.Fatalf("failed to register format pattern: %v", err)
	}

	// Verify it can be retrieved
	gen, ok := m.Get("CUSTOM_DATE")
	if !ok {
		t.Fatal("CUSTOM_DATE generator not found")
	}

	result := gen.Generate("")
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, result)
	if !matched {
		t.Errorf("expected YYYY-MM-DD format, got %s", result)
	}

	// Register another pattern with auto-detection
	err = m.RegisterFormatPattern(FormatPatternConfig{
		Name:   "CUSTOM_ID",
		Format: "ID-####-AAA",
	})
	if err != nil {
		t.Fatalf("failed to register format pattern: %v", err)
	}

	gen2, ok := m.Get("CUSTOM_ID")
	if !ok {
		t.Fatal("CUSTOM_ID generator not found")
	}

	result2 := gen2.Generate("")
	matched2, _ := regexp.MatchString(`^ID-\d{4}-[A-Z]{3}$`, result2)
	if !matched2 {
		t.Errorf("expected ID-NNNN-XXX format, got %s", result2)
	}
}

// TestIPv4Generator tests IPv4 address generation
func TestIPv4Generator(t *testing.T) {
	g := NewIPv4Generator()

	if g.Name() != "IPV4_ADDRESS" {
		t.Errorf("expected name IPV4_ADDRESS, got %s", g.Name())
	}

	t.Run("generates valid IPv4 format", func(t *testing.T) {
		ipv4Pattern := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
		for i := 0; i < 100; i++ {
			result := g.Generate("")
			if !ipv4Pattern.MatchString(result) {
				t.Errorf("expected IPv4 format, got %s", result)
			}
		}
	})

	t.Run("generates valid octet values", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			result := g.Generate("")
			parts := strings.Split(result, ".")
			if len(parts) != 4 {
				t.Errorf("expected 4 octets, got %d in %s", len(parts), result)
				continue
			}
			for _, part := range parts {
				val, err := strconv.Atoi(part)
				if err != nil {
					t.Errorf("invalid octet %q in %s: %v", part, result, err)
					continue
				}
				if val < 0 || val > 255 {
					t.Errorf("invalid octet value %d in %s", val, result)
				}
			}
		}
	})
}

// TestIPv6Generator tests IPv6 address generation
func TestIPv6Generator(t *testing.T) {
	g := NewIPv6Generator()

	if g.Name() != "IPV6_ADDRESS" {
		t.Errorf("expected name IPV6_ADDRESS, got %s", g.Name())
	}

	t.Run("generates valid IPv6 full format", func(t *testing.T) {
		result := g.Generate("2001:0db8:85a3:0000:0000:8a2e:0370:7334")
		// Should be 8 groups separated by colons
		parts := strings.Split(result, ":")
		if len(parts) < 3 {
			t.Errorf("expected IPv6 format with colons, got %s", result)
		}
	})

	t.Run("generates compressed format for compressed input", func(t *testing.T) {
		result := g.Generate("2001:db8::1")
		// Result should contain valid hex characters and colons
		matched, _ := regexp.MatchString(`^[0-9a-fA-F:]+$`, result)
		if !matched {
			t.Errorf("expected IPv6 format, got %s", result)
		}
	})

	t.Run("preserves case", func(t *testing.T) {
		// Uppercase input should produce uppercase output
		resultUpper := g.Generate("2001:0DB8:85A3:0000:0000:8A2E:0370:7334")
		if strings.ContainsAny(resultUpper, "abcdef") && !strings.ContainsAny(resultUpper, "ABCDEF") {
			t.Errorf("uppercase input should preserve case, got lowercase: %s", resultUpper)
		}
	})
}

// TestHostnameGenerator tests hostname generation
func TestHostnameGenerator(t *testing.T) {
	d := data.Load()
	g := NewHostnameGenerator(d)

	if g.Name() != "HOSTNAME" {
		t.Errorf("expected name HOSTNAME, got %s", g.Name())
	}

	t.Run("generates simple hostname", func(t *testing.T) {
		result := g.Generate("webserver")
		// Should be a simple alphanumeric hostname
		matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, result)
		if !matched {
			t.Errorf("expected simple hostname, got %s", result)
		}
	})

	t.Run("generates FQDN for FQDN input", func(t *testing.T) {
		result := g.Generate("server01.example.com")
		// Should contain a dot (FQDN)
		if !strings.Contains(result, ".") {
			t.Errorf("expected FQDN with dot, got %s", result)
		}
	})

	t.Run("generates hostname with number for numbered input", func(t *testing.T) {
		result := g.Generate("web01")
		// Should contain digits
		hasDigit := false
		for _, c := range result {
			if c >= '0' && c <= '9' {
				hasDigit = true
				break
			}
		}
		if !hasDigit {
			t.Errorf("expected hostname with number, got %s", result)
		}
	})
}

// TestUSZipGenerator tests US ZIP code generation
func TestUSZipGenerator(t *testing.T) {
	g := NewUSZipGenerator()

	if g.Name() != "US_ZIP" {
		t.Errorf("expected name US_ZIP, got %s", g.Name())
	}

	t.Run("generates 5-digit ZIP", func(t *testing.T) {
		result := g.Generate("12345")
		matched, _ := regexp.MatchString(`^\d{5}$`, result)
		if !matched {
			t.Errorf("expected 5-digit ZIP, got %s", result)
		}
	})

	t.Run("generates ZIP+4 format", func(t *testing.T) {
		result := g.Generate("12345-6789")
		matched, _ := regexp.MatchString(`^\d{5}-\d{4}$`, result)
		if !matched {
			t.Errorf("expected ZIP+4 format, got %s", result)
		}
	})
}

// TestCityGenerator tests city name generation
func TestCityGenerator(t *testing.T) {
	cd := countries.Load()
	g := NewCityGenerator(cd)

	if g.Name() != "CITY" {
		t.Errorf("expected name CITY, got %s", g.Name())
	}

	t.Run("generates city name", func(t *testing.T) {
		result := g.Generate("New York")
		if result == "" {
			t.Error("expected non-empty city name")
		}
	})

	t.Run("preserves uppercase", func(t *testing.T) {
		result := g.Generate("NEW YORK")
		if strings.ToUpper(result) != result {
			t.Errorf("expected uppercase city, got %s", result)
		}
	})

	t.Run("preserves lowercase", func(t *testing.T) {
		result := g.Generate("new york")
		if strings.ToLower(result) != result {
			t.Errorf("expected lowercase city, got %s", result)
		}
	})
}

// TestUKPostcodeGenerator tests UK postcode generation
func TestUKPostcodeGenerator(t *testing.T) {
	g := NewUKPostcodeGenerator()

	if g.Name() != "UK_POSTCODE" {
		t.Errorf("expected name UK_POSTCODE, got %s", g.Name())
	}

	t.Run("generates valid UK postcode with space", func(t *testing.T) {
		result := g.Generate("SW1A 1AA")
		// UK postcodes: 2-4 chars outward + space + 3 chars inward
		matched, _ := regexp.MatchString(`^[A-Z]{1,2}\d[A-Z\d]? \d[A-Z]{2}$`, result)
		if !matched {
			t.Errorf("expected UK postcode format, got %s", result)
		}
	})

	t.Run("generates UK postcode without space", func(t *testing.T) {
		result := g.Generate("SW1A1AA")
		if strings.Contains(result, " ") {
			t.Errorf("expected no space in postcode, got %s", result)
		}
	})
}

// TestCAPostcodeGenerator tests Canadian postcode generation
func TestCAPostcodeGenerator(t *testing.T) {
	g := NewCAPostcodeGenerator()

	if g.Name() != "CA_POSTCODE" {
		t.Errorf("expected name CA_POSTCODE, got %s", g.Name())
	}

	t.Run("generates valid Canadian postcode with space", func(t *testing.T) {
		result := g.Generate("K1A 0B1")
		// Canadian postcodes: A9A 9A9 format
		matched, _ := regexp.MatchString(`^[A-Z]\d[A-Z] \d[A-Z]\d$`, result)
		if !matched {
			t.Errorf("expected Canadian postcode format A9A 9A9, got %s", result)
		}
	})

	t.Run("generates Canadian postcode without space", func(t *testing.T) {
		result := g.Generate("K1A0B1")
		if strings.Contains(result, " ") {
			t.Errorf("expected no space in postcode, got %s", result)
		}
		matched, _ := regexp.MatchString(`^[A-Z]\d[A-Z]\d[A-Z]\d$`, result)
		if !matched {
			t.Errorf("expected Canadian postcode format A9A9A9, got %s", result)
		}
	})
}

// TestWorldwidePostcodeGenerator tests worldwide postcode generation
func TestWorldwidePostcodeGenerator(t *testing.T) {
	g := NewWorldwidePostcodeGenerator()

	if g.Name() != "WORLDWIDE_POSTCODE" {
		t.Errorf("expected name WORLDWIDE_POSTCODE, got %s", g.Name())
	}

	t.Run("generates US format for digit-only input", func(t *testing.T) {
		result := g.Generate("12345")
		matched, _ := regexp.MatchString(`^\d{5}$`, result)
		if !matched {
			t.Errorf("expected US ZIP format for digit input, got %s", result)
		}
	})

	t.Run("generates UK format for UK-style input", func(t *testing.T) {
		result := g.Generate("SW1A 1AA")
		// Should generate a UK-style postcode
		matched, _ := regexp.MatchString(`^[A-Z]{1,2}\d[A-Z\d]? ?\d[A-Z]{2}$`, result)
		if !matched {
			t.Errorf("expected UK postcode format, got %s", result)
		}
	})

	t.Run("generates Canadian format for Canadian-style input", func(t *testing.T) {
		result := g.Generate("K1A 0B1")
		// Should generate a Canadian-style postcode
		matched, _ := regexp.MatchString(`^[A-Z]\d[A-Z] ?\d[A-Z]\d$`, result)
		if !matched {
			t.Errorf("expected Canadian postcode format, got %s", result)
		}
	})
}
