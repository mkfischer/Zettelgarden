package srs

import (
	"fmt"
	"go-backend/models"
	"go-backend/server"
	"log"
)

func InitCardAsFlashcard(s *server.Server, userID, cardPK int) error {
	query := `
        UPDATE cards SET flashcard_state = 'new', flashcard_due = NOW()
        WHERE id = $1 AND user_id = $2 AND is_flashcard = TRUE
        `
	_, err := s.DB.Exec(query, cardPK, userID)
	if err != nil {
		log.Printf("init flashcard err %v", err)
		return err
	}
	return nil
}

func GetFlashcard(s *server.Server, userID, cardPK int) (models.Flashcard, error) {

	var card models.Flashcard

	err := s.DB.QueryRow(`
	SELECT
	id, card_id, user_id, title, body, created_at, updated_at,
        flashcard_state, flashcard_reps, flashcard_lapses, flashcard_last_review,
        flashcard_due
	FROM cards 
	WHERE is_deleted = FALSE AND card_id = $1 AND user_id = $2
	`, cardPK, userID).Scan(
		&card.ID,
		&card.CardID,
		&card.UserID,
		&card.Title,
		&card.Body,
		&card.CreatedAt,
		&card.UpdatedAt,
		&card.State,
		&card.Reps,
		&card.Lapses,
		&card.LastReview,
		&card.Due,
	)
	if err != nil {
		log.Printf("err %v", err)
		return models.Flashcard{}, fmt.Errorf("something went wrong")
	}
	return card, nil
}

func Next(s *server.Server, userID int) (int, error) {
	var count int
	var result int

	err := s.DB.QueryRow(`
        SELECT count(id) FROM cards
        WHERE is_flashcard = TRUE and user_id = $1 AND 
        DATE(flashcard_due) <= DATE(NOW())
`, userID).Scan(&count)
	if err != nil {
		log.Printf("get next flashcard err %v", err)
		return 0, fmt.Errorf("unable to access card")
	}
	if count < 1 {
		return -1, nil
	}

	err = s.DB.QueryRow(`
        SELECT id FROM cards
        WHERE is_flashcard = TRUE AND user_id = $1 AND
        DATE(flashcard_due) <= DATE(NOW())
        ORDER BY RANDOM()
        LIMIT 1
	`, userID).Scan(&result)
	if err != nil {
		log.Printf("get next flashcard err %v", err)
		return 0, fmt.Errorf("unable to access card")
	}

	return result, nil

}

func getNewFlashcardState(card models.Flashcard, rating models.Rating) (string, error) {
	if card.State == "New" || card.State == "" {
		return "Learning", nil

	} else if card.State == "Learning" || card.State == "Relearning" {
		if rating == models.Again {
			return "Learning", nil
		} else {
			return "Review", nil
		}
	} else if card.State == "Review" {
		if rating == models.Again {
			return "Relearning", nil
		} else {
			return "Review", nil
		}

	} else {
		return "", fmt.Errorf("card state inconsistent, got %v", card.State)

	}
}

func UpdateCardFromReview(s *server.Server, userID int, cardPK int, rating models.Rating) error {

	flashcard, err := GetFlashcard(s, userID, cardPK)
	if err != nil {
		log.Printf("get flashcard error %v", err)
		return err
	}
	flashcard_due := 1
	flashcard_reps_inc := 1
	flashcard_lapses_inc := 0

	log.Printf("tee hee")
	if rating == models.Again {
		log.Printf("?")
		flashcard_due = 0
		flashcard_lapses_inc = 1
	}
	flashcard_state, err := getNewFlashcardState(flashcard, rating)
	if err != nil {
		log.Printf("get flashcard error %v", err)
		return err
	}

	query := `
        UPDATE cards
        SET
          flashcard_due = NOW() + MAKE_INTERVAL(days => $1),
          flashcard_reps = flashcard_reps + $2,
flashcard_lapses = flashcard_lapses + $3,
          flashcard_last_review = NOW(),
flashcard_state = $4

        WHERE id = $5 AND user_id = $6
        `
	_, err = s.DB.Exec(
		query,
		flashcard_due,
		flashcard_reps_inc,
		flashcard_lapses_inc,
		flashcard_state,
		cardPK,
		userID,
	)
	if err != nil {
		log.Printf("record flashcard review err %v", err)
		return fmt.Errorf("unable to update card: %v", err.Error())

	}
	return nil
}

func RecordFlashcardReview(s *server.Server, userID int, cardPK int, rating models.Rating) error {
	query := `INSERT INTO flashcard_reviews (card_pk, user_id, rating) VALUES ($1, $2, $3)`

	_, err := s.DB.Exec(query, cardPK, userID, rating)
	if err != nil {
		log.Printf("record flashcard review err %v", err)
		return fmt.Errorf("unable to record flashcard review: %v", err.Error())

	}
	return nil
}
