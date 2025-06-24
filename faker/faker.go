package faker

import (
	oneFake "github.com/manveru/faker"
	anotherFake "github.com/pioz/faker"
)

const (
	TypeFirstName   = "first_name"
	TypeLastName    = "last_name"
	TypeName        = "name"
	TypePhone       = "phone"
	TypeEmail       = "email"
	TypeCompanyName = "company"
	TypeAddress     = "address"
	TypeStreet      = "street_address"
	TypeCity        = "city"
	TypeZipCode     = "zip_code"
	TypeIPv4        = "ipv4"
	TypeURL         = "url"
	TypeLorem       = "lorem"
	TypeFixed       = "fixed"
	TypeString      = "string"
	TypeNull        = "null"
)

var (
	externalFakeGenerator *oneFake.Faker
)

func init() {
	externalFakeGenerator, _ = oneFake.New("en")
}

type FakeGenerator interface {
	GetData() interface{}
}

type FakeFirstName struct{}
type FakeLastName struct{}
type FakeName struct{}
type FakePhone struct{}
type FakeEmail struct{}
type FakeCompanyName struct{}
type FakeAddress struct{}
type FakeStreetAddress struct{}
type FakeCity struct{}
type FakeZipCode struct{}
type FakeIPv4 struct{}
type FakeURL struct{}
type FakeLorem struct{}

type FakeFixed struct {
	Value string
}

type FakeString struct {
	Length int
}

type FakeNull struct{}

func New(fakeConfig map[string]interface{}) FakeGenerator {
	if fakeType, ok := fakeConfig["type"]; ok {

		switch fakeType {
		case TypeFirstName:
			return &FakeFirstName{}
		case TypeLastName:
			return &FakeLastName{}
		case TypeName:
			return &FakeName{}

		case TypePhone:
			return &FakePhone{}
		case TypeEmail:
			return &FakeEmail{}

		case TypeCompanyName:
			return &FakeCompanyName{}
		case TypeAddress:
			return &FakeAddress{}
		case TypeStreet:
			return &FakeStreetAddress{}
		case TypeCity:
			return &FakeCity{}
		case TypeZipCode:
			return &FakeZipCode{}

		case TypeIPv4:
			return &FakeIPv4{}
		case TypeURL:
			return &FakeURL{}
		case TypeLorem:
			return &FakeLorem{}

		case TypeFixed:
			var value string
			if fixedValue, ok := fakeConfig["string"]; ok {
				value = fixedValue.(string)
			}
			return &FakeFixed{
				Value: value,
			}
		case TypeString:
			var length int
			if stringLength, ok := fakeConfig["length"]; ok {
				length = stringLength.(int)
			}
			return &FakeString{
				Length: length,
			}
		case TypeNull:
			return &FakeNull{}
		}
	}
	return nil
}

func (ff *FakeFixed) GetData() interface{} {
	return ff.Value
}

func (fs *FakeString) GetData() interface{} {
	return anotherFake.StringWithSize(fs.Length)
}

func (ff *FakeFirstName) GetData() interface{} {
	return anotherFake.FirstName()
}

func (ff *FakeLastName) GetData() interface{} {
	return anotherFake.LastName()
}

func (ff *FakeName) GetData() interface{} {
	return anotherFake.FullName()
}

func (ff *FakePhone) GetData() interface{} {
	return externalFakeGenerator.PhoneNumber()
}

func (ff *FakeEmail) GetData() interface{} {
	// Add more dispersion as email should be unique in some cases
	return anotherFake.StringWithSize(5) + externalFakeGenerator.Email()
}

func (ff *FakeCompanyName) GetData() interface{} {
	return externalFakeGenerator.CompanyName()
}

func (ff *FakeAddress) GetData() interface{} {
	return anotherFake.AddressFull()
}

func (ff *FakeStreetAddress) GetData() interface{} {
	return anotherFake.AddressSecondaryAddress()
}

func (ff *FakeCity) GetData() interface{} {
	return anotherFake.AddressCity()
}

func (ff *FakeZipCode) GetData() interface{} {
	return anotherFake.AddressZip()
}

func (ff *FakeIPv4) GetData() interface{} {
	return externalFakeGenerator.IPv4Address()
}

func (ff *FakeURL) GetData() interface{} {
	return externalFakeGenerator.URL()
}

func (ff *FakeLorem) GetData() interface{} {
	return anotherFake.Sentence()
}

func (fn *FakeNull) GetData() interface{} {
	return nil
}
