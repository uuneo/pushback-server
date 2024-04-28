package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pushbackServer/middleware/ipban"
	"strings"
)

// 允许的图片类型
var allowedImageTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/gif":     true,
	"image/webp":    true,
	"image/bmp":     true,
	"image/svg+xml": true,
}

// UploadController 处理图片上传请求
// 支持GET和POST两种请求方式:
// GET: 返回上传页面
// POST: 上传图片并保存
func UploadController(c *gin.Context) {
	// 验证管理员权限
	admin, ok := c.Get("admin")

	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "upload.html", gin.H{})
		return
	}

	if !ok || !admin.(bool) {
		log.Println("Unauthorized upload attempt")
		manager := ipban.GetManager()
		manager.BanIP(c.ClientIP())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 获取文件名
	fileName := c.PostForm("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	// 创建上传目录
	uploadDir := "./images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}
	// 验证文件是否在uploads目录下
	if isTrue, _ := isFileInDirectory(fileName, uploadDir); isTrue {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found in uploads directory"})
		return

	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload failed: " + err.Error()})
		return
	}

	// 验证文件类型
	if !allowedImageTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only image files are allowed"})
		return
	}

	// 生成安全的文件名
	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = filepath.Ext(file.Filename)
	}
	if ext == "" {
		ext = ".jpg" // 默认扩展名
	}

	// 使用时间戳和随机字符串生成唯一文件名
	safeFileName := fmt.Sprintf("%s%s",
		strings.TrimSuffix(fileName, ext),
		strings.ToLower(ext))

	// 完整的文件路径
	filePath := filepath.Join(uploadDir, safeFileName)

	// 保存文件
	if err = c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file: " + err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"filename": safeFileName,
		"path":     filePath,
		"size":     file.Size,
		"type":     file.Header.Get("Content-Type"),
	})
}

func GetImage(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	// 构建文件路径
	filePath := filepath.Join("./images", fileName)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// 返回文件
	c.File(filePath)

}

func isFileInDirectory(dirPath, fileName string) (bool, error) {
	// 对目录路径进行规范化处理
	dirPath = filepath.Clean(dirPath)

	// 检查目录是否存在
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // 目录不存在时，直接返回文件不在目录中
		}
		return false, fmt.Errorf("检查目录状态出错: %w", err)
	}

	// 确认路径指向的是一个目录
	if !dirInfo.IsDir() {
		return false, fmt.Errorf("路径 %q 不是一个目录", dirPath)
	}

	// 构建文件的完整路径
	filePath := filepath.Join(dirPath, fileName)

	// 检查文件是否存在
	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // 文件不存在
		}
		return false, fmt.Errorf("检查文件状态出错: %w", err)
	}

	return true, nil
}
