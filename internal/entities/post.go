package entities

type Post struct {
	ID        uint64 `db:"id"`
	Content   string `db:"content"`
	Title     string `db:"title"`
	UserID    uint64 `db:"user_id"`
	CreatedAt string `db:"created_at"`
}
