package show

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/iooikaak/frame/helper"

	"github.com/iooikaak/microService1/database/es"

	"github.com/iooikaak/frame/util"
	"github.com/iooikaak/frame/xlog"
	"github.com/olivere/elastic/v7"
)

func (d *Dao) InsertShowList(ctx context.Context, show *es.ShowListEsModel) (err error) {
	err = helper.Retry5BackOff(func() error {
		_, err = d.es.Index().Index(es.ShowListElasticSearchIndex).Type(es.ShowListElasticSearchType).
			BodyJson(show).Refresh("true").Do(ctx)
		return err
	})
	if err != nil {
		xlog.Errorf("InsertShowList failed err: %v", err)
	}
	return
}

func (d *Dao) UpdateShowListByQueryAttr(ctx context.Context, stdShowId string, lastShowTime time.Time) (err error) {
	q := elastic.NewTermQuery("std_show_id", stdShowId)
	lastShowTimeStr := lastShowTime.Format(util.TimeLayout)
	scriptArr := []string{
		"ctx._source.last_show_time=params.last_show_time",
	}
	params := map[string]interface{}{
		"last_show_time": lastShowTimeStr,
	}
	err = helper.Retry5BackOff(func() error {
		_, err = d.es.UpdateByQuery().Index(es.ShowListElasticSearchIndex).Type(es.ShowListElasticSearchType).Query(q).
			Refresh("true").
			Script(elastic.NewScriptInline(strings.Join(scriptArr, ";")).Params(params)).ProceedOnVersionConflict().Do(ctx)
		return err
	})
	if err != nil {
		return
	}
	return
}

func (d *Dao) GetShowList(ctx context.Context, id string) (showList *es.ShowListEsModel, err error) {
	r, err := d.es.Get().Index(es.ShowListElasticSearchIndex).Type(es.ShowListElasticSearchType).Id(id).Do(ctx)
	showList = &es.ShowListEsModel{}
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return showList, nil
		}
		return
	}
	if r == nil || !r.Found {
		return showList, nil
	}

	data, err := r.Source.MarshalJSON()
	if err != nil {
		return
	}
	err = json.Unmarshal(data, showList)
	return
}

func (d *Dao) GetShowLists(ctx context.Context, params *es.ShowListSearchEsModel) (data *elastic.SearchResult, err error) {
	boolQuery := elastic.NewBoolQuery()
	if params.IdentityRequiredType > 0 {
		boolQuery.Must(elastic.NewTermQuery("identity_required_type", params.IdentityRequiredType))
	}

	if len(params.BizArtistIds) > 0 {
		boolQuery.Must(elastic.NewMatchQuery("biz_artist_ids", params.BizArtistIds))
	}

	boolQuery.Must(elastic.NewTermQuery("is_show_std_content", params.IsShowStdContent))
	if len(params.LastShowTime) > 0 {
		boolQuery.Must(elastic.NewRangeQuery("last_show_time").Gte(params.LastShowTime))
	}

	query := d.es.Search().Index(es.ShowListElasticSearchIndex).
		Type(es.ShowListElasticSearchType).
		From(params.From).
		Size(params.PageSize)

	if len(params.SortBy) > 0 {
		ascendingDirection := true
		switch params.SortDirection {
		case "asc":
			ascendingDirection = true
		case "desc":
			ascendingDirection = false
		}
		query.Sort(params.SortBy, ascendingDirection)
	}
	data, err = query.From(params.From).Size(params.PageSize).Query(boolQuery).Do(ctx)
	if err != nil {
		return
	}
	if data == nil {
		return
	}
	return
}
