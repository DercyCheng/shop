package valueobject

import (
	"path"
	"strings"
)

// ImageType 图片类型枚举
type ImageType string

const (
	ImageTypeProduct  ImageType = "product"  // 商品图片
	ImageTypeThumbnail ImageType = "thumbnail" // 缩略图
	ImageTypeBanner   ImageType = "banner"   // 轮播图
	ImageTypeLogo     ImageType = "logo"     // 品牌Logo
)

// Image 图片值对象
type Image struct {
	URL         string    // 图片URL
	Type        ImageType // 图片类型
	Title       string    // 图片标题/描述
	Size        int64     // 文件大小(字节)
	Width       int       // 宽度(像素)
	Height      int       // 高度(像素)
	ContentType string    // 内容类型，如image/jpeg
}

// NewImage 创建图片值对象
func NewImage(url string, imageType ImageType) *Image {
	return &Image{
		URL:         url,
		Type:        imageType,
		ContentType: inferContentType(url),
	}
}

// inferContentType 根据URL推断内容类型
func inferContentType(url string) string {
	ext := strings.ToLower(path.Ext(url))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "image/jpeg" // 默认类型
	}
}

// IsValidImage 验证图片有效性
func (i *Image) IsValidImage() bool {
	if i == nil || i.URL == "" {
		return false
	}
	
	// 检查URL格式
	if !strings.HasPrefix(i.URL, "http://") && !strings.HasPrefix(i.URL, "https://") {
		return false
	}
	
	// 检查文件扩展名
	ext := strings.ToLower(path.Ext(i.URL))
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".svg":  true,
	}
	
	return validExts[ext]
}

// IsSameImage 判断是否是同一图片
func (i *Image) IsSameImage(other *Image) bool {
	if i == nil || other == nil {
		return false
	}
	
	return i.URL == other.URL
}

// WithDimensions 设置图片尺寸
func (i *Image) WithDimensions(width, height int) *Image {
	if i == nil {
		return nil
	}
	
	i.Width = width
	i.Height = height
	return i
}

// WithSize 设置图片文件大小
func (i *Image) WithSize(size int64) *Image {
	if i == nil {
		return nil
	}
	
	i.Size = size
	return i
}

// WithTitle 设置图片标题
func (i *Image) WithTitle(title string) *Image {
	if i == nil {
		return nil
	}
	
	i.Title = title
	return i
}
