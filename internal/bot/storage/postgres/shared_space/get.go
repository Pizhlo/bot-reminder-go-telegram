package sharedspace

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	user "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
)

func (db *sharedSpaceRepo) GetAllByUserID(ctx context.Context, userID int64) ([]model.SharedSpace, error) {
	spaces := make([]model.SharedSpace, 0)

	rows, err := db.db.QueryContext(ctx, `select shared_spaces.shared_spaces.id, space_number, name, created
	from shared_spaces.shared_spaces
		join users.users on users.id = shared_spaces.shared_spaces.creator
		join shared_spaces.shared_spaces_view on shared_spaces.shared_spaces_view.id = shared_spaces.shared_spaces.id
		where shared_spaces.shared_spaces.creator = (select id from users.users where tg_id = $1)
		order by created ASC`, userID)

	if err != nil {
		return nil, fmt.Errorf("error while getting all shared spaces from DB by user ID %d: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		space := model.SharedSpace{}

		err := rows.Scan(&space.ID, &space.ViewID, &space.Name, &space.Created)
		if err != nil {
			return nil, fmt.Errorf("error scanning shared space: %+v", err)
		}

		spaces = append(spaces, space)
	}

	for i, space := range spaces {
		participants, err := db.getAllParticipants(ctx, space.ID)
		if err != nil {
			return nil, err
		}

		users := make([]user.User, 0)

		for _, u := range participants {
			users = append(users, user.User{
				TGID: u,
			})
		}

		spaces[i].Participants = users
	}

	return spaces, nil
}

func (db *sharedSpaceRepo) getAllParticipants(ctx context.Context, spaceID int) ([]int64, error) {
	res := make([]int64, 0)

	rows, err := db.db.QueryContext(ctx, "select user_id from shared_spaces.participants where space_id = $1", spaceID)
	if err != nil {
		return nil, fmt.Errorf("error while getting all users for shared space from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		id := int64(0)
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error scanning user ID while searching all space's participants: %+v", err)
		}

		res = append(res, id)
	}

	return res, nil
}
