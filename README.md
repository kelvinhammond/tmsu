Overview
========

TMSU is an application that allows you to organise your files by associating
them with tags. It provides a tool for managing these tags and a virtual
file-system to allow tag based access to your files.

TMSU's virtual file system does not store your files: it merely provides an
alternative, tagged based view of your files stored elsewhere in the file-
system. That way you have the freedom to choose the most suitable file system
for storage whilst still benefiting from tag based access.

Usage
=====

A command overview and details on how to use each command are available via the
integrated help:

    $ tmsu help

Full documentation is maintained online on the wiki:

  * <http://bitbucket.org/oniony/tmsu/wiki>

Downloading
===========

Binary builds for a limited number of architectures and operating system
combinations are available:

  * <http://bitbucket.org/oniony/tmsu/downloads>

You will need to ensure that both FUSE and Sqlite3 are installed for the
program to function.

Compiling
=========

The following steps are for compiling from source.

1. Installing Go

    TMSU is written in the Go programming language. To compile from source you must
    first install Go and the packages that TMSU depends upon. You can get that from
    the Go website:

    * <http://www.golang.org/>

    Go can be installed per the instructions on the Go website or it may be
    available in the package management system that comes with your operating
    system.

    TMSU is currently built against the Go weekly builds. See VERSIONS for the
    latest version that is known to work with the dependent packages.

2. Install the dependent packages.

        $ go get github.com/mattn/go-sqlite3
        $ go get github.com/hanwen/go-fuse/fuse

3. Clone the TMSU respository:

        $ hg clone https://bitbucket.org/oniony/tmsu

4. Make the project

        $ cd tmsu
        $ make

    This will compile to 'bin/tmsu' within the working directory.

5. Install the project

        $ sudo make install

    This will install TMSU to '/usr/bin/tmsu'.

    It will also install the Zsh completion to '/usr/share/zsh/site-functions'.

    To change the paths used override the environment variables in the Makefile.

About
=====

  * Website: <http://www.tmsu.org/>
  * Project: <http://bitbucket.org/oniony/tmsu>
  * Wiki: <http://bitbucket.org/oniony/tmsu/wiki>
  * Mailing list: <http://groups.google.com/group/tmsu>

TMSU is written in Go.

  * <http://www.golang.org/>

TMSU itself is written and maintained by Paul Ruane <paul@tmsu.org>, however
much of the functionality it provides is made possible by the Fuse and Sqlite3
libraries, their Go bindings and, of course, the Go language standard library.

Release Notes
=============

tip
---

This version changes the behaviour of the 'files' command to not automatic-
ally discover untagged files in the file-system that inherit tags. Instead the
--recursive option can be used when this behaviour is desirable:

For additional flexibility, the 'tag' and 'untag' command now also have a
--recursive option allowing a directory's contents to be likewise tagged or
untagged.

These two changes allow TMSU to be used in two different ways: dynamically
discovering directory contents using 'files --recursive' or statically
adding directory contents to the database for faster retrieval.

IMPORTANT: Please back up your database then upgrade it using the upgrade
script. The 'repair' step may take a while to run as every file is reexamined
to populate the new columns.

IMPORTANT: If you have been following 'tip' there is a separate upgrade
script called `tip_to_0.1.0.sql`. Due to a bug in the previous version
please re-run this even if you have previously upgraded otherwise you may end
up with duplicate entries in the `file_tag` table.

    $ cp ~/.tmsu/default.db ~/.tmsu/default.db.bak  # back up
    $ sqlite3 -init misc/db-upgrades/0.0.9_to_0.1.0.sql ~/.tmsu/default.db
    $ tmsu repair

  * 'files' command no longer finds files that inherit a tag from the file-
    system on-the-fly unless run with the --recursive option.
  * --recursive (-r) option on 'tag' and 'untag' for tagging/untagging
    directory contents recursively.
  * 'files' command now has --directory (-d) and --file (-f) options to limit
    output to just files or directories.
  * 'files' command now has --branch (-b) and --leaf (-l) options to omit items
    within matching directories and to omit parent directories of matching items
    respectively.
  * Improved command-line parsing: now supports global options, short options
    and mixed option ordering.
  * 'repair' command rewritten to fix bugs.
  * Fingerprints for directories no longer calculated. (Recursively tag
    files instead to detect file duplicates.)
  * Removed the 'export' command. (Sqlite tooling has better facilities.)
  * Tags containing '/' are no longer legal.
  * File lists are now shown sorted alphanumerically.
  * The 'tag' command no longer identifies modified files. (Use 'repair'
    instead.)
  * The 'mount' command now has a '--allow-other' option which allows other
    users to access the mounted file-system.
  * Updated Zsh completion.
  * Improved error messages.
  * Improved unit-test coverage.
  * Minor bug fixes.

v0.0.9
------

  * Fixed bug which caused process hosting the virtual file-system to crash if
    a non-existant tag directory is 'stat'ed.
  * Untagged files now inherit parent directory tags.

v0.0.8
------

Files can now be tagged within tagged directories. Files within tagged
directories will inherit the directory's tags.

  * Fixed bug with 'untag' command when non-existant tag is specified.
  * Updated with respect to go-fuse API change. 
  * 'mount' command now lists mount points if invoked without arguments.
  * Improved 'mount' command help.
  * 'rename' command now validates destination tag name.
  * Fixed bug with 'unmount --all' returning an error if there are no mounts.
  * Removed dependency upon 'pgrep'; now accesses proc-fs directotly for mount
    information.
  * 'files' command will now show files that inherit the specified tags
    ('--explicit' option turns this off).
  * 'tags' command will now shows inherited tags ('--explicit' turns this off).
  * 'status' command now reports inherited tags.
  * 'stats' command formatting updated.
  * Other minor bug fixes.

v0.0.7
------

  Files larger than 5MB now use a different fingerprinting algorithm where the
  fingerprint is produced by taking a 500KB of the start, middle and end of the
  file. This should dramatically improve performance, especially on slow file-
  systems.

  NOTE: it is advisable to run 'repair' after upgrading to this version to update
        large files with fingerprints produced using the new algorithm.

v0.0.6
------

IMPORTANT: This version adds a column to one of the database tables to record
the file's modification timestamep. The new code will not work with an
existing TMSU database until it has been upgraded.

To upgrade the database, run the upgrade script using the Sqlite3 tooling:

    $ sqlite3 -init misc/db-upgrade/0.0.5_to_0.0.6.sql

It is also advisable to run 'repair' after upgrading which will populate the
new column and also fix the directory fingerprints which have a new algorithm.

  * Upgraded to Go 1.
  * Added 'repair' command to fix up database when files are moved or modified.
  * Added 'version' command.
  * Added 'copy' command (duplicates a tag).
  * Added modification timestamp which is used in preference to fingerprint.
  * Command output now includes 'tmsu' to make it clear where output is coming
    from when piping.
  * Relative paths now calculated more accurately.
  * Zsh completion now supports tags containing colons.
  * 'status' command performance and functionality improvements.
  * Added directory fingerprinting.
  * Added 'repair' command to Zsh completion.
  * The 'files' command now allows tags to be excluded by prefixing them with a
    minus, e.g. -jazz.

- - -

Copyright 2011 Paul Ruane

Copying and distribution of this file, with or without modification,
are permitted in any medium without royalty provided the copyright
notice and this notice are preserved.  This file is offered as-is,
without any warranty.
