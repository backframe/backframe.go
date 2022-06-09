package parser

/**
Sections
	Providers - Auth
		Provider
	Interfaces - APIS
		Interface [REST, GQL, GRPC]
	Integrations - Middleware & third party
		Integration
	Database - db config
		...config
*/

/*
	Specfile
		: VersionDeclaration
		: SectionBlock
		;
*/

type Spec struct {
	version  int
	sections map[SectionT]Section
}

func (s *Spec) SetVersion(val int) *Spec {
	s.version = val
	return s
}

func (s *Spec) AddSection(t SectionT) *Spec {
	s.sections[t] = NewSection(t)
	return s
}

func (s *Spec) RemoveSection(t SectionT) *Spec {
	delete(s.sections, t)
	return s
}

type SectionT byte

const (
	Providers SectionT = iota
	Interfaces
	Integrations
	Databases
)

type Section struct {
	sectionType SectionT
}

func NewSection(t SectionT) Section {
	return Section{
		sectionType: t,
	}
}

type ProviderT byte

const (
	Google ProviderT = iota
	Twitter
	Github
	EmailAndPassword
	PhoneNumber
	Facebook
)

type Provider struct {
	providerType ProviderT
}

func NewProvider(t ProviderT) Provider {
	return Provider{
		providerType: t,
	}
}

type InterfaceT byte

const (
	Rest InterfaceT = iota
	Graphql
	Grpc
)

type Interface struct {
	interfaceType InterfaceT
	routes        []string
	isVersioned   bool
	resources     []Resource
}

func NewInterface(t InterfaceT, versioned bool) Interface {
	return Interface{
		interfaceType: t,
		routes:        []string{},
		isVersioned:   versioned,
		resources:     []Resource{},
	}
}

func (i *Interface) AddRoute(val string) *Interface {
	i.routes = append(i.routes, val)
	return i
}

func (i *Interface) AddResource(name string) *Interface {
	i.resources = append(i.resources, NewResource(name))
	return i
}

type Resource struct {
	name        string
	schema_path string
	methods     []Method
}

func NewResource(name string) Resource {
	return Resource{
		name:    name,
		methods: []Method{},
	}
}

func (r *Resource) AddMethod(m MethodT, sec bool) *Resource {
	r.methods = append(r.methods, NewMethod(m, sec))
	return r
}

func (r *Resource) DefineSchemaPath(val string) *Resource {
	r.schema_path = val
	return r
}

type MethodT byte

const (
	GET MethodT = iota
	POST
	PUT
	DELETE
)

type Method struct {
	methodType   MethodT
	secured      bool
	subResources []Resource
	subFields    []string
}

func NewMethod(t MethodT, sec bool) Method {
	return Method{
		methodType:   t,
		secured:      sec,
		subResources: []Resource{},
		subFields:    []string{},
	}
}

type Database struct{}
type Integration struct{}
