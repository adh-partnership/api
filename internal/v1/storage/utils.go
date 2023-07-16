package storage

var base = "https://cdn.zdvartcc.org/"

func GenerateURL(slug string) string {
	return base + slug
}

func GetSlugFromURL(url string) string {
	return url[len(base+"uploads/"):]
}
