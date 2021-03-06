package xsd

import (
	"fmt"

	"github.com/sezzle/sezzle-go-xml"
)

type InnerSchema struct {
	TargetNamespace    string        `xml:"targetNamespace,attr"`
	ElementFormDefault string        `xml:"elementFormDefault,attr"`
	Version            string        `xml:"version,attr"`
	ComplexTypes       []ComplexType `xml:"http://www.w3.org/2001/XMLSchema complexType"`
	SimpleTypes        []SimpleType  `xml:"http://www.w3.org/2001/XMLSchema simpleType"`
	Elements           []Element     `xml:"http://www.w3.org/2001/XMLSchema element"`
}

type Schema struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema schema"`
	Aliases map[string]string
	InnerSchema
}

func (s *Schema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	err := d.DecodeElement(&s.InnerSchema, &start)
	if err != nil {
		return err
	}

	s.XMLName = start.Name
	s.Aliases = map[string]string{}

	for _, attr := range start.Attr {
		s.Aliases[attr.Name.Local] = attr.Value
	}
	return nil
}

func (s *Schema) Namespace() string {
	return s.TargetNamespace
}

func (s *Schema) GetAlias(alias string) (space string) {
	return s.Aliases[alias]
}

// EncodeElement : Begins encoding to XML from the top level body element, calling Encode and EncodeType recursively on the
// nested elements until there are no more to be encoded.
func (s *Schema) EncodeElement(name string, enc *xml.Encoder, sr SchemaRepository, params map[string]interface{}, useNamespace, keepUsingNamespace bool, path ...string) error {
	// Starts encoding the top level xml element
	for _, elem := range s.Elements {
		if elem.Name == name {
			return elem.Encode(enc, sr, s, params, useNamespace, keepUsingNamespace, path...)
		}
	}

	return fmt.Errorf("did not find element '%s'", name)
}

func (s *Schema) EncodeType(name string, enc *xml.Encoder, sr SchemaRepository, params map[string]interface{}, useNamespace, keepUsingNamespace bool, path ...string) error {
	for _, cmplx := range s.ComplexTypes {
		if cmplx.Name == name {
			return cmplx.Encode(enc, sr, s, params, useNamespace, keepUsingNamespace, path...)
		}
	}

	for _, smpl := range s.SimpleTypes {
		if smpl.Name == name {
			return smpl.Encode(enc, sr, s, params, useNamespace, keepUsingNamespace, path...)
		}
	}

	return fmt.Errorf("did not find type '%s'", name)
}
