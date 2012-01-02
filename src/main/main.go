/*
Copyright 2011 Paul Ruane.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"os"
)

var commands map[string]Command

func main() {
	commandArray := []Command{
		DeleteCommand{},
		DupesCommand{},
		ExportCommand{},
		FilesCommand{},
		HelpCommand{},
		MergeCommand{},
		MountCommand{},
		RenameCommand{},
		StatsCommand{},
		StatusCommand{},
		TagCommand{},
		TagsCommand{},
		UnmountCommand{},
		UntagCommand{},
		VersionCommand{},
		VfsCommand{},
	}

	commands = make(map[string]Command, len(commandArray))
	for _, command := range commandArray {
		commands[command.Name()] = command
	}

	var commandName string
	if len(os.Args) > 1 {
		commandName = os.Args[1]
	} else {
		commandName = "help"
	}

	command := commands[commandName]
	if command == nil {
		fmt.Printf("No such command, '%v'.\n", commandName)
		os.Exit(1)
	}

	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	} else {
		args = []string{}
	}

	err := command.Exec(args)
	if err != nil {
	    fmt.Fprintln(os.Stderr, err.Error())
	    os.Exit(1)
	}
}
