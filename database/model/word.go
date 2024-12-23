package model

import (
	"context"
	"database/sql"
)

type WordDto struct {
	ID            int64          `json:"id"`
	Word          string         `json:"word"`
	Transcription sql.NullString `json:"transcription"`
	Meaning       sql.NullString `json:"meaning"`
	Example       sql.NullString `json:"example"`
	WordLevel     sql.NullString `json:"word_level"`
	Translations  sql.NullString `json:"translations"`
	Frequency     sql.NullInt16  `json:"frequency"`
}

type WordStorage struct {
	db *sql.DB
}

func (s *WordStorage) GetWord(ctx context.Context, text string) (*WordDto, *DatabaseError) {
	// Implement getting logic here
	query := `SELECT id, word, transcription, meaning, example, word_level, translation FROM words WHERE word = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	wordDto := &WordDto{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		text,
	).Scan(
		&wordDto.ID,
		&wordDto.Word,
		&wordDto.Transcription,
		&wordDto.Meaning,
		&wordDto.Example,
		&wordDto.WordLevel,
		&wordDto.Translations,
	)

	if err != nil {
		return nil, ProcessErrorFromDatabase(err)
	}

	return wordDto, nil
}

func (s *WordStorage) SaveWords(ctx context.Context, wordDto *WordDto) (*WordDto, *DatabaseError) {
	// Implement saving logic here
	query := `INSERT INTO words (word) VALUES ($1) RETURNING id`
	ctx, cancel := context.WithTimeout(ctx, QueryRowTimeout)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		wordDto.Word,
	).Scan(
		&wordDto.ID,
	)

	if err != nil {
		return nil, ProcessErrorFromDatabase(err)
	}

	return wordDto, nil
}

func (s *WordStorage) SaveWordWithBookConnection(ctx context.Context, book *BookDto, wordDto *WordDto) *DatabaseError {
	query := `INSERT INTO books_words (book_id, word_id, frequency) VALUES ($1, $2, $3)`

	_, err := s.db.ExecContext(
		ctx,
		query,
		book.ID,
		wordDto.ID,
		book.WordMap[wordDto.Word],
	)

	if err != nil {
		return ProcessErrorFromDatabase(err)
	}

	return nil
}