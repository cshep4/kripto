package com.cshep4.kripto.receiptemailer.email

import com.amazonaws.services.lambda.runtime.Context
import com.cshep4.kripto.receiptemailer.model.Trade
import com.cshep4.kripto.receiptemailer.result.EmailResult
import com.sendgrid.Method
import com.sendgrid.Request
import com.sendgrid.SendGrid
import com.sendgrid.helpers.mail.Mail
import com.sendgrid.helpers.mail.objects.Content
import com.sendgrid.helpers.mail.objects.Email


class Emailer(private val to: String, private val sendGrid: SendGrid) {
    fun send(context: Context?, trade: Trade): EmailResult {
        val logger = context?.logger

        val from = Email("receipt@kripto.com", "Kripto Platform")
        val subject = "Kripto Trade Receipt"
        val to = Email(to)
        val content = Content("text/html", Template.build(trade))
        val mail = Mail(from, subject, to, content)

        val request = Request()

        return try {
            request.method = Method.POST
            request.endpoint = "mail/send"
            request.body = mail.build()

            val res = sendGrid.api(request)

            logger?.log("response: ${res.statusCode}. body: ${res.body}")

            EmailResult.Success
        } catch (e: Throwable) {
            EmailResult.Error(e.message!!, e)
        }

    }
}