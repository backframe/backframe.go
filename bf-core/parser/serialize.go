package parser

var interfaceTmpl = `
interface {{name}} {
	isVersioned = {{Versioned}}

	if isVersioned {
		${{versions}}
	} else {
		${{resources}}
	}
}
`

var resourceTmpl = `
resource {{type}} {
	schema = {{schema}}

	${{methods}}
}
`
var methodTmpl = `
method {{name}} {
	secured = {{secured}}
	pubfields = {{pubfields}}
	subresources = {{subresources}}
}
`
