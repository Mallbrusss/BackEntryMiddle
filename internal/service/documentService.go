package service

import (
	"fmt"
	"internal/models"
	"internal/repository"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

type DocumentService struct {
	docRepo   *repository.DocumentRepository
	uploadDir string
	wg        sync.WaitGroup
}

func NewDocumentService(dockRepo *repository.DocumentRepository, uploadDir string) *DocumentService {
	return &DocumentService{
		docRepo:   dockRepo,
		uploadDir: uploadDir,
		wg:        sync.WaitGroup{},
	}
}

// TODO: Мб подавать на вход структуру?? Подумать
// TODO: подумать, над fileName или как хранить одинаковые файлы
func (ds *DocumentService) UploadDocument(document *models.Document, fileData []byte, grant []string) (*models.Document, error) {
	ds.wg.Add(2)

	errorCh := make(chan error, 2)
	filePathCh := make(chan string, 1)

	go func() {
		defer ds.wg.Done()
		fileName := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(document.Name))
		filePath := filepath.Join(ds.uploadDir, fileName)

		if err := os.WriteFile(filePath, fileData, 0644); err != nil {
			errorCh <- fmt.Errorf("error write file: %w", err)
			return
		}
		filePathCh <- filePath
	}()

	go func() {
		defer ds.wg.Done()
		err := ds.docRepo.CreateDocument(document, grant)
		if err != nil {
			errorCh <- fmt.Errorf("error save to database: %w", err)
		}
	}()

	ds.wg.Wait()
	close(errorCh)
	close(filePathCh)

	var fError error
	for err := range errorCh {
		if fError == nil {
			fError = err
		}
	}

	if fError != nil {
		select {
		case filepath := <-filePathCh:
			removeErr := os.Remove(filepath)
			if removeErr != nil && !os.IsNotExist(removeErr) {
				return nil, fmt.Errorf("error delete file after error: %v, %v", fError, removeErr)
			}
		default:
		}
		return nil, fError
	}

	document.FilePath = <-filePathCh

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
