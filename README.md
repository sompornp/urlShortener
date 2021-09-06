# URL Shortener

## What is it?

A simple program to generate shorten url. It is implemented in Golang and saved data in sqlite db file.

### Features
- Generate shortened URL
- Support blacklist urls
- Hit count from visiting generated shortened URL
- Search data
- Delete data

### Requirements

- Go (should be 1.16+)

### How to run

1. git clone https://github.com/sompornp/urlShortener
2. go mod download
3. Rename or copy .env.example to .env. Take a look on each config value in the file or description in later section.
4. go run main.go

## APIs

| Path | Method | Description |
| ---- | ------ | ----------- |
| /ping | GET |  Health check api |
| /encode | POST | URL shortened api. Example curl: <br/> ``` curl -v -X POST localhost:8080/encode -H 'Content-Type: application/json'  --data-raw '{"Url": "http://www.ku.ac.th","Expire": 1630990800}'```|
| /admin | GET | Query URL shortened data in db. Example curls: <br/> ``` curl -v http://localhost:8080/admin -H 'Authorization: testit' ``` <br/><br/> ``` curl -v http://localhost:8080/admin?keyword=http://www.ku.th -H 'Authorization: testit' ``` <br/><br/> ``` curl -v http://localhost:8080/admin?Shortcode=ojFaTA&Keyword=http://www.ku.ac.th -H 'Authorization: testit' ```|
| /admin | DELETE | Delete URL shortened record in db. Example curl: <br/> ``` curl -v -X DELETE http://localhost:8080/admin -H 'Authorization: testit' --data-raw '{"Shortcode": "qHMLYY"}' ``` |
| /{shortLink} | GET | Given shortLink generated from /encode api and redirect to target URL. Example curl: <br/> ``` curl -v http://localhost:8080/VxFCtx ``` |

## ENV Configs

### DB_FILE

Sqlite file, this field is required

### BLACKLIST_URLS

Blacklist urls, supporting regular expression. For regular expression url, prefix with 'regex:'.

```
BLACKLIST_URLS=regex: http://news.bbc.co.uk, http://www.ku.ac.th, http://wwww.
```
This field is optional

### REFERER

Hosting url at which shortening
url is referred to. If omitted, default to http://localhost:8080
```
REFERER=http://localhost:8080
```

### ADMIN_TOKEN

Token for admin related apis. If omitted, will auto-generate. See it on console when first started
```
ADMIN_TOKEN=testit
```