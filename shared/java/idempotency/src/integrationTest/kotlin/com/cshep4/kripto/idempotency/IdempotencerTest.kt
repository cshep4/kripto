package com.cshep4.kripto.idempotency

import com.cshep4.kripto.idempotency.result.IdempotencyResult
import com.natpryce.hamkrest.assertion.assertThat
import com.natpryce.hamkrest.equalTo
import org.bson.Document
import org.junit.After
import org.junit.Test
import org.litote.kmongo.KMongo

internal class IdempotencerTest {
    companion object {
        const val DB_NAME = "db"
        const val KEY = "🔑"
        const val MONGO_URI = "mongodb://localhost:27017"
    }

    private val client = KMongo.createClient(MONGO_URI)

    private val idempotencer = Idempotencer(DB_NAME, client)

    @Test
    fun `'check' returns false if idempotency key not used'`() {
        val expected: IdempotencyResult = IdempotencyResult.Success(false)

        val res1 = idempotencer.check(KEY)
        val res2 = idempotencer.check("2")

        assertThat(res1, equalTo(expected))
        assertThat(res2, equalTo(expected))
    }

    @Test
    fun `'check' returns true if idempotency key used'`() {
        val expected1: IdempotencyResult = IdempotencyResult.Success(false)
        val expected2: IdempotencyResult = IdempotencyResult.Success(true)

        val res1 = idempotencer.check(KEY)
        val res2 = idempotencer.check(KEY)

        assertThat(res1, equalTo(expected1))
        assertThat(res2, equalTo(expected2))
    }

    @After
    fun cleanup() {
        client.getDatabase(DB_NAME)
                .getCollection("idempotency")
                .deleteMany(Document())
    }
}
