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

package vfs

import (
	"fmt"
	"github.com/hanwen/go-fuse/fuse"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"tmsu/log"
	"tmsu/storage"
	"tmsu/storage/database"
)

type FuseVfs struct {
	fuse.DefaultFileSystem

	store     *storage.Storage
	mountPath string
	state     *fuse.MountState
}

func MountVfs(databasePath string, mountPath string, allowOther bool) (*FuseVfs, error) {
	fuseVfs := FuseVfs{}
	pathNodeFs := fuse.NewPathNodeFs(&fuseVfs, nil)
	conn := fuse.NewFileSystemConnector(pathNodeFs, nil)
	state := fuse.NewMountState(conn)

	mountOptions := fuse.MountOptions{AllowOther: allowOther}
	err := state.Mount(mountPath, &mountOptions)
	if err != nil {
		return nil, fmt.Errorf("could not mount virtual filesystem at '%v': %v", mountPath, err)
	}

	store, err := storage.OpenAt(databasePath)
	if err != nil {
		return nil, fmt.Errorf("could not open database at '%v': %v", databasePath, err)
	}

	fuseVfs.store = store
	fuseVfs.mountPath = mountPath
	fuseVfs.state = state

	return &fuseVfs, nil
}

func (vfs FuseVfs) Unmount() {
	vfs.state.Unmount()
}

func (vfs FuseVfs) Loop() {
	vfs.state.Loop()
}

func (vfs FuseVfs) Access(name string, mode uint32, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN Access(%v, %v)", name, mode)
	defer log.Infof("END Access(%v, %v)", name, mode)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Chmod(name string, mode uint32, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN Chmod(%v, %v)", name, mode)
	defer log.Infof("BEGIN Chmod(%v, %v)", name, mode)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Chown(name string, uid uint32, gid uint32, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN Chown(%v, %v, %v)", name, uid, gid)
	defer log.Infof("BEGIN Chown(%v, %v)", name, uid, gid)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Create(name string, flags uint32, mode uint32, context *fuse.Context) (fuse.File, fuse.Status) {
	log.Infof("BEGIN Create(%v, %v, %v)", name, flags, mode)
	defer log.Infof("BEGIN Create(%v, %v)", name, flags, mode)

	return nil, fuse.ENOSYS
}

func (vfs FuseVfs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Infof("BEGIN GetAttr(%v)", name)
	defer log.Infof("END GetAttr(%v)", name)

	switch name {
	case "":
		fallthrough
	case "tags":
		return vfs.getTagsAttr()
	}

	path := vfs.splitPath(name)

	switch path[0] {
	case "tags":
		return vfs.getTaggedEntryAttr(path[1:])
	}

	return nil, fuse.ENOENT
}

func (vfs FuseVfs) GetXAttr(name string, attr string, context *fuse.Context) ([]byte, fuse.Status) {
	log.Infof("BEGIN GetXAttr(%v, %v)", name, attr)
	defer log.Infof("END GetAttr(%v, %v)", name, attr)

	return nil, fuse.ENOSYS
}

func (vfs FuseVfs) Link(oldName string, newName string, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN Link(%v, %v)", oldName, newName)
	defer log.Infof("END Link(%v, %v)", oldName, newName)

	return fuse.ENOSYS
}

func (vfs FuseVfs) ListXAttr(name string, context *fuse.Context) ([]string, fuse.Status) {
	log.Infof("BEGIN ListXAttr(%v)", name)
	defer log.Infof("END ListXAttr(%v)", name)

	return nil, fuse.ENOSYS
}

func (vfs FuseVfs) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN Mkdir(%v)", name)
	defer log.Infof("END Mkdir(%v)", name)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN Mknod(%v)", name)
	defer log.Infof("END Mknod(%v)", name)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Open(name string, flags uint32, context *fuse.Context) (fuse.File, fuse.Status) {
	log.Infof("BEGIN Open(%v)", name)
	defer log.Infof("END Open(%v)", name)

	return nil, fuse.ENOSYS
}

func (vfs FuseVfs) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	log.Infof("BEGIN OpenDir(%v)", name)
	defer log.Infof("END OpenDir(%v)", name)

	switch name {
	case "":
		return vfs.topDirectories()
	case "tags":
		return vfs.tagDirectories()
	}

	path := vfs.splitPath(name)
	switch path[0] {
	case "tags":
		return vfs.openTaggedEntryDir(path[1:])
	}

	return nil, fuse.ENOENT
}

func (vfs FuseVfs) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
	log.Infof("BEGIN Readlink(%v)", name)
	defer log.Infof("END Readlink(%v)", name)

	path := vfs.splitPath(name)
	switch path[0] {
	case "tags":
		return vfs.readTaggedEntryLink(path[1:])
	}

	return "", fuse.ENOENT
}

func (vfs FuseVfs) RemoveXAttr(name string, attr string, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN RemoveXAttr(%v, %v)", name, attr)
	defer log.Infof("END RemoveXAttr(%v, %v)", name, attr)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Rename(oldName string, newName string, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN Rename(%v, %v)", oldName, newName)
	defer log.Infof("END Rename(%v, %v)", oldName, newName)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Rmdir(name string, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN Rmdir(%v)", name)
	defer log.Infof("END Rmdir(%v)", name)

	return fuse.ENOSYS
}

func (vfs FuseVfs) SetXAttr(name string, attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN SetXAttr(%v, %v)", name, attr)
	defer log.Infof("END SetXAttr(%v, %v)", name, attr)

	return fuse.ENOSYS
}

func (vfs FuseVfs) StatFs(name string) *fuse.StatfsOut {
	log.Infof("BEGIN StatFs(%v)", name)
	defer log.Infof("END StatFs(%v)", name)

	return &fuse.StatfsOut{}
}

func (vfs FuseVfs) Symlink(value string, linkName string, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN Symlink(%v, %v)", value, linkName)
	defer log.Infof("END Symlink(%v, %v)", value, linkName)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Truncate(name string, offset uint64, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN Truncate(%v)", name)
	defer log.Infof("END Truncate(%v)", name)

	return fuse.ENOSYS
}

func (vfs FuseVfs) Unlink(name string, context *fuse.Context) fuse.Status {
	log.Infof("BEGIN Unlink(%v)", name)
	defer log.Infof("END Unlink(%v)", name)

	fileId, err := vfs.parseFileId(name)
	if err != nil {
		log.Fatal("Could not unlink: ", err)
	}

	if fileId == 0 {
		// cannot unlink tag directories
		return fuse.EPERM
	}

	path := vfs.splitPath(name)
	tagNames := path[1 : len(path)-1]

	for _, tagName := range tagNames {
		tag, err := vfs.store.Db.TagByName(tagName)
		if err != nil {
			log.Fatal(err)
		}
		if tag == nil {
			log.Fatalf("Could not retrieve tag '%v'.", tagName)
		}

		err = vfs.store.RemoveFileTag(fileId, tag.Id)
		if err != nil {
			log.Fatal(err)
		}
	}

	return fuse.OK
}

// non-exported

func (vfs FuseVfs) splitPath(path string) []string {
	return strings.Split(path, string(filepath.Separator))
}

func (vfs FuseVfs) parseFileId(name string) (uint, error) {
	parts := strings.Split(name, ".")
	count := len(parts)

	if count == 1 {
		return 0, nil
	}

	id, err := Atoui(parts[count-2])
	if err != nil {
		id, err = Atoui(parts[count-1])
		if err != nil {
			return 0, nil
		}
	}

	return id, nil
}

func (vfs FuseVfs) topDirectories() ([]fuse.DirEntry, fuse.Status) {
	log.Infof("BEGIN topDirectories")
	defer log.Infof("END topDirectories")

	entries := make([]fuse.DirEntry, 0, 1)
	entries = append(entries, fuse.DirEntry{Name: "tags", Mode: fuse.S_IFDIR})
	return entries, fuse.OK
}

func (vfs FuseVfs) tagDirectories() ([]fuse.DirEntry, fuse.Status) {
	log.Infof("BEGIN tagDirectories")
	defer log.Infof("END tagDirectories")

	tags, err := vfs.store.Db.Tags()
	if err != nil {
		log.Fatalf("Could not retrieve tags: %v", err)
	}

	entries := make([]fuse.DirEntry, len(tags))
	for index, tag := range tags {
		entries[index] = fuse.DirEntry{Name: tag.Name, Mode: fuse.S_IFDIR}
	}

	return entries, fuse.OK
}

func (vfs FuseVfs) getTagsAttr() (*fuse.Attr, fuse.Status) {
	log.Infof("BEGIN getTagsAttr")
	defer log.Infof("END getTagsAttr")

	tagCount, err := vfs.store.Db.TagCount()
	if err != nil {
		log.Fatalf("Could not get tag count: %v", err)
	}

	now := time.Now()
	return &fuse.Attr{Mode: fuse.S_IFDIR | 0755, Nlink: 2, Size: uint64(tagCount), Mtime: uint64(now.Unix()), Mtimensec: uint32(now.Nanosecond())}, fuse.OK
}

func (vfs FuseVfs) getTaggedEntryAttr(path []string) (*fuse.Attr, fuse.Status) {
	log.Infof("BEGIN getTaggedEntryAttr(%v)", path)
	defer log.Infof("END getTaggedEntryAttr(%v)", path)

	pathLength := len(path)
	name := path[pathLength-1]

	fileId, err := vfs.parseFileId(name)
	if err != nil {
		return nil, fuse.ENOENT
	}

	if fileId == 0 {
		// tag directory
		tagIds, err := vfs.tagNamesToIds(path)
		if err != nil {
			log.Fatalf("Could not lookup tag IDs: %v.", err)
		}
		if tagIds == nil {
			return nil, fuse.ENOENT
		}

		//TODO slow
		//		fileCount, err := vfs.store.FileCountWithTags(tagIds)
		//		if err != nil {
		//			log.Fatalf("Could not retrieve count of files with tags: %v. (%v)", path, err)
		//		}   
		fileCount := 0

		now := time.Now()
		return &fuse.Attr{Mode: fuse.S_IFDIR | 0755, Nlink: 2, Size: uint64(fileCount), Mtime: uint64(now.Unix()), Mtimensec: uint32(now.Nanosecond())}, fuse.OK
	}

	file, err := vfs.store.File(fileId)
	if err != nil {
		log.Fatalf("Could not retrieve file #%v: %v", fileId, err)
	}
	if file == nil {
		return &fuse.Attr{Mode: fuse.S_IFREG}, fuse.ENOENT
	}

	fileInfo, err := os.Stat(file.Path())
	var size int64
	var modTime time.Time
	if err == nil {
		size = fileInfo.Size()
		modTime = fileInfo.ModTime()
	} else {
		size = 0
		modTime = time.Time{}
	}

	return &fuse.Attr{Mode: fuse.S_IFLNK | 0755, Size: uint64(size), Mtime: uint64(modTime.Unix()), Mtimensec: uint32(modTime.Nanosecond())}, fuse.OK
}

func (vfs FuseVfs) openTaggedEntryDir(path []string) ([]fuse.DirEntry, fuse.Status) {
	log.Infof("BEGIN openTaggedEntryDir(%v)", path)
	defer log.Infof("END openTaggedEntryDir(%v)", path)

	tagIds, err := vfs.tagNamesToIds(path)
	if err != nil {
		log.Fatalf("Could not lookup tag IDs: %v.", err)
	}
	if tagIds == nil {
		return nil, fuse.ENOENT
	}

	furtherTagIds, err := vfs.store.TagsForTags(tagIds)
	if err != nil {
		log.Fatalf("Could not retrieve tags for tags: %v", err)
	}

	files, err := vfs.store.FilesWithTags(tagIds, []uint{})
	if err != nil {
		log.Fatalf("Could not retrieve tagged files: %v", err)
	}

	entries := make([]fuse.DirEntry, 0, len(files)+len(furtherTagIds))
	for _, tag := range furtherTagIds {
		entries = append(entries, fuse.DirEntry{Name: tag.Name, Mode: fuse.S_IFDIR | 0755})
	}
	for _, file := range files {
		linkName := vfs.getLinkName(file)
		entries = append(entries, fuse.DirEntry{Name: linkName, Mode: fuse.S_IFLNK})
	}

	return entries, fuse.OK
}

func (vfs FuseVfs) readTaggedEntryLink(path []string) (string, fuse.Status) {
	log.Infof("BEGIN readTaggedEntryLink(%v)", path)
	defer log.Infof("END readTaggedEntryLink(%v)", path)

	name := path[len(path)-1]

	fileId, err := vfs.parseFileId(name)
	if err != nil {
		log.Fatalf("Could not parse file identifier: %v", err)
	}
	if fileId == 0 {
		return "", fuse.ENOENT
	}

	file, err := vfs.store.File(fileId)
	if err != nil {
		log.Fatalf("Could not find file %v in database.", fileId)
	}

	return file.Path(), fuse.OK
}

func (vfs FuseVfs) getLinkName(file *database.File) string {
	extension := filepath.Ext(file.Path())
	fileName := filepath.Base(file.Path())
	linkName := fileName[0 : len(fileName)-len(extension)]
	suffix := "." + Uitoa(file.Id) + extension

	if len(linkName)+len(suffix) > 255 {
		linkName = linkName[0 : 255-len(suffix)]
	}

	return linkName + suffix
}

func (vfs FuseVfs) tagNamesToIds(tagNames []string) ([]uint, error) {
	tagIds := make([]uint, len(tagNames))

	for index, tagName := range tagNames {
		tag, err := vfs.store.Db.TagByName(tagName)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve tag '%v': %v", tagName, err)
		}
		if tag == nil {
			return nil, nil
		}

		tagIds[index] = tag.Id
	}

	return tagIds, nil
}

func Uitoa(ui uint) string {
	return strconv.FormatUint(uint64(ui), 10)
}

func Atoui(str string) (uint, error) {
	ui64, err := strconv.ParseUint(str, 10, 0)
	return uint(ui64), err
}
