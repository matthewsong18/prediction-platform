package bets

import (
	"database/sql"
	"errors"
	"fmt"
)

type libSQLRepository struct {
	db *sql.DB
}

func NewLibSQLRepository(db *sql.DB) BetRepository {
	return &libSQLRepository{
		db: db,
	}
}

func (repo libSQLRepository) Save(bet *bet) error {
	query := `INSERT INTO bets (poll_id, user_id, selected_option_index, bet_status) VALUES (?, ?, ?, ?)`

	preparedStatement, preparedErr := repo.db.Prepare(query)
	if preparedErr != nil {
		return fmt.Errorf("error while preparing save bet statement: %w", preparedErr)
	}

	_, execErr := preparedStatement.Exec(bet.PollID, bet.UserID, bet.SelectedOptionIndex, bet.BetStatus)
	if execErr != nil {
		return fmt.Errorf("error while executing save bet statement: %w", execErr)
	}

	return nil
}

func (repo libSQLRepository) GetByPollIdAndUserId(pollID string, userID string) (*bet, error) {
	query := "SELECT poll_id, user_id, selected_option_index, bet_status FROM bets WHERE poll_id = ? AND user_id = ?"
	preparedStatement, preparedErr := repo.db.Prepare(query)
	if preparedErr != nil {
		return nil, fmt.Errorf("error while preparing get bet by poll_id and user_id statement: %w", preparedErr)
	}

	var bet bet
	row := preparedStatement.QueryRow(pollID, userID)
	if scanErr := row.Scan(&bet.PollID, &bet.UserID, &bet.SelectedOptionIndex, &bet.BetStatus); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return nil, ErrBetNotFound
		}
		return nil, fmt.Errorf("error while scanning bet: %w", scanErr)
	}

	return &bet, nil
}

func (repo libSQLRepository) GetBetsFromUser(userID string) ([]*bet, error) {
	query := "SELECT poll_id, user_id, selected_option_index, bet_status FROM bets WHERE user_id = ?"
	preparedStatement, preparedErr := repo.db.Prepare(query)
	if preparedErr != nil {
		return nil, fmt.Errorf("error while preparing get bets from user statement: %w", preparedErr)
	}

	rows, queryErr := preparedStatement.Query(userID)
	if queryErr != nil {
		return nil, fmt.Errorf("error while querying bets from user: %w", queryErr)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmtErr := fmt.Errorf("failed to close bets rows: %w\n", err)
			fmt.Println(fmtErr)
		}
	}(rows)

	var bets []*bet
	for rows.Next() {
		var bet bet
		if scanErr := rows.Scan(&bet.PollID, &bet.UserID, &bet.SelectedOptionIndex, &bet.BetStatus); scanErr != nil {
			return nil, fmt.Errorf("error while scanning bet: %w", scanErr)
		}
		bets = append(bets, &bet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %w", err)
	}

	return bets, nil
}

func (repo libSQLRepository) GetBetsByPollId(pollID string) ([]*bet, error) {
	query := "SELECT poll_id, user_id, selected_option_index, bet_status FROM bets WHERE poll_id = ?"
	preparedStatement, preparedErr := repo.db.Prepare(query)
	if preparedErr != nil {
		return nil, fmt.Errorf("error while preparing get bets by poll_id")
	}

	rows, queryErr := preparedStatement.Query(pollID)
	if queryErr != nil {
		return nil, fmt.Errorf("error while querying bets by poll_id")
	}

	var bets []*bet
	for rows.Next() {
		var bet bet
		if scanErr := rows.Scan(&bet.PollID, &bet.UserID, &bet.SelectedOptionIndex, &bet.BetStatus); scanErr != nil {
			return nil, fmt.Errorf("error while scanning bet: %w", scanErr)
		}

		bets = append(bets, &bet)

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error while iterating over rows: %w", err)
		}
	}
	return bets, nil
}

func (repo libSQLRepository) UpdateBet(bet *bet) error {
	query := "UPDATE bets SET selected_option_index = ?, bet_status = ? WHERE poll_id = ? AND user_id = ?"
	preparedStatement, preparedErr := repo.db.Prepare(query)
	if preparedErr != nil {
		return fmt.Errorf("error while preparing update bet statement: %w", preparedErr)
	}

	result, execErr := preparedStatement.Exec(bet.SelectedOptionIndex, bet.BetStatus, bet.PollID, bet.UserID)
	if execErr != nil {
		return fmt.Errorf("error while executing update bet statement: %w", execErr)
	}

	rowsAffected, rowErr := result.RowsAffected()
	if rowErr != nil {
		return fmt.Errorf("error while getting rows affected: %w", rowErr)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were affected by the update operation")
	}

	return nil
}

func (repo *libSQLRepository) GetAllUserStats() ([]*UserStats, error) {
	query := `
		SELECT
			user_id,
			SUM(CASE WHEN bet_status = 1 THEN 1 ELSE 0 END) AS wins,
			SUM(CASE WHEN bet_status = 2 THEN 1 ELSE 0 END) AS losses
		FROM bets
		GROUP BY user_id;
	`
	rows, queryErr := repo.db.Query(query)
	if queryErr != nil {
		return nil, fmt.Errorf("error while executing get all user stats: %w", queryErr)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmtErr := fmt.Errorf("failed to close bets rows: %w", err)
			fmt.Println(fmtErr)
		}
	}(rows)

	var allUserStats []*UserStats
	for rows.Next() {
		userStat := &UserStats{}
		if scanErr := rows.Scan(&userStat.UserID, &userStat.Wins, &userStat.Losses); scanErr != nil {
			return nil, fmt.Errorf("error while scanning bet: %w", scanErr)
		}

		userStat.Total = userStat.Wins + userStat.Losses
		if userStat.Total > 0 {
			userStat.WinLossRatio = float64(userStat.Wins) / float64(userStat.Losses)
		} else {
			userStat.WinLossRatio = 0.0
		}

		allUserStats = append(allUserStats, userStat)

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error while iterating over rows: %w", err)
		}
	}

	return allUserStats, nil
}

var _ BetRepository = (*libSQLRepository)(nil)
