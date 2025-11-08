package crawler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type Crawler struct {
	Client   *http.Client
	Visited  map[string]bool
	SaveDir  string
	MaxDepth int
}

// Сrawl рекурсивно проходит все странички до depth, и скачивает их вместе с ресурсами
func (cr *Crawler) Crawl(baseURL string, depth int) ([]string, error) {
	// Если достигли максимальной глубины -> возвращаем пустоту
	if depth <= 0 {
		return nil, nil
	}
	// Если пришли на уже посещенную -> возвращаем пустоту
	if cr.Visited[baseURL] {
		return nil, nil
	}
	cr.Visited[baseURL] = true

	// вычисляем локальный путь (если уже скачано — пропускаем)
	filePath := PathFromURL(baseURL, cr.SaveDir)
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("Пропуск (уже скачан): %s -> %s\n", baseURL, filePath)
		return []string{baseURL}, nil
	}

	// GET страницы
	resp, err := cr.Client.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("код состояния %d для %s", resp.StatusCode, baseURL)
	}

	// читаем всё resp.Body в память
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения тела %s: %v", baseURL, err)
	}
	// парсим HTML
	doc, err := html.Parse(bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML для %s: %v", baseURL, err)
	}
	base, _ := url.Parse(baseURL)

	// Модифицируем HTML
	RewriteLinksAndSaveResources(doc, base, cr)

	// Сохраняем модифицированный HTML в filePath
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("не удалось создать папки для %s: %v", filePath, err)
	}
	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать файл %s: %v", filePath, err)
	}
	defer f.Close()

	if err := html.Render(f, doc); err != nil {
		return nil, fmt.Errorf("ошибка записи HTML %s: %v", filePath, err)
	}
	fmt.Printf("Сохранено: %s\n", filePath)

	// уже в модифицированной страничке ищем ссылки на другие странички по тегу "a"
	links := FindLinks(doc, base, []string{"a"})

	// для дочерних страниц
	childPages := []string{}

	for _, link := range links {
		if strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "tel:") {
			continue
		}
		u, err := url.Parse(link)
		if err != nil {
			continue
		}
		// ограничиваем обход 1им хостом
		if u.Host != base.Host {
			continue
		}
		child, err := cr.Crawl(u.String(), depth-1)
		if err != nil {
			log.Printf("Ошибка при обходе %s: %v", u.String(), err)
			continue
		}
		childPages = append(childPages, child...)
	}

	allLinks := append([]string{baseURL}, childPages...)

	return allLinks, nil
}

func (cr *Crawler) getResourceBody(link string) (io.ReadCloser, error) {
	resp, err := cr.Client.Get(link)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("ошибка получения ресурса: статус %d", resp.StatusCode)
	}
	return resp.Body, nil
}
