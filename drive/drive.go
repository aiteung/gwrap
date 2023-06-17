package drive

import (
	"google.golang.org/api/drive/v2"
	"net/http"
)

type GoogleDrive struct {
	srv *drive.Service
}

func NewGoogleDrive(srvDrive *drive.Service) (res GoogleDrive) {
	res = GoogleDrive{srv: srvDrive}
	return
}

func (gd *GoogleDrive) Service() *drive.Service {
	return gd.srv
}

func (gd *GoogleDrive) CreateDuplicate(fileId, filename, desc string) (fileDupId string, err error) {
	file, err := gd.srv.Files.Copy(fileId, &drive.File{Title: filename, Description: desc}).Do()
	if err != nil {
		return
	}
	switch file.Shared {
	case true:
		fileDupId = file.DriveId
		return
	case false:
		fileDupId = file.Id
		return
	}
	return
}

func (gd *GoogleDrive) DeleteFiles(fileId ...string) (fileStatus map[string]bool, err error) {
	fileStat := make(map[string]bool, len(fileId))
	for _, v := range fileId {
		er := gd.srv.Files.Delete(v).Do()
		if er != nil {
			fileStat[v] = false
			continue
		}
		fileStat[v] = true
	}
	fileStatus = fileStat

	return
}

func (gd *GoogleDrive) DownloadFile(fileId, mimeType string) (res *http.Response, err error) {
	res, err = gd.srv.Files.Export(fileId, mimeType).Download()
	return
}

func (gd *GoogleDrive) GetURI(fileId string) (url string, err error) {
	res, err := gd.srv.Files.Get(fileId).Do()
	url = res.WebViewLink
	if url == "" {
		url = res.WebContentLink
	}
	if url == "" {
		url = res.AlternateLink
	}
	return
}

// TODO: UploadFile
//func (gd *GoogleDrive) UploadFile(fileId string) (err error) {
//	resx, err := gd.srv.Files.Create(fileId).Do()
//	fmt.Printf("URI: %+v\n", resx)
//	return
//}
