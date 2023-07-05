package drive

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JPratama7/safe"
	"github.com/aiteung/gwrap"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func DriveNewServices(FileCred string, FileToken string) (driveService GoogleDrive) {
	cfgGoogle, err := gwrap.NewGoogleConfig(FileCred, drive.DriveScope, drive.DriveReadonlyScope)
	if err != nil {
		return
	}
	client := gwrap.GetClient(cfgGoogle, FileToken)
	curCtx := context.Background()
	srvDrive := safe.AsResult(drive.NewService(curCtx, option.WithHTTPClient(client)))
	if srvDrive.IsErr() {
		log.Print("Google Drive or Docs Service Unavailable")
	}
	driveService = NewGoogleDrive(srvDrive.Unwrap())
	return
}

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

func (gd *GoogleDrive) CreateDuplicate(fileId, filename, desc string, permission *drive.Permission) (fileDupId string, err error) {
	file, err := gd.srv.Files.Copy(fileId, &drive.File{Name: filename, Description: desc}).Do()

	if permission == nil {
		permission = &drive.Permission{
			Type: "anyone",
			Role: "writer",
		}
	}

	if err != nil {
		return
	}
	switch file.Shared {
	case true:
		fileDupId = file.DriveId
		return
	case false:
		fileDupId = file.Id
		_, err = gd.srv.Permissions.Create(fileDupId, permission).Do()
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

// https://stackoverflow.com/questions/70311191/access-to-the-provided-image-was-forbidden-even-though-it-was-uploaded-from-the
func (gd *GoogleDrive) GetURI(fileId string) (url string, err error) {
	res, err := gd.srv.Files.Get(fileId).Fields("id", "name", "webViewLink", "webContentLink", "thumbnailLink").Do()
	url = res.ThumbnailLink
	if url == "" {
		url = res.WebViewLink
	}
	return
}

func (gd *GoogleDrive) UploadFile(fileName, mimeType, filePath string, permission *drive.Permission) (fileId string, err error) {
	fileData := drive.File{
		Name:        fileName,
		CreatedTime: time.Now().Format(time.RFC3339),
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

	fileRes, err := gd.srv.Files.Create(&fileData).Media(open).Do()
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

	_, err = gd.srv.Permissions.Create(fileId, permission_).Do()
	return
}
func (gd *GoogleDrive) UploadFileReader(fileName, mimeType string, fileReader io.ReadSeeker, permission *drive.Permission) (fileId string, err error) {
	fileData := drive.File{
		Name:        fileName,
		CreatedTime: time.Now().Format(time.RFC3339),
		MimeType:    mimeType,
	}

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

	fileRes, err := gd.srv.Files.Create(&fileData).Media(fileReader).Do()
	if err != nil {
		return
	}

	fileId = fileRes.Id
	if fileId == "" {
		fileId = fileRes.DriveId
	}

	if permission == nil {
		permission = &drive.Permission{
			Type: "anyone",
			Role: "reader",
		}
	}

	_, err = gd.srv.Permissions.Create(fileId, permission).Do()
	return
}
