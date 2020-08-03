package repository

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jinzhu/gorm"
	"local.packages/config"
	"local.packages/entity"
	"local.packages/validator"
)

// FileRepository ...
type FileRepository struct {
	db *gorm.DB
}

// NewFileRepository is constructor for FileRepository.
func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

func selectFileColumn(db *gorm.DB) *gorm.DB {
	return db.Select("files.id, files.display_name, files.url, files.content_type, files.card_id")
}

// ValidateUID validates whether a cardID received as args was created by the login user.
func (r *FileRepository) ValidateUID(cid, uid uint) []validator.ValidationError {
	var b entity.Board

	if r.db.Joins("Join lists ON boards.id = lists.board_id").
		Joins("Join cards ON lists.id = cards.list_id").
		Select("user_id").
		Where("cards.id = ?", cid).
		Where("boards.user_id = ?", uid).
		First(&b).
		RecordNotFound() {
		return validator.NewValidationErrors("invalid parameters")
	}

	return nil
}

// Find returns a record of File that found by id.
func (r *FileRepository) Find(id, uid uint) (*entity.File, []validator.ValidationError) {
	var f entity.File

	if r.db.Joins("Join cards ON files.card_id = cards.id").
		Joins("Join lists ON cards.list_id = lists.id").
		Joins("Join boards ON lists.board_id = boards.id").
		Where("boards.user_id = ?", uid).
		First(&f, id).
		RecordNotFound() {
		return &f, validator.NewValidationErrors("invalid parameters")
	}

	return &f, nil
}

// Upload upload file to S3
func (r *FileRepository) Upload(fh *multipart.FileHeader, cid uint) *entity.File {
	u := s3manager.NewUploader(config.AWSSession())

	ct := fh.Header.Get("Content-Type")

	f, err := fh.Open()

	if err != nil {
		log.Printf("failed open file: %v", err)
		return nil
	}

	defer f.Close()

	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		log.Printf("failed create hash: %v", err)
		return nil
	}

	hash := base64.URLEncoding.EncodeToString(b)

	key := fmt.Sprintf("%s-%s", hash[:8], fh.Filename)

	uo, err := u.Upload(&s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Bucket:      aws.String(config.Config.AWS.Bucket),
		Key:         aws.String(fmt.Sprintf("%d/%s", cid, key)),
		Body:        f,
		ContentType: aws.String(ct),
	})

	if err != nil {
		log.Printf("failed upload file: %v", err)
		return nil
	}

	return &entity.File{
		DisplayName: fh.Filename,
		Key:         key,
		URL:         uo.Location,
		ContentType: ct,
		CardID:      cid,
	}
}

// Create insert a new record to a files table.
func (r *FileRepository) Create(f *entity.File) []validator.ValidationError {
	if err := r.db.Create(f).Error; err != nil {
		return validator.FormattedMySQLError(err)
	}

	return nil
}

// Delete delete a record from a files table.
func (r *FileRepository) Delete(f *entity.File) []validator.ValidationError {
	if rslt := r.db.Delete(f); rslt.RowsAffected == 0 {
		log.Printf("fail to delete file: %v", rslt.Error)
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

// DeleteObject deletes an object from S3 bucket.
func (r *FileRepository) DeleteObject(key string) error {
	svc := s3.New(config.AWSSession())

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(config.Config.AWS.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		log.Printf("failed delete file object: %v", err)
		return err
	}

	return nil
}

// GetAll returns slice of File's record.
func (r *FileRepository) GetAll(bid, uid uint) *[]entity.File {
	var fs []entity.File

	r.db.Scopes(selectFileColumn).
		Joins("Join cards ON files.card_id = cards.id").
		Joins("Join lists ON cards.list_id = lists.id").
		Joins("Join boards ON lists.board_id = boards.id").
		Where("boards.id = ?", bid).
		Where("boards.user_id = ?", uid).
		Find(&fs)

	return &fs
}
