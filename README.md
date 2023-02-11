# Implementing Full Text Search (FTS) in Postgres


![cover 0](https://miro.medium.com/v2/resize:fit:786/format:webp/1*m9BuaqJQ-FOOX4ifFQglng.jpeg)

> Text search is a common feature many applications have. We usually use this feature on Instagram, Twitter, or any social media to find someone by their account username, first name, or last name. We use this feature to find an address in Google Maps, (Apple) Maps, Waze, or any Map Navigation

> Implementing the Full Text Search (FTS) in Postgres is similar to using the LIKE operator, i.e., text query-document matching. In FTS, the text matching is conducted by the operator @@ given <doc_text> as the document and <query_text> as the query, e.g., <doc_text> @@ <query_text>.

# Tech Stack

- Go
- Fiber
- PostgreSQL
- Pgx Pool
- Fiber Web Socket
- React JS
- Tailwind CSS

## Screnshoots

![ss 1](https://github.com/SemmiDev/gosh/blob/main/images/screnshoot.png)
![ss 2](https://github.com/SemmiDev/gosh/blob/main/images/screnshoot-2.png)

## How to Run?

```bash
➤ Clone this repo
➤ cd gosh

# setup postgres with dummy data (see: postgres/init.sql)
➤ docker-compose up -d

# install dependencies
➤ go mod tidy

# run backend
➤ go run main.go

# change to frontend
➤ cd gosh

# install dependencies
➤ yarn

# run frontend
➤ yarn dev

# open in browser (http://localhost:5173/)
```

## References

[Fast Text Search to Boost User Experience in Kampus Merdeka Platform](https://medium.com/govtech-edu/fast-text-search-to-boost-user-experience-in-kampus-merdeka-platform-a3a444522754)
