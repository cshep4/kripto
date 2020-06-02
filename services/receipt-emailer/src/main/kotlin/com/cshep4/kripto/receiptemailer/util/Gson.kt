package com.cshep4.kripto.receiptemailer.util

import com.google.gson.*
import com.google.gson.Gson
import java.lang.reflect.Type
import java.time.LocalDateTime
import java.time.ZonedDateTime

object Gson {
    fun build(): Gson {
        return GsonBuilder().registerTypeAdapter(LocalDateTime::class.java,
                object : JsonDeserializer<LocalDateTime> {
                    override fun deserialize(json: JsonElement?, typeOfT: Type?, context: JsonDeserializationContext?): LocalDateTime {
                        return ZonedDateTime.parse(json?.asJsonPrimitive?.asString).toLocalDateTime()
                    }

                }).create()
    }
}