package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/xid"
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

func generateS3Key(workspaceId, threadId, messageId, filename string) string {
	return fmt.Sprintf("%s/%s/%s/attachments/%s", workspaceId, threadId, messageId, filename)
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

// ProcessMessageAttachment handles the processing and uploading of file attachments to S3.
// It takes a base64-encoded file content, content type, filename and other metadata,
// decodes the content, uploads it to S3 and returns a MessageAttachment object with
// the result details and error if any.
//
// Parameters:
// - ctx: Context for the operation
// - workspaceId: ID of the workspace the attachment belongs to
// - threadId: ID of the thread the attachment belongs to
// - messageId: ID of the message the attachment belongs to
// - base64Content: Base64 encoded file content
// - contentType: MIME type of the file
// - filename: Optional filename (will be generated if empty)
// - s3Client: S3 client configuration
//
// Returns:
// A MessageAttachment object containing the upload results or error details if failed and error object.
func ProcessMessageAttachment(
	ctx context.Context, workspaceId, threadId, messageId,
	base64Content, contentType, filename string, s3Client store.S3Config) (models.MessageAttachment, error) {
	if base64Content == "" {
		return models.MessageAttachment{
			HasError: true,
			Error:    "base64Content cannot be empty",
		}, errors.New("base64Content cannot be empty")
	}

	now := time.Now().UTC()
	attachment := models.MessageAttachment{
		AttachmentId: xid.New().String(),
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

	s3Key := generateS3Key(workspaceId, threadId, messageId, attachment.Name)

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

	attachment.ContentKey = s3Key
	attachment.ContentUrl = generateS3URL(s3Client.BaseEndpoint, s3Client.BucketName, s3Key)
	return attachment, nil
}
