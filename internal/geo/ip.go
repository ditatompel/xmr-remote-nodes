package geo

import (
	"errors"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// IPInfo represents IP address information from Maxmind's GeoLite2 database
type IPInfo struct {
	IP                  string  `json:"ip"`
	IsAnonymousProxy    bool    `json:"is_anonymous_proxy"`
	IsSatelliteProvider bool    `json:"is_satellite_provider"`
	City                string  `json:"city"`
	ContinentName       string  `json:"continent_name"`
	ContinentCode       string  `json:"continent_code"`
	IsInEuropeanUnion   bool    `json:"is_in_european_union"`
	CountryName         string  `json:"country_name"`
	CountryCode         string  `json:"country_code"`
	TimeZone            string  `json:"timezone"`
	Latitude            float64 `json:"latitude"`
	Longitude           float64 `json:"longitude"`
	AccuracyRadius      uint16  `json:"accuracy_radius"`
	ASNOrg              string  `json:"asn_org"`
	ASN                 uint    `json:"asn"`
}

// Info returns GeoIP information from given IP address
func Info(ipAddr string) (*IPInfo, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return nil, errors.New("Invalid IP address")
	}
	dbCity, err := geoip2.Open("./assets/geoip/GeoLite2-City.mmdb")
	if err != nil {
		return nil, errors.New("Cannot open GeoIP City database")
	}
	defer dbCity.Close()

	dbAsn, err := geoip2.Open("./assets/geoip/GeoLite2-ASN.mmdb")
	if err != nil {
		return nil, errors.New("Cannot open GeoIP ASN database")
	}
	defer dbAsn.Close()

	city, err := dbCity.City(ip)
	if err != nil {
		return nil, errors.New("Cannot read GeoIP City database")
	}

	asn, err := dbAsn.ASN(ip)
	if err != nil {
		return nil, errors.New("Cannot read GeoIP ASN database")
	}

	qip := IPInfo{
		IP:                  ipAddr,
		IsAnonymousProxy:    city.Traits.IsAnonymousProxy,
		IsSatelliteProvider: city.Traits.IsSatelliteProvider,
		City:                city.City.Names["en"],
		ContinentName:       city.Continent.Names["en"],
		ContinentCode:       city.Continent.Code,
		IsInEuropeanUnion:   city.Country.IsInEuropeanUnion,
		CountryName:         city.Country.Names["en"],
		CountryCode:         city.Country.IsoCode,
		TimeZone:            city.Location.TimeZone,
		Latitude:            city.Location.Latitude,
		Longitude:           city.Location.Longitude,
		AccuracyRadius:      city.Location.AccuracyRadius,
		ASNOrg:              asn.AutonomousSystemOrganization,
		ASN:                 asn.AutonomousSystemNumber,
	}

	return &qip, nil
}
