package snmp

import (
    "log"
    "fmt"

    "github.com/pkg/errors"
    g "github.com/soniah/gosnmp"
    "github.com/elastic/beats/v7/metricbeat/mb"
)

// SNMP type represents a full SNMP agent
type SNMP struct {
    client    *g.GoSNMP
    results   []string
}

func NewSNMP(base mb.BaseMetricSet) (*SNMP, error) {
    config := defaultConfig()
    err := base.Module().UnpackConfig(&config)
    if err != nil {
        return nil, err
    }
    
    ver, err := parseVersion(config.Version)
    if err != nil {
        return nil, err
    }

    return &SNMP{
        client: &g.GoSNMP{
            Target: config.Host,
            Port: uint16(config.Port),
            Community: config.Community,
            Version: ver,
            Timeout: config.Timeout,
            Logger: nil,
            MaxOids: config.MaxOids,
        },
    }, nil
}

func (s *SNMP) Get(oidslice []string) (*g.SnmpPacket, error) {
    err := s.client.Connect()
    result, err := s.client.Get(oidslice)
    if err != nil {
        log.Fatalf("Error while executing SNMP Get: %v", err)
    }
    defer s.client.Conn.Close()
    
    return result, nil
}

func (s *SNMP) Walk(oid string) error{
    err := s.client.Connect()
    if err != nil {
        log.Fatalf("Error while connecting to SNMP port: %v", err)
    }

    err = s.client.Walk(oid, s.walkf)
    if err != nil {
        log.Fatalf("Error while executing SNMP Walk: %v", err)
    }

    defer s.client.Conn.Close()

    return err
}

func (s *SNMP) BulkWalk(oid string) error {
    err := s.client.Connect()
    if err != nil {
        log.Fatalf("Error while connecting to SNMP port: %v", err)
    }

    err = s.client.BulkWalk(oid, s.walkf)
    if err != nil {
        log.Fatalf("Error while executing SNMP Walk: %v", err)
    }

    defer s.client.Conn.Close()

    return err
}

func (s *SNMP) BulkWalkAll(oid string) (map[string][]g.SnmpPDU, error) {
    var results []g.SnmpPDU
    err := s.client.Connect()
    if err != nil {
        log.Fatalf("Error while connecting to SNMP port: %v", err)
    }

    results, err = s.client.BulkWalkAll(oid)
    if err != nil {
        log.Fatalf("Error while executing SNMP Walk: %v", err)
    }

    resultsArray := make(map[string][]g.SnmpPDU, len(results))
    for _, entry := range results {
        resultsArray[entry.Name[len(entry.Name)-1:]] = append(resultsArray[entry.Name[len(entry.Name)-1:]], entry)
    }

    defer s.client.Conn.Close()

    return resultsArray, err
}

func (s *SNMP) walkf(pdu g.SnmpPDU) error {
    fmt.Printf("%s = ", pdu.Name)

    switch pdu.Type {
    case g.OctetString:
        b := pdu.Value.([]byte)
        fmt.Printf("STRING: %s\n", string(b))
    default:
        fmt.Printf("TYPE %d: %d\n", pdu.Type, g.ToBigInt(pdu.Value))
    }
    return nil
}

func parseVersion(v string) (g.SnmpVersion, error) {
    if v == "v2c" {
        return g.Version2c, nil
    }

    if v == "1" {
        return g.Version1, nil
    }

    if v == "v3" {
        return g.Version3, nil
    }

    return g.Version2c, errors.New("Unknown SNMP version configured")
}

func (s *SNMP) ToInt(pdu g.SnmpPDU) int64 {
    bint := g.ToBigInt(pdu.Value)
    return bint.Int64()
}

func (s *SNMP) ToUint(pdu g.SnmpPDU) uint64 {
    bint := g.ToBigInt(pdu.Value)
    return bint.Uint64()
}

