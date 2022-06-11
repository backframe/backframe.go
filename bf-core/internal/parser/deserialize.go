package parser

func Deserialize(stream []byte) []Block {
	parser := NewParser()
	return parser.Parse(string(stream))
}

func AddBlock(ast []Block) []Block {
	var newBlocks = []Block{}
	for _, a := range ast {
		// sections
		if a.id == "Providers" {
			b := Block{
				_type: "PROVIDER",
				id:    "Twitter",
			}
			a.blocks = append(a.blocks, b)
		}

		if a.id == "Interfaces" {
			b := Block{
				_type: "INTERFACE",
				id:    "Grpc",
			}
			a.blocks = append(a.blocks, b)
		}

		newBlocks = append(newBlocks, a)
	}
	return newBlocks
}
