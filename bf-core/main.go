package main

import (
	"fmt"
	"io/fs"
	"os"

	"backframe.io/backframe/bf-core/internal/parser"
	"backframe.io/backframe/bf-core/internal/serde"
)

func main() {
	p := parser.Parser{}
	ast := p.Parse(`  
	
	# the providers section
	section Interfaces {
		interface Rest {
			isVersioned = true;
			
			version 1 {
				resource User {
					schema = "models/user.schema";
	
					method GET {
						# secured = true;
					
						pubfields = [
						"name", 
						"position",
						];
					}
				}
				
				resource Book {
					schema = "models/book.schema";
	
					method POST {}
				}
			}


		}

		interface Graphql {}
	}

	section Providers {
		provider Google {}

		provider Facebook {}
	}

	
	`)

	fmt.Printf("%v\n", ast)
	contents := serde.Serialize(ast)

	os.WriteFile("Specfile", contents, fs.FileMode(os.O_RDWR))
}
