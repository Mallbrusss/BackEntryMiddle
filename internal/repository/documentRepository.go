package repository

import (
	"fmt"
	"log"

	"github.com/Mallbrusss/BackEntryMiddle/models"

	"gorm.io/gorm"
)

type DocumentRepositoryInterface interface {
	CreateDocument(document *models.Document, grant []string) error
	GetDocuments(login string, filter map[string]any, limit int) ([]models.Document, error)
	GetDocumentByID(documentID, login string) (*models.Document, error)
	DeleteDocument(document *models.Document) error
	GetDocumentAccessByID(documentID string) ([]models.DocumentAccess, error)
	FindByToken(token string) (*models.User, error)
	FindByLogin(login string) (*models.User, error)
	IsPermission(documentId string, user *models.User) (bool, error)
}

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (dr *DocumentRepository) CreateDocument(document *models.Document, grant []string) error {
	tx := dr.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	err := tx.Create(document).Error
	if err != nil {
		return err
	}
	var accesses []models.DocumentAccess
	for _, login := range grant {
		accesses = append(accesses, models.DocumentAccess{
			DocID: document.ID,
			Login: login,
		})
	}

	if err := tx.Create(&accesses).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (dr *DocumentRepository) GetDocuments(login string, filter map[string]any, limit int) ([]models.Document, error) {
	var documents []models.Document
	user, err := dr.FindByLogin(login)
	if err != nil {
		return nil, err
	}

	query := dr.db.Model(&models.Document{}).
		Joins("LEFT JOIN document_accesses AS da ON documents.id = da.doc_id")

	if !user.IsAdmin {
		query = query.Where("documents.public = TRUE OR da.login = ?", login)
	}

	for k, v := range filter {
		query = query.Where(fmt.Sprintf("documents.%s = ?", k), v)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err = query.Order("documents.name ASC, documents.created_at DESC").Find(&documents).Error
	if err != nil {
		return nil, err
	}

	if len(documents) > 0 {
		var grants []models.DocumentAccess
		var documentIDs []string
		for _, doc := range documents {
			documentIDs = append(documentIDs, doc.ID)
		}

		err = dr.db.Where("doc_id IN (?)", documentIDs).Find(&grants).Error
		if err != nil {
			return nil, err
		}

		for i := range documents {
			var logins []string
			for _, grant := range grants {
				if documents[i].ID == grant.DocID {
					logins = append(logins, grant.Login)
				}
			}
			documents[i].Grant = logins
		}
	}

	log.Println(query.Statement.SQL.String())
	return documents, nil
}

func (dr *DocumentRepository) GetDocumentAccessByID(documentID string) ([]models.DocumentAccess, error) {
	var grants []models.DocumentAccess

	err := dr.db.Model(&models.DocumentAccess{}).
		Where("doc_id = ?", documentID).
		Find(&grants).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching document accesses: %w", err)
	}

	return grants, nil
}

func (dr *DocumentRepository) GetDocumentByID(documentID, login string) (*models.Document, error) {
	var document models.Document

	user, err := dr.FindByLogin(login)
	if err != nil {
		return nil, err
	}

	query := dr.db.Model(&models.Document{}).
		Joins("LEFT JOIN document_accesses AS da ON documents.id = da.doc_id")

	if !user.IsAdmin {
		query = query.Where("documents.public = TRUE OR da.login = ?", login)
	}

	err = query.Where("documents.id = ?", documentID).First(&document).Error
	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (dr *DocumentRepository) DeleteDocument(document *models.Document) error {
	err := dr.db.Where("doc_id = ?", document.ID).Delete(&models.DocumentAccess{}).Error
	if err != nil {
		return err
	}
	return dr.db.Delete(document).Error
}

func (dr *DocumentRepository) CreateAccess(access *models.DocumentAccess) error {
	return dr.db.Create(access).Error
}

func (dr *DocumentRepository) FindByToken(token string) (*models.User, error) {
	var user models.User
	err := dr.db.
		Where("token = ?", token).
		First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (dr *DocumentRepository) FindByLogin(login string) (*models.User, error) {
	var user models.User

	err := dr.db.
		Where("login = ?", login).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (dr *DocumentRepository) IsPermission(documentId string, user *models.User) (bool, error) {
	if user.IsAdmin {
		return true, nil
	}

	var count int64
	err := dr.db.Model(&models.DocumentAccess{}).
		Where("doc_id = ? AND login = ?", documentId, user.Login).
		Count(&count)

	if err.Error != nil {
		return false, err.Error
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
