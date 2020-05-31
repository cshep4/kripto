package com.cshep4.kripto.receiptemailer.email

import com.cshep4.kripto.receiptemailer.model.Trade
import com.cshep4.kripto.receiptemailer.result.EmailResult
import com.natpryce.hamkrest.assertion.assertThat
import com.natpryce.hamkrest.equalTo
import com.sendgrid.SendGrid
import org.junit.Test
import java.time.LocalDateTime

internal class EmailerIntegrationTest {
    @Test
    fun send() {
        val sendGrid = SendGrid(System.getenv("SEND_GRID_API_KEY"))
        val emailer = Emailer("chris_shepherd2@hotmail.com", sendGrid)

        val trade = Trade(
                id = "aa368788-bb4f-40c0-b80f-afcfdaf18574",
                side = "buy",
                productId = "BTC-GBP",
                funds = "9.95024875",
                settled = true,
                createdAt = LocalDateTime.now(),
                fillFees = "0.049751102976",
                filledSize = "0.00125952",
                executedValue = "9.9502205952"
        )

        val res = emailer.send(null, trade)
        val expectedResult: EmailResult = EmailResult.Success

        assertThat(res, equalTo(expectedResult))
    }
}