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
	"os"
	"testing"
	"time"
	"tmsu/cli"
	"tmsu/fingerprint"
	"tmsu/storage"
)

func TestDeleteSuccessful(test *testing.T) {
	// set-up

	databasePath := configureDatabase()
	defer os.Remove(databasePath)

	store, err := storage.Open()
	if err != nil {
		test.Fatal(err)
	}
	defer store.Close()

	fileD, err := store.AddFile("/tmp/d", fingerprint.Fingerprint("abc"), time.Now(), 123, false)
	if err != nil {
		test.Fatal(err)
	}

	fileF, err := store.AddFile("/tmp/f", fingerprint.Fingerprint("abc"), time.Now(), 123, false)
	if err != nil {
		test.Fatal(err)
	}

	fileB, err := store.AddFile("/tmp/b", fingerprint.Fingerprint("abc"), time.Now(), 123, false)
	if err != nil {
		test.Fatal(err)
	}

	tagDeathrow, err := store.AddTag("deathrow")
	if err != nil {
		test.Fatal(err)
	}

	tagFreeman, err := store.AddTag("freeman")
	if err != nil {
		test.Fatal(err)
	}

	if _, err := store.AddFileTag(fileD.Id, tagDeathrow.Id); err != nil {
		test.Fatal(err)
	}

	if _, err := store.AddFileTag(fileF.Id, tagFreeman.Id); err != nil {
		test.Fatal(err)
	}

	if _, err := store.AddFileTag(fileB.Id, tagDeathrow.Id); err != nil {
		test.Fatal(err)
	}
	if _, err := store.AddFileTag(fileB.Id, tagFreeman.Id); err != nil {
		test.Fatal(err)
	}

	command := DeleteCommand{false}

	// test

	if err := command.Exec(cli.Options{}, []string{"deathrow"}); err != nil {
		test.Fatal(err)
	}

	// validate

	tagDeathrow, err = store.TagByName("deathrow")
	if err != nil {
		test.Fatal(err)
	}
	if tagDeathrow != nil {
		test.Fatal("Deleted tag still exists.")
	}

	fileTagsD, err := store.FileTagsByFileId(fileD.Id)
	if err != nil {
		test.Fatal(err)
	}
	if len(fileTagsD) != 0 {
		test.Fatal("Expected no file-tags for file 'd'.")
	}

	fileTagsF, err := store.FileTagsByFileId(fileF.Id)
	if err != nil {
		test.Fatal(err)
	}
	if len(fileTagsF) != 1 {
		test.Fatal("Expected one file-tag for file 'f'.")
	}
	if fileTagsF[0].TagId != tagFreeman.Id {
		test.Fatal("Expected file-tag for tag 'freeman'.")
	}

	fileTagsB, err := store.FileTagsByFileId(fileB.Id)
	if err != nil {
		test.Fatal(err)
	}
	if len(fileTagsB) != 1 {
		test.Fatal("Expected one file-tag for file 'b'.")
	}
	if fileTagsB[0].TagId != tagFreeman.Id {
		test.Fatal("Expected file-tag for tag 'freeman'.")
	}
}

func TestDeleteNonExistentTag(test *testing.T) {
	// set-up

	databasePath := configureDatabase()
	defer os.Remove(databasePath)

	command := DeleteCommand{false}

	// test

	err := command.Exec(cli.Options{}, []string{"deleteme"})

	// validate

	if err == nil {
		test.Fatal("Non-existent from tag was not identified.")
	}
}
