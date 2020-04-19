package main

import (
	"github.com/c8112002/news-crawler/db"
	"github.com/c8112002/news-crawler/db/seeds"
	"github.com/c8112002/news-crawler/utils"
)


func main()  {

	if err := utils.LoadEnv(); err != nil {
		panic(err.Error())
	}

	d, err := db.New(utils.GetEnv())
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		if err := d.Close(); err != nil {
			panic(err.Error())
		}
	}()

	seeds.MakeSites(d)
	seeds.MakeTags(d)
}