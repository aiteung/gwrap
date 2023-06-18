package drive

import (
	"fmt"
	"google.golang.org/api/drive/v2"
	"io"
	"net/http"
	"os"
	"time"
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

func (gd *GoogleDrive) UploadFile(fileName, mimeType, filePath string, permission *drive.Permission) (fileId string, err error) {
	fileData := drive.File{
		Title:       fileName,
		CreatedDate: time.Now().Format(time.RFC3339),
	}

	open, err := os.Open(filePath)
	if err != nil {
		return
	}

	buffer := make([]byte, 512)
	_, err = open.Read(buffer)
	if err != nil {
		err = fmt.Errorf("could not read file: %v", err)
		return
	}
	_, err = open.Seek(0, io.SeekStart)
	if err != nil {
		err = fmt.Errorf("could not revert file offset: %v", err)
		return
	}

	switch mimeType {
	case "":
		fileData.MimeType = http.DetectContentType(buffer)
	default:
		fileData.MimeType = mimeType
	}

	fileRes, err := gd.srv.Files.Insert(&fileData).Media(open).Do()
	if err != nil {
		return
	}

	fileId = fileRes.Id
	if fileId == "" {
		fileId = fileRes.DriveId
	}

	permission_ := permission
	if permission == nil {
		permission_ = &drive.Permission{
			Type: "anyone",
			Role: "reader",
		}
	}

	_, err = gd.srv.Permissions.Insert(fileId, permission_).Do()
	return
}
func (gd *GoogleDrive) UploadFileReader(fileName, mimeType string, fileReader io.ReadSeeker, permission *drive.Permission) (fileId string, err error) {
	fileData := drive.File{
		Title:       fileName,
		CreatedDate: time.Now().Format(time.RFC3339),
	}
	fileData.MimeType = mimeType

	if mimeType == "" {
		buffer := make([]byte, 512)
		if _, er := fileReader.Read(buffer); er != nil {
			err = fmt.Errorf("could not read file: %v", er)
			return
		}
		if _, er := fileReader.Seek(0, io.SeekStart); er != nil {
			err = fmt.Errorf("could not revert file offset: %v", er)
			return
		}
		fileData.MimeType = http.DetectContentType(buffer)
	}

	fileRes, err := gd.srv.Files.Insert(&fileData).Media(fileReader).Do()
	if err != nil {
		return
	}

	fileId = fileRes.Id
	if fileId == "" {
		fileId = fileRes.DriveId
	}

	permission_ := permission
	if permission == nil {
		permission_ = &drive.Permission{
			Type: "anyone",
			Role: "reader",
		}
	}

	_, err = gd.srv.Permissions.Insert(fileId, permission_).Do()
	return
}
