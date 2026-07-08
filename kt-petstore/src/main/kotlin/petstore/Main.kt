package petstore

import petstore.db.Database
import petstore.db.SnowflakeConfig
import petstore.model.*
import petstore.repository.*
import petstore.service.*
import java.time.LocalDate

fun main() {
    val config = SnowflakeConfig(
        account = System.getenv("SNOWFLAKE_ACCOUNT"),
        user = System.getenv("SNOWFLAKE_USER"),
        password = System.getenv("SNOWFLAKE_PASSWORD"),
        database = "TEST",
        schema = "PETSTORE",
        warehouse = System.getenv("SNOWFLAKE_WAREHOUSE"),
        role = System.getenv("SNOWFLAKE_ROLE"),
    )

    Database(config).use { db ->
        val customerRepo = CustomerRepository(db)
        val breedRepo = BreedRepository(db)
        val shelterRepo = ShelterRepository(db)
        val catRepo = CatRepository(db)
        val cartRepo = CartRepository(db)
        val cartItemRepo = CartItemRepository(db)
        val orderRepo = OrderRepository(db)
        val orderItemRepo = OrderItemRepository(db)
        val inventoryRepo = InventoryRepository(db)
        val auditRepo = AuditLogRepository(db)

        val inventoryService = InventoryService(inventoryRepo)
        val auditService = AuditService(auditRepo)
        val catService = CatService(catRepo, breedRepo, shelterRepo, inventoryService)
        val cartService = CartService(cartRepo, cartItemRepo, inventoryService)
        val orderService = OrderService(orderRepo, orderItemRepo, catRepo, inventoryService, auditService)
        val customerService = CustomerService(customerRepo)
        val breedService = BreedService(breedRepo)
        val shelterService = ShelterService(shelterRepo)

        println("=== Cat Adoption Store ===")

        val shelter = shelterService.create(
            name = "Whiskers Rescue",
            address = "123 Purr Lane, Catherton",
            phone = "555-0123",
            email = "adopt@whiskers.org",
        )
        println("Created shelter: $shelter")

        val breed = breedService.create(
            name = "Maine Coon",
            description = "Large, gentle, fluffy",
            originCountry = "USA",
            avgLifespanYears = 13,
        )
        println("Created breed: $breed")

        val breed2 = breedService.create(
            name = "Siamese",
            description = "Vocal, sleek, social",
            originCountry = "Thailand",
            avgLifespanYears = 15,
        )
        println("Created breed: $breed2")

        val customer = customerService.create(
            fullName = "Alice Johnson",
            email = "alice@example.com",
            phone = "555-9999",
            address = "456 Felidae St, Meowville",
        )
        println("Created customer: $customer")

        val cat1 = catService.create(
            name = "Luna",
            breedId = breed.id,
            shelterId = shelter.id,
            birthDate = LocalDate.of(2023, 5, 14),
            priceCents = 25000,
            currency = "USD",
            description = "Friendly, loves cuddles",
        )
        println("Created cat: ${cat1.name} (${cat1.status})")

        val cat2 = catService.create(
            name = "Simba",
            breedId = breed2.id,
            shelterId = shelter.id,
            birthDate = LocalDate.of(2022, 8, 1),
            priceCents = 30000,
            currency = "USD",
            description = "Very talkative, loves warmth",
        )
        println("Created cat: ${cat2.name} (${cat2.status})")

        val cat3 = catService.create(
            name = "Luna Bell",
            breedId = breed.id,
            shelterId = shelter.id,
            birthDate = LocalDate.of(2024, 1, 20),
            priceCents = 20000,
            currency = "USD",
            description = "Playful kitten",
        )
        println("Created cat: ${cat3.name} (${cat3.status})")

        println("\nSearch results for 'Luna':")
        catService.search("Luna").forEach { println("  - ${it.name} (id=${it.id})") }

        println("\nInventory for ${cat1.name}: ${inventoryService.get(cat1.id)}")

        val cart = cartService.create(customer.id)
        cartService.addItem(cart.id, cat1.id, 1)
        cartService.addItem(cart.id, cat3.id, 1)
        println("\nCart ${cart.id} items:")
        cartService.items(cart.id).forEach { println("  - catId=${it.catId} qty=${it.quantity}") }

        val order = orderService.place(
            customerId = customer.id,
            items = listOf(cat1.id to 1, cat3.id to 1),
            currency = "USD",
            shipAddress = customer.address,
            billAddress = customer.address,
            notes = "Please deliver in the morning",
        )
        println("\nPlaced order: ${order.id} status=${order.status} total=${order.totalCents}")
        println("Order items:")
        orderService.items(order.id).forEach { println("  - catId=${it.catId} qty=${it.quantity} price=${it.priceCents}") }

        val paid = orderService.markPaid(order.id)
        println("Order paid: status=${paid.status}")

        val shipped = orderService.ship(order.id)
        println("Order shipped: status=${shipped.status} shippedAt=${shipped.shippedAt}")

        catService.markAdopted(cat1.id)
        catService.markAdopted(cat3.id)
        println("\n${cat1.name} status: ${catService.get(cat1.id).status}")
        println("${cat3.name} status: ${catService.get(cat3.id).status}")

        println("\nRecent audit log:")
        auditService.recent(10).forEach { e ->
            println("  - [${e.createdAt}] ${e.actor} ${e.action} -> ${e.target}")
        }

        println("\nAll breeds:")
        breedService.list().forEach { println("  - ${it.name} (${it.originCountry})" ) }

        println("\nAll shelters:")
        shelterService.list().forEach { println("  - ${it.name} (${it.address})" ) }

        println("\nCustomer orders:")
        customerService.list().forEach { c ->
            println("  Customer: ${c.fullName}")
            orderService.listByCustomer(c.id).forEach { o ->
                println("    Order ${o.id}: ${o.status} total=${o.totalCents} ${o.currency}")
            }
        }

        cartService.clear(cart.id)
        println("\nCart cleared: ${cart.id}")

        println("\n=== Done ===")
    }
}
