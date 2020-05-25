package com.cshep4.kripto.idempotency.result

sealed class IdempotencyResult {
    data class Success(val exists: Boolean) : IdempotencyResult()
    data class Error(val message: String, val cause: Throwable? = null) : IdempotencyResult()
}