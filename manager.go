package sfm

import (
	"errors"
	"github.com/xbsoftware/wfs"
	"github.com/xbsoftware/wfs-local"
	"io"
	"sort"
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
		Include: func(name string) bool {
			for _, e := range ext {
				if !strings.ContainsAny(e, ".") {
					return true
				}
				if strings.HasSuffix(name, strings.ToLower(e)) {
					return true
				}
			}
			return false
		},
	}
	files, err := s.Drive.List(dir, &config)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(files, func(i, j int) bool {
		if files[i].Type == files[j].Type {
			return strings.Compare(files[i].Name, files[j].Name) == -1
		}
		// sort folder first
		var aOrder, bOrder int
		if files[i].Type == "folder" {
			aOrder = 0
		} else {
			aOrder = 1
		}
		if files[j].Type == "folder" {
			bOrder = 0
		} else {
			bOrder = 1
		}
		return aOrder < bOrder
	})
	return files, nil
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
