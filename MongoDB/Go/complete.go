package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	uri            = "mongodb://localhost:27017"
	dbName         = "example_db"
	collectionName = "example_collection"
)

type Reaction struct {
	UserID string `bson:"userId"`
	Emoji  string `bson:"emoji"`
}

type ExampleDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Age       int                `bson:"age"`
	Reactions []Reaction         `bson:"reactions,omitempty"`
}

var collection *mongo.Collection

func connect() {
	clientOptions := options.Client().ApplyURI(uri).SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected successfully to MongoDB")

	collection = client.Database(dbName).Collection(collectionName)
}

// CRUD Operations
func createDocument(document ExampleDocument) {
	_, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Document inserted")
}

func createDocuments(documents []ExampleDocument) {
	var docs []interface{}
	for _, doc := range documents {
		docs = append(docs, doc)
	}
	_, err := collection.InsertMany(context.TODO(), docs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Documents inserted")
}

func readDocument(filter bson.M) {
	var result ExampleDocument
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Read document:", result)
}

func readDocuments(filter bson.M) {
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	var results []ExampleDocument
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Read documents:", results)
}

func countDocuments(filter bson.M) {
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Count documents:", count)
}

func updateDocument(filter bson.M, update bson.M) {
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated document")
}

func updateDocuments(filter bson.M, update bson.M) {
	_, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated documents")
}

func replaceDocument(filter bson.M, replacement ExampleDocument) {
	_, err := collection.ReplaceOne(context.TODO(), filter, replacement)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Replaced document")
}

func deleteDocument(filter bson.M) {
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted document")
}

func deleteDocuments(filter bson.M) {
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted documents")
}

// Query Operators
func queryOperators() {
	filter := bson.M{
		"age":       bson.M{"$gte": 18, "$lte": 30},
		"name":      bson.M{"$in": bson.A{"Alice", "Bob"}},
		"$or":       bson.A{bson.M{"age": bson.M{"$lt": 25}}, bson.M{"name": "Charlie"}},
		"$and":      bson.A{bson.M{"age": bson.M{"$gt": 20}}, bson.M{"name": bson.M{"$ne": "Dave"}}},
		"reactions": bson.M{"$exists": true},
	}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	var results []ExampleDocument
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Query operator results:", results)
}

// Update Operators
func updateOperators() {
	filter := bson.M{"name": "Alice"}
	update := bson.M{
		"$set":    bson.M{"age": 26},
		"$unset":  bson.M{"reactions": ""},
		"$inc":    bson.M{"age": 1},
		"$rename": bson.M{"age": "years"},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Update operators result")
}

// Array Update Operators
func arrayUpdateOperators() {
	filter := bson.M{"name": "Alice"}
	update := bson.M{
		"$push":     bson.M{"reactions": Reaction{UserID: "user3", Emoji: "😃"}},
		"$pull":     bson.M{"reactions": bson.M{"userId": "user2"}},
		"$addToSet": bson.M{"reactions": Reaction{UserID: "user4", Emoji: "😎"}},
		"$pop":      bson.M{"reactions": 1}, // removes the last item
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Array update operators result")
}

// Aggregation
func aggregation() {
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"age", bson.D{{"$gte", 18}}}}}},
		bson.D{{"$group", bson.D{{"_id", "$age"}, {"count", bson.D{{"$sum", 1}}}}}},
		bson.D{{"$sort", bson.D{{"count", -1}}}},
		bson.D{{"$limit", 5}},
		bson.D{{"$skip", 1}},
		bson.D{{"$project", bson.D{{"age", "$_id"}, {"count", 1}, {"_id", 0}}}},
		bson.D{{"$unwind", "$reactions"}},
		bson.D{{"$lookup", bson.D{{"from", "another_collection"}, {"localField", "name"}, {"foreignField", "name"}, {"as", "related_docs"}}}},
		bson.D{{"$addFields", bson.D{{"additionalField", "new value"}}}},
		bson.D{{"$replaceRoot", bson.D{{"newRoot", "$related_docs"}}}},
	}
	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Aggregation results:", results)
}

// Indexing
func indexing() {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{"name", 1}},
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created index")

	indexModels := []mongo.IndexModel{
		{Keys: bson.D{{"age", 1}}},
		{Keys: bson.D{{"name", 1}, {"age", 1}}},
	}
	_, err = collection.Indexes().CreateMany(context.TODO(), indexModels)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created indexes")

	cursor, err := collection.Indexes().List(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	var indexes []bson.M
	if err := cursor.All(context.TODO(), &indexes); err != nil {
		log.Fatal(err)
	}
	fmt.Println("List indexes:", indexes)

	_, err = collection.Indexes().DropOne(context.TODO(), "name_1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Dropped index")

	_, err = collection.Indexes().DropAll(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Dropped all indexes")
}

// Miscellaneous Operations
func miscellaneous() {
	models := []mongo.WriteModel{
		mongo.NewInsertOneModel().SetDocument(ExampleDocument{Name: "Eve", Age: 22}),
		mongo.NewUpdateOneModel().SetFilter(bson.M{"name": "Alice"}).SetUpdate(bson.M{"$set": bson.M{"age": 29}}),
		mongo.NewDeleteOneModel().SetFilter(bson.M{"name": "Bob"}),
	}
	bulkResult, err := collection.BulkWrite(context.TODO(), models)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bulk write result:", bulkResult)

	distinctValues, err := collection.Distinct(context.TODO(), "name", bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Distinct values:", distinctValues)

	updatedDoc := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"name": "Eve"},
		bson.M{"$set": bson.M{"age": 23}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var updatedResult ExampleDocument
	if err := updatedDoc.Decode(&updatedResult); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Find one and update result:", updatedResult)

	deletedDoc := collection.FindOneAndDelete(context.TODO(), bson.M{"name": "Eve"})
	var deletedResult ExampleDocument
	if err := deletedDoc.Decode(&deletedResult); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Find one and delete result:", deletedResult)

	replacedDoc := collection.FindOneAndReplace(
		context.TODO(),
		bson.M{"name": "Alice"},
		ExampleDocument{Name: "Alice", Age: 30},
	)
	var replacedResult ExampleDocument
	if err := replacedDoc.Decode(&replacedResult); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Find one and replace result:", replacedResult)
}

func main() {
	connect()

	// Run CRUD operations
	createDocument(ExampleDocument{Name: "Alice", Age: 25})
	createDocuments([]ExampleDocument{
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
	})
	readDocument(bson.M{"name": "Alice"})
	readDocuments(bson.M{})
	countDocuments(bson.M{"age": bson.M{"$gte": 25}})
	updateDocument(bson.M{"name": "Alice"}, bson.M{"$set": bson.M{"age": 26}})
	updateDocuments(bson.M{"age": bson.M{"$gt": 25}}, bson.M{"$set": bson.M{"age": 27}})
	replaceDocument(bson.M{"name": "Bob"}, ExampleDocument{Name: "Bob", Age: 31})
	deleteDocument(bson.M{"name": "Charlie"})
	deleteDocuments(bson.M{"age": 27})

	// Run Query Operators
	queryOperators()

	// Run Update Operators
	updateOperators()

	// Run Array Update Operators
	arrayUpdateOperators()

	// Run Aggregation
	aggregation()

	// Run Indexing
	indexing()

	// Run Miscellaneous Operations
	miscellaneous()
}
