package com.cshep4.kripto.idempotency

import com.cshep4.kripto.idempotency.model.Idempotency
import com.cshep4.kripto.idempotency.result.IdempotencyResult
import com.mongodb.client.MongoClient
import com.mongodb.client.model.IndexOptions
import org.litote.kmongo.ensureIndex
import org.litote.kmongo.findOneById
import java.time.LocalDateTime
import java.util.concurrent.TimeUnit

class Idempotencer(database: String, client: MongoClient) {
    private val collection = client
            .getDatabase(database)
            .getCollection("idempotency", Idempotency::class.java)

    init {
        val opts = IndexOptions()
                .unique(true)
                .name("createdAtIdx")
                .background(true)
                .expireAfter(4, TimeUnit.DAYS)

        collection.ensureIndex("{'createdAt':1}", opts)
    }

    fun check(key: String): IdempotencyResult {
        return try {
            val doc = collection.findOneById(key)

            if (doc != null) {
                return IdempotencyResult.Success(true)
            }

            collection.insertOne(Idempotency(_id = key, createdAt = LocalDateTime.now()))

            IdempotencyResult.Success(false)
        } catch (e: Throwable) {
            IdempotencyResult.Error(e.message!!, e)
        }
    }
}