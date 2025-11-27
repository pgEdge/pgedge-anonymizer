/*-------------------------------------------------------------------------
 *
 * pgEdge Anonymizer
 *
 * Portions copyright (c) 2025, pgEdge, Inc.
 * This software is released under The PostgreSQL License
 *
 *-------------------------------------------------------------------------
 */

// Package generator provides data generators for anonymization patterns.
package generator

import (
	"github.com/pgedge/pgedge-anonymizer/internal/generator/data"
	"github.com/pgedge/pgedge-anonymizer/internal/generator/data/countries"
)

// Manager coordinates generators and provides access to the data set.
type Manager struct {
	registry    *Registry
	data        *data.DataSet
	countryData *countries.CountryDataSet
}

// FormatPatternConfig holds configuration for creating a format-based generator.
type FormatPatternConfig struct {
	Name    string // Pattern name (becomes generator name)
	Format  string // Format string
	Type    string // "date", "number", or "mask" (auto-detected if empty)
	Min     int64  // Min value for number type
	Max     int64  // Max value for number type
	MinYear int    // Min year for date type
	MaxYear int    // Max year for date type
}

// NewManager creates a new generator manager with all built-in generators.
func NewManager() *Manager {
	dataset := data.Load()
	countryDataset := countries.Load()
	registry := NewRegistry()

	m := &Manager{
		registry:    registry,
		data:        dataset,
		countryData: countryDataset,
	}

	// Register all built-in generators
	m.registerBuiltins()

	return m
}

// registerBuiltins registers all built-in generators.
func (m *Manager) registerBuiltins() {
	// Phone generators (legacy/generic)
	m.registry.Register(NewUSPhoneGenerator())
	m.registry.Register(NewUKPhoneGenerator())
	m.registry.Register(NewInternationalPhoneGenerator())
	m.registry.Register(NewWorldwidePhoneGenerator())

	// Country-specific phone generators
	m.registry.Register(NewAUPhoneGenerator())
	m.registry.Register(NewCAPhoneGenerator())
	m.registry.Register(NewDEPhoneGenerator())
	m.registry.Register(NewESPhoneGenerator())
	m.registry.Register(NewFIPhoneGenerator())
	m.registry.Register(NewFRPhoneGenerator())
	m.registry.Register(NewIEPhoneGenerator())
	m.registry.Register(NewINPhoneGenerator())
	m.registry.Register(NewITPhoneGenerator())
	m.registry.Register(NewJPPhoneGenerator())
	m.registry.Register(NewKRPhoneGenerator())
	m.registry.Register(NewMXPhoneGenerator())
	m.registry.Register(NewNOPhoneGenerator())
	m.registry.Register(NewNZPhoneGenerator())
	m.registry.Register(NewPKPhoneGenerator())
	m.registry.Register(NewSEPhoneGenerator())
	m.registry.Register(NewSGPhoneGenerator())

	// Person data generators (legacy/generic)
	m.registry.Register(NewNameGenerator(m.data))
	m.registry.Register(NewFirstNameGenerator(m.data))
	m.registry.Register(NewLastNameGenerator(m.data))
	m.registry.Register(NewEmailGenerator(m.data))
	m.registry.Register(NewAddressGenerator(m.countryData))
	m.registry.Register(NewCityGenerator(m.countryData))

	// Country-specific name generators
	for _, code := range countries.AllCountries {
		if data := m.countryData.Countries[code]; data != nil {
			m.registry.Register(NewCountryFirstNameGenerator(code, data))
			m.registry.Register(NewCountryLastNameGenerator(code, data))
			m.registry.Register(NewCountryFullNameGenerator(code, data))
			m.registry.Register(NewCountryCityGenerator(code, data))
		}
	}

	// Worldwide name generators
	m.registry.Register(NewWorldwideFirstNameGenerator(m.countryData))
	m.registry.Register(NewWorldwideLastNameGenerator(m.countryData))
	m.registry.Register(NewWorldwideNameGenerator(m.countryData))
	m.registry.Register(NewWorldwideCityGenerator(m.countryData))

	// Postal code generators (legacy/generic)
	m.registry.Register(NewUSZipGenerator())
	m.registry.Register(NewUKPostcodeGenerator())
	m.registry.Register(NewCAPostcodeGenerator())
	m.registry.Register(NewWorldwidePostcodeGenerator())

	// Country-specific postcode generators
	m.registry.Register(NewAUPostcodeGenerator())
	m.registry.Register(NewDEPostcodeGenerator())
	m.registry.Register(NewESPostcodeGenerator())
	m.registry.Register(NewFIPostcodeGenerator())
	m.registry.Register(NewFRPostcodeGenerator())
	m.registry.Register(NewIEPostcodeGenerator())
	m.registry.Register(NewINPostcodeGenerator())
	m.registry.Register(NewITPostcodeGenerator())
	m.registry.Register(NewJPPostcodeGenerator())
	m.registry.Register(NewKRPostcodeGenerator())
	m.registry.Register(NewMXPostcodeGenerator())
	m.registry.Register(NewNOPostcodeGenerator())
	m.registry.Register(NewNZPostcodeGenerator())
	m.registry.Register(NewPKPostcodeGenerator())
	m.registry.Register(NewSEPostcodeGenerator())
	m.registry.Register(NewSGPostcodeGenerator())

	// Country-specific address generators
	if data := m.countryData.Countries[countries.AU]; data != nil {
		m.registry.Register(NewAUAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.CA]; data != nil {
		m.registry.Register(NewCAAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.DE]; data != nil {
		m.registry.Register(NewDEAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.ES]; data != nil {
		m.registry.Register(NewESAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.FI]; data != nil {
		m.registry.Register(NewFIAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.FR]; data != nil {
		m.registry.Register(NewFRAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.IE]; data != nil {
		m.registry.Register(NewIEAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.IN]; data != nil {
		m.registry.Register(NewINAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.IT]; data != nil {
		m.registry.Register(NewITAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.JP]; data != nil {
		m.registry.Register(NewJPAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.KR]; data != nil {
		m.registry.Register(NewKRAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.MX]; data != nil {
		m.registry.Register(NewMXAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.NO]; data != nil {
		m.registry.Register(NewNOAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.NZ]; data != nil {
		m.registry.Register(NewNZAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.PK]; data != nil {
		m.registry.Register(NewPKAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.SE]; data != nil {
		m.registry.Register(NewSEAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.SG]; data != nil {
		m.registry.Register(NewSGAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.UK]; data != nil {
		m.registry.Register(NewUKAddressGenerator(data))
	}
	if data := m.countryData.Countries[countries.US]; data != nil {
		m.registry.Register(NewUSAddressGenerator(data))
	}

	// Worldwide address generator
	m.registry.Register(NewWorldwideAddressGenerator(m.countryData))

	// Financial generators
	m.registry.Register(NewCreditCardGenerator())
	m.registry.Register(NewCreditCardExpiryGenerator())
	m.registry.Register(NewCreditCardCVVGenerator())

	// ID number generators (legacy/generic)
	m.registry.Register(NewSSNGenerator())
	m.registry.Register(NewUKNIGenerator())
	m.registry.Register(NewUKNHSGenerator())
	m.registry.Register(NewPassportGenerator())

	// Country-specific ID number generators
	m.registry.Register(NewAUTFNGenerator())
	m.registry.Register(NewCASINGenerator())
	m.registry.Register(NewDESteurIDGenerator())
	m.registry.Register(NewESNIFGenerator())
	m.registry.Register(NewFIHETUGenerator())
	m.registry.Register(NewFRNIRGenerator())
	m.registry.Register(NewIEPPSGenerator())
	m.registry.Register(NewINAadhaarGenerator())
	m.registry.Register(NewINPANGenerator())
	m.registry.Register(NewITCFGenerator())
	m.registry.Register(NewJPMyNumberGenerator())
	m.registry.Register(NewKRRRNGenerator())
	m.registry.Register(NewMXCURPGenerator())
	m.registry.Register(NewNOFNRGenerator())
	m.registry.Register(NewNZIRDGenerator())
	m.registry.Register(NewPKCNICGenerator())
	m.registry.Register(NewSEPNRGenerator())
	m.registry.Register(NewSGNRICGenerator())
	m.registry.Register(NewUSSSNGenerator())

	// Date generators
	m.registry.Register(NewDOBGenerator())
	m.registry.Register(NewDOBOver13Generator())
	m.registry.Register(NewDOBOver16Generator())
	m.registry.Register(NewDOBOver18Generator())
	m.registry.Register(NewDOBOver21Generator())

	// Text generators
	m.registry.Register(NewLoremGenerator(m.data))

	// Network generators
	m.registry.Register(NewIPv4Generator())
	m.registry.Register(NewIPv6Generator())
	m.registry.Register(NewHostnameGenerator(m.data))
}

// Get retrieves a generator by name.
func (m *Manager) Get(name string) (Generator, bool) {
	return m.registry.Get(name)
}

// List returns all registered generator names.
func (m *Manager) List() []string {
	return m.registry.List()
}

// Data returns the embedded data set.
func (m *Manager) Data() *data.DataSet {
	return m.data
}

// RegisterFormatPattern creates and registers a format-based generator.
func (m *Manager) RegisterFormatPattern(cfg FormatPatternConfig) error {
	// Determine format type
	var formatType FormatType
	switch cfg.Type {
	case "date":
		formatType = FormatTypeDate
	case "number":
		formatType = FormatTypeNumber
	case "mask":
		formatType = FormatTypeMask
	case "":
		// Auto-detect
		formatType = DetectFormatType(cfg.Format)
	default:
		formatType = FormatTypeMask
	}

	// Create format config
	formatConfig := FormatConfig{
		Format:  cfg.Format,
		Type:    formatType,
		Min:     cfg.Min,
		Max:     cfg.Max,
		MinYear: cfg.MinYear,
		MaxYear: cfg.MaxYear,
	}

	// Create and register the generator
	gen := NewFormatGenerator(cfg.Name, formatConfig)
	m.registry.Register(gen)

	return nil
}
