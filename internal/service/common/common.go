package common

import (
	"context"
	"dndroller/internal/model"
	"dndroller/internal/repo"
	"dndroller/internal/repo/user"
)

func GetSets(rp *repo.Repo, id int) ([]*repo.Set, error) {
	ctx := context.Background()
	user, err := rp.User.Query().Where(user.ChatID(int64(id))).Only(ctx)
	if repo.IsNotFound(err) {
		user, err = rp.User.Create().SetChatID(int64(id)).Save(ctx)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return user.QuerySet().All(ctx)
}

func CreateSet(repo *repo.Repo, id int) (*repo.Set, error) {
	ctx := context.Background()
	user, err := repo.User.Query().Where(user.ChatID(int64(id))).Only(ctx)
	if err != nil {
		return nil, err
	}
	return repo.Set.Create().SetOwner(user).SetData(model.NewDiceSetEmpty()).Save(ctx)
}

func DeleteSet(repo *repo.Repo, id int) (*repo.Set, error) {
	ctx := context.Background()
	set, err := repo.Set.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	err = repo.Set.DeleteOneID(set.ID).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return set, nil
}
