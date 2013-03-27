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
	"time"
	"tmsu/cli"
	"tmsu/fingerprint"
	"tmsu/log"
	"tmsu/storage"
)

func TestTagsForSingleFile(test *testing.T) {
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

	file, err := store.AddFile("/tmp/tmsu/a", fingerprint.Fingerprint("123"), time.Now(), 0, false)
	if err != nil {
		test.Fatal(err)
	}

	appleTag, err := store.AddTag("apple")
	if err != nil {
		test.Fatal(err)
	}

	bananaTag, err := store.AddTag("banana")
	if err != nil {
		test.Fatal(err)
	}

	_, err = store.AddFileTag(file.Id, appleTag.Id)
	if err != nil {
		test.Fatal(err)
	}

	_, err = store.AddFileTag(file.Id, bananaTag.Id)
	if err != nil {
		test.Fatal(err)
	}

	tagsCommand := TagsCommand{false}

	// test

	if err := tagsCommand.Exec(cli.Options{}, []string{"/tmp/tmsu/a"}); err != nil {
		test.Fatal(err)
	}

	// verify

	log.Outfile.Seek(0, 0)

	bytes, err := ioutil.ReadAll(log.Outfile)
	compareOutput(test, "apple\nbanana\n", string(bytes))
}

func TestTagsForMultipleFiles(test *testing.T) {
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

	aFile, err := store.AddFile("/tmp/tmsu/a", fingerprint.Fingerprint("123"), time.Now(), 0, false)
	if err != nil {
		test.Fatal(err)
	}

	bFile, err := store.AddFile("/tmp/tmsu/b", fingerprint.Fingerprint("123"), time.Now(), 0, false)
	if err != nil {
		test.Fatal(err)
	}

	appleTag, err := store.AddTag("apple")
	if err != nil {
		test.Fatal(err)
	}

	bananaTag, err := store.AddTag("banana")
	if err != nil {
		test.Fatal(err)
	}

	_, err = store.AddFileTag(aFile.Id, appleTag.Id)
	if err != nil {
		test.Fatal(err)
	}

	_, err = store.AddFileTag(aFile.Id, bananaTag.Id)
	if err != nil {
		test.Fatal(err)
	}

	_, err = store.AddFileTag(bFile.Id, appleTag.Id)
	if err != nil {
		test.Fatal(err)
	}

	tagsCommand := TagsCommand{false}

	// test

	if err := tagsCommand.Exec(cli.Options{}, []string{"/tmp/tmsu/a", "/tmp/tmsu/b"}); err != nil {
		test.Fatal(err)
	}

	// verify

	log.Outfile.Seek(0, 0)

	bytes, err := ioutil.ReadAll(log.Outfile)
	compareOutput(test, "/tmp/tmsu/a: apple banana\n/tmp/tmsu/b: apple\n", string(bytes))
}

func TestAllTags(test *testing.T) {
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

	_, err = store.AddTag("apple")
	if err != nil {
		test.Fatal(err)
	}

	_, err = store.AddTag("banana")
	if err != nil {
		test.Fatal(err)
	}

	tagsCommand := TagsCommand{false}

	// test

	if err := tagsCommand.Exec(cli.Options{cli.Option{"--all", "-a", "", false, ""}}, []string{}); err != nil {
		test.Fatal(err)
	}

	// verify

	log.Outfile.Seek(0, 0)

	bytes, err := ioutil.ReadAll(log.Outfile)
	compareOutput(test, "apple\nbanana\n", string(bytes))
}
