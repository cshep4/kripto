package com.cshep4.kripto.accountretriever.handler

import com.amazonaws.services.lambda.runtime.Context
import com.amazonaws.services.lambda.runtime.RequestHandler
import com.amazonaws.services.lambda.runtime.events.SQSEvent
import com.cshep4.kripto.idempotency.Idempotencer
import com.google.gson.Gson
import com.mongodb.client.MongoClient
import io.lumigo.handlers.LumigoRequestExecutor
import org.litote.kmongo.KMongo
import java.util.function.Supplier
import kotlin.system.exitProcess

class Handler : RequestHandler<SQSEvent, Unit> {
    private val mongoUri = getEnv("MONGO_URI")
    private val client = KMongo.createClient(mongoUri)

    private val idempotencer = Idempotencer("account", client)

    override fun handleRequest(event: SQSEvent, context: Context?) {

        val supplier: Supplier<Unit> = Supplier<Unit> {
            val logger = context?.logger

            logger?.log("request payload: " + Gson().toJson(event))
        }
        return LumigoRequestExecutor.execute(event, context, supplier)
    }

    private fun getEnv(key: String): String {
        val env = System.getenv(key)
        if (env != null) {
            return env
        }

        println("missing environment variable: $key")
        exitProcess(1)
    }
}