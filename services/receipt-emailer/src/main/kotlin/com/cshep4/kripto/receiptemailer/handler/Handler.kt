package com.cshep4.kripto.receiptemailer.handler

import com.amazonaws.services.lambda.runtime.Context
import com.amazonaws.services.lambda.runtime.RequestHandler
import com.amazonaws.services.lambda.runtime.events.SQSEvent
import com.cshep4.kripto.idempotency.Idempotencer
import com.cshep4.kripto.idempotency.result.IdempotencyResult
import com.cshep4.kripto.receiptemailer.email.Emailer
import com.cshep4.kripto.receiptemailer.model.Trade
import com.cshep4.kripto.receiptemailer.result.EmailResult
import com.cshep4.kripto.receiptemailer.util.Gson
import com.sendgrid.SendGrid
import org.litote.kmongo.KMongo
import kotlin.system.exitProcess


class Handler : RequestHandler<SQSEvent, Unit> {
    private val sendGrid = SendGrid(getEnv("SEND_GRID_API_KEY"))
    private val emailer = Emailer(getEnv("RECEIPT_RECIPIENT"), sendGrid)

    private val client = KMongo.createClient(getEnv("MONGO_URI"))
    private val idempotencer = Idempotencer("receipt", client)

    private val gson = Gson.build()

    override fun handleRequest(event: SQSEvent, context: Context?) {
        val logger = context?.logger

        event.records.forEach {
            val trade = gson.fromJson(it.body, Trade::class.java)

            when (val i = idempotencer.check(trade.id)) {
                is IdempotencyResult.Error -> throw Exception(i.message, i.cause)
                is IdempotencyResult.Success -> {
                    if (i.exists) {
                        logger?.log("msg_already_processed - trade: " + gson.toJson(trade))
                        return@forEach
                    }
                }
            }

            when (val e = emailer.send(context, trade)) {
                is EmailResult.Error -> throw throw Exception(e.message, e.cause)
            }
        }
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