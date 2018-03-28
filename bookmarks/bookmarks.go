package bookmarks

import (
	"time"

	"github.com/fchoquet/bookmarks/oembed"
	"github.com/fchoquet/bookmarks/pager"
	"github.com/jmoiron/sqlx"
	"gopkg.in/go-playground/validator.v9"
)

// Bookmark represents a bookmark
type Bookmark struct {
	ID int `json:"id;omitempty" db:"id"`

	// Shared properties
	// Some properties are required at the domain level but not in the API
	// Missing properties are populated using oEmbed
	URL        string     `json:"url" db:"url" validate:"required,max=255"`
	Title      string     `json:"title" db:"title" validate:"required,max=100"`
	AuthorName string     `json:"author_name" db:"author_name" validate:"required,max=100"`
	AddedDate  *time.Time `json:"added_date" db:"added_date"`

	// These properties are specific to the target link and might not make sense
	// in a given context. Let's keep them in a single struct as long as their
	// number is not overwhelming. They don't worth any extra complexity for now
	Width    int `json:"width;omitempty" db:"width"`
	Height   int `json:"height;omitempty" db:"height"`
	Duration int `json:"duration;omitempty" db:"duration"`

	Keywords []Keyword `json:"keywords"`

	// TODO: Bookmarks should be attached to a user. I'm not sure if it's in the
	// scope of this exercise
	UserID int
}

// Repository stores bookmarks to a permanent storage
type Repository interface {
	// List returns a list of bookmarks. Can be filtered.
	// It also returns the total number of bookmarks (useful with pagination)
	List(fitler Filter) ([]*Bookmark, int, error)

	// Load loads a unique bookmark by its URL. returns nil if not found
	ByID(id int) (*Bookmark, error)

	// Insert creates a new bookmark. Returns an error if already exists
	Insert(*Bookmark) (*Bookmark, error)

	// Update updates an existing bookmark's keywords
	UpdateKeywords(id int, keywords []Keyword) error

	// Delete delets an existing bookmark
	Delete(id int) error
}

// Filter allows filtering of Bookmarks
type Filter struct {
	ID    *int
	Pager pager.Pager
}

// NewRepository returns a default Repository implementation
// Let's not use anything more fancy than sqlx
// Raw sql is enough given the extreme simplicity of the queries
func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

type repository struct {
	db *sqlx.DB
}

func (rep *repository) List(filter Filter) ([]*Bookmark, int, error) {
	sql := `SELECT * FROM bookmarks`
	args := struct {
		ID int
	}{}

	if filter.ID != nil {
		sql += ` WHERE id = :id`
		args.ID = *filter.ID
	}

	if filter.Pager == nil {
		filter.Pager = pager.NoPager()
	}

	rows, err := rep.db.NamedQuery(sql, args)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	bookmarks := []*Bookmark{}
	index := -1
	for rows.Next() {
		index++
		if !filter.Pager.IsVisible(index) {
			continue
		}

		var b Bookmark
		if err = rows.StructScan(&b); err != nil {
			return nil, 0, err
		}

		keywords, err := loadKeywords(rep.db, b.ID)
		if err != nil {
			return nil, 0, err
		}

		b.Keywords = keywords

		bookmarks = append(bookmarks, &b)
	}

	return bookmarks, index + 1, nil
}

func (rep *repository) ByID(id int) (*Bookmark, error) {
	bookmarks, _, err := rep.List(Filter{ID: &id})
	if err != nil || len(bookmarks) == 0 {
		return nil, err
	}

	return bookmarks[0], nil
}

func (rep *repository) Insert(b *Bookmark) (*Bookmark, error) {
	if err := validator.New().Struct(b); err != nil {
		return nil, err
	}

	// generates an added date if none passed
	if b.AddedDate == nil {
		now := time.Now()
		b.AddedDate = &now
	}

	// start a transaction
	tx, err := rep.db.Beginx()
	if err != nil {
		return nil, err
	}

	newB, err := insert(tx, b)
	if err != nil {
		// there's little we can do if Rollback fails, so let's ignore this case
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return newB, nil
}

func insert(tx *sqlx.Tx, b *Bookmark) (*Bookmark, error) {
	// the primary key on url will ensure that the record does not exist
	sql := `
INSERT INTO bookmarks (
    url, title, author_name, added_date, width, height, duration
) VALUES (
    :url, :title, :author_name, :added_date, :width, :height, :duration
)
`
	res, err := tx.NamedExec(sql, b)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	b.ID = int(id)

	// Now insert keywords
	if err := saveKeywords(tx, b.ID, b.Keywords); err != nil {
		return nil, err
	}

	return b, nil
}

func (rep *repository) UpdateKeywords(id int, keywords []Keyword) error {
	tx, err := rep.db.Beginx()
	if err != nil {
		return err
	}

	if err := saveKeywords(tx, id, keywords); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (rep *repository) Delete(id int) error {
	tx, err := rep.db.Beginx()
	if err != nil {
		return err
	}

	if err := delete(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func delete(tx *sqlx.Tx, id int) error {
	if err := deleteKwAssociations(tx, id); err != nil {
		return err
	}

	sql := `DELETE FROM bookmarks WHERE id = ?`
	_, err := tx.Exec(sql, id)
	return err
}

// FromOembed decorates a bookmark with oEmbed information
// it does not overwrite existing properties
// this is obvioulsy arguable, but I'm not sure of the expectation here
func FromOembed(b *Bookmark, link *oembed.Link) *Bookmark {
	if link == nil {
		return b
	}

	if b.Title == "" {
		b.Title = link.Title
	}
	if b.AuthorName == "" {
		b.AuthorName = link.AuthorName
	}
	if b.Width == 0 {
		b.Width = int(link.Width)
	}
	if b.Height == 0 {
		b.Height = int(link.Height)
	}
	if b.Duration == 0 {
		b.Duration = link.Duration
	}

	return b
}
