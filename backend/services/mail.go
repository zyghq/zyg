package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/zyghq/zyg/adapters/store"
	"github.com/zyghq/zyg/models"
)

//goland:noinspection ALL
func getExtensionFromContentType(contentType string) string {
	extensions := map[string]string{
		"text/plain":         ".txt",
		"text/html":          ".html",
		"image/jpeg":         ".jpg",
		"image/png":          ".png",
		"image/gif":          ".gif",
		"application/pdf":    ".pdf",
		"application/msword": ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"application/vnd.ms-excel": ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
		"application/zip": ".zip",
	}

	if ext, ok := extensions[contentType]; ok {
		return ext
	}
	return ".bin"
}

func generateS3Key(workspaceId, threadId, activityID, filename string) string {
	return fmt.Sprintf("%s/%s/%s/attachments/%s", workspaceId, threadId, activityID, filename)
}

func generateS3URL(endpoint, bucket, key string) string {
	return fmt.Sprintf("%s/%s/%s", endpoint, bucket, key)
}

func getFilename(filename, contentType string) string {
	if filename != "" {
		return filename
	}

	// Generate filename based on timestamp and content type
	ext := getExtensionFromContentType(contentType)
	return fmt.Sprintf("Attachment_%s%s", time.Now().UTC().Format("20060102_150405"), ext)
}

// removeDataURLPrefix removes the data URL prefix (e.g., "data:image/jpeg;base64,") from a base64-encoded string
// If no prefix is found, returns the original string unchanged
// This ensures the base64 string can be properly decoded by removing metadata prefixes
func removeDataURLPrefix(base64String string) string {
	if idx := strings.Index(base64String, ","); idx != -1 {
		return base64String[idx+1:]
	}
	return base64String
}

// ProcessMessageAttachment processes the base64 media content for the activity
func ProcessMessageAttachment(
	ctx context.Context, workspaceId, threadId, activityID,
	base64Content, contentType, filename string, s3Client store.S3Config) (models.ActivityAttachment, error) {
	if base64Content == "" {
		return models.ActivityAttachment{
			HasError: true,
			Error:    "base64Content cannot be empty",
		}, errors.New("base64Content cannot be empty")
	}

	now := time.Now().UTC()
	attachmentId := (&models.ActivityAttachment{}).GenId()
	attachment := models.ActivityAttachment{
		AttachmentId: attachmentId,
		ActivityID:   activityID,
		Name:         getFilename(filename, contentType),
		ContentType:  contentType,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Decode the base64Content, stripping any data URL prefix
	decodedData, err := base64.StdEncoding.DecodeString(removeDataURLPrefix(base64Content))
	if err != nil {
		attachment.HasError = true
		attachment.Error = fmt.Sprintf("failed to decode base64 content: %v", err)
		return attachment, errors.New("failed to decode base64 content")
	}

	s3Key := generateS3Key(workspaceId, threadId, attachment.ActivityID, attachment.Name)

	// Upload decoded data to S3
	_, err = s3Client.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s3Client.BucketName),
		Key:         aws.String(s3Key),
		Body:        bytes.NewReader(decodedData),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		attachment.HasError = true
		attachment.Error = fmt.Sprintf("failed to upload attachment: %v", err)
		return attachment, errors.New("failed to upload attachment")
	}

	hash := md5.Sum(decodedData)
	md5String := fmt.Sprintf("%x", hash)

	attachment.MD5Hash = md5String
	attachment.ContentKey = s3Key
	attachment.ContentUrl = generateS3URL(s3Client.BaseEndpoint, s3Client.BucketName, s3Key)
	return attachment, nil
}
