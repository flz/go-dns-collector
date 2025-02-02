package transformers

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/dmachard/go-dnscollector/dnsutils"
	"github.com/dmachard/go-logger"
)

func TestSuspicious_Json(t *testing.T) {
	// enable feature
	config := dnsutils.GetFakeConfigTransformers()

	// get fake
	dm := dnsutils.GetFakeDnsMessage()
	dm.Init()

	// init subproccesor
	suspicious := NewSuspiciousSubprocessor(config, logger.New(false), "test")
	suspicious.InitDnsMessage(&dm)

	// expected json
	refJson := `
			{
				"suspicious": {
					"score":0,
					"malformed-pkt":false,
					"large-pkt":false,
					"long-domain":false,
					"slow-domain":false,
					"unallowed-chars":false,
					"uncommon-qtypes":false,
					"excessive-number-labels":false
				}
			}
			`

	var dmMap map[string]interface{}
	err := json.Unmarshal([]byte(dm.ToJson()), &dmMap)
	if err != nil {
		t.Fatalf("could not unmarshal dm json: %s\n", err)
	}

	var refMap map[string]interface{}
	err = json.Unmarshal([]byte(refJson), &refMap)
	if err != nil {
		t.Fatalf("could not unmarshal ref json: %s\n", err)
	}

	if _, ok := dmMap["suspicious"]; !ok {
		t.Fatalf("transformer key is missing")
	}

	if !reflect.DeepEqual(dmMap["suspicious"], refMap["suspicious"]) {
		t.Errorf("json format different from reference")
	}
}

func TestSuspicious_MalformedPacket(t *testing.T) {
	// config
	config := dnsutils.GetFakeConfigTransformers()
	config.Suspicious.Enable = true

	// init subproccesor
	suspicious := NewSuspiciousSubprocessor(config, logger.New(false), "test")

	// malformed DNS message
	dm := dnsutils.GetFakeDnsMessage()
	dm.DNS.MalformedPacket = true

	// init dns message with additional part
	suspicious.InitDnsMessage(&dm)

	suspicious.CheckIfSuspicious(&dm)

	if dm.Suspicious.Score != 1.0 {
		t.Errorf("suspicious score should be equal to 1.0")
	}

	if dm.Suspicious.MalformedPacket != true {
		t.Errorf("suspicious malformed packet flag should be equal to true")
	}
}

func TestSuspicious_LongDomain(t *testing.T) {
	// config
	config := dnsutils.GetFakeConfigTransformers()
	config.Suspicious.Enable = true
	config.Suspicious.ThresholdQnameLen = 4

	// init subproccesor
	suspicious := NewSuspiciousSubprocessor(config, logger.New(false), "test")

	// malformed DNS message
	dm := dnsutils.GetFakeDnsMessage()
	dm.DNS.Qname = "longdomain.com"

	// init dns message with additional part
	suspicious.InitDnsMessage(&dm)

	suspicious.CheckIfSuspicious(&dm)

	if dm.Suspicious.Score != 1.0 {
		t.Errorf("suspicious score should be equal to 1.0")
	}

	if dm.Suspicious.LongDomain != true {
		t.Errorf("suspicious long domain flag should be equal to true")
	}
}

func TestSuspiciousLargePacket(t *testing.T) {
	// config
	config := dnsutils.GetFakeConfigTransformers()
	config.Suspicious.Enable = true
	config.Suspicious.ThresholdPacketLen = 4

	// init subproccesor
	suspicious := NewSuspiciousSubprocessor(config, logger.New(false), "test")

	// malformed DNS message
	dm := dnsutils.GetFakeDnsMessage()
	dm.DNS.Length = 50

	// init dns message with additional part
	suspicious.InitDnsMessage(&dm)

	suspicious.CheckIfSuspicious(&dm)

	if dm.Suspicious.Score != 1.0 {
		t.Errorf("suspicious score should be equal to 1.0")
	}

	if dm.Suspicious.LargePacket != true {
		t.Errorf("suspicious large packet flag should be equal to true")
	}
}

func TestSuspiciousUncommonQtype(t *testing.T) {
	// config
	config := dnsutils.GetFakeConfigTransformers()
	config.Suspicious.Enable = true

	// init subproccesor
	suspicious := NewSuspiciousSubprocessor(config, logger.New(false), "test")

	// malformed DNS message
	dm := dnsutils.GetFakeDnsMessage()
	dm.DNS.Qtype = "LOC"

	// init dns message with additional part
	suspicious.InitDnsMessage(&dm)

	suspicious.CheckIfSuspicious(&dm)

	if dm.Suspicious.Score != 1.0 {
		t.Errorf("suspicious score should be equal to 1.0")
	}

	if dm.Suspicious.UncommonQtypes != true {
		t.Errorf("suspicious uncommon qtype flag should be equal to true")
	}
}

func TestSuspiciousExceedMaxLabels(t *testing.T) {
	// config
	config := dnsutils.GetFakeConfigTransformers()
	config.Suspicious.Enable = true
	config.Suspicious.ThresholdMaxLabels = 2

	// init subproccesor
	suspicious := NewSuspiciousSubprocessor(config, logger.New(false), "test")

	// malformed DNS message
	dm := dnsutils.GetFakeDnsMessage()
	dm.DNS.Qname = "test.sub.dnscollector.com"

	// init dns message with additional part
	suspicious.InitDnsMessage(&dm)

	suspicious.CheckIfSuspicious(&dm)

	if dm.Suspicious.Score != 1.0 {
		t.Errorf("suspicious score should be equal to 1.0")
	}

	if dm.Suspicious.ExcessiveNumberLabels != true {
		t.Errorf("suspicious excessive number labels flag should be equal to true")
	}
}

func TestSuspiciousUnallowedChars(t *testing.T) {
	// config
	config := dnsutils.GetFakeConfigTransformers()
	config.Suspicious.Enable = true

	// init subproccesor
	suspicious := NewSuspiciousSubprocessor(config, logger.New(false), "test")

	// malformed DNS message
	dm := dnsutils.GetFakeDnsMessage()
	dm.DNS.Qname = "AAAAAA==.dnscollector.com"

	// init dns message with additional part
	suspicious.InitDnsMessage(&dm)

	suspicious.CheckIfSuspicious(&dm)

	if dm.Suspicious.Score != 1.0 {
		t.Errorf("suspicious score should be equal to 1.0")
	}

	if dm.Suspicious.UnallowedChars != true {
		t.Errorf("suspicious unallowed chars flag should be equal to true")
	}
}
