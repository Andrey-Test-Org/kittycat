package com.kittycat.dashboard

data class KittyEvent(
    val catId: String,
    val kind: String,
    val weightG: Int,
    val atMillis: Long
)

data class KittyStats(
    val catId: String,
    val totalEvents: Int,
    val avgWeightG: Double,
    val lastSeenMillis: Long,
    val kinds: List<String>
)

object Dashboard {
    fun aggregate(events: List<KittyEvent>): Map<String, KittyStats> {
        return events.groupBy { it.catId }.mapValues { (catId, group) ->
            val avg = group.sumOf { it.weightG }.toDouble() / group.size
            val last = group.maxOf { it.atMillis }
            val kinds = group.map { it.kind }.distinct()
            KittyStats(catId, group.size, avg, last, kinds)
        }
    }
}
