/*
Copyright 2011-2013 Paul Ruane.

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

package commands

import (
	"io/ioutil"
	"os"
	"testing"
	"tmsu/cli"
	"tmsu/log"
	"tmsu/storage"
)

func TestRepairMovedFile(test *testing.T) {
	// set-up

	databasePath := configureDatabase()
	defer os.Remove(databasePath)

	store, err := storage.Open()
	if err != nil {
		test.Fatal(err)
	}
	defer store.Close()

	if err := createFile("/tmp/tmsu/a", "hello"); err != nil {
		test.Fatal(err)
	}
	defer os.Remove("/tmp/tmsu/a")

	tagCommand := TagCommand{false, false}
	if err := tagCommand.Exec(cli.Options{}, []string{"/tmp/tmsu/a", "a"}); err != nil {
		test.Fatal(err)
	}

	if err := os.Rename("/tmp/tmsu/a", "/tmp/tmsu/b"); err != nil {
		test.Fatal(err)
	}

	command := RepairCommand{false, false}

	// test

	if err := command.Exec(cli.Options{}, []string{"/tmp/tmsu"}); err != nil {
		test.Fatal(err)
	}

	// validate

	files, err := store.Files()
	if err != nil {
		test.Fatal(err)
	}

	if len(files) != 1 {
		test.Fatalf("Expected one file but are %v", len(files))
	}

	if files[0].Path() != "/tmp/tmsu/b" {
		test.Fatalf("File rename was not repaired.")
	}
}

func TestRepairModifiedFile(test *testing.T) {
	// set-up

	databasePath := configureDatabase()
	defer os.Remove(databasePath)

	store, err := storage.Open()
	if err != nil {
		test.Fatal(err)
	}
	defer store.Close()

	if err := createFile("/tmp/tmsu/a", "hello"); err != nil {
		test.Fatal(err)
	}
	defer os.Remove("/tmp/tmsu/a")

	tagCommand := TagCommand{false, false}
	if err := tagCommand.Exec(cli.Options{}, []string{"/tmp/tmsu/a", "a"}); err != nil {
		test.Fatal(err)
	}

	if err := createFile("/tmp/tmsu/a", "banana"); err != nil {
		test.Fatal(err)
	}

	command := RepairCommand{false, false}

	// test

	if err := command.Exec(cli.Options{}, []string{"/tmp/tmsu"}); err != nil {
		test.Fatal(err)
	}

	// validate

	files, err := store.Files()
	if err != nil {
		test.Fatal(err)
	}

	if len(files) != 1 {
		test.Fatalf("Expected one file but are %v", len(files))
	}

	if files[0].Size != 6 {
		test.Fatalf("File modification was not repaired.")
	}
}

func TestReportsMissingFiles(test *testing.T) {
	// set-up

	databasePath := configureDatabase()
	defer os.Remove(databasePath)

	outPath, errPath, err := configureOutput()
	if err != nil {
		test.Fatal(err)
	}
	defer os.Remove(outPath)
	defer os.Remove(errPath)

	store, err := storage.Open()
	if err != nil {
		test.Fatal(err)
	}
	defer store.Close()

	if err := createFile("/tmp/tmsu/a", "hello"); err != nil {
		test.Fatal(err)
	}

	tagCommand := TagCommand{false, false}
	if err := tagCommand.Exec(cli.Options{}, []string{"/tmp/tmsu/a", "a"}); err != nil {
		test.Fatal(err)
	}

	if err := os.Remove("/tmp/tmsu/a"); err != nil {
		test.Fatal(err)
	}

	command := RepairCommand{false, false}

	// test

	if err := command.Exec(cli.Options{}, []string{}); err != nil {
		test.Fatal(err)
	}

	// validate

	log.Outfile.Seek(0, 0)

	bytes, err := ioutil.ReadAll(log.Outfile)
	compareOutput(test, "tmsu: New tag 'a'.\ntmsu: /tmp/tmsu/a: missing\n", string(bytes))
}
