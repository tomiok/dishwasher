## Architectural decisions

## Architecture Overview

### Components

#### 1. Vector Store (Internal)
- Holds all vectors in memory as float32 slices
- Maintains ID-to-index mapping for O(1) lookups
- Tracks metadata as key-value pairs
- Thread-safe with RWMutex

#### 2. HNSW Index (Internal)
- Hierarchical Navigable Small World graph
- Builds approximate nearest neighbor index
- Enables fast O(log n) search instead of O(n)
- Multi-layer structure

#### 3. File Storage (Internal)
- Binary file format for persistence
- Serializes vectors, metadata, and index
- Supports atomic writes and crash recovery
- Version management for migrations

#### 4. Database API (Public)
- Open(path, opts) - initialize DB
- Insert(id, vector, metadata) - add vectors
- Search(vector, options) - find K nearest neighbors
- Close() - flush to disk
- Stats() - report size and performance

### Data Flow

**Insert Operation:**
User → VectorDB.Insert()
→ VectorStore.Add() (in-memory)
→ HNSWIndex.Insert() (update graph)
→ FileStorage.AppendWAL() (durability)
→ Response

**Search Operation:**
User → VectorDB.Search(query_vector)
→ HNSWIndex.Search() (traverse graph)
→ Calculate distances
→ Apply filters
→ Return top-K results

### Memory Layout

Vectors stored as:
- Array of float32 arrays: [][]float32
- Each vector = dimension * 4 bytes
- Example: 1M vectors × 768 dimensions = ~3GB RAM

ID Index stored as:
- Map[string]int: ID → vector index
- Overhead: ~100 bytes per entry
- Example: 1M vectors = ~100MB

### Thread Safety
- RWMutex protects all mutations
- Searches allowed concurrently (read lock)
- Inserts/deletes exclusive (write lock)