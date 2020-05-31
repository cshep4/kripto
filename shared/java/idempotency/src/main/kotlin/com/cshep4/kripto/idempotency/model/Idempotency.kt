package com.cshep4.kripto.idempotency.model

import java.time.LocalDateTime

data class Idempotency(val _id: String, val createdAt: LocalDateTime)