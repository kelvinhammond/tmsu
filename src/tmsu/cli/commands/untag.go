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
	"fmt"
	"path/filepath"
	"strings"
	"tmsu/cli"
	"tmsu/log"
	"tmsu/storage"
	"tmsu/storage/database"
)

type UntagCommand struct {
	verbose   bool
	recursive bool
}

func (UntagCommand) Name() cli.CommandName {
	return "untag"
}

func (UntagCommand) Synopsis() string {
	return "Remove tags from files"
}

func (UntagCommand) Description() string {
	return `tmsu untag FILE TAG...
tmsu untag --all FILE...
tmsu untag --tags "TAG..." FILE...

Disassociates FILE with the TAGs specified.`
}

func (UntagCommand) Options() cli.Options {
	return cli.Options{{"--all", "-a", "strip each file of all tags"},
		{"--tags", "-t", "the set of tags to remove"},
		{"--recursive", "-r", "recursively remove tags from directory contents"}}
}

func (command UntagCommand) Exec(options cli.Options, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no arguments specified.")
	}

	command.verbose = options.HasOption("--verbose")
	command.recursive = options.HasOption("--recursive")

	store, err := storage.Open()
	if err != nil {
		return fmt.Errorf("could not open storage: %v", err)
	}
	defer store.Close()

	if options.HasOption("--all") {
		if len(args) < 1 {
			return fmt.Errorf("files to untag must be specified.")
		}

		paths := args

		if err := command.untagPathsAll(store, paths); err != nil {
			return err
		}
	} else if options.HasOption("--tags") {
		if len(args) < 2 {
			return fmt.Errorf("quoted set of tags and at least one file to untag must be specified.")
		}

		tagNames := strings.Fields(args[0])
		paths := args[1:]

		tagIds, err := command.lookupTagIds(store, tagNames)
		if err != nil {
			return err
		}

		if err := command.untagPaths(store, paths, tagIds); err != nil {
			return err
		}
	} else {
		if len(args) < 2 {
			return fmt.Errorf("tag to remove and files to untag must be specified.")
		}

		path := args[0]
		tagNames := args[1:]

		tagIds, err := command.lookupTagIds(store, tagNames)
		if err != nil {
			return err
		}

		if err := command.untagPath(store, path, tagIds); err != nil {
			return err
		}
	}

	return nil
}

func (command UntagCommand) lookupTagIds(store *storage.Storage, names []string) ([]uint, error) {
	tags, err := store.TagsByNames(names)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve tags: %v", err)
	}

	for _, name := range names {
		if !tags.Any(func(tag *database.Tag) bool { return tag.Name == name }) {
			return nil, fmt.Errorf("no such tag '%v'", name)
		}
	}

	tagIds := make([]uint, len(tags))
	for index, tag := range tags {
		tagIds[index] = tag.Id
	}

	return tagIds, nil
}

func (command UntagCommand) untagPathsAll(store *storage.Storage, paths []string) error {
	for _, path := range paths {
		if err := command.untagPathAll(store, path); err != nil {
			return err
		}
	}

	return nil
}

func (command UntagCommand) untagPaths(store *storage.Storage, paths []string, tagIds []uint) error {
	for _, path := range paths {
		if err := command.untagPath(store, path, tagIds); err != nil {
			return err
		}
	}

	return nil
}

func (command UntagCommand) untagPathAll(store *storage.Storage, path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("%v: could not get absolute path: %v", path, err)
	}

	file, err := store.FileByPath(absPath)
	if err != nil {
		return fmt.Errorf("%v: could not retrieve file: %v", path, err)
	}
	if file == nil {
		return fmt.Errorf("%v: file is not tagged.", path)
	}

	if command.verbose {
		log.Infof("%v: removing all tags.", file.Path())
	}

	if err := store.RemoveFileTagsByFileId(file.Id); err != nil {
		return fmt.Errorf("%v: could not remove file's tags: %v", file.Path(), err)
	}

	if err := command.removeUntaggedFile(store, file); err != nil {
		return err
	}

	if command.recursive {
		childFiles, err := store.FilesByDirectory(file.Path())
		if err != nil {
			return fmt.Errorf("%v: could not retrieve files for directory: %v", file.Path())
		}

		for _, childFile := range childFiles {
			if err := store.RemoveFileTagsByFileId(childFile.Id); err != nil {
				return fmt.Errorf("%v: could not remove file's tags: %v", childFile.Path(), err)
			}

			if err := command.removeUntaggedFile(store, childFile); err != nil {
				return err
			}

		}
	}

	return nil
}

func (command UntagCommand) untagPath(store *storage.Storage, path string, tagIds []uint) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("%v: could not get absolute path: %v", path, err)
	}

	file, err := store.FileByPath(absPath)
	if err != nil {
		return fmt.Errorf("%v: could not retrieve file: %v", path, err)
	}
	if file == nil {
		return fmt.Errorf("%v: file is not tagged.", path)
	}

	for _, tagId := range tagIds {
		if command.verbose {
			log.Infof("%v: unapplying tag #%v.", file.Path(), tagId)
		}

		if err := store.RemoveFileTag(file.Id, tagId); err != nil {
			return fmt.Errorf("%v: could not remove tag #%v: %v", file.Path(), tagId, err)
		}
	}

	if err := command.removeUntaggedFile(store, file); err != nil {
		return err
	}

	if command.recursive {
		childFiles, err := store.FilesByDirectory(file.Path())
		if err != nil {
			return fmt.Errorf("%v: could not retrieve files for directory: %v", file.Path())
		}

		for _, childFile := range childFiles {
			for _, tagId := range tagIds {
				if command.verbose {
					log.Infof("%v: unapplying tag #%v.", childFile.Path(), tagId)
				}

				if err := store.RemoveFileTag(childFile.Id, tagId); err != nil {
					return fmt.Errorf("%v: could not remove tag #%v: %v", childFile.Path(), tagId, err)
				}
			}

			if err := command.removeUntaggedFile(store, childFile); err != nil {
				return err
			}
		}
	}

	return nil
}

func (command UntagCommand) removeUntaggedFile(store *storage.Storage, file *database.File) error {
	if command.verbose {
		log.Infof("%v: identifying whether file is tagged.", file.Path())
	}

	filetagCount, err := store.FileTagCountByFileId(file.Id)
	if err != nil {
		return fmt.Errorf("%v: could not get tag count: %v", file.Path(), err)
	}

	if filetagCount == 0 {
		if command.verbose {
			log.Infof("%v: removing untagged file.", file.Path())
		}

		err = store.RemoveFile(file.Id)
		if err != nil {
			return fmt.Errorf("%v: could not remove file: %v", file.Path(), err)
		}
	}

	return nil
}
