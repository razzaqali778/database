import { MongoClient, Db, Collection, Document, ObjectId } from 'mongodb'

const url = 'mongodb://localhost:27017'
const dbName = 'example_db'
const collectionName = 'example_collection'

let db: Db
let collection: Collection<Document>

// Interfaces to define document structure
interface Reaction {
  userId: string
  emoji: string
}

interface Post extends Document {
  name: string
  reactions: Reaction[]
}

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

async function createDocument(document: Document) {
  try {
    const result = await collection.insertOne(document)
    console.log('Document inserted with _id: ', result.insertedId)
  } catch (error) {
    console.error('Create document error', error)
  }
}

async function readDocument(query: Document) {
  try {
    const document = await collection.findOne(query)
    return document
  } catch (error) {
    console.error('Read document error', error)
  }
}

async function updateDocument(query: Document, update: Document) {
  try {
    const result = await collection.updateOne(query, { $set: update })
    console.log(
      `${result.matchedCount} document(s) matched the filter, updated ${result.modifiedCount} document(s)`
    )
  } catch (error) {
    console.error('Update document error', error)
  }
}

async function deleteDocument(query: Document) {
  try {
    const result = await collection.deleteOne(query)
    console.log(`Deleted ${result.deletedCount} document(s)`)
  } catch (error) {
    console.error('Delete document error', error)
  }
}

async function addReaction(postId: string, userId: string, emoji: string) {
  try {
    const post = await collection.findOne({ _id: new ObjectId(postId) })
    if (!post) {
      throw new Error('Post not found')
    }

    const existingReaction = post.reactions.find(
      (reaction: { userId: string }) => reaction.userId === userId
    )

    if (existingReaction) {
      await collection.updateOne(
        { _id: new ObjectId(postId), 'reactions.userId': userId },
        { $set: { 'reactions.$.emoji': emoji } }
      )
      console.log(`Updated reaction for user ${userId} on post ${postId}`)
    } else {
      await collection.updateOne(
        { _id: new ObjectId(postId) },
        { $push: { reactions: { userId, emoji } } }
      )
      console.log(
        `Added reaction ${emoji} for user ${userId} on post ${postId}`
      )
    }
  } catch (error) {
    console.error('Add reaction error', error)
  }
}

async function removeReaction(postId: string, userId: string) {
  try {
    await collection.updateOne(
      { _id: new ObjectId(postId) },
      { $pull: { reactions: { userId } } }
    )
    console.log(`Removed reaction for user ${userId} from post ${postId}`)
  } catch (error) {
    console.error('Remove reaction error', error)
  }
}

// Usage example
;(async () => {
  await connect()
  const document = { name: 'Post 1', reactions: [] }
  await createDocument(document)
  const postId = document._id

  await addReaction(postId.toString(), 'user1', '👍')
  await addReaction(postId.toString(), 'user2', '❤️')
  await addReaction(postId.toString(), 'user1', '😂') // This should update user1's reaction to 😂

  console.log(await readDocument({ _id: new ObjectId(postId) }))

  await removeReaction(postId.toString(), 'user1')

  console.log(await readDocument({ _id: new ObjectId(postId) }))

  await deleteDocument({ _id: new ObjectId(postId) })

  // Close the connection
  await db.client.close()
})()
