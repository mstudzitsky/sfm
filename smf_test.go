package sfm

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

func TestNewSiteFileManager(t *testing.T) {
	manager, err := NewSiteFileManager("/upload")
	if err != nil {
		t.Errorf("Error creating new site file manager: %v", err)
	}
	if manager.RootPath != "/upload" {
		t.Errorf("RootPath is %s, expecting %s", manager.RootPath, "/upload")
	}
	stats1, stats2, err := manager.Drive.Stats()
	if err != nil {
		t.Errorf("Eroor get drive stats: %v", err)
	}
	fmt.Printf("Stats: 1 - %d , 2- %d", stats1, stats2)
}

func TestSiteFileManager_Ls(t *testing.T) {
	var testFolderContents = []string{"id: /test/test1, name: test1, type: folder, size: 0", "id: /test/IMAG0003.jpg, name: IMAG0003.jpg, type: image, size: 2061645", "id: /test/IMAG0004.JPG, name: IMAG0004.JPG, type: image, size: 1943197", "id: /test/brands-logo-1.png, name: brands-logo-1.png, type: image, size: 41086"}
	manager, err := NewSiteFileManager("/upload")
	if err != nil {
		t.Errorf("Error creating new site file manager: %v", err)
	}
	files, err := manager.Ls("/test", "jpg", "png")
	fmt.Println("=================================================")
	for i, f := range files {
		check := fmt.Sprintf("id: %s, name: %s, type: %s, size: %d", f.ID, f.Name, f.Type, f.Size)
		if check != testFolderContents[i] {
			t.Errorf("Error in line %d: expecting: %s , got: %s", i, check, testFolderContents[i])
		}
		fmt.Println(check)
	}
	fmt.Println("=================================================")
}

func BenchmarkSiteFileManager_Ls(b *testing.B) {
	manager, err := NewSiteFileManager("/upload")
	if err != nil {
		b.Errorf("Error creating new site file manager: %v", err)
	}
	for i := 0; i < b.N; i++ {
		_, _ = manager.Ls("/test", "jpg", "png")
	}
}

func TestSiteFileManager_Mkdir(t *testing.T) {
	manager, err := NewSiteFileManager("/upload")
	if err != nil {
		t.Errorf("Error creating new site file manager: %v", err)
	}
	folder, err := manager.MkDir("/", "created1")
	if err != nil {
		t.Errorf("Cannot create folder %v", err)
		return
	}
	search, err := manager.Drive.Search("/", "created1", nil)
	if err != nil {
		t.Errorf("Error searching new folder %v", err)
		return
	}
	if len(search) > 0 {
		if search[0].Name != folder.Name {
			t.Errorf("Expected folder name: %s, got: %s", folder.Name, search[0].Name)
		}
		_ = manager.Drive.Remove(folder.ID)
	} else {
		t.Errorf("New folder %s not found", folder.Name)
	}
	check := fmt.Sprintf("id: %s, name: %s, type: %s, size: %d", folder.ID, folder.Name, folder.Type, folder.Size)
	fmt.Println(check)
}

func TestSiteFileManager_Create(t *testing.T) {
	manager, err := NewSiteFileManager("/upload")
	if err != nil {
		t.Errorf("Error creating new site file manager: %v", err)
		return
	}
	testFileInfo, err := manager.Drive.Search("/", "brands-logo-10.png")
	if err != nil {
		t.Errorf("Error searching test file %v", err)
		return
	}
	testFile, err := manager.Drive.Read(testFileInfo[0].ID + testFileInfo[0].Name)
	if err != nil {
		t.Errorf("Error reading test file: %s, %v", testFileInfo[0].ID, err)
		return
	}
	file, err := manager.Create("/", "tes_new_file.png", testFile)
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
		return
	}
	newFileInfo, err := manager.Drive.Info(file.ID)
	if err != nil {
		t.Errorf("Error get info of new file: %v", err)
		return
	}
	if file.ID != newFileInfo.ID || file.Type != newFileInfo.Type || file.Size != newFileInfo.Size || file.Name != newFileInfo.Name {
		fmt.Printf("expecting id: %s, type %s, size %d, name %s\n", file.ID, file.Type, file.Size, file.Name)
		t.Errorf("got id: %s, type: %s, size: %d, name %s", newFileInfo.ID, newFileInfo.Type, newFileInfo.Size, newFileInfo.Name)
		return
	}
	fmt.Printf("Created id: %s, type %s, size %d, name %s\n", newFileInfo.ID, newFileInfo.Type, newFileInfo.Size, newFileInfo.Name)
	_ = manager.Drive.Remove(newFileInfo.ID)
}

func TestSiteFileManager_Delete(t *testing.T) {
	manager, err := NewSiteFileManager("/upload")
	if err != nil {
		t.Errorf("Error creating new site file manager: %v", err)
		return
	}
	testFile, err := manager.Drive.Read("/brands-logo-10.png")
	file, err := manager.Create("/", "tes_new_file.png", testFile)
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
		return
	}
	err = manager.Delete("/", file.Name)
	if err != nil {
		t.Errorf("Error deleting test file: %v", err)
		return
	}
	testFileInfo, err := manager.Drive.Search("/", file.Name)
	if err != nil {
		t.Errorf("Error searching test file %v", err)
		return
	}
	if len(testFileInfo) > 0 {
		t.Errorf("Error test file is not detected %s", testFileInfo[0].ID)
	}
}

func TestSiteFileManager_GetImage(t *testing.T) {
	manager, err := NewSiteFileManager("/upload")
	if err != nil {
		t.Errorf("Error creating new site file manager: %v", err)
		return
	}
	img, format, err := manager.GetImage("/217.jpeg", 50, 50, true)
	if err != nil {
		t.Errorf("can not get image: %v", err)
		return
	}
	fmt.Printf("Image format: %s\n", format)
	out, _ := os.Create(manager.RootPath + "/test/encoded/" + "new_file." + format)
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			t.Errorf("can not close file: %v", err)
		}
		//err = os.Remove(manager.RootPath + "/test/encoded/" + "new_file." + format)
		//if err != nil {
		//	t.Errorf("can not delete file: %v", err)
		//}
	}(out)
	switch format {
	case "png", "webp":
		err = png.Encode(out, *img)
		break
	case "jpeg", "jpg":
		err = jpeg.Encode(out, *img, &jpeg.Options{Quality: 90})
		break
		//case "webp":
		//	err = webp.Encode(out, *img, &encoder.Options{Lossless: true})
		//	break
	}
	if err != nil {
		t.Errorf("Error encoding image: %v", err)
	}
}

func BenchmarkSiteFileManager_GetImage(b *testing.B) {
	manager, err := NewSiteFileManager("/upload")
	if err != nil {
		b.Errorf("Error creating new site file manager: %v", err)
		return
	}
	img, format, err := manager.GetImage("/217.jpeg", 20, 0, true)
	if err != nil {
		b.Errorf("can not get image: %v", err)
		return
	}
	fmt.Printf("Image format: %s\n", format)
	out, _ := os.Create(manager.RootPath + "/test/encoded/" + "new_file." + format)
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			b.Errorf("can not close file: %v", err)
		}
		err = os.Remove(manager.RootPath + "/test/encoded/" + "new_file." + format)
		if err != nil {
			b.Errorf("can not delete file: %v", err)
		}
	}(out)
	switch format {
	case "png", "webp":
		err = png.Encode(out, *img)
		break
	case "jpeg", "jpg":
		err = jpeg.Encode(out, *img, &jpeg.Options{Quality: 90})
		break
		//case "webp":
		//	err = webp.Encode(out, *img, &encoder.Options{Lossless: true})
		//	break
	}
	if err != nil {
		b.Errorf("Error encoding image: %v", err)
	}
}
