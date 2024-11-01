package modelentities

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Blog struct {
	Id        pgtype.Int4
	Title     pgtype.Text
	Content   pgtype.Text
	Category  pgtype.Text
	Tags      pgtype.Text
	CreatedAt pgtype.Int8
	UpdatedAt pgtype.Int8
}
