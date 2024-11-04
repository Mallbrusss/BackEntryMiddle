package service

import (
	"fmt"
	"github.com/Mallbrusss/BackEntryMiddle/models"
	"github.com/Mallbrusss/BackEntryMiddle/internal/repository"
	"mime"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

type DocumentService struct {
	docRepo   *repository.DocumentRepository
	uploadDir string
}

func NewDocumentService(dockRepo *repository.DocumentRepository, uploadDir string) *DocumentService {
	return &DocumentService{
		docRepo:   dockRepo,
		uploadDir: uploadDir,
	}
}

func (ds *DocumentService) getFileExtensions(mimeType string) string {
	ext, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(ext) == 0 {
		return ".bin"
	}
	return ext[0]
}

func (ds *DocumentService) UploadDocument(document *models.Document, fileData []byte, grant []string) (*models.Document, error) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	errorCh := make(chan error, 2)
	// filePathCh := make(chan string, 1)

	ext := ds.getFileExtensions(document.Mime)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(ds.uploadDir, fileName)

	go func() {
		defer wg.Done()
		if err := os.WriteFile(filePath, fileData, 0644); err != nil {
			errorCh <- fmt.Errorf("error write file: %w", err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		document.ID = uuid.New().String()
		err := ds.docRepo.CreateDocument(document, grant)
		if err != nil {
			errorCh <- fmt.Errorf("error save to database: %w", err)
		}
	}()

	wg.Wait()
	close(errorCh)

	var fError error
	for err := range errorCh {
		if fError == nil {
			fError = err
		}
	}

	if fError != nil {
			removeErr := os.Remove(filePath)
			if removeErr != nil && !os.IsNotExist(removeErr) {
				return nil, fmt.Errorf("error: %v, error deleting file: %v", fError, removeErr)
			}
		return nil, fError
	}

	document.FilePath = filePath

	return document, nil
}

func (ds *DocumentService) DeleteDocument(documentID string) error {
	document, err := ds.docRepo.GetDocumentByID(documentID, "")
	if err != nil {
		return err
	}

	if err := os.Remove(document.FilePath); err != nil && !os.IsExist(err) {
		return err
	}

	return ds.docRepo.DeleteDocument(document)
}
