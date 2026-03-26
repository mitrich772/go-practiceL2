package crawler

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func contain[T comparable](elem T, slice []T) bool {
	for _, v := range slice {
		if elem == v {
			return true
		}
	}
	return false
}

// FindLinks ищет href в тегах и отправляет их в канал.
// tags слайс тегов в которых искать ссылки
func FindLinks(node *html.Node, base *url.URL, tags []string, links chan string) {
	// Создаем единую WaitGroup для всех дочерних элементов
	var wg sync.WaitGroup

	// Внутренняя функция для рекурсивного обхода дерева
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n == nil {
			return
		}

		if n.Type == html.ElementNode && contain(n.Data, tags) {
			for _, v := range n.Attr {
				if v.Key == "href" || v.Key == "src" {
					// Нашли ссылку парсим ее в url, добавляем абсолютный путь
					u, err := url.Parse(v.Val)
					if err != nil {
						continue
					}
					if base != nil {
						u = base.ResolveReference(u)
					}
					// Если нет схемы то добавлям базовую
					if u.Scheme == "" {
						u.Scheme = base.Scheme
					}
					if u.Scheme == "http" || u.Scheme == "https" {
						links <- u.String()
					}
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			wg.Add(1)

			go func(c *html.Node) {
				defer wg.Done()
				walk(c)
			}(child)
		}
	}

	walk(node)

	wg.Wait()

	close(links)
}

// RewriteLinksAndSaveResources проходит по HTML-дереву, находит ссылки (href/src),
// скачивает ресурсы и переписывает ссылки на локальные пути.
func RewriteLinksAndSaveResources(node *html.Node, base *url.URL, cr *Crawler) {
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n == nil {
			return
		}

		if n.Type == html.ElementNode {
			for i, a := range n.Attr {
				if a.Key != "href" && a.Key != "src" {
					continue
				}

				u := resolveURL(a.Val, base)
				if u == nil {
					continue
				}

				ext := filepath.Ext(u.Path)
				if isResource(ext) {
					newVal := processResourceLink(u, base, cr)
					if newVal != "" {
						n.Attr[i].Val = newVal
					}
				} else {
					n.Attr[i].Val = u.String()
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(node)
}

// resolveURL превращает ссылку из атрибута в абсолютный URL.
func resolveURL(raw string, base *url.URL) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		return nil
	}
	if base != nil {
		u = base.ResolveReference(u)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil
	}
	if u.Scheme == "mailto" || u.Scheme == "tel" {
		return nil
	}
	return u
}

// isResource возвращает true, если это не HTML-страница, а ресурс (css, js, img и т.д.)
func isResource(ext string) bool {
	if ext == "" {
		return false
	}
	ext = strings.ToLower(ext)
	return ext != ".html" && ext != ".htm"
}

// processResourceLink скачивает файл, если его нет, и возвращает относительный путь,
// чтобы заменить значение атрибута href/src.
// u - URl нашего ресурса
// base - текущий URL на котором найден ресурс
func processResourceLink(u, base *url.URL, cr *Crawler) string {
	localPath := PathFromURL(u.String(), cr.SaveDir)

	// если файла нет — качаем
	if _, err := os.Stat(localPath); err != nil {
		rc, err := cr.getResourceBody(u.String())
		if err != nil {
			return ""
		}
		defer rc.Close()

		if err := SaveFile(localPath, rc); err != nil {
			return ""
		}
		fmt.Printf("Ресурс сохранён: %s\n", localPath)
	}

	// делаем относительный путь от текущей страницы до ресурса
	currentPageLocal := PathFromURL(base.String(), cr.SaveDir)
	rel, err := filepath.Rel(filepath.Dir(currentPageLocal), localPath)
	if err != nil {
		return localPath // fallback: абсолютный путь
	}
	return filepath.ToSlash(rel)
}
