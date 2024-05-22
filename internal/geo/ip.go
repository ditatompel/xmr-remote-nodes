package geo

import (
	"errors"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type IPInfo struct {
	Ip                  string  `json:"ip"`
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
	AsnOrg              string  `json:"asn_org"`
	Asn                 uint    `json:"asn"`
}

func IpInfo(ipAddr string) (*IPInfo, error) {
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

	cityRecord, err := dbCity.City(ip)
	if err != nil {
		return nil, errors.New("Cannot read GeoIP City database")
	}

	asnRecord, err := dbAsn.ASN(ip)
	if err != nil {
		return nil, errors.New("Cannot read GeoIP ASN database")
	}

	qip := IPInfo{
		Ip:                  ipAddr,
		IsAnonymousProxy:    cityRecord.Traits.IsAnonymousProxy,
		IsSatelliteProvider: cityRecord.Traits.IsSatelliteProvider,
		City:                cityRecord.City.Names["en"],
		ContinentName:       cityRecord.Continent.Names["en"],
		ContinentCode:       cityRecord.Continent.Code,
		IsInEuropeanUnion:   cityRecord.Country.IsInEuropeanUnion,
		CountryName:         cityRecord.Country.Names["en"],
		CountryCode:         cityRecord.Country.IsoCode,
		TimeZone:            cityRecord.Location.TimeZone,
		Latitude:            cityRecord.Location.Latitude,
		Longitude:           cityRecord.Location.Longitude,
		AccuracyRadius:      cityRecord.Location.AccuracyRadius,
		AsnOrg:              asnRecord.AutonomousSystemOrganization,
		Asn:                 asnRecord.AutonomousSystemNumber,
	}

	return &qip, nil
}
