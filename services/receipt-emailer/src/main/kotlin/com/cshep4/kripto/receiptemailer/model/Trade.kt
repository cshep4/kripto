package com.cshep4.kripto.receiptemailer.model

import java.time.LocalDateTime

data class Trade(
        val id: String,
        val side: String,
        val productId: String,
        val funds: String,
        val settled: Boolean,
        val createdAt: LocalDateTime,
        val fillFees: String,
        val filledSize: String,
        val executedValue: String
)