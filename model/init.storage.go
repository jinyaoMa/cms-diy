package model

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

type UTF16 *uint16
type Size uint64
type Branch struct {
	pointer   UTF16
	APath     string
	All       Size
	Used      Size
	Available Size
	Free      Size
}
type Branches []*Branch
type Storage struct {
	branches  Branches
	All       Size
	Used      Size
	Available Size
	Free      Size
}
type ScannedFileCallback func(apath string, rpath string, depth int, workspace string, fileInfo os.FileInfo)

const (
	B  Size = 1
	KB Size = 1024 * B
	MB Size = 1024 * KB
	GB Size = 1024 * MB
	TB Size = 1024 * GB
)

var (
	h   *windows.DLL
	cmd *windows.Proc

	storage Storage
)

func init_storage() {
	h = windows.MustLoadDLL("kernel32.dll")
	cmd = h.MustFindProc("GetDiskFreeSpaceExW")

	for _, apath := range StorageBranches {
		ptr, err := windows.UTF16PtrFromString(apath)
		if err != nil {
			panic("Failed to locate workspace")
		}
		storage.branches = append(storage.branches, &Branch{
			pointer: ptr,
			APath:   filepath.Join(apath, "."),
		})
	}

	storage.Update()
}

func (s *Storage) Update() {
	s.Free = 0
	s.All = 0
	s.Available = 0

	checkedVolumes := ""
	for _, branch := range s.branches {
		var free, all, available Size
		_, _, _ = cmd.Call(uintptr(unsafe.Pointer(branch.pointer)),
			uintptr(unsafe.Pointer(&free)),
			uintptr(unsafe.Pointer(&all)),
			uintptr(unsafe.Pointer(&available)))
		branch.Free = free
		branch.All = all
		branch.Available = available
		branch.Used = branch.All - branch.Free

		volume := filepath.VolumeName(branch.APath)
		if !strings.Contains(checkedVolumes, volume) {
			checkedVolumes += volume
			s.Free += branch.Free
			s.All += branch.All
			s.Available += branch.Available
		}
	}
	s.Used = s.All - s.Free

	if WorkspaceLimit < s.Available {
		s.Available = WorkspaceLimit
	}
}

func IsDirectoryNamedCache(str string) bool {
	return strings.Compare(str, DIRECTORY_CACHE_NAME) == 0
}

func IsDirectoryNamedRecycle(str string) bool {
	return strings.Compare(str, DIRECTORY_RECYCLE_NAME) == 0
}

func GenerateRecycleAPath(file File) string {
	return filepath.Join(file.Workspace, DIRECTORY_RECYCLE_NAME, parseNumberToString(file.ID)+".$"+file.Name)
}

func IsFileNameCharValid(str string) bool {
	return !strings.ContainsAny(str, "\\/:*?\"<>|")
}

func IsPathCharValid(str string) bool {
	matched, err := regexp.MatchString("\\.[\\/]", str)
	return !(err != nil || matched || strings.ContainsAny(str, ":*?\"<>|"))
}

func GetValidUserWorkspace(userAccount string, size Size) (validUserWorkspace string, ok bool) {
	storage.Update()
	if storage.Available >= size {
		for _, branch := range storage.branches {
			if branch.Available >= size {
				validUserWorkspace = filepath.Join(branch.APath, userAccount)
				ok = true
				return
			}
		}
	}
	validUserWorkspace = ""
	ok = false
	return
}

func NewUserSpace(userAccount string, fn ScannedFileCallback, isForcedScan bool) error {
	if !IsFileNameCharValid(userAccount) {
		return newError("User account format not allowed")
	}

	for _, branch := range storage.branches {
		targetPath := filepath.Join(branch.APath, userAccount)
		err := os.Mkdir(targetPath, os.ModeDir)
		if os.IsExist(err) {
			if isForcedScan && ScanUserFiles(targetPath, fn) != nil {
				return newError("User files initializing error")
			}
		} else if err != nil {
			return err
		}
		err = os.Mkdir(filepath.Join(targetPath, DIRECTORY_RECYCLE_NAME), os.ModeDir)
		if !(err == nil || os.IsExist(err)) {
			return newError("User '" + DIRECTORY_RECYCLE_NAME + "' initializing error")
		}
		err = os.Mkdir(filepath.Join(targetPath, DIRECTORY_CACHE_NAME), os.ModeDir)
		if !(err == nil || os.IsExist(err)) {
			return newError("User '" + DIRECTORY_CACHE_NAME + "' initializing error")
		}
	}

	return nil
}

func ScanUserFiles(userSpace string, fn ScannedFileCallback) error {
	return filepath.Walk(userSpace, func(apath string, fileInfo os.FileInfo, err error) error {
		if fileInfo == nil {
			return err
		}

		rpath, errRel := filepath.Rel(userSpace, apath)
		if errRel != nil {
			return errRel
		}
		if !(rpath == "." || IsDirectoryNamedRecycle(fileInfo.Name()) || IsDirectoryNamedCache(fileInfo.Name())) {
			depth := strings.Count(rpath, string(filepath.Separator))
			fn(apath, rpath, depth, userSpace, fileInfo)
		}

		return nil
	})
}
