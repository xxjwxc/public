package myelastic

import (
	"context"
	"encoding/json"
	"log"

	"github.com/xxjwxc/public/errors"

	"reflect"
	"strings"
	"time"

	"github.com/olivere/elastic"
)

//
type MyElastic struct {
	Client *elastic.Client
	Err    error
	Ctx    context.Context
}

//
func OnInitES(url string) MyElastic {
	var es MyElastic
	es.Ctx = context.Background()
	es.Client, es.Err = elastic.NewClient(elastic.SetURL(url))
	if es.Err != nil {
		log.Println(es.Err)
		//mylog.Error(es.Err)
		//panic(es.Err)
	}

	return es
}

//func (es *MyElastic) Model(refs interface{}) *MyElastic {
//	if reflect.ValueOf(refs).Type().Kind() != reflect.Ptr {
//		mylog.Println("Model: attempt to Model into a pointer")
//		panic(0)
//	}

//	es.Element = refs
//	return es
//}

/*
创建索引（相当于数据库）
mapping 如果为空("")则表示不创建模型
*/
func (es *MyElastic) CreateIndex(index_name, mapping string) (result bool) {
	es.Err = nil
	exists, err := es.Client.IndexExists(index_name).Do(es.Ctx)
	if err != nil {
		es.Err = err
		log.Println(es.Err)
		return false
	}

	if !exists {
		var re *elastic.IndicesCreateResult
		if len(mapping) == 0 {
			re, es.Err = es.Client.CreateIndex(index_name).Do(es.Ctx)
		} else {
			re, es.Err = es.Client.CreateIndex(index_name).BodyString(mapping).Do(es.Ctx)
		}

		if es.Err != nil {
			log.Println(es.Err)
			return false
		}

		return re.Acknowledged
	}

	return false
}

/*
	排序查询
	返回json数据集合
*/
func (es *MyElastic) SortQuery(index_name string, builder []elastic.Sorter, query []elastic.Query) (bool, []string) {

	searchResult := es.Client.Search().Index(index_name)

	if len(builder) > 0 {
		for _, v := range builder {
			searchResult = searchResult.SortBy(v)
		}
	}
	if len(query) > 0 {
		for _, v := range query {
			searchResult = searchResult.Query(v)
		}
	}
	es_result, err := searchResult.Do(es.Ctx) // execute
	if err != nil {
		log.Println(es.Err)
		return false, nil
	}
	//log.Println("Found a total of %d entity\n", es_result.TotalHits())

	if es_result.Hits.TotalHits > 0 {
		var result []string
		//log.Println("Found a total of %d entity\n", searchResult.Hits.TotalHits)
		for _, hit := range es_result.Hits.Hits {

			result = append(result, string(*hit.Source))

		}
		return true, result
	} else {
		// No hits
		return true, nil
	}
}

/*
   排序查询
   返回原始Hit
   builder：排序
   agg：聚合 类似group_by sum
   query：查询
*/
func (es *MyElastic) SortQueryReturnHits(index_name string, from, size int, builder []elastic.Sorter, query []elastic.Query) (bool, []*elastic.SearchHit) {

	searchResult := es.Client.Search().Index(index_name)

	if len(builder) > 0 {
		for _, v := range builder {
			searchResult = searchResult.SortBy(v)
		}
	}
	if len(query) > 0 {
		for _, v := range query {
			searchResult = searchResult.Query(v)
		}
	}
	if size > 0 {
		searchResult = searchResult.From(from)
		searchResult = searchResult.Size(size)
	}
	es_result, err := searchResult.Do(es.Ctx) // execute
	if err != nil {
		log.Println(es.Err)
		return false, nil
	}

	//	log.Println("wwwwww", es_result.Aggregations)
	if es_result.Hits.TotalHits > 0 {
		return true, es_result.Hits.Hits
	} else {
		return true, nil
	}
}

/*
添加记录,覆盖添加
*/
func (es *MyElastic) Add(index_name, type_name, id string, data interface{}) (result bool) {
	result = false
	// Index a tweet (using JSON serialization)
	if len(id) > 0 {
		_, es.Err = es.Client.Index().
			Index(index_name).
			Type(type_name).
			Id(id).
			BodyJson(data).
			Do(es.Ctx)
	} else {
		_, es.Err = es.Client.Index().
			Index(index_name).
			Type(type_name).
			BodyJson(data).
			Do(es.Ctx)
	}

	if es.Err != nil {
		log.Println(es.Err)
		return false
	}
	_, es.Err = es.Client.Flush().Index(index_name).Do(es.Ctx)
	if es.Err != nil {
		log.Println(es.Err)
		return false
	}
	return true
}

/*
添加记录,覆盖添加
index_name
type_name
query interface{} //查询条件
out *[]Param //查询结果
*/
func (es *MyElastic) SearchMap(index_name, type_name string, query interface{}, out *[]map[string]interface{}) (result bool) {
	es_search := es.Client.Search()
	if len(type_name) > 0 {
		es_search = es_search.Type(type_name)
	}
	if len(index_name) > 0 {
		es_search = es_search.Index(index_name)
	}
	var es_result *elastic.SearchResult
	es_result, es.Err = es_search.Source(query).Do(es.Ctx)
	if es.Err != nil {
		log.Println(es.Err)
		return false
	}
	if es_result.Hits == nil {
		log.Println(errors.New("expected SearchResult.Hits != nil; got nil"))
		return false
	}

	for _, hit := range es_result.Hits.Hits {
		tmp := make(map[string]interface{})
		err := json.Unmarshal(*hit.Source, &tmp)
		if err != nil {
			log.Println(es.Err)
		} else {
			*out = append(*out, tmp)
		}
	}

	return true
}

/*
添加记录,覆盖添加
index_name
type_name
query interface{} //查询条件
out *[]Param //查询结果
*/
func (es *MyElastic) Search(index_name, type_name string, query interface{}, out interface{}) (result bool) {

	sliceValue := reflect.Indirect(reflect.ValueOf(out))

	if sliceValue.Kind() != reflect.Slice {
		log.Println(errors.New("needs a pointer to a slice"))
		return false
	}

	sliceElementType := sliceValue.Type().Elem()
	es_search := es.Client.Search()
	if len(type_name) > 0 {
		es_search = es_search.Type(type_name)
	}
	if len(index_name) > 0 {
		es_search = es_search.Index(index_name)
	}
	var es_result *elastic.SearchResult
	es_result, es.Err = es_search.Source(query).Do(es.Ctx)
	if es.Err != nil {
		log.Println(es.Err)
		return false
	}
	if es_result.Hits == nil {
		log.Println(errors.New("expected SearchResult.Hits != nil; got nil"))
		return false
	}

	for _, hit := range es_result.Hits.Hits {
		newValue := reflect.New(sliceElementType)

		item := make(map[string]interface{})
		err := json.Unmarshal(*hit.Source, &item)
		//fmt.Println(string(*hit.Source))

		err = scanMapIntoStruct(newValue.Interface(), item)
		if err != nil {
			log.Println(err)
		}

		if err != nil {
			log.Println(err)
		} else {
			sliceValue.Set(reflect.Append(sliceValue, reflect.Indirect(reflect.ValueOf(newValue.Interface()))))
			//out = append(out, tmp)
		}
	}

	return true
}

func scanMapIntoStruct(obj interface{}, objMap map[string]interface{}) error {
	dataStruct := reflect.Indirect(reflect.ValueOf(obj))
	if dataStruct.Kind() != reflect.Struct {
		return errors.New("expected a pointer to a struct")
	}

	dataStructType := dataStruct.Type()

	for i := 0; i < dataStructType.NumField(); i++ {
		field := dataStructType.Field(i)
		fieldv := dataStruct.Field(i)

		err := scanMapElement(fieldv, field, objMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func scanMapElement(fieldv reflect.Value, field reflect.StructField, objMap map[string]interface{}) error {

	objFieldName := field.Name
	bb := field.Tag
	sqlTag := bb.Get("json")

	if sqlTag == "-" || reflect.ValueOf(bb).String() == "-" {
		return nil
	}

	sqlTags := strings.Split(sqlTag, ",")
	sqlFieldName := objFieldName
	if len(sqlTags[0]) > 0 {
		sqlFieldName = sqlTags[0]
	}

	data, ok := objMap[sqlFieldName]
	if !ok || data == nil {
		return nil
	}

	//	fmt.Println("================")
	//	fmt.Println(field.Type.Kind())
	//	fmt.Println(sqlFieldName)
	var v interface{}
	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x := int(data.(float64))
		v = x
	case reflect.Slice:
		if fieldv.Type().String() == "[]uint8" {
			x := []byte(data.(string))
			v = x
		} else if fieldv.Type().String() == "[]string" {
			mp := data.([]interface{})
			var ss []string
			for _, v := range mp {
				ss = append(ss, v.(string))
			}
			v = ss
		} else if fieldv.Type().String() == "[]int" {
			mp := data.([]interface{})
			var ss []int
			for _, v := range mp {
				ss = append(ss, int(v.(float64)))
			}
			v = ss

		} else {
			v = data
		}

	case reflect.Struct:
		if fieldv.Type().String() == "time.Time" {
			x, err := time.Parse("2006-01-02 15:04:05", data.(string))
			if err != nil {
				x, err = time.Parse("2006-01-02 15:04:05.000 -0700", data.(string))
				if err != nil {
					if err != nil {
						x, err = time.Parse("2006-01-02T15:04:05.999999999Z07:00", data.(string))
						if err != nil {
							return errors.New("unsupported time format: " + data.(string))
						}
					}
				}
			}

			v = x
		} else {
			v = data
		}

	default:
		v = data
	}

	fieldv.Set(reflect.ValueOf(v))
	//	fmt.Println("================")

	return nil
}
