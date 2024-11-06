package sharedspace

import (
	"context"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	user "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
)

func (db *sharedSpaceRepo) GetAllByUserID(ctx context.Context, userID int64) ([]model.SharedSpace, error) {
	spaces := make([]model.SharedSpace, 0)

	rows, err := db.db.QueryContext(ctx, `select shared_spaces.shared_spaces.id, space_number, name, created
	from shared_spaces.shared_spaces
		join shared_spaces.shared_spaces_view on shared_spaces.shared_spaces_view.id = shared_spaces.shared_spaces.id
		where shared_spaces.shared_spaces.creator = (select id from users.users where tg_id = $1)
		order by created ASC;`, userID)

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

	if len(spaces) == 0 {
		return nil, api_errors.ErrSharedSpacesNotFound
	}

	for i, space := range spaces {
		participants, err := db.getAllParticipants(ctx, space.ID)
		if err != nil {
			return nil, err
		}

		spaces[i].Participants = participants
	}

	return spaces, nil
}

func (db *sharedSpaceRepo) getAllParticipants(ctx context.Context, spaceID int) ([]user.User, error) {
	res := make([]user.User, 0)

	rows, err := db.db.QueryContext(ctx, `select tg_id, username from shared_spaces.participants 
join users.users on users.users.id = shared_spaces.participants.user_id
where space_id = $1;`, spaceID)
	if err != nil {
		return nil, fmt.Errorf("error while getting all users for shared space from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := user.User{}

		err := rows.Scan(&user.TGID, &user.UsernameSQL)
		if err != nil {
			return nil, fmt.Errorf("error scanning user ID while searching all space's participants: %+v", err)
		}

		res = append(res, user)
	}

	return res, nil
}
