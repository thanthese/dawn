package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	var userName string
	err = dbpool.QueryRow(ctx, "select CURRENT_USER").Scan(&userName)
	if err != nil {
		log.Fatalf("QueryRow failed: %v\n", err)
	}
	fmt.Printf("Current user is '%s'\n", userName)

	albs, err := albumsByArtist(dbpool, ctx, "John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albs)

	err = addAlbum(dbpool, ctx, Album{
		Title:  "Abbey Road",
		Artist: "Beatles",
		Price:  10})
	if err != nil {
		log.Fatal(err)
	}

	count, err := rowsCount(dbpool, ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of rows in album: %v", count)
}

func rowsCount(dbpool *pgxpool.Pool, ctx context.Context) (count int, err error) {
	err = dbpool.QueryRow(ctx, "select count(*) from album;").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("couldn't get count: %w", err)
	}
	return count, nil
}

func albumsByArtist(dbpool *pgxpool.Pool, ctx context.Context, name string) ([]Album, error) {
	rows, err := dbpool.Query(ctx, "select * from album where artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("make connection, albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Album])
}

func addAlbum(dbpool *pgxpool.Pool, ctx context.Context, alb Album) error {
	query := `INSERT INTO album (title, artist, price) VALUES (@title, @artist, @price)`
	args := pgx.NamedArgs{
		"title":  alb.Title,
		"artist": alb.Artist,
		"price":  alb.Price,
	}
	_, err := dbpool.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	return nil
}
