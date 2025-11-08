package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/errybase/bar"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

var (
	sqldb, _ = sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	db       = bun.NewDB(sqldb, sqlitedialect.New())
	ctx      = context.Background()
)

func init() {
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	db.ResetModel(ctx,
		(*User)(nil),
		(*Story)(nil),
	)
}

func main() {
	// create users
	users := initUsers()
	if err := bar.Model(&users).Create(ctx, db); err != nil {
		panic(err)
	} else {
		fmt.Println("created users:", users)
	}
	user1, user2 := users[0], users[1]

	// create stories
	stories := initStories(users)
	if err := bar.Model(&stories).Create(ctx, db); err != nil {
		panic(err)
	} else {
		fmt.Println("created stories:", stories)
	}
	story1, story2 := stories[0], stories[1]

	// get story of user 1
	if s, err := user1.RelStory().Get(ctx, db); err != nil {
		panic(err)
	} else {
		fmt.Println("story of user 1:", s)
	}

	// set story of user 1 to story 2
	if err := user1.RelStory().Set(ctx, db, &story2); err != nil {
		panic(err)
	} else {
		fmt.Println("updated story of user 1:", story2)
	}

	// create new story for user 1
	story3 := Story{Title: "Fun story"}
	if err := user1.RelStory().Create(ctx, db, &story3); err != nil {
		panic(err)
	} else {
		fmt.Println("new story of user 1:", story3)
	}

	// update story 3 of user 1
	story3.Title = "Love story"
	if err := user1.RelStory().Update(ctx, db, &story3); err != nil {
		panic(err)
	} else {
		fmt.Println("updated story 3 of user 1:", story3)
	}

	// delete story 3 of user 1
	if err := user1.RelStory().Delete(ctx, db, &story3); err != nil {
		panic(err)
	} else {
		fmt.Println("deleted story 3 of user 1:", story3)
	}

	// get associated author of story 1
	if u, err := story1.RelAuthor().Get(ctx, db); err != nil {
		panic(err)
	} else {
		fmt.Println("author of story 1:", u)
	}

	// set author of story 1 to user 2
	if err := story1.RelAuthor().Set(ctx, db, user2); err != nil {
		panic(err)
	} else {
		fmt.Println("updated story 1:", story1)
	}

	// create new user associated with story
	story1, user3 := stories[0], User{Name: "member"}
	if err := story1.RelAuthor().Create(ctx, db, &user3); err != nil {
		panic(err)
	} else {
		fmt.Printf("updated story: %s; created user: %s", story1, user3)
	}
}

func initUsers() []User {
	return []User{
		{
			Name:   "admin",
			Emails: []string{"admin1@admin", "admin2@admin"},
		},
		{
			Name:   "root",
			Emails: []string{"root1@root", "root2@root"},
		},
	}
}

func initStories(users []User) []Story {
	return []Story{
		{
			Title:    "Cool story",
			AuthorID: users[0].ID,
		},
		{
			Title:    "Sad story",
			AuthorID: users[1].ID,
		},
	}
}

type User struct {
	bar.BaseModel

	Name   string
	Emails []string

	Story *Story `bun:"rel:has-one,join:id=author_id"`
}

func (u *User) RelStory() bar.HasOne[Story] {
	return bar.HasOne[Story]{
		Model:        u,
		RelationName: "Story",
	}
}

func (u User) String() string {
	return fmt.Sprintf("User<id=%d, name=%s, emails=%v, createdAt=%v, updated=%v>", u.ID, u.Name, u.Emails, u.CreatedAt, u.UpdatedAt)
}

type Story struct {
	bar.BaseModel

	Title    string
	AuthorID int64
	Author   User `bun:"rel:belongs-to,join:author_id=id"`
}

func (s *Story) RelAuthor() bar.BelongsTo[User] {
	return bar.BelongsTo[User]{
		Model:        s,
		RelationName: "Author",
	}
}

func (s Story) String() string {
	return fmt.Sprintf("Story<id=%d, title=%s, authorId=%d, createdAt=%v, updatedAt=%v>", s.ID, s.Title, s.AuthorID, s.CreatedAt, s.UpdatedAt)
}
