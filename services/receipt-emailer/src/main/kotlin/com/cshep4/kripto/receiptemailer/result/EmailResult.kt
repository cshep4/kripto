package com.cshep4.kripto.receiptemailer.result

sealed class EmailResult {
    object Success : EmailResult()
    data class Error(val message: String, val cause: Throwable? = null) : EmailResult()
}