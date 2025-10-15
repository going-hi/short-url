package link

import (
	"database/sql"
)

type LinkRepository struct {
	Db *sql.DB
}

func (rep *LinkRepository) Create(url, code string, userId int) (*Link, error) {
	query := `INSERT INTO (url, code, userId) VALUES ($1, $2, $3) RETURNING id, url, code`

	link := &Link{}
	err := rep.Db.QueryRow(query, url, code, userId).Scan(&link.Id, &link.Url, &link.Code)

	if err != nil {
		return nil, err
	}

	return link, nil
}

func (rep *LinkRepository) FindByCode(code string) (*Link, error) {
	query := `SELECT id, code, url, clicks, userId FROM links WHERE code = $1 LIMIT 1`

	link := &Link{}
	err := rep.Db.QueryRow(query, code).Scan(
		&link.Id,
		&link.Code,
		&link.Url,
		&link.UserId,
	)

	if err != nil {
		return nil, err
	}

	return link, nil
}

func (rep *LinkRepository) FindById(id int) (*Link, error) {
	query := `SELECT id, code, url, clicks, userId FROM links WHERE id = $1 LIMIT 1`

	link := &Link{}
	err := rep.Db.QueryRow(query, id).Scan(
		&link.Id,
		&link.Code,
		&link.Url,
		&link.UserId,
	)

	if err != nil {
		return nil, err
	}

	return link, nil
}

func (rep *LinkRepository) UpdateClick(id int) error {
	query := `UPDATE links SET clicks = clicks + 1 WHERE id = $1`

	_, err := rep.Db.Exec(query, id)
	return err
}

func (rep *LinkRepository) Delete(id int) error {
	query := `DELETE FROM links WHERE id = $1`

	_, err := rep.Db.Exec(query, id)
	return err
}

func (rep *LinkRepository) FindAllByUserId(userId int) ([]*Link, error) {
	query := `SELECT id, code, url, clicks, userId FROM links WHERE userId = $1`

	rows, err := rep.Db.Query(query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var links []*Link

	for rows.Next() {
		l := &Link {}
		err := rows.Scan(&l.Id, &l.Code, &l.Url, &l.Clicks, &l.UserId)
		if err != nil {
			return nil, err
		}
		links = append(links, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}
