package sfm

import (
	"github.com/xbsoftware/wfs"
	"github.com/xbsoftware/wfs-local"
	"strings"
)

type SiteFileManager struct {
	RootPath    string
	UploadLimit uint
	Drive       wfs.Drive
}

func NewSiteFileManager(RootPath string, UploadLimit uint) (*SiteFileManager, error) {
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
		RootPath:    RootPath,
		UploadLimit: UploadLimit,
		Drive:       drive,
	}, nil
}

func (s *SiteFileManager) ls(dir string, ext ...string) ([]wfs.File, error) {
	config := wfs.ListConfig{
		SkipFiles:  false,
		SubFolders: true,
		Nested:     true,
		Exclude:    func(name string) bool { return strings.HasPrefix(name, ".") },
		Include: func(name string) bool {
			for _, e := range ext {
				if strings.HasSuffix(name, e) {
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
	return files, nil
}
