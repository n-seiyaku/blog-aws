package image

const MaxSize int64 = 5 * 1024 * 1024 // 5 MB

var allowedTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

func IsAllowedContentType(contentType string) bool {
	_, ok := allowedTypes[contentType]
	return ok
}

func Extension(contentType string) (string, bool) {
	ext, ok := allowedTypes[contentType]
	return ext, ok
}

func IsValidSize(size int64) bool {
	return size > 0 && size <= MaxSize
}
