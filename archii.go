package archii

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/rs/xid"
)

func fileName(article *readability.Article) string {
	return fmt.Sprintf("%s.%s.org", strings.Map(func(r rune) rune {
		switch r {
		case ' ', '/', ':', '*', '|', '?', '"', '\'', '<', '>', '\\':
			return '_'
		}
		return r
	}, article.Title), xid.New())
}

func RunApp(url string, dir string) error {
	article, err := readability.FromURL(url, 20*time.Second)
	if err != nil {
		return err
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("#+title: %s\n", article.Title))
	builder.WriteString(fmt.Sprintf("#+url: %s\n", url))
	builder.WriteString(fmt.Sprintf("#+author: %s\n", article.Byline))
	builder.WriteString(fmt.Sprintf("#+date_saved: %s\n", time.Now().Format(time.RFC3339)))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("* %s\n%s", article.Title, article.TextContent))

	path := fmt.Sprintf("%s/%s", dir, fileName(&article))
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	buffer := bufio.NewWriter(out)
	buffer.WriteString(builder.String())
	buffer.Flush()

	fmt.Println(builder.String())
	fmt.Printf("Saved to: %s\n", path)
	return nil
}
