package crawler

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// pathFromURL делает безопасное имя файла из URL и его путь от raw
// если нет нужного пути создает папки, гарантирует что возвращаемый путь существует
func PathFromURL(raw, saveDir string) string {
	u, err := url.Parse(raw)
	if err != nil {
		u = &url.URL{Path: raw}
	}

	host := u.Host
	if host == "" {
		host = "unknown"
	}

	path := u.Path
	if path == "" || path == "/" {
		path = "/index"
	}

	// Убираем лишние символы
	path = strings.TrimSuffix(path, "/")
	path = strings.ReplaceAll(path, ":", "_")
	path = strings.ReplaceAll(path, "?", "_")
	path = strings.ReplaceAll(path, "&", "_")

	// Полный путь с подпапками
	fullPath := filepath.Join(saveDir, host, path)

	// Создаём подпапки, если не хватает
	dir := filepath.Dir(fullPath)
	os.MkdirAll(dir, 0755)

	// Добавляем .html, если нет расширения
	ext := filepath.Ext(fullPath)
	if ext == "" {
		fullPath += ".html"
	}
	return fullPath
}

// SaveFile сохраняет поток r по указанному пути, создавая нужные папки.
// Для ресурсов
func SaveFile(path string, r io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}
