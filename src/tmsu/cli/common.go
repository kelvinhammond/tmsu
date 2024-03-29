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

package cli

import (
	"errors"
)

func ValidateTagNames(tagNames []string) error {
	for _, tagName := range tagNames {
		if err := ValidateTagName(tagName); err != nil {
			return err
		}
	}

	return nil
}

func ValidateTagName(tagName string) error {
	if tagName == "." || tagName == ".." {
		return errors.New("Tag name cannot be '.' or '..'.")
	}

	if tagName[0] == '-' {
		return errors.New("Tag names cannot start with '-'.")
	}

	for _, ch := range tagName {
		switch ch {
		case ',':
			return errors.New("tag names cannot contain ','.")
		case '=':
			return errors.New("tag names cannot contain '='.")
		case ' ':
			return errors.New("tag names cannot contain ' '.")
		case '/':
			return errors.New("tag names cannot contain '/'.")
		}
	}

	return nil
}
