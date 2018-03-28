package bookmarks

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Keyword represents a Keyword linked to a Bookmark
type Keyword string

// keyword => db ID
type keywordsMap map[Keyword]int

func loadKeywords(db *sqlx.DB, id int) ([]Keyword, error) {
	sql := `
SELECT kw.name
FROM bookmark_keywords bkw
INNER JOIN keywords kw ON kw.id = bkw.keyword_id
WHERE bkw.bookmark_id = ?
`
	rows, err := db.Query(sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keywords := []Keyword{}
	for rows.Next() {
		var kw string
		if err = rows.Scan(&kw); err != nil {
			return nil, err
		}
		keywords = append(keywords, Keyword(kw))
	}

	return keywords, nil
}

func saveKeywords(tx *sqlx.Tx, bookmarkID int, keywords []Keyword) error {
	// fist remove old associations
	if err := deleteKwAssociations(tx, bookmarkID); err != nil {
		return err
	}

	if len(keywords) == 0 {
		// nothing more to do
		return nil
	}

	// find which keywords already exist in DB to reuse them
	kwMap, err := getKeywordsMap(tx, keywords)
	if err != nil {
		return err
	}

	kwIDs := []int{}
	for kw, id := range kwMap {
		switch {
		case id == 0:
			// create new keywords. We can't bulk insert because we want to get the IDs back
			newID, err := createKeyword(tx, kw)
			if err != nil {
				return err
			}
			kwIDs = append(kwIDs, newID)
		default:
			kwIDs = append(kwIDs, id)
		}
	}

	// finally create associations
	if err := saveKwAssociations(tx, bookmarkID, kwIDs); err != nil {
		return err
	}

	return nil
}

// return known keywords with their IDs
// new keywords are assigned the ID 0
func getKeywordsMap(tx *sqlx.Tx, keywords []Keyword) (keywordsMap, error) {
	sql := `SELECT id, name FROM keywords WHERE name IN (?)`

	// sqlx is not able to treat a []Keyword as a []string so we have to build it manually
	// there is no smarter way to do that in go (except using unsafe but it's another story)
	// this is the downside of creating a distinct Keyword types
	// but using clear types always pays back eventually
	kws := make([]string, 0, len(keywords))
	for _, kw := range keywords {
		kws = append(kws, string(kw))
	}

	// sqlx management of IN is odd
	// but this is the recommended way: http://jmoiron.github.io/sqlx/#namedParams
	query, args, err := sqlx.In(sql, kws)
	if err != nil {
		return nil, err
	}

	query = tx.Rebind(query)
	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kwmap := keywordsMap{}
	for rows.Next() {
		var (
			id   int
			name string
		)
		if err = rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		kwmap[Keyword(name)] = id
	}

	for _, kw := range keywords {
		if _, ok := kwmap[kw]; !ok {
			kwmap[kw] = 0
		}
	}

	return kwmap, nil
}

func createKeyword(tx *sqlx.Tx, kw Keyword) (int, error) {
	sql := `INSERT INTO keywords (name) VALUES (?)`
	res, err := tx.Exec(sql, string(kw))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func saveKwAssociations(tx *sqlx.Tx, bookmarkID int, kwIDs []int) error {
	if len(kwIDs) == 0 {
		// nothing left to do
		return nil
	}

	// let's build a query for bulk insert
	sql := `INSERT INTO bookmark_keywords (bookmark_id, keyword_id) VALUES %s`
	placeHolders := make([]string, 0, len(kwIDs))
	args := make([]interface{}, 0, len(kwIDs)*2)
	for _, id := range kwIDs {
		placeHolders = append(placeHolders, "(?, ?)")
		args = append(args, bookmarkID, id)
	}
	_, err := tx.Exec(fmt.Sprintf(sql, strings.Join(placeHolders, ",")), args...)
	return err
}

func deleteKwAssociations(tx *sqlx.Tx, bookmarkID int) error {
	sql := `DELETE FROM bookmark_keywords WHERE bookmark_id = ?`
	_, err := tx.Exec(sql, bookmarkID)
	return err
}
