package main

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"strings"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Create: Index Document
	createDocument(es)

	// Create: Bulk Index Documents
	bulkIndexDocuments(es)

	// Read: Get Document
	getDocument(es)

	// Read: Search Documents
	searchDocuments(es)

	// Read: Count Documents
	countDocuments(es)

	// Update: Update Document
	updateDocument(es)

	// Update: Upsert Document
	upsertDocument(es)

	// Delete: Delete Document
	deleteDocument(es)

	// Delete: Delete by Query
	deleteByQuery(es)
}

func createDocument(es *elasticsearch.Client) {
	doc := `{"field": "value"}`
	res, err := es.Index(
		"index",
		strings.NewReader(doc),
		es.Index.WithDocumentID("id"),
		es.Index.WithRefresh("true"),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}

func bulkIndexDocuments(es *elasticsearch.Client) {
	bulk := strings.NewReader(`
    { "index": { "_id": "1" } }
    { "field": "value1" }
    { "index": { "_id": "2" } }
    { "field": "value2" }
    `)
	res, err := es.Bulk(bulk, es.Bulk.WithIndex("index"), es.Bulk.WithRefresh("true"))
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}

func getDocument(es *elasticsearch.Client) {
	res, err := es.Get("index", "id")
	if err != nil {
		log.Fatalf("Error getting document: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}

func searchDocuments(es *elasticsearch.Client) {
	query := `{
        "query": {
            "match": {
                "field": "value"
            }
        }
    }`
	res, err := es.Search(
		es.Search.WithIndex("index"),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error searching documents: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}

func countDocuments(es *elasticsearch.Client) {
	query := `{
        "query": {
            "match_all": {}
        }
    }`
	res, err := es.Count(
		es.Count.WithIndex("index"),
		es.Count.WithBody(strings.NewReader(query)),
		es.Count.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error counting documents: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}

func updateDocument(es *elasticsearch.Client) {
	update := `{
        "doc": {
            "field": "new_value"
        }
    }`
	res, err := es.Update(
		"index",
		"id",
		strings.NewReader(update),
		es.Update.WithRefresh("true"),
	)
	if err != nil {
		log.Fatalf("Error updating document: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}

func upsertDocument(es *elasticsearch.Client) {
	upsert := `{
        "doc": { "field": "new_value" },
        "doc_as_upsert": true
    }`
	res, err := es.Update(
		"index",
		"id",
		strings.NewReader(upsert),
		es.Update.WithRefresh("true"),
	)
	if err != nil {
		log.Fatalf("Error upserting document: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}

func deleteDocument(es *elasticsearch.Client) {
	res, err := es.Delete(
		"index",
		"id",
		es.Delete.WithRefresh("true"),
	)
	if err != nil {
		log.Fatalf("Error deleting document: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}

func deleteByQuery(es *elasticsearch.Client) {
	query := `{
        "query": {
            "match": {
                "field": "value"
            }
        }
    }`
	res, err := es.DeleteByQuery(
		[]string{"index"},
		strings.NewReader(query),
		es.DeleteByQuery.WithRefresh("true"),
	)
	if err != nil {
		log.Fatalf("Error deleting by query: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}
