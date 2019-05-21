package render

import (
	"io"
	"text/template"
	"time"

	"github.com/wtask/sitemap/internal/sitemap"
)

const (
	xmlMap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
{{- range . }}
	{{- if .URI }}
	<url>
		<loc>{{ .URI.String }}</loc>
		{{- with .DocumentMeta }}
			{{- if not .Modified.IsZero }}
		<lastmod>{{ .Modified.Format "2006-01-02T15:04:05Z07:00" }}</lastmod>
			{{- end }}
		{{- end }}
	</url>
	{{- end}}
{{- end }}
</urlset>
`
	xmlIndex = `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
{{- $time := .Modified }}
{{- range .Files }}	
	{{- with . }}
	<sitemap>
		<loc>{{ . }}</loc>
			{{- if not $time.IsZero }}
		<lastmod>{{ $time.Format "2006-01-02T15:04:05Z07:00" }}</lastmod>
			{{- end }}
	</sitemap>
	{{- end }}
{{- end }}
</sitemapindex>
`
)

var xml *template.Template

func init() {
	xml = template.Must(template.New("map").Parse(xmlMap))
	xml = template.Must(xml.New("index").Parse(xmlIndex))
}

func XMLMap(writer io.Writer, m []sitemap.MapItem) error {
	return xml.Lookup("map").Execute(writer, m)
}

func XMLIndex(writer io.Writer, modified time.Time, fileURI []string) error {
	data := struct {
		Modified time.Time
		Files    []string
	}{
		modified,
		fileURI,
	}
	return xml.Lookup("index").Execute(writer, data)
}
