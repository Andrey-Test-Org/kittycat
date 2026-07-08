package petstore.service

import petstore.model.*
import petstore.repository.*
import java.time.Instant
import java.util.UUID

class CustomerService(private val repo: CustomerRepository) {

    fun create(fullName: String, email: String, phone: String?, address: String): Customer {
        require(fullName.isNotBlank()) { "fullName must not be blank" }
        require(email.isNotBlank()) { "email must not be blank" }
        val customer = Customer(
            id = newId(),
            fullName = fullName,
            email = email,
            phone = phone,
            address = address,
        )
        return repo.create(customer)
    }

    fun get(id: String): Customer =
        repo.findById(id) ?: throw NotFoundException("customer $id not found")

    fun list(offset: Int = 0, limit: Int = 25): List<Customer> = repo.list(offset, limit)

    fun update(id: String, fullName: String?, email: String?, phone: String?, address: String?): Customer {
        val existing = get(id)
        val updated = existing.copy(
            fullName = fullName ?: existing.fullName,
            email = email ?: existing.email,
            phone = phone ?: existing.phone,
            address = address ?: existing.address,
        )
        return repo.update(updated)
    }

    fun delete(id: String) = repo.delete(id)
}

class BreedService(private val repo: BreedRepository) {

    fun create(name: String, description: String?, originCountry: String, avgLifespanYears: Int): Breed {
        require(name.isNotBlank()) { "name must not be blank" }
        require(avgLifespanYears > 0) { "avgLifespanYears must be positive" }
        return repo.create(Breed(
            id = newId(),
            name = name,
            description = description,
            originCountry = originCountry,
            avgLifespanYears = avgLifespanYears,
        ))
    }

    fun get(id: String): Breed =
        repo.findById(id) ?: throw NotFoundException("breed $id not found")

    fun list(): List<Breed> = repo.list()
}

class ShelterService(private val repo: ShelterRepository) {

    fun create(name: String, address: String, phone: String?, email: String?): Shelter {
        require(name.isNotBlank()) { "name must not be blank" }
        return repo.create(Shelter(
            id = newId(),
            name = name,
            address = address,
            phone = phone,
            email = email,
        ))
    }

    fun get(id: String): Shelter =
        repo.findById(id) ?: throw NotFoundException("shelter $id not found")

    fun list(): List<Shelter> = repo.list()
}

class InventoryService(private val repo: InventoryRepository) {

    fun stock(catId: String, available: Int) {
        require(available >= 0) { "available must be >= 0" }
        val existing = repo.findByCat(catId) ?: Inventory(catId, 0, 0)
        repo.upsert(existing.copy(available = available, updatedAt = Instant.now()))
    }

    fun reserve(catId: String, qty: Int) {
        require(qty > 0) { "qty must be > 0" }
        val inv = repo.findByCat(catId) ?: Inventory(catId, 0, 0)
        require(inv.available >= qty) { "not enough stock: have ${inv.available}, need $qty" }
        repo.upsert(inv.copy(
            available = inv.available - qty,
            reserved = inv.reserved + qty,
            updatedAt = Instant.now(),
        ))
    }

    fun release(catId: String, qty: Int) {
        require(qty > 0) { "qty must be > 0" }
        val inv = repo.findByCat(catId) ?: return
        val newReserved = maxOf(0, inv.reserved - qty)
        val released = inv.reserved - newReserved
        repo.upsert(inv.copy(
            reserved = newReserved,
            available = inv.available + released,
            updatedAt = Instant.now(),
        ))
    }

    fun fulfil(catId: String, qty: Int) {
        require(qty > 0) { "qty must be > 0" }
        val inv = repo.findByCat(catId) ?: return
        val newReserved = maxOf(0, inv.reserved - qty)
        repo.upsert(inv.copy(
            reserved = newReserved,
            updatedAt = Instant.now(),
        ))
    }

    fun get(catId: String): Inventory =
        repo.findByCat(catId) ?: Inventory(catId, 0, 0)
}

class CatService(
    private val repo: CatRepository,
    private val breedRepo: BreedRepository,
    private val shelterRepo: ShelterRepository,
    private val inventoryService: InventoryService,
) {

    fun create(
        name: String,
        breedId: String,
        shelterId: String,
        birthDate: java.time.LocalDate,
        priceCents: Long,
        currency: String,
        description: String?,
    ): Cat {
        require(name.isNotBlank()) { "name must not be blank" }
        require(priceCents > 0) { "priceCents must be > 0" }
        require(breedRepo.findById(breedId) != null) { "breed $breedId not found" }
        require(shelterRepo.findById(shelterId) != null) { "shelter $shelterId not found" }
        val cat = Cat(
            id = newId(),
            name = name,
            breedId = breedId,
            shelterId = shelterId,
            birthDate = birthDate,
            priceCents = priceCents,
            currency = currency,
            description = description,
        )
        val created = repo.create(cat)
        inventoryService.stock(created.id, 1)
        return created
    }

    fun get(id: String): Cat =
        repo.findById(id) ?: throw NotFoundException("cat $id not found")

    fun list(offset: Int = 0, limit: Int = 25): List<Cat> = repo.list(offset, limit)

    fun search(query: String, limit: Int = 25): List<Cat> = repo.searchByName(query, limit)

    fun markAdopted(id: String) = repo.updateStatus(id, CatStatus.ADOPTED)

    fun markRetired(id: String) = repo.updateStatus(id, CatStatus.RETIRED)

    fun delete(id: String) = repo.delete(id)
}

class CartService(
    private val cartRepo: CartRepository,
    private val cartItemRepo: CartItemRepository,
    private val inventoryService: InventoryService,
) {

    fun create(customerId: String): Cart {
        return cartRepo.create(Cart(
            id = newId(),
            customerId = customerId,
        ))
    }

    fun get(id: String): Cart =
        cartRepo.findById(id) ?: throw NotFoundException("cart $id not found")

    fun addItem(cartId: String, catId: String, quantity: Int): CartItem {
        require(quantity > 0) { "quantity must be > 0" }
        val inv = inventoryService.get(catId)
        require(inv.available >= quantity) { "not enough stock for cat $catId" }
        return cartItemRepo.create(CartItem(
            id = newId(),
            cartId = cartId,
            catId = catId,
            quantity = quantity,
        ))
    }

    fun removeItem(cartId: String, catId: String) =
        cartItemRepo.delete(cartId, catId)

    fun items(cartId: String): List<CartItem> =
        cartItemRepo.findByCart(cartId)

    fun clear(cartId: String) {
        cartItemRepo.deleteAllByCart(cartId)
        cartRepo.delete(cartId)
    }
}

class AuditService(private val repo: AuditLogRepository) {

    fun record(actor: String, action: String, target: String): AuditLog {
        return repo.append(AuditLog(
            id = newId(),
            actor = actor,
            action = action,
            target = target,
        ))
    }

    fun recent(limit: Int = 50): List<AuditLog> = repo.list(limit)
}

class OrderService(
    private val orderRepo: OrderRepository,
    private val orderItemRepo: OrderItemRepository,
    private val catRepo: CatRepository,
    private val inventoryService: InventoryService,
    private val auditService: AuditService,
) {

    fun place(
        customerId: String,
        items: List<Pair<String, Int>>,
        currency: String,
        shipAddress: String,
        billAddress: String,
        notes: String? = null,
    ): Order {
        require(items.isNotEmpty()) { "order must have at least one item" }
        var total = 0L
        val orderId = newId()
        val now = Instant.now()
        for ((catId, qty) in items) {
            inventoryService.reserve(catId, qty)
        }
        for ((catId, qty) in items) {
            val cat = catRepo.findById(catId) ?: throw NotFoundException("cat $catId not found")
            require(cat.currency == currency) { "currency mismatch: cat $catId uses ${cat.currency}, order uses $currency" }
            val pricePerUnit = cat.priceCents
            orderItemRepo.create(OrderItem(
                id = newId(),
                orderId = orderId,
                catId = catId,
                quantity = qty,
                priceCents = pricePerUnit,
                currency = currency,
            ))
            total += pricePerUnit * qty
        }
        val order = Order(
            id = orderId,
            customerId = customerId,
            status = OrderStatus.PENDING,
            totalCents = total,
            currency = currency,
            shipAddress = shipAddress,
            billAddress = billAddress,
            notes = notes,
            placedAt = now,
        )
        val saved = orderRepo.create(order)
        auditService.record(customerId, "order.place", orderId)
        return saved
    }

    fun markPaid(id: String): Order {
        val order = get(id)
        require(order.status == OrderStatus.PENDING) { "order must be PENDING to pay" }
        return orderRepo.update(order.copy(status = OrderStatus.PAID))
    }

    fun ship(id: String): Order {
        val order = get(id)
        require(order.status == OrderStatus.PAID) { "order must be PAID to ship" }
        for (item in orderItemRepo.findByOrder(id)) {
            inventoryService.fulfil(item.catId, item.quantity)
        }
        return orderRepo.update(order.copy(
            status = OrderStatus.SHIPPED,
            shippedAt = Instant.now(),
        ))
    }

    fun cancel(id: String): Order {
        val order = get(id)
        require(order.status in listOf(OrderStatus.PENDING, OrderStatus.PAID)) { "order cannot be cancelled in state ${order.status}" }
        for (item in orderItemRepo.findByOrder(id)) {
            inventoryService.release(item.catId, item.quantity)
        }
        return orderRepo.update(order.copy(
            status = OrderStatus.CANCELLED,
            cancelledAt = Instant.now(),
        ))
    }

    fun get(id: String): Order =
        orderRepo.findById(id) ?: throw NotFoundException("order $id not found")

    fun listByCustomer(customerId: String, offset: Int = 0, limit: Int = 25): List<Order> =
        orderRepo.listByCustomer(customerId, offset, limit)

    fun items(orderId: String): List<OrderItem> =
        orderItemRepo.findByOrder(orderId)
}

private fun newId(): String = UUID.randomUUID().toString().replace("-", "")
