package petstore.db

import java.sql.Connection
import java.sql.DriverManager
import java.sql.ResultSet
import java.sql.Timestamp
import java.time.Instant
import java.time.LocalDate

data class SnowflakeConfig(
    val account: String,
    val user: String,
    val password: String,
    val database: String = "TEST",
    val schema: String = "PETSTORE",
    val warehouse: String,
    val role: String,
)

class Database(private val config: SnowflakeConfig) : AutoCloseable {

    private val connection: Connection = DriverManager.getConnection(jdbcUrl(), config.user, config.password)

    private fun jdbcUrl(): String {
        return "jdbc:snowflake://${config.account}.snowflakecomputing.com" +
            "/?db=${config.database}&schema=${config.schema}&warehouse=${config.warehouse}&role=${config.role}"
    }

    fun connection(): Connection = connection

    override fun close() {
        if (!connection.isClosed) connection.close()
    }
}

fun ResultSet.getInstant(col: String): Instant? =
    getTimestamp(col)?.let { it.toInstant() }

fun ResultSet.getInstant(col: Int): Instant? =
    getTimestamp(col)?.let { it.toInstant() }

fun ResultSet.getLocalDate(col: String): LocalDate? =
    getDate(col)?.let { it.toLocalDate() }

fun Instant.toSqlTimestamp(): Timestamp = Timestamp.from(this)
