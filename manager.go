package sfm

import (
	"errors"
	"github.com/xbsoftware/wfs"
	"github.com/xbsoftware/wfs-local"
	"io"
	"slices"
	"strings"
)

type SiteFileManager struct {
	RootPath string
	Drive    wfs.Drive
}

func NewSiteFileManager(RootPath string) (*SiteFileManager, error) {
	policy := wfs.Policy(&wfs.AllowPolicy{})
	driveConfig := &wfs.DriveConfig{
		Verbose: false,
		Policy:  &policy,
	}
	drive, err := local.NewLocalDrive(RootPath, driveConfig)
	if err != nil {
		return nil, err
	}
	return &SiteFileManager{
		RootPath: RootPath,
		Drive:    drive,
	}, nil
}

func (s *SiteFileManager) Ls(dir string, ext ...string) ([]wfs.File, error) {
	config := wfs.ListConfig{
		SkipFiles:  false,
		SubFolders: false,
		Nested:     false,
		Exclude:    func(name string) bool { return strings.HasPrefix(name, ".") },
	}
	files, err := s.Drive.List(dir, &config)
	if err != nil {
		return nil, err
	}
	files = slices.DeleteFunc(files, func(f wfs.File) bool {
		if f.Type == "folder" {
			return false
		}
		for _, e := range ext {
			if strings.HasSuffix(f.Name, strings.ToLower(e)) {
				return false
			}
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	if !slices.IsSortedFunc(files, cmpFiles) {
		slices.SortFunc(files, cmpFiles)
	}
	return files, nil
}

func cmpFiles(f1, f2 wfs.File) int {
	if f1.Type == f2.Type {
		return strings.Compare(f1.Name, f2.Name)
	} else if f1.Type == "folder" {
		return -1
	} else if f2.Type == "folder" {
		return 1
	}
	return 0
}

func (s *SiteFileManager) Create(folder string, name string, file io.Reader) (*wfs.File, error) {
	fileId, err := s.Drive.Make(folder, name, false)
	if err != nil {
		return nil, err
	}

	err = s.Drive.Write(fileId, file)
	if err != nil {
		_ = s.Drive.Remove(fileId)
		return nil, err
	}
	info, err := s.Drive.Info(fileId)
	if err != nil {
		_ = s.Drive.Remove(fileId)
		return nil, err
	}
	return &info, nil
}

func (s *SiteFileManager) MkDir(folder string, name string) (*wfs.File, error) {
	fileId, err := s.Drive.Make(folder, name, true)
	if err != nil {
		return nil, err
	}
	info, err := s.Drive.Info(fileId)
	if err != nil {
		_ = s.Drive.Remove(fileId)
		return nil, err
	}
	return &info, nil
}

func (s *SiteFileManager) Delete(folder string, name string) error {
	config := wfs.ListConfig{
		SkipFiles:  false,
		SubFolders: false,
		Nested:     false,
		Exclude:    func(name string) bool { return strings.HasPrefix(name, ".") },
	}
	file, err := s.Drive.Search(folder, name, &config)
	if err != nil {
		return err
	}
	if len(file) == 0 {
		return errors.New("no such file")
	}
	err = s.Drive.Remove(file[0].ID + file[0].Name)
	if err != nil {
		return err
	}
	return nil
}
