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

	"github.com/pgedge/pgedge-anonymizer/internal/generator/data/countries"
)

// CountryAddressGenerator generates addresses for a specific country.
type CountryAddressGenerator struct {
	BaseGenerator
	cities      []string
	streetTypes []string
	format      func(num int, street, streetType, city, postcode string) string
	postcodeGen Generator
}

// NewUSAddressGenerator creates a US address generator.
func NewUSAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "US_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"St", "Ave", "Blvd", "Dr", "Ln", "Rd", "Way", "Ct", "Pl"},
		format: func(num int, street, streetType, city, postcode string) string {
			return fmt.Sprintf("%d %s %s, %s %s", num, street, streetType, city, postcode)
		},
		postcodeGen: NewUSZipGenerator(),
	}
}

// NewUKAddressGenerator creates a UK address generator.
func NewUKAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "UK_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"Street", "Road", "Avenue", "Lane", "Close", "Drive", "Way", "Gardens", "Crescent"},
		format: func(num int, street, streetType, city, postcode string) string {
			return fmt.Sprintf("%d %s %s, %s, %s", num, street, streetType, city, postcode)
		},
		postcodeGen: NewUKPostcodeGenerator(),
	}
}

// NewCAAddressGenerator creates a Canadian address generator.
func NewCAAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "CA_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"St", "Ave", "Blvd", "Dr", "Rd", "Way", "Cres", "Pl"},
		format: func(num int, street, streetType, city, postcode string) string {
			return fmt.Sprintf("%d %s %s, %s %s", num, street, streetType, city, postcode)
		},
		postcodeGen: NewCAPostcodeGenerator(),
	}
}

// NewAUAddressGenerator creates an Australian address generator.
func NewAUAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "AU_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"Street", "Road", "Avenue", "Drive", "Court", "Place", "Crescent", "Parade"},
		format: func(num int, street, streetType, city, postcode string) string {
			return fmt.Sprintf("%d %s %s, %s %s", num, street, streetType, city, postcode)
		},
		postcodeGen: NewAUPostcodeGenerator(),
	}
}

// NewDEAddressGenerator creates a German address generator.
func NewDEAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "DE_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"straße", "weg", "platz", "allee", "ring", "gasse"},
		format: func(num int, street, streetType, city, postcode string) string {
			// German format: Streetname + number, postcode city
			return fmt.Sprintf("%s%s %d, %s %s", street, streetType, num, postcode, city)
		},
		postcodeGen: NewDEPostcodeGenerator(),
	}
}

// NewESAddressGenerator creates a Spanish address generator.
func NewESAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "ES_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"Calle", "Avenida", "Plaza", "Paseo", "Carrer", "Carretera"},
		format: func(num int, street, streetType, city, postcode string) string {
			// Spanish format: Street type Street name, number, postcode city
			return fmt.Sprintf("%s %s, %d, %s %s", streetType, street, num, postcode, city)
		},
		postcodeGen: NewESPostcodeGenerator(),
	}
}

// NewFIAddressGenerator creates a Finnish address generator.
func NewFIAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "FI_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"katu", "tie", "polku", "kuja", "puisto"},
		format: func(num int, street, streetType, city, postcode string) string {
			// Finnish format: Streetname + suffix number, postcode city
			return fmt.Sprintf("%s%s %d, %s %s", street, streetType, num, postcode, city)
		},
		postcodeGen: NewFIPostcodeGenerator(),
	}
}

// NewFRAddressGenerator creates a French address generator.
func NewFRAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "FR_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"Rue", "Avenue", "Boulevard", "Place", "Chemin", "Allée"},
		format: func(num int, street, streetType, city, postcode string) string {
			// French format: number street type street name, postcode city
			return fmt.Sprintf("%d %s %s, %s %s", num, streetType, street, postcode, city)
		},
		postcodeGen: NewFRPostcodeGenerator(),
	}
}

// NewIEAddressGenerator creates an Irish address generator.
func NewIEAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "IE_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"Street", "Road", "Avenue", "Lane", "Drive", "Park", "Close", "Grove"},
		format: func(num int, street, streetType, city, postcode string) string {
			return fmt.Sprintf("%d %s %s, %s, %s", num, street, streetType, city, postcode)
		},
		postcodeGen: NewIEPostcodeGenerator(),
	}
}

// NewINAddressGenerator creates an Indian address generator.
func NewINAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "IN_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"Road", "Street", "Marg", "Nagar", "Colony", "Lane", "Gali"},
		format: func(num int, street, streetType, city, postcode string) string {
			return fmt.Sprintf("%d %s %s, %s - %s", num, street, streetType, city, postcode)
		},
		postcodeGen: NewINPostcodeGenerator(),
	}
}

// NewITAddressGenerator creates an Italian address generator.
func NewITAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "IT_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"Via", "Viale", "Piazza", "Corso", "Largo", "Vicolo"},
		format: func(num int, street, streetType, city, postcode string) string {
			// Italian format: Street type Street name, number, postcode city
			return fmt.Sprintf("%s %s, %d, %s %s", streetType, street, num, postcode, city)
		},
		postcodeGen: NewITPostcodeGenerator(),
	}
}

// NewJPAddressGenerator creates a Japanese address generator.
func NewJPAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "JP_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{},
		format: func(num int, street, streetType, city, postcode string) string {
			// Japanese format: postcode city district-block-number
			block := 1 + randomInt(30)
			lot := 1 + randomInt(20)
			return fmt.Sprintf("〒%s %s %d-%d-%d", postcode, city, block, lot, num)
		},
		postcodeGen: NewJPPostcodeGenerator(),
	}
}

// NewKRAddressGenerator creates a South Korean address generator.
func NewKRAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "KR_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"로", "길", "대로"},
		format: func(num int, street, streetType, city, postcode string) string {
			// Korean format: city street+type number (postcode)
			return fmt.Sprintf("%s %s%s %d (%s)", city, street, streetType, num, postcode)
		},
		postcodeGen: NewKRPostcodeGenerator(),
	}
}

// NewMXAddressGenerator creates a Mexican address generator.
func NewMXAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "MX_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"Calle", "Avenida", "Boulevard", "Calzada", "Privada"},
		format: func(num int, street, streetType, city, postcode string) string {
			// Mexican format: Street type Street name #number, postcode city
			return fmt.Sprintf("%s %s #%d, %s %s", streetType, street, num, postcode, city)
		},
		postcodeGen: NewMXPostcodeGenerator(),
	}
}

// NewNOAddressGenerator creates a Norwegian address generator.
func NewNOAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "NO_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"gate", "vei", "veien", "plass"},
		format: func(num int, street, streetType, city, postcode string) string {
			// Norwegian format: Streetname + suffix number, postcode city
			return fmt.Sprintf("%s%s %d, %s %s", street, streetType, num, postcode, city)
		},
		postcodeGen: NewNOPostcodeGenerator(),
	}
}

// NewNZAddressGenerator creates a New Zealand address generator.
func NewNZAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "NZ_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"Street", "Road", "Avenue", "Drive", "Place", "Terrace", "Crescent"},
		format: func(num int, street, streetType, city, postcode string) string {
			return fmt.Sprintf("%d %s %s, %s %s", num, street, streetType, city, postcode)
		},
		postcodeGen: NewNZPostcodeGenerator(),
	}
}

// NewPKAddressGenerator creates a Pakistani address generator.
func NewPKAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "PK_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"Road", "Street", "Colony", "Block", "Sector"},
		format: func(num int, street, streetType, city, postcode string) string {
			return fmt.Sprintf("%d %s %s, %s - %s", num, street, streetType, city, postcode)
		},
		postcodeGen: NewPKPostcodeGenerator(),
	}
}

// NewSEAddressGenerator creates a Swedish address generator.
func NewSEAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "SE_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"gatan", "vägen", "torget", "platsen"},
		format: func(num int, street, streetType, city, postcode string) string {
			// Swedish format: Streetname + suffix number, postcode city
			return fmt.Sprintf("%s%s %d, %s %s", street, streetType, num, postcode, city)
		},
		postcodeGen: NewSEPostcodeGenerator(),
	}
}

// NewSGAddressGenerator creates a Singaporean address generator.
func NewSGAddressGenerator(data *countries.CountryData) *CountryAddressGenerator {
	return &CountryAddressGenerator{
		BaseGenerator: BaseGenerator{name: "SG_ADDRESS"},
		cities:        data.Cities,
		streetTypes:   []string{"Road", "Street", "Avenue", "Drive", "Lane", "Crescent", "Way"},
		format: func(num int, street, streetType, city, postcode string) string {
			// Singapore format: Block number Street name, Singapore postcode
			return fmt.Sprintf("Blk %d %s %s, Singapore %s", num, street, streetType, postcode)
		},
		postcodeGen: NewSGPostcodeGenerator(),
	}
}

// Generic street names used for countries
var genericStreetNames = []string{
	"Main", "Oak", "Maple", "Park", "Lake", "Hill", "River",
	"Forest", "Garden", "Central", "North", "South", "East", "West",
}

// WorldwideAddressGenerator generates addresses from any country.
type WorldwideAddressGenerator struct {
	BaseGenerator
	generators []Generator
}

// NewWorldwideAddressGenerator creates a worldwide address generator.
func NewWorldwideAddressGenerator(data *countries.CountryDataSet) *WorldwideAddressGenerator {
	return &WorldwideAddressGenerator{
		BaseGenerator: BaseGenerator{name: "WORLDWIDE_ADDRESS"},
		generators: []Generator{
			NewUSAddressGenerator(data.Get(countries.US)),
			NewUKAddressGenerator(data.Get(countries.UK)),
			NewCAAddressGenerator(data.Get(countries.CA)),
			NewAUAddressGenerator(data.Get(countries.AU)),
			NewDEAddressGenerator(data.Get(countries.DE)),
			NewESAddressGenerator(data.Get(countries.ES)),
			NewFIAddressGenerator(data.Get(countries.FI)),
			NewFRAddressGenerator(data.Get(countries.FR)),
			NewIEAddressGenerator(data.Get(countries.IE)),
			NewINAddressGenerator(data.Get(countries.IN)),
			NewITAddressGenerator(data.Get(countries.IT)),
			NewJPAddressGenerator(data.Get(countries.JP)),
			NewKRAddressGenerator(data.Get(countries.KR)),
			NewMXAddressGenerator(data.Get(countries.MX)),
			NewNOAddressGenerator(data.Get(countries.NO)),
			NewNZAddressGenerator(data.Get(countries.NZ)),
			NewPKAddressGenerator(data.Get(countries.PK)),
			NewSEAddressGenerator(data.Get(countries.SE)),
			NewSGAddressGenerator(data.Get(countries.SG)),
		},
	}
}

// Generate produces an address from a randomly selected country.
func (g *WorldwideAddressGenerator) Generate(input string) string {
	// Pick a random country generator
	gen := g.generators[randomInt(len(g.generators))]
	return gen.Generate(input)
}

// Generate produces a street address for the country.
func (g *CountryAddressGenerator) Generate(input string) string {
	streetNum := 1 + randomInt(999)
	streetName := randomString(genericStreetNames)
	streetType := ""
	if len(g.streetTypes) > 0 {
		streetType = randomString(g.streetTypes)
	}
	city := randomString(g.cities)
	postcode := g.postcodeGen.Generate(input)

	result := g.format(streetNum, streetName, streetType, city, postcode)

	// Preserve case if needed
	if strings.ToUpper(input) == input && len(input) > 1 {
		return strings.ToUpper(result)
	}
	if strings.ToLower(input) == input && len(input) > 1 {
		return strings.ToLower(result)
	}
	return result
}
