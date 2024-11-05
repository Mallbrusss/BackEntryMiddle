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
	query := dr.db.Model(&models.Document{}).
		Joins("LEFT JOIN document_accesses ON documents.id = document_accesses.doc_id").
		Where("documents.public = TRUE OR document_accesses.login = ?", login)

	for k, v := range filter {
		query = query.Where(fmt.Sprintf("documents.%s = ?", k), v)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("name ASC, created_at DESC").Find(&documents).Error

	var grants []models.DocumentAccess
	if len(documents) > 0 {
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

	return documents, err
}

func (dr *DocumentRepository) GetDocumentByID(documentID, login string) (*models.Document, error) {
	var document models.Document

	err := dr.db.Model(&models.Document{}).
		Joins("LEFT JOIN document_accesses ON documents.id = document_accesses.doc_id").
		Where("documents.id = ? AND (documents.public = TRUE OR document_accesses.login = ?)", documentID, login).
		First(&document).Error
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
		return nil, err // Вернем ошибку, если токен не найден
	}

	return &user, nil
}

func (dr *DocumentRepository) FindByLogin(login string) (*models.User, error) {
	var user models.User

	err := dr.db.
		Where("login = ?", login).
		First(&user).Error

	if err != nil {
		return nil, err // Вернем ошибку, если пользователь не найден
	}

	return &user, nil
}
