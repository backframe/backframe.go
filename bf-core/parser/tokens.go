package parser

/***************************** Root specfile config *******************************/
type BfSpec struct {
	Sections []Section
}

// TODO: Specfile methods

/***************************** Section block methods *******************************/
type SectionType byte

type Section struct {
	_type SectionType
	Body  SectionBody
}

type SectionBody struct {
	Providers    []Provider
	Interfaces   []Interface
	Database     Database
	Integrations Integration
}

const (
	Providers SectionType = iota
	Interfaces
	Databases
	Integrations
)

func ConsumeSection(p *Parser) Section {
	p._eat("SECTION")
	id := p._eat("IDENTIFIER")
	p._eat("{")

	var _t SectionType
	switch id.value {
	case "Providers":
		_t = Providers
	case "Integrations":
		_t = Integrations
	case "Interfaces":
		_t = Interfaces
	case "Databases":
		_t = Databases
	}

	body := ConsumeSectionContents(p)
	p._eat("}")

	return Section{
		_type: _t,
		Body:  body,
	}
}

func ConsumeSectionContents(p *Parser) SectionBody {
	body := SectionBody{}

	return body
}

/***************************** Provider block *******************************/
type Provider struct {
	_type   string
	options ProviderOptions
}

type ProviderOptions struct{}

func ConsumeProviderBlock(p *Parser) []Provider {
	list := []Provider{}

	for {
		if p._lookahead.value == "}" {
			break
		}
		list = append(list, ConsumeProvider(p))
	}

	return list
}

func ConsumeProvider(p *Parser) Provider {
	p._eat("PROVIDER")
	id := p._eat("IDENTIFIER")
	p._eat("{")
	// TODO: Parse provider options
	p._eat("}")

	return Provider{
		_type: id.value,
	}
}

/***************************** Interface block *******************************/
type InterfaceType byte

type Interface struct {
	_type     InterfaceType
	Routes    []string
	Resources []Resource
	Endpoint  string
	Versioned bool
}

const (
	Rest InterfaceType = iota
	Graphql
	Grpc
)

type Resource struct {
	Name    string
	Schema  string
	Methods []Method
}

func ConsumeResourceBlock(p *Parser) []Resource {
	list := []Resource{}

	for {
		if p._lookahead.value == "}" {
			break
		}
		var r Resource
		p._eat("RESOURCE")
		id := p._eat("IDENTIFIER")
		p._eat("{")
		ConsumeResourceVariables(p, &r)
		methods := ConsumeMethodBlock(p)
		p._eat("}")
		r.Name = id.value
		r.Methods = methods
		list = append(list, r)
	}

	return list
}

func ConsumeResourceVariables(p *Parser, r *Resource) {
	for {
		isClosed := p._lookahead.value == "}" || p._lookahead._type == "METHOD"
		if isClosed {
			break
		}

		for {
			isClosed := p._lookahead.value == "}" || p._lookahead._type == "METHOD"
			if p._lookahead.value == ";" || isClosed {
				break
			}
			id := p._eat("IDENTIFIER")
			p._eat("ASSIGNMENT_OP")
			v := p._eat("STRING")

			if id.value == "schema" {
				r.Schema = v.value
			}
			p._eat(";")
		}

	}
}

type Method struct {
	_type        string
	Secured      bool
	SubResources []Resource
	PublicFields []string
}

func ConsumeMethodBlock(p *Parser) []Method {
	list := []Method{}

	for {
		if p._lookahead.value == "}" {
			break
		}
		m := ConsumeMethod(p)
		list = append(list, m)
	}

	return list
}

func ConsumeMethod(p *Parser) Method {
	p._eat("METHOD")
	id := p._eat("IDENTIFIER")
	p._eat("{")

	secured := false
	subResources := []Resource{}
	publicFields := []string{}

	for {
		if p._lookahead.value == "}" {
			break
		}

		for {
			if p._lookahead.value == ";" || p._lookahead.value == "}" {
				break
			}
			v := p._eat("IDENTIFIER")
			p._eat("ASSIGNMENT_OP")

			if v.value == "secured" {
				i := p._eat("BOOLEAN")
				if i.value == "true" {
					secured = true
				}
			}

			if v.value == "pubfields" {
				values := ConsumeArray(p)
				publicFields = append(publicFields, values...)
			}

			// Check at runtime whether subresource is valid
			if v.value == "subresources" {
				values := ConsumeArray(p)
				for _, v := range values {
					r := Resource{
						Name: v,
					}
					subResources = append(subResources, r)
				}
			}

			p._eat(";")
		}
	}

	p._eat("}")
	return Method{
		_type:        id.value,
		Secured:      secured,
		SubResources: subResources,
		PublicFields: publicFields,
	}
}

func ConsumeInterfaceBlock(p *Parser) []Interface {
	list := []Interface{}

	for {
		if p._lookahead.value == "}" {
			break
		}
		list = append(list, ConsumeInterface(p))
	}

	return list
}

func ConsumeInterfaceVariables(p *Parser, i *Interface) {
	for {
		if p._lookahead.value == "}" || p._lookahead._type == "VERSION" || p._lookahead._type == "RESOURCE" {
			break
		}

		for {
			isClosed := p._lookahead.value == "}" || p._lookahead._type == "VERSION" || p._lookahead._type == "RESOURCE"
			if isClosed || p._lookahead.value == ";" {
				id := p._eat("IDENTIFIER")
				p._eat("ASSIGNMENT_OP")

				if id.value == "isVersioned" {
					val := p._eat("BOOLEAN")
					if val.value == "true" {
						i.Versioned = true
					}
				}

				p._eat(";")
			}
		}
	}
}

func ConsumeInterface(p *Parser) Interface {
	p._eat("INTERFACE")
	id := p._eat("IDENTIFIER")
	p._eat("{")

	var api Interface
	var t InterfaceType

	switch id.value {
	case "Rest":
		t = Rest
	case "Graphql":
		t = Graphql
	case "Grpc":
		t = Grpc
	}

	api._type = t
	// Check for variables
	ConsumeInterfaceVariables(p, &api)

	if api.Versioned {
		// Consume Version blocks which then consume Resource blocks
	} else {
		api.Resources = ConsumeResourceBlock(p)
	}

	p._eat("}")
	return api

}

/***************************** Database Block *******************************/
type Database struct {
	Type             string
	DbUser           string
	DbPassword       string
	ServerURI        string
	ConnectionString string
}

/***************************** Integration Block *******************************/
type Integration struct{}

/***************************** Utility Methods *******************************/
func ConsumeArray(p *Parser) []string {
	list := []string{}

	p._eat("[")

	for {
		if p._lookahead.value == "]" {
			break
		}
		val := p._eat("STRING")
		p._eat("COMMA")
		list = append(list, val.value)
	}

	p._eat("]")
	return list
}
