package storage

const base = "https://cdn.denartcc.org/uploads/"

func GenerateURL(slug string) string {
	return base + slug
}

func GetSlugFromURL(url string) string {
	return url[len(base):]
}
