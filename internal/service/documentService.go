package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/Mallbrusss/BackEntryMiddle/internal/repository"
	"github.com/Mallbrusss/BackEntryMiddle/models"
	"github.com/redis/go-redis/v9"

	"github.com/google/uuid"
)

type DocumentService struct {
	docRepo   *repository.DocumentRepository
	rdb       *redis.Client
	uploadDir string
}

func NewDocumentService(dockRepo *repository.DocumentRepository, rdb *redis.Client, uploadDir string) *DocumentService {
	return &DocumentService{
		docRepo:   dockRepo,
		rdb:       rdb,
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

	errorCh := make(chan error, 3)

	ext := ds.getFileExtensions(document.Mime)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(ds.uploadDir, fileName)
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := os.WriteFile(filePath, fileData, 0644); err != nil {
			errorCh <- fmt.Errorf("error write file: %w", err)
			return
		}
	}()
	wg.Add(1)

	go func() {
		defer wg.Done()
		document.ID = uuid.New().String()
		document.FilePath = filePath
		err := ds.docRepo.CreateDocument(document, grant)
		if err != nil {
			errorCh <- fmt.Errorf("error save to database: %w", err)
			return
		}

		jsonData, err := json.Marshal(document)
		if err != nil {
			errorCh <- fmt.Errorf("error marshaling doc to json: %w", err)
			return
		}

		err = ds.rdb.Set(context.Background(), document.ID, jsonData, 10*time.Second).Err()
		if err != nil {
			errorCh <- fmt.Errorf("error saving document to Redis: %w", err)
		}

	}()

	wg.Wait()
	close(errorCh)

	for err := range errorCh {
		if err != nil {
			removeErr := os.Remove(filePath)
			if removeErr != nil && !os.IsNotExist(removeErr) {
				return nil, fmt.Errorf("error: %v, error deleting file: %v", err, removeErr)
			}
			return nil, err
		}
	}

	return document, nil
}

func (ds *DocumentService) DeleteDocument(documentID, login string) error {
	document, err := ds.docRepo.GetDocumentByID(documentID, login)
	if err != nil {
		log.Println("File not found")
		return err
	}

	user, err := ds.docRepo.FindByLogin(login)
	if err != nil {
		log.Println("User not found")
		return err
	}

	if user.IsAdmin {
		return ds.DeleteDocumentFromSystem(document)
	}

	hasPermission, err := ds.docRepo.IsPermission(documentID, user)
	if err != nil {
		log.Printf("Error checking user permissions: %v", err)
		return err
	}
	if !hasPermission {
		return fmt.Errorf("user does not have permission to delete this document")
	}

	return ds.DeleteDocumentFromSystem(document)
}

func (ds *DocumentService) DeleteDocumentFromSystem(document *models.Document) error {
	var wg sync.WaitGroup
	errorCh := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if document.File && document.FilePath != "" {
			if err := os.Remove(document.FilePath); err != nil && !os.IsExist(err) {
				errorCh <- fmt.Errorf("error deleting file: %w", err)
				return
			}

		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := ds.docRepo.DeleteDocument(document)
		if err != nil {
			errorCh <- fmt.Errorf("error deleting document from database: %w", err)
			return
		}

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx := context.Background()
		err := ds.rdb.Del(ctx, document.ID).Err()
		if err != nil {
			errorCh <- fmt.Errorf("error deleting document from Redis cache: %w", err)
			return
		}

	}()

	wg.Wait()
	close(errorCh)

	for err := range errorCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func (ds *DocumentService) GetDocumentByID(documentID, login string) (*models.Document, error) {
	ctx := context.Background()

	cachedDoc, err := ds.rdb.Get(ctx, documentID).Result()
	if err == redis.Nil {
		document, err := ds.docRepo.GetDocumentByID(documentID, login)
		if err != nil {
			return nil, err
		}

		grants, err := ds.docRepo.GetDocumentAccessByID(documentID)
		if err != nil {
			log.Printf("error getting document accesses: %v", err)
			return nil, err
		}

		for _, grant := range grants {
			document.Grant = append(document.Grant, grant.Login)
		}

		cacheDocument := models.CacheDocument{
			ID:        document.ID,
			Name:      document.Name,
			Mime:      document.Mime,
			FilePath:  document.FilePath,
			File:      document.File,
			Public:    document.Public,
			Token:     document.Token,
			CreatedAt: document.CreatedAt,
			Grant:     document.Grant,
		}

		jsonData, err := json.Marshal(cacheDocument)
		if err != nil {
			log.Println("error marshaling document to JSON: %w", err)
			return nil, err
		}

		err = ds.rdb.Set(ctx, documentID, jsonData, 10*time.Second).Err()
		if err != nil {

			log.Println("error saving document to Redis: %w")
			return nil, err
		}

		return document, nil
	} else if err != nil {
		log.Println("error saving document to Redis: %w", err)
		return nil, err
	}

	var cachedDocument models.CacheDocument
	err = json.Unmarshal([]byte(cachedDoc), &cachedDocument)
	if err != nil {
		log.Println("error unmarshaling document from JSON: %w", err)
		return nil, err
	}

	document := &models.Document{
		ID:        cachedDocument.ID,
		Name:      cachedDocument.Name,
		Mime:      cachedDocument.Mime,
		FilePath:  cachedDocument.FilePath,
		File:      cachedDocument.File,
		Public:    cachedDocument.Public,
		Token:     cachedDocument.Token,
		CreatedAt: cachedDocument.CreatedAt,
		Grant:     cachedDocument.Grant,
	}

	return document, nil
}

func (ds *DocumentService) GetDocuments(login string, filter map[string]any, limit int) ([]models.Document, error) {
	ctx := context.Background()

	cacheKey := fmt.Sprintf("documents:%s:%v:%d", login, filter, limit)
	cachedDocs, err := ds.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		documents, err := ds.docRepo.GetDocuments(login, filter, limit)
		if err != nil {
			return nil, err
		}

		jsonData, err := json.Marshal(documents)
		if err != nil {
			return nil, fmt.Errorf("error marshaling documents to JSON: %w", err)
		}

		err = ds.rdb.Set(ctx, cacheKey, jsonData, 10*time.Second).Err()
		if err != nil {
			return nil, fmt.Errorf("error saving documents to Redis: %w", err)
		}

		sort.Slice(documents, func(i, j int) bool {
			if documents[i].Name == documents[j].Name {
				return documents[i].CreatedAt.Before(documents[j].CreatedAt)
			}
			return documents[i].Name < documents[j].Name
		})
		return documents, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error retrieving documents from Redis: %w", err)
	}

	var documents []models.Document
	err = json.Unmarshal([]byte(cachedDocs), &documents)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling documents from JSON: %w", err)
	}

	return documents, nil
}

func (ds *DocumentService) GetUserByToken(token string) (*models.User, error) {
	session, err := ds.docRepo.FindByToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	user, err := ds.docRepo.FindByLogin(session.Login)
	if err != nil {
		return nil, err
	}
	return user, nil
}
