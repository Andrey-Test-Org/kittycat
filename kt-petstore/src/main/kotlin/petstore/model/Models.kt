package petstore.model

import java.time.Instant

data class Customer(
    val id: String,
    val fullName: String,
    val email: String,
    val phone: String?,
    val address: String,
    val createdAt: Instant = Instant.now(),
    val updatedAt: Instant = Instant.now(),
)

data class Breed(
    val id: String,
    val name: String,
    val description: String?,
    val originCountry: String,
    val avgLifespanYears: Int,
    val createdAt: Instant = Instant.now(),
    val updatedAt: Instant = Instant.now(),
)

data class Shelter(
    val id: String,
    val name: String,
    val address: String,
    val phone: String?,
    val email: String?,
    val createdAt: Instant = Instant.now(),
    val updatedAt: Instant = Instant.now(),
)

enum class CatStatus { AVAILABLE, RESERVED, ADOPTED, RETIRED }

data class Cat(
    val id: String,
    val name: String,
    val breedId: String,
    val shelterId: String,
    val birthDate: java.time.LocalDate,
    val priceCents: Long,
    val currency: String,
    val status: CatStatus = CatStatus.AVAILABLE,
    val description: String?,
    val createdAt: Instant = Instant.now(),
    val updatedAt: Instant = Instant.now(),
)

data class Cart(
    val id: String,
    val customerId: String,
    val createdAt: Instant = Instant.now(),
    val updatedAt: Instant = Instant.now(),
)

data class CartItem(
    val id: String,
    val cartId: String,
    val catId: String,
    val quantity: Int,
    val createdAt: Instant = Instant.now(),
)

enum class OrderStatus { PENDING, PAID, SHIPPED, CANCELLED }

data class Order(
    val id: String,
    val customerId: String,
    val status: OrderStatus = OrderStatus.PENDING,
    val totalCents: Long,
    val currency: String,
    val shipAddress: String,
    val billAddress: String,
    val notes: String? = null,
    val createdAt: Instant = Instant.now(),
    val updatedAt: Instant = Instant.now(),
    val placedAt: Instant? = null,
    val shippedAt: Instant? = null,
    val cancelledAt: Instant? = null,
)

data class OrderItem(
    val id: String,
    val orderId: String,
    val catId: String,
    val quantity: Int,
    val priceCents: Long,
    val currency: String,
    val createdAt: Instant = Instant.now(),
)

data class Inventory(
    val catId: String,
    val available: Int,
    val reserved: Int,
    val updatedAt: Instant = Instant.now(),
)

data class AuditLog(
    val id: String,
    val actor: String,
    val action: String,
    val target: String,
    val createdAt: Instant = Instant.now(),
)
