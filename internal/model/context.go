package model

// ContextKey 定義專屬的型別，避免 Context 鍵值衝突
type ContextKey string

const (
	// ImageAttachmentKey 用來在 Context 中傳遞圖片
	ImageAttachmentKey ContextKey = "latest_uploaded_image"
)

// FileAttachment 暫存使用者上傳的檔案資訊
type FileAttachment struct {
	MimeType string
	Data     []byte
}
