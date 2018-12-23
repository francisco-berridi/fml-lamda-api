package net.fitmylife.core.domain.activity

import net.fitmylife.core.enums.*
import java.time.DayOfWeek
import java.time.Duration


class WorkoutSchedule {

    val routine: Map<DayOfWeek, Pair<WorkoutCategory, Workout?>> = HashMap()
    val averagaeCaloriesBurned: Double? = null
}
