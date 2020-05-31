package com.cshep4.kripto.idempotency

import com.cshep4.kripto.idempotency.model.Idempotency
import com.cshep4.kripto.idempotency.result.IdempotencyResult
import com.mongodb.MongoException
import com.mongodb.client.MongoClient
import com.mongodb.client.MongoCollection
import com.mongodb.client.MongoDatabase
import com.mongodb.client.result.InsertOneResult
import com.natpryce.hamkrest.assertion.assertThat
import com.natpryce.hamkrest.equalTo
import io.mockk.every
import io.mockk.mockk
import io.mockk.verify
import org.junit.Test
import org.litote.kmongo.ensureIndex
import org.litote.kmongo.findOneById
import java.time.LocalDateTime

internal class IdempotencerTest {
    companion object {
        const val DB_NAME = "db"
        const val KEY = "🔑"
    }

    @Test
    fun `TTL index is ensured on init`() {
        val (collection, _) = initIdempotencer()

        verify { collection.ensureIndex("{'createdAt':1}", any()) }
    }

    @Test
    fun `'check' returns true if doc exists`() {
        val (collection, idempotencer) = initIdempotencer()

        every { collection.findOneById(KEY) } returns Idempotency(_id = KEY, createdAt = LocalDateTime.now())

        val res = idempotencer.check(KEY)

        val expectedResult: IdempotencyResult = IdempotencyResult.Success(true)

        assertThat(res, equalTo(expectedResult))
        verify(exactly = 0) { collection.insertOne(any()) }
    }

    @Test
    fun `'check' returns false if doc does not exist`() {
        val (collection, idempotencer) = initIdempotencer()

        every { collection.findOneById(KEY) } returns null
        every { collection.insertOne(any()) } returns InsertOneResult.unacknowledged()

        val res = idempotencer.check(KEY)

        val expectedResult: IdempotencyResult = IdempotencyResult.Success(false)

        assertThat(res, equalTo(expectedResult))
        verify { collection.insertOne(any()) }
    }

    @Test
    fun `'check' returns error if error storing key`() {
        val (collection, idempotencer) = initIdempotencer()

        val me = MongoException("error")

        every { collection.findOneById(KEY) } returns null
        every { collection.insertOne(any()) } throws me

        val res = idempotencer.check(KEY)

        val expectedResult: IdempotencyResult = IdempotencyResult.Error(me.message!!, me)

        assertThat(res, equalTo(expectedResult))
        verify { collection.insertOne(any()) }
    }

    private fun initIdempotencer(): Pair<MongoCollection<Idempotency>, Idempotencer> {
        val client = mockk<MongoClient>()
        val db = mockk<MongoDatabase>()
        val collection = mockk<MongoCollection<Idempotency>>()

        every { client.getDatabase(DB_NAME) } returns db
        every { db.getCollection("idempotency", Idempotency::class.java) } returns collection
        every { collection.ensureIndex("{'createdAt':1}", any()) } returns "index"

        val idempotencer = Idempotencer(DB_NAME, client)
        return Pair(collection, idempotencer)
    }
}
