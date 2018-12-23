package net.fitmylife.core.domain.body

import net.fitmylife.core.enums.Gender
import net.fitmylife.core.enums.ObesityCategory

data class BodyStats(
        val height: Double,
        val weight: Double,
        val gender: Gender,
        val age: Int,
        val RMR: Double,
        val BMI: Double,
        val maxHeartRate: Int,
        val obesityCategory: ObesityCategory,
        val mifflinStJeorCalculation: Double,
        val averageHRDuringWorkouts: Int
)