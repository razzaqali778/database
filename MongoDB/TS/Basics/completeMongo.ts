import { MongoClient, Db, Collection, Document, ObjectId } from 'mongodb'

const url = 'mongodb://localhost:27017'
const dbName = 'example_db'
const collectionName = 'example_collection'

// Interfaces to define document structure
interface Reaction {
  userId: string
  emoji: string
}

interface ExampleDocument extends Document {
  name: string
  age: number
  reactions?: Reaction[]
}

let db: Db
let collection: Collection<ExampleDocument>

async function connect() {
  const client = new MongoClient(url, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
  })
  try {
    await client.connect()
    console.log('Connected successfully to MongoDB')
    db = client.db(dbName)
    collection = db.collection(collectionName)
  } catch (error) {
    console.error('Connection error', error)
  }
}

// CRUD Operations
async function createDocument(document: ExampleDocument) {
  try {
    const result = await collection.insertOne(document)
    console.log('Document inserted with _id:', result.insertedId)
  } catch (error) {
    console.error('Create document error:', error)
  }
}

async function createDocuments(documents: ExampleDocument[]) {
  try {
    const result = await collection.insertMany(documents)
    console.log('Documents inserted with _ids:', result.insertedIds)
  } catch (error) {
    console.error('Create documents error:', error)
  }
}

//** In TypeScript, Partial<Type> is a utility type that constructs a type with all properties of Type set to optional. This means that when you create an object of this type, you don't need to provide all the properties that the original type has—only the ones you want to include.* */

async function readDocument(query: Partial<ExampleDocument>) {
  try {
    const document = await collection.findOne(query)
    return document
  } catch (error) {
    console.error('Read document error:', error)
  }
}

async function readDocuments(query: Partial<ExampleDocument>) {
  try {
    const documents = await collection.find(query).toArray()
    return documents
  } catch (error) {
    console.error('Read documents error:', error)
  }
}

async function countDocuments(query: Partial<ExampleDocument>) {
  try {
    const count = await collection.countDocuments(query)
    return count
  } catch (error) {
    console.error('Count documents error:', error)
  }
}

async function updateDocument(
  query: Partial<ExampleDocument>,
  update: Partial<ExampleDocument>
) {
  try {
    const result = await collection.updateOne(query, { $set: update })
    console.log('Updated document count:', result.modifiedCount)
  } catch (error) {
    console.error('Update document error:', error)
  }
}

async function updateDocuments(
  query: Partial<ExampleDocument>,
  update: Partial<ExampleDocument>
) {
  try {
    const result = await collection.updateMany(query, { $set: update })
    console.log('Updated documents count:', result.modifiedCount)
  } catch (error) {
    console.error('Update documents error:', error)
  }
}

async function replaceDocument(
  query: Partial<ExampleDocument>,
  replacement: ExampleDocument
) {
  try {
    const result = await collection.replaceOne(query, replacement)
    console.log('Replaced document count:', result.modifiedCount)
  } catch (error) {
    console.error('Replace document error:', error)
  }
}

async function deleteDocument(query: Partial<ExampleDocument>) {
  try {
    const result = await collection.deleteOne(query)
    console.log('Deleted document count:', result.deletedCount)
  } catch (error) {
    console.error('Delete document error:', error)
  }
}

async function deleteDocuments(query: Partial<ExampleDocument>) {
  try {
    const result = await collection.deleteMany(query)
    console.log('Deleted documents count:', result.deletedCount)
  } catch (error) {
    console.error('Delete documents error:', error)
  }
}

// Query Operators
async function queryOperators() {
  try {
    const documents = await collection
      .find({
        age: { $gte: 18, $lte: 30 },
        name: { $in: ['Alice', 'Bob'] },
        $or: [{ age: { $lt: 25 } }, { name: { $eq: 'Charlie' } }],
        $and: [{ age: { $gt: 20 } }, { name: { $ne: 'Dave' } }],
        reactions: { $exists: true },
      })
      .toArray()
    console.log('Query operator results:', documents)
  } catch (error) {
    console.error('Query operators error:', error)
  }
}

// Update Operators
async function updateOperators() {
  try {
    const result = await collection.updateOne(
      { name: 'Alice' },
      {
        $set: { age: 26 },
        $unset: { reactions: '' },
        $inc: { age: 1 },
        $rename: { age: 'years' },
      }
    )
    console.log('Update operators result:', result)
  } catch (error) {
    console.error('Update operators error:', error)
  }
}

// Array Update Operators
async function arrayUpdateOperators() {
  try {
    const result = await collection.updateOne(
      { name: 'Alice' },
      {
        $push: { reactions: { userId: 'user3', emoji: '😃' } as Reaction },
        $pull: { reactions: { userId: 'user2' } as Reaction },
        $addToSet: { reactions: { userId: 'user4', emoji: '😎' } as Reaction },
        $pop: { reactions: 1 }, // removes the last item
      }
    )
    console.log('Array update operators result:', result)
  } catch (error) {
    console.error('Array update operators error:', error)
  }
}
// Aggregation
async function aggregation() {
  try {
    const pipeline = [
      { $match: { age: { $gte: 18 } } },
      { $group: { _id: '$age', count: { $sum: 1 } } },
      { $sort: { count: -1 } },
      { $limit: 5 },
      { $skip: 1 },
      { $project: { age: '$_id', count: 1, _id: 0 } },
      { $unwind: '$reactions' },
      {
        $lookup: {
          from: 'another_collection',
          localField: 'name',
          foreignField: 'name',
          as: 'related_docs',
        },
      },
      { $addFields: { additionalField: 'new value' } },
      { $replaceRoot: { newRoot: '$related_docs' } },
    ]
    const documents = await collection.aggregate(pipeline).toArray()
    console.log('Aggregation results:', documents)
  } catch (error) {
    console.error('Aggregation error:', error)
  }
}

// Indexing
async function indexing() {
  try {
    const result1 = await collection.createIndex({ name: 1 })
    console.log('Created index:', result1)

    const result2 = await collection.createIndexes([
      { key: { age: 1 } },
      { key: { name: 1, age: 1 } },
    ])
    console.log('Created indexes:', result2)

    const indexes = await collection.listIndexes().toArray()
    console.log('List indexes:', indexes)

    const result3 = await collection.dropIndex('name_1')
    console.log('Dropped index:', result3)

    const result4 = await collection.dropIndexes()
    console.log('Dropped all indexes:', result4)
  } catch (error) {
    console.error('Indexing error:', error)
  }
}

// Miscellaneous Operations
async function miscellaneous() {
  try {
    const bulkResult = await collection.bulkWrite([
      { insertOne: { document: { name: 'Eve', age: 22 } } },
      {
        updateOne: { filter: { name: 'Alice' }, update: { $set: { age: 29 } } },
      },
      { deleteOne: { filter: { name: 'Bob' } } },
    ])
    console.log('Bulk write result:', bulkResult)

    const distinctValues = await collection.distinct('name')
    console.log('Distinct values:', distinctValues)

    const updatedDoc = await collection.findOneAndUpdate(
      { name: 'Eve' },
      { $set: { age: 23 } },
      { returnDocument: 'after' }
    )
    console.log('Find one and update result:', updatedDoc)

    const deletedDoc = await collection.findOneAndDelete({ name: 'Eve' })
    console.log('Find one and delete result:', deletedDoc)

    const replacedDoc = await collection.findOneAndReplace(
      { name: 'Alice' },
      { name: 'Alice', age: 30 }
    )
    console.log('Find one and replace result:', replacedDoc)
  } catch (error) {
    console.error('Miscellaneous operations error:', error)
  }
}

// Main function to run the examples
;(async () => {
  await connect()

  // Run CRUD operations
  await createDocument({ name: 'Alice', age: 25 })
  await createDocuments([
    { name: 'Bob', age: 30 },
    { name: 'Charlie', age: 35 },
  ])
  console.log(await readDocument({ name: 'Alice' }))
  console.log(await readDocuments({}))
  console.log(await countDocuments({ age: { $gte: 25 } }))
  await updateDocument({ name: 'Alice' }, { age: 26 })
  await updateDocuments({ age: { $gt: 25 } }, { age: 27 })
  await replaceDocument({ name: 'Bob' }, { name: 'Bob', age: 31 })
  await deleteDocument({ name: 'Charlie' })
  await deleteDocuments({ age: 27 })

  // Run Query Operators
  await queryOperators()

  // Run Update Operators
  await updateOperators()

  // Run Array Update Operators
  await arrayUpdateOperators()

  // Run Aggregation
  await aggregation()

  // Run Indexing
  await indexing()

  // Run Miscellaneous Operations
  await miscellaneous()

  // Close the connection
  await db.client.close()
})()
