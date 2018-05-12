package main

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/icrowley/fake"
	"encoding/json"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	RealName string `json:"real_name"`
}

func Populate(number int) error {
	client, err := elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"))
	if err != nil {
		return err
	}

	idxExists, err := client.IndexExists("users").Do(context.Background())
	if err != nil {
		return err
	}
	if !idxExists {
		client.CreateIndex("users").Do(context.Background())
	}

	for i := 0; i < number; i++ {
		user := User{
			Username: fake.UserName(),
			Email: fake.EmailAddress(),
			RealName: fake.FullName(),
		}
		_, err = client.Index().
			Index("users").
			Type("doc").
			BodyJson(user).
			Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func Search(term string, from, size int) ([]*User, error) {
	client, err := elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"))
	if err != err {
		return nil, err
	}
	q := elastic.NewMultiMatchQuery(term, "username", "email", "real_name").Fuzziness("AUTO:2,5")
	res, err := client.Search().
		Index("users").
		Query(q).
		From(from).
		Size(size).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	users := make([]*User, 0)

	for _, hit := range res.Hits.Hits {
		var user User
		err := json.Unmarshal(*hit.Source, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
