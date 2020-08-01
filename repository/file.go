package repository

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
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
		Name:        fh.Filename,
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
