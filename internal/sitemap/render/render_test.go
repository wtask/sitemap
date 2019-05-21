package render

import (
	"fmt"
	"os"
	"time"

	"github.com/wtask/sitemap/internal/sitemap"
)

var MoscowTZ *time.Location

func init() {
	var err error
	MoscowTZ, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}
}

func ExampleEmptyXMLMap() {
	err := XMLMap(os.Stdout, []sitemap.MapItem{})
	if err != nil {
		panic(fmt.Errorf("Unexpected error: %s", err))
	}

	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	// </urlset>
}

func ExampleXMLMap() {
	uri := func(source string) *sitemap.URI {
		u, err := sitemap.NewURI(source)
		if err != nil {
			panic(err)
		}
		return u
	}
	err := XMLMap(
		os.Stdout,
		[]sitemap.MapItem{
			sitemap.MapItem{
				uri("http://localhost/"),
				&sitemap.DocumentMeta{Modified: time.Time{}},
			},
			sitemap.MapItem{
				uri("http://localhost/homepage.html"),
				&sitemap.DocumentMeta{Modified: time.Date(2019, 5, 21, 23, 26, 0, 0, time.UTC)},
			},
			sitemap.MapItem{
				uri("http://localhost/protocol.html"),
				&sitemap.DocumentMeta{Modified: time.Date(2019, 5, 21, 23, 26, 0, 0, MoscowTZ)},
			},
			sitemap.MapItem{
				//  should be no output
				nil,
				&sitemap.DocumentMeta{Modified: time.Date(2019, 5, 21, 23, 26, 0, 0, time.UTC)},
			},
			sitemap.MapItem{
				uri("http://localhost/faq.html"),
				nil,
			},
		},
	)
	if err != nil {
		panic(fmt.Errorf("Unexpected error: %s", err))
	}

	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	//	<url>
	//		<loc>http://localhost/</loc>
	//	</url>
	//	<url>
	//		<loc>http://localhost/homepage.html</loc>
	//		<lastmod>2019-05-21T23:26:00Z</lastmod>
	//	</url>
	//	<url>
	//		<loc>http://localhost/protocol.html</loc>
	//		<lastmod>2019-05-21T23:26:00+03:00</lastmod>
	//	</url>
	//	<url>
	//		<loc>http://localhost/faq.html</loc>
	//	</url>
	// </urlset>
}

func ExampleEmptyXMLIndex() {
	err := XMLIndex(os.Stdout, time.Time{}, nil)
	if err != nil {
		panic(err)
	}

	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	// </sitemapindex>
}

func ExampleXMLIndexWithoutTime() {
	err := XMLIndex(
		os.Stdout,
		time.Time{},
		[]string{
			"http://www.example.com/sitemap1.xml.gz",
			"",
			"http://www.example.com/sitemap2.xml.gz",
		},
	)
	if err != nil {
		panic(err)
	}

	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	//	<sitemap>
	//		<loc>http://www.example.com/sitemap1.xml.gz</loc>
	//	</sitemap>
	//	<sitemap>
	//		<loc>http://www.example.com/sitemap2.xml.gz</loc>
	//	</sitemap>
	// </sitemapindex>
}

func ExampleXMLIndexWithTime() {
	err := XMLIndex(
		os.Stdout,
		time.Date(2019, 5, 21, 23, 26, 0, 0, MoscowTZ),
		[]string{
			"http://www.example.com/sitemap1.xml.gz",
			"",
			"http://www.example.com/sitemap2.xml.gz",
		},
	)
	if err != nil {
		panic(err)
	}

	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	//	<sitemap>
	//		<loc>http://www.example.com/sitemap1.xml.gz</loc>
	//		<lastmod>2019-05-21T23:26:00+03:00</lastmod>
	//	</sitemap>
	//	<sitemap>
	//		<loc>http://www.example.com/sitemap2.xml.gz</loc>
	//		<lastmod>2019-05-21T23:26:00+03:00</lastmod>
	//	</sitemap>
	// </sitemapindex>
}

func ExampleXMLIndexWithUTCTime() {
	err := XMLIndex(
		os.Stdout,
		time.Date(2019, 5, 21, 23, 26, 0, 0, time.UTC),
		[]string{
			"http://www.example.com/sitemap1.xml.gz",
			"",
			"http://www.example.com/sitemap2.xml.gz",
		},
	)
	if err != nil {
		panic(err)
	}

	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	//	<sitemap>
	//		<loc>http://www.example.com/sitemap1.xml.gz</loc>
	//		<lastmod>2019-05-21T23:26:00Z</lastmod>
	//	</sitemap>
	//	<sitemap>
	//		<loc>http://www.example.com/sitemap2.xml.gz</loc>
	//		<lastmod>2019-05-21T23:26:00Z</lastmod>
	//	</sitemap>
	// </sitemapindex>
}
