package petstore.repository

import petstore.db.Database
import petstore.db.getInstant
import petstore.db.getLocalDate
import petstore.db.toSqlTimestamp
import petstore.model.*
import java.sql.SQLException
import java.time.Instant

class NotFoundException(message: String) : RuntimeException(message)

class CustomerRepository(private val db: Database) {

    fun create(c: Customer): Customer {
        val sql = """
            INSERT INTO customers (id, full_name, email, phone, address, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        """.trimIndent()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, c.id)
            ps.setString(2, c.fullName)
            ps.setString(3, c.email)
            ps.setString(4, c.phone)
            ps.setString(5, c.address)
            ps.setTimestamp(6, c.createdAt.toSqlTimestamp())
            ps.setTimestamp(7, c.updatedAt.toSqlTimestamp())
            ps.executeUpdate()
        }
        return c
    }

    fun findById(id: String): Customer? {
        val sql = "SELECT id, full_name, email, phone, address, created_at, updated_at FROM customers WHERE id = ?"
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, id)
            ps.executeQuery().use { rs ->
                if (rs.next()) return rsToCustomer(rs)
            }
        }
        return null
    }

    fun list(offset: Int, limit: Int): List<Customer> {
        val sql = "SELECT id, full_name, email, phone, address, created_at, updated_at FROM customers ORDER BY created_at DESC LIMIT ? OFFSET ?"
        val results = mutableListOf<Customer>()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setInt(1, limit)
            ps.setInt(2, offset)
            ps.executeQuery().use { rs ->
                while (rs.next()) results.add(rsToCustomer(rs))
            }
        }
        return results
    }

    fun update(c: Customer): Customer {
        val sql = """
            UPDATE customers
            SET full_name = ?, email = ?, phone = ?, address = ?, updated_at = ?
            WHERE id = ?
        """.trimIndent()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, c.fullName)
            ps.setString(2, c.email)
            ps.setString(3, c.phone)
            ps.setString(4, c.address)
            ps.setTimestamp(5, Instant.now().toSqlTimestamp())
            ps.setString(6, c.id)
            val rows = ps.executeUpdate()
            if (rows == 0) throw NotFoundException("customer ${c.id} not found")
        }
        return c
    }

    fun delete(id: String) {
        db.connection().prepareStatement("DELETE FROM customers WHERE id = ?").use { ps ->
            ps.setString(1, id)
            ps.executeUpdate()
        }
    }

    private fun rsToCustomer(rs: java.sql.ResultSet) = Customer(
        id = rs.getString("id"),
        fullName = rs.getString("full_name"),
        email = rs.getString("email"),
        phone = rs.getString("phone"),
        address = rs.getString("address"),
        createdAt = rs.getInstant("created_at")!!,
        updatedAt = rs.getInstant("updated_at")!!,
    )
}

class BreedRepository(private val db: Database) {

    fun create(b: Breed): Breed {
        val sql = """
            INSERT INTO breeds (id, name, description, origin_country, avg_lifespan_years, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        """.trimIndent()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, b.id)
            ps.setString(2, b.name)
            ps.setString(3, b.description)
            ps.setString(4, b.originCountry)
            ps.setInt(5, b.avgLifespanYears)
            ps.setTimestamp(6, b.createdAt.toSqlTimestamp())
            ps.setTimestamp(7, b.updatedAt.toSqlTimestamp())
            ps.executeUpdate()
        }
        return b
    }

    fun findById(id: String): Breed? {
        val sql = "SELECT id, name, description, origin_country, avg_lifespan_years, created_at, updated_at FROM breeds WHERE id = ?"
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, id)
            ps.executeQuery().use { rs ->
                if (rs.next()) return rsToBreed(rs)
            }
        }
        return null
    }

    fun list(): List<Breed> {
        val sql = "SELECT id, name, description, origin_country, avg_lifespan_years, created_at, updated_at FROM breeds ORDER BY name"
        val results = mutableListOf<Breed>()
        db.connection().prepareStatement(sql).use { ps ->
            ps.executeQuery().use { rs ->
                while (rs.next()) results.add(rsToBreed(rs))
            }
        }
        return results
    }

    private fun rsToBreed(rs: java.sql.ResultSet) = Breed(
        id = rs.getString("id"),
        name = rs.getString("name"),
        description = rs.getString("description"),
        originCountry = rs.getString("origin_country"),
        avgLifespanYears = rs.getInt("avg_lifespan_years"),
        createdAt = rs.getInstant("created_at")!!,
        updatedAt = rs.getInstant("updated_at")!!,
    )
}

class ShelterRepository(private val db: Database) {

    fun create(s: Shelter): Shelter {
        val sql = """
            INSERT INTO shelters (id, name, address, phone, email, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        """.trimIndent()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, s.id)
            ps.setString(2, s.name)
            ps.setString(3, s.address)
            ps.setString(4, s.phone)
            ps.setString(5, s.email)
            ps.setTimestamp(6, s.createdAt.toSqlTimestamp())
            ps.setTimestamp(7, s.updatedAt.toSqlTimestamp())
            ps.executeUpdate()
        }
        return s
    }

    fun findById(id: String): Shelter? {
        val sql = "SELECT id, name, address, phone, email, created_at, updated_at FROM shelters WHERE id = ?"
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, id)
            ps.executeQuery().use { rs ->
                if (rs.next()) return rsToShelter(rs)
            }
        }
        return null
    }

    fun list(): List<Shelter> {
        val sql = "SELECT id, name, address, phone, email, created_at, updated_at FROM shelters ORDER BY name"
        val results = mutableListOf<Shelter>()
        db.connection().prepareStatement(sql).use { ps ->
            ps.executeQuery().use { rs ->
                while (rs.next()) results.add(rsToShelter(rs))
            }
        }
        return results
    }

    private fun rsToShelter(rs: java.sql.ResultSet) = Shelter(
        id = rs.getString("id"),
        name = rs.getString("name"),
        address = rs.getString("address"),
        phone = rs.getString("phone"),
        email = rs.getString("email"),
        createdAt = rs.getInstant("created_at")!!,
        updatedAt = rs.getInstant("updated_at")!!,
    )
}

class CatRepository(private val db: Database) {

    fun create(c: Cat): Cat {
        val sql = """
            INSERT INTO cats (id, name, breed_id, shelter_id, birth_date, price_cents, currency, status, description, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """.trimIndent()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, c.id)
            ps.setString(2, c.name)
            ps.setString(3, c.breedId)
            ps.setString(4, c.shelterId)
            ps.setDate(5, java.sql.Date.valueOf(c.birthDate))
            ps.setLong(6, c.priceCents)
            ps.setString(7, c.currency)
            ps.setString(8, c.status.name)
            ps.setString(9, c.description)
            ps.setTimestamp(10, c.createdAt.toSqlTimestamp())
            ps.setTimestamp(11, c.updatedAt.toSqlTimestamp())
            ps.executeUpdate()
        }
        return c
    }

    fun findById(id: String): Cat? {
        val sql = "SELECT id, name, breed_id, shelter_id, birth_date, price_cents, currency, status, description, created_at, updated_at FROM cats WHERE id = ?"
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, id)
            ps.executeQuery().use { rs ->
                if (rs.next()) return rsToCat(rs)
            }
        }
        return null
    }

    fun list(offset: Int, limit: Int): List<Cat> {
        val sql = "SELECT id, name, breed_id, shelter_id, birth_date, price_cents, currency, status, description, created_at, updated_at FROM cats ORDER BY created_at DESC LIMIT ? OFFSET ?"
        val results = mutableListOf<Cat>()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setInt(1, limit)
            ps.setInt(2, offset)
            ps.executeQuery().use { rs ->
                while (rs.next()) results.add(rsToCat(rs))
            }
        }
        return results
    }

    fun searchByName(query: String, limit: Int): List<Cat> {
        val sql = "SELECT id, name, breed_id, shelter_id, birth_date, price_cents, currency, status, description, created_at, updated_at FROM cats WHERE name ILIKE ? ORDER BY name LIMIT ?"
        val results = mutableListOf<Cat>()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, "%$query%")
            ps.setInt(2, limit)
            ps.executeQuery().use { rs ->
                while (rs.next()) results.add(rsToCat(rs))
            }
        }
        return results
    }

    fun updateStatus(id: String, status: CatStatus) {
        val sql = "UPDATE cats SET status = ?, updated_at = ? WHERE id = ?"
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, status.name)
            ps.setTimestamp(2, Instant.now().toSqlTimestamp())
            ps.setString(3, id)
            val rows = ps.executeUpdate()
            if (rows == 0) throw NotFoundException("cat $id not found")
        }
    }

    fun delete(id: String) {
        db.connection().prepareStatement("DELETE FROM cats WHERE id = ?").use { ps ->
            ps.setString(1, id)
            ps.executeUpdate()
        }
    }

    private fun rsToCat(rs: java.sql.ResultSet) = Cat(
        id = rs.getString("id"),
        name = rs.getString("name"),
        breedId = rs.getString("breed_id"),
        shelterId = rs.getString("shelter_id"),
        birthDate = rs.getLocalDate("birth_date")!!,
        priceCents = rs.getLong("price_cents"),
        currency = rs.getString("currency"),
        status = CatStatus.valueOf(rs.getString("status")),
        description = rs.getString("description"),
        createdAt = rs.getInstant("created_at")!!,
        updatedAt = rs.getInstant("updated_at")!!,
    )
}

class CartRepository(private val db: Database) {

    fun create(cart: Cart): Cart {
        val sql = "INSERT INTO carts (id, customer_id, created_at, updated_at) VALUES (?, ?, ?, ?)"
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, cart.id)
            ps.setString(2, cart.customerId)
            ps.setTimestamp(3, cart.createdAt.toSqlTimestamp())
            ps.setTimestamp(4, cart.updatedAt.toSqlTimestamp())
            ps.executeUpdate()
        }
        return cart
    }

    fun findById(id: String): Cart? {
        val sql = "SELECT id, customer_id, created_at, updated_at FROM carts WHERE id = ?"
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, id)
            ps.executeQuery().use { rs ->
                if (rs.next()) return rsToCart(rs)
            }
        }
        return null
    }

    fun delete(id: String) {
        db.connection().prepareStatement("DELETE FROM carts WHERE id = ?").use { ps ->
            ps.setString(1, id)
            ps.executeUpdate()
        }
    }

    private fun rsToCart(rs: java.sql.ResultSet) = Cart(
        id = rs.getString("id"),
        customerId = rs.getString("customer_id"),
        createdAt = rs.getInstant("created_at")!!,
        updatedAt = rs.getInstant("updated_at")!!,
    )
}

class CartItemRepository(private val db: Database) {

    fun create(item: CartItem): CartItem {
        val sql = "INSERT INTO cart_items (id, cart_id, cat_id, quantity, created_at) VALUES (?, ?, ?, ?, ?)"
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, item.id)
            ps.setString(2, item.cartId)
            ps.setString(3, item.catId)
            ps.setInt(4, item.quantity)
            ps.setTimestamp(5, item.createdAt.toSqlTimestamp())
            ps.executeUpdate()
        }
        return item
    }

    fun findByCart(cartId: String): List<CartItem> {
        val sql = "SELECT id, cart_id, cat_id, quantity, created_at FROM cart_items WHERE cart_id = ? ORDER BY created_at"
        val results = mutableListOf<CartItem>()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, cartId)
            ps.executeQuery().use { rs ->
                while (rs.next()) results.add(rsToCartItem(rs))
            }
        }
        return results
    }

    fun delete(cartId: String, catId: String) {
        db.connection().prepareStatement("DELETE FROM cart_items WHERE cart_id = ? AND cat_id = ?").use { ps ->
            ps.setString(1, cartId)
            ps.setString(2, catId)
            ps.executeUpdate()
        }
    }

    fun deleteAllByCart(cartId: String) {
        db.connection().prepareStatement("DELETE FROM cart_items WHERE cart_id = ?").use { ps ->
            ps.setString(1, cartId)
            ps.executeUpdate()
        }
    }

    private fun rsToCartItem(rs: java.sql.ResultSet) = CartItem(
        id = rs.getString("id"),
        cartId = rs.getString("cart_id"),
        catId = rs.getString("cat_id"),
        quantity = rs.getInt("quantity"),
        createdAt = rs.getInstant("created_at")!!,
    )
}

class OrderRepository(private val db: Database) {

    fun create(o: Order): Order {
        val sql = """
            INSERT INTO orders (id, customer_id, status, total_cents, currency, ship_address, bill_address, notes, created_at, updated_at, placed_at, shipped_at, cancelled_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """.trimIndent()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, o.id)
            ps.setString(2, o.customerId)
            ps.setString(3, o.status.name)
            ps.setLong(4, o.totalCents)
            ps.setString(5, o.currency)
            ps.setString(6, o.shipAddress)
            ps.setString(7, o.billAddress)
            ps.setString(8, o.notes)
            ps.setTimestamp(9, o.createdAt.toSqlTimestamp())
            ps.setTimestamp(10, o.updatedAt.toSqlTimestamp())
            ps.setTimestamp(11, o.placedAt?.toSqlTimestamp())
            ps.setTimestamp(12, o.shippedAt?.toSqlTimestamp())
            ps.setTimestamp(13, o.cancelledAt?.toSqlTimestamp())
            ps.executeUpdate()
        }
        return o
    }

    fun findById(id: String): Order? {
        val sql = "SELECT id, customer_id, status, total_cents, currency, ship_address, bill_address, notes, created_at, updated_at, placed_at, shipped_at, cancelled_at FROM orders WHERE id = ?"
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, id)
            ps.executeQuery().use { rs ->
                if (rs.next()) return rsToOrder(rs)
            }
        }
        return null
    }

    fun listByCustomer(customerId: String, offset: Int, limit: Int): List<Order> {
        val sql = "SELECT id, customer_id, status, total_cents, currency, ship_address, bill_address, notes, created_at, updated_at, placed_at, shipped_at, cancelled_at FROM orders WHERE customer_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?"
        val results = mutableListOf<Order>()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, customerId)
            ps.setInt(2, limit)
            ps.setInt(3, offset)
            ps.executeQuery().use { rs ->
                while (rs.next()) results.add(rsToOrder(rs))
            }
        }
        return results
    }

    fun update(o: Order): Order {
        val sql = """
            UPDATE orders
            SET status = ?, total_cents = ?, updated_at = ?, placed_at = ?, shipped_at = ?, cancelled_at = ?
            WHERE id = ?
        """.trimIndent()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, o.status.name)
            ps.setLong(2, o.totalCents)
            ps.setTimestamp(3, Instant.now().toSqlTimestamp())
            ps.setTimestamp(4, o.placedAt?.toSqlTimestamp())
            ps.setTimestamp(5, o.shippedAt?.toSqlTimestamp())
            ps.setTimestamp(6, o.cancelledAt?.toSqlTimestamp())
            ps.setString(7, o.id)
            val rows = ps.executeUpdate()
            if (rows == 0) throw NotFoundException("order ${o.id} not found")
        }
        return o
    }

    private fun rsToOrder(rs: java.sql.ResultSet) = Order(
        id = rs.getString("id"),
        customerId = rs.getString("customer_id"),
        status = OrderStatus.valueOf(rs.getString("status")),
        totalCents = rs.getLong("total_cents"),
        currency = rs.getString("currency"),
        shipAddress = rs.getString("ship_address"),
        billAddress = rs.getString("bill_address"),
        notes = rs.getString("notes"),
        createdAt = rs.getInstant("created_at")!!,
        updatedAt = rs.getInstant("updated_at")!!,
        placedAt = rs.getInstant("placed_at"),
        shippedAt = rs.getInstant("shipped_at"),
        cancelledAt = rs.getInstant("cancelled_at"),
    )
}

class OrderItemRepository(private val db: Database) {

    fun create(item: OrderItem): OrderItem {
        val sql = "INSERT INTO order_items (id, order_id, cat_id, quantity, price_cents, currency, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, item.id)
            ps.setString(2, item.orderId)
            ps.setString(3, item.catId)
            ps.setInt(4, item.quantity)
            ps.setLong(5, item.priceCents)
            ps.setString(6, item.currency)
            ps.setTimestamp(7, item.createdAt.toSqlTimestamp())
            ps.executeUpdate()
        }
        return item
    }

    fun findByOrder(orderId: String): List<OrderItem> {
        val sql = "SELECT id, order_id, cat_id, quantity, price_cents, currency, created_at FROM order_items WHERE order_id = ? ORDER BY created_at"
        val results = mutableListOf<OrderItem>()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, orderId)
            ps.executeQuery().use { rs ->
                while (rs.next()) results.add(rsToOrderItem(rs))
            }
        }
        return results
    }

    private fun rsToOrderItem(rs: java.sql.ResultSet) = OrderItem(
        id = rs.getString("id"),
        orderId = rs.getString("order_id"),
        catId = rs.getString("cat_id"),
        quantity = rs.getInt("quantity"),
        priceCents = rs.getLong("price_cents"),
        currency = rs.getString("currency"),
        createdAt = rs.getInstant("created_at")!!,
    )
}

class InventoryRepository(private val db: Database) {

    fun upsert(inv: Inventory) {
        val sql = """
            MERGE INTO inventory t
            USING (SELECT ? AS cat_id, ? AS available, ? AS reserved, ? AS updated_at) s
            ON t.cat_id = s.cat_id
            WHEN MATCHED THEN UPDATE SET available = s.available, reserved = s.reserved, updated_at = s.updated_at
            WHEN NOT MATCHED THEN INSERT (cat_id, available, reserved, updated_at) VALUES (s.cat_id, s.available, s.reserved, s.updated_at)
        """.trimIndent()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, inv.catId)
            ps.setInt(2, inv.available)
            ps.setInt(3, inv.reserved)
            ps.setTimestamp(4, inv.updatedAt.toSqlTimestamp())
            ps.executeUpdate()
        }
    }

    fun findByCat(catId: String): Inventory? {
        val sql = "SELECT cat_id, available, reserved, updated_at FROM inventory WHERE cat_id = ?"
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, catId)
            ps.executeQuery().use { rs ->
                if (rs.next()) return rsToInventory(rs)
            }
        }
        return null
    }

    private fun rsToInventory(rs: java.sql.ResultSet) = Inventory(
        catId = rs.getString("cat_id"),
        available = rs.getInt("available"),
        reserved = rs.getInt("reserved"),
        updatedAt = rs.getInstant("updated_at")!!,
    )
}

class AuditLogRepository(private val db: Database) {

    fun append(entry: AuditLog): AuditLog {
        val sql = "INSERT INTO audit_log (id, actor, action, target, created_at) VALUES (?, ?, ?, ?, ?)"
        db.connection().prepareStatement(sql).use { ps ->
            ps.setString(1, entry.id)
            ps.setString(2, entry.actor)
            ps.setString(3, entry.action)
            ps.setString(4, entry.target)
            ps.setTimestamp(5, entry.createdAt.toSqlTimestamp())
            ps.executeUpdate()
        }
        return entry
    }

    fun list(limit: Int): List<AuditLog> {
        val sql = "SELECT id, actor, action, target, created_at FROM audit_log ORDER BY created_at DESC LIMIT ?"
        val results = mutableListOf<AuditLog>()
        db.connection().prepareStatement(sql).use { ps ->
            ps.setInt(1, limit)
            ps.executeQuery().use { rs ->
                while (rs.next()) results.add(rsToAuditLog(rs))
            }
        }
        return results
    }

    private fun rsToAuditLog(rs: java.sql.ResultSet) = AuditLog(
        id = rs.getString("id"),
        actor = rs.getString("actor"),
        action = rs.getString("action"),
        target = rs.getString("target"),
        createdAt = rs.getInstant("created_at")!!,
    )
}
