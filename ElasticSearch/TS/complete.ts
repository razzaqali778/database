import { Client } from '@elastic/elasticsearch'

const client = new Client({ node: 'http://localhost:9200' })

// Utility function to log results
const logResult = (operation: string, result: any) => {
  console.log(`${operation}:`, result)
}

// Create: Index Document
const createDocument = async () => {
  const doc = { field: 'value' }
  const res = await client.index({
    index: 'index',
    id: 'id',
    body: doc,
    refresh: 'true',
  })
  logResult('createDocument', res)
}

// Create: Bulk Index Documents
const bulkIndexDocuments = async () => {
  const bulk = [
    { index: { _id: '1' } },
    { field: 'value1' },
    { index: { _id: '2' } },
    { field: 'value2' },
  ]
  const res = await client.bulk({
    index: 'index',
    body: bulk,
    refresh: 'true',
  })
  logResult('bulkIndexDocuments', res)
}

// Read: Get Document
const getDocument = async () => {
  const res = await client.get({
    index: 'index',
    id: 'id',
  })
  logResult('getDocument', res)
}

// Read: Search Documents
const searchDocuments = async () => {
  const query = {
    query: {
      match: {
        field: 'value',
      },
    },
  }
  const res = await client.search({
    index: 'index',
    body: query,
    pretty: true,
  })
  logResult('searchDocuments', res)
}

// Read: Count Documents
const countDocuments = async () => {
  const query = {
    query: {
      match_all: {},
    },
  }
  const res = await client.count({
    index: 'index',
    body: query,
    pretty: true,
  })
  logResult('countDocuments', res)
}

// Update: Update Document
const updateDocument = async () => {
  const update = {
    doc: {
      field: 'new_value',
    },
  }
  const res = await client.update({
    index: 'index',
    id: 'id',
    body: update,
    refresh: 'true',
  })
  logResult('updateDocument', res)
}

// Update: Upsert Document
const upsertDocument = async () => {
  const upsert = {
    doc: { field: 'new_value' },
    doc_as_upsert: true,
  }
  const res = await client.update({
    index: 'index',
    id: 'id',
    body: upsert,
    refresh: 'true',
  })
  logResult('upsertDocument', res)
}

// Delete: Delete Document
const deleteDocument = async () => {
  const res = await client.delete({
    index: 'index',
    id: 'id',
    refresh: 'true',
  })
  logResult('deleteDocument', res)
}

// Delete: Delete by Query
const deleteByQuery = async () => {
  const query = {
    query: {
      match: {
        field: 'value',
      },
    },
  }
  const res = await client.deleteByQuery({
    index: 'index',
    body: query,
    refresh: 'true',
  })
  logResult('deleteByQuery', res)
}

// Main function to run the examples
const main = async () => {
  await createDocument()
  await bulkIndexDocuments()
  await getDocument()
  await searchDocuments()
  await countDocuments()
  await updateDocument()
  await upsertDocument()
  await deleteDocument()
  await deleteByQuery()
}

main().catch(console.error)
