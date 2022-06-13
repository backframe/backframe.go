package main

import (
	"backframe.io/backframe/bf-cli/cmd"
	"github.com/ndaba1/gommander"
)

func main() {
	app := gommander.App()

	app.Help("A framework for rapid development and deployment of APIs").
		Name("bf-go").
		Author("The backframe team").
		Version("0.1.0")

	app.SubCommand("new").
		Alias("n").
		Help("Creates a new backframe project in the specified directory").
		AddArgument(
			gommander.NewArgument("<app-name>").
				Help("The name of the new project").
				Default("server"),
		).
		Flag("-g --git", "Whether or not to initialize a git repository").
		Flag("-d --default", "Skip prompts and use default preset").
		Flag("-f --force", "Override target directory if it exists").
		Option("-p --preset <path>", "Pass the path to custom bfconfig.json").
		Action(cmd.New)

	app.SubCommand("serve").
		Alias("s").
		Help("Serve the project in the current working directory").
		AddOption(
			gommander.NewOption("port").
				Short('p').
				Help("The port to serve the app on").
				AddArgument(
					gommander.NewArgument("<port-no>").Default("9000"),
				),
		).
		Action(cmd.Serve)

	app.SubCommand("watch").
		Alias("w").
		Help("Start the app server in watch mode").
		AddOption(
			gommander.NewOption("port").
				Short('p').
				Help("The port to serve the app on").
				AddArgument(
					gommander.NewArgument("<port-no>").Default("9000"),
				),
		).
		AddOption(
			gommander.NewOption("exclude").
				Short('e').
				Help("Pass a list of files not to watch for changes").
				AddArgument(
					gommander.NewArgument("<files>").Variadic(true),
				),
		).
		Action(cmd.Watch)

	app.SubCommand("build").
		Alias("b").
		Help("Build a new server from a defined backframe config").
		Flag("-q --quiet", "Suppress output printing when building").
		AddOption(
			gommander.NewOption("config").
				Short('c').
				Help("The path to the backframe.json config file").
				Required(true).
				AddArgument(
					gommander.NewArgument("<path>").Default("./backframe.json"),
				),
		).
		Action(cmd.Build)

	app.SubCommand("generate").
		Alias("g").
		AddArgument(
			gommander.
				NewArgument("<api-type>").
				Help("The type of API to generate files for").
				ValidateWith([]string{"REST", "GRAPHQL", "GRPC"}).
				Default("REST"),
		).
		Help("Uses the cli to generate new api endpoints by prompting for values").
		Flag("-v --version", "Specify whether or not to version the endpoints").
		Action(cmd.Generate)

	app.SubCommand("add").
		Alias("a").
		Argument("<package-name>", "The name of the package to add").
		Help("Adds and invokes a plugin/middleware to a project").
		UsageStr("bf-go add <package-name> -- [pluginOptions]").
		Action(cmd.Add)

	app.SubCommand("deploy").
		Alias("d").
		Help("Deploy the backframe server on backframe.io").
		Action(cmd.Deploy)

	app.SubCommand("login").
		Help("Login to your backframe user account").
		Action(cmd.Login)

	rest := app.SubCommand("rest").Alias("r").Help("Manage Rest API functionality")

	rest.SubCommand("generate").
		Alias("g").
		Help("Generate a Rest API for current backframe project")

	rest.SubCommand("routes").Help("List all the rest routes of the project")

	db := app.SubCommand("db").Alias("database").Help("Manage database functionality")

	db.SubCommand("migrate").
		Help("Handle database migrations").
		AddArgument(
			gommander.
				NewArgument("<type>").
				Help("The type of the migration").
				ValidateWith([]string{"UP", "DOWN", "ROLLBACK"}),
		)

	db.SubCommand("create").
		Help("Create a new database migration").
		AddArgument(
			gommander.
				NewArgument("<name>").
				Help("The name of the new migration"),
		)

	app.Set(gommander.IncludeHelpSubcommand, true)

	app.Parse()
}
