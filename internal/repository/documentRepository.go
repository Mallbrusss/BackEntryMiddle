package repository

import (
	"fmt"
	"github.com/Mallbrusss/BackEntryMiddle/models"

	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (dr *DocumentRepository) CreateDocument(document *models.Document, grant []string) error {
	err := dr.db.Create(document).Error
	if err != nil {
		return err
	}
	for _, login := range grant {
		access := models.DocumentAccess{
			ID:    document.ID,
			Login: login,
		}
		if err := dr.db.Create(&access).Error; err != nil {
			return err
		}
	}
	return nil
}

func (dr *DocumentRepository) GetDocuments(login string, filter map[string]any, limit int) ([]models.Document, error) {
	var documents []models.Document
	query := dr.db.Model(&models.Document{}).
		Joins("LEFT JOIN document_access ON document.id = document_access.document_id").
		Where("documents.public = TRUE OR document_access.login = ?", login)

	for k, v := range filter {
		query = query.Where(fmt.Sprintf("documents.%s = ?", k), v)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&documents).Error
	return documents, err
}

func (dr *DocumentRepository) GetDocumentByID(documentID, login string) (*models.Document, error) {
	var document models.Document

	err := dr.db.Model(&models.Document{}).
		Joins("LEFT JOIN document_access ON document.id = document_access.document_id").
		Where("documents.id = ? AND (documents.public = TRUE OR document_access.login = ?)", documentID).
		First(&document).Error
	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (dr *DocumentRepository) DeleteDocument(document *models.Document) error {
	err := dr.db.Where("document_id = ?", document.ID).Delete(&models.DocumentAccess{}).Error
	if err != nil {
		return err
	}
	return dr.db.Delete(document).Error
}

func (dr *DocumentRepository) CreateAccess(access *models.DocumentAccess) error {
	return dr.db.Create(access).Error
}

//TODO: Добавить транзакции