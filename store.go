package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/btlike/repository"
	"github.com/xgfone/gobt/g"
	"gopkg.in/olivere/elastic.v3"
)

func GetTorrentByInfohashFromDB(client repository.Repository, infohash string) (repository.Torrent, error) {
	if t, err := client.GetTorrentByInfohash(infohash); err != nil || t.Infohash != infohash {
		return nil, err
	} else {
		return t, nil
	}
}

func GetTorrentByInfohashFromSE(client *elastic.Client, infohash string) (*TorrentSearch, error) {
	indexType := strings.ToLower(string(infohash[0]))
	result, err := client.Get().Index("torrent").Type(indexType).Id(infohash).Do()
	if err != nil || result == nil || !result.Found || result.Source == nil {
		return nil, fmt.Errorf("can't get the result")
	}

	var data TorrentSearch
	err = json.Unmarshal(*result.Source, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func SearchKeyword(client *elastic.Client, key string, from, size int) ([]repository.Torrent, int, error) {
	query := elastic.NewQueryStringQuery(key)
	result, err := client.Search().Index("torrent").Query(query).From(from).Size(size).Do()
	if err != nil {
		return nil, 0, err
	}

	results := make([]repository.Torrent, 0)
	if result.Hits.TotalHits > 0 {
		for _, hit := range result.Hits.Hits {
			if t, err := GetTorrentByInfohashFromDB(g.Repository, hit.Id); err == nil {
				results = append(results, t)
			}
		}
	}

	return results, int(result.Hits.TotalHits), nil
}
