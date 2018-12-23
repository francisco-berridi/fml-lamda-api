package net.fitmylife.core.analytics

import net.fitmylife.core.domain.body.BodyStats
import net.fitmylife.core.enums.Gender
import net.fitmylife.core.enums.ObesityCategory
import net.fitmylife.core.formulas.Formulas
import java.time.LocalDate
import java.time.Period

class BodyAnalytics(private val gender: Gender, private val dateOfBirth: LocalDate, private val height: Double, private val weight: Double) {

    fun calculateRMR(): Double {
        return Formulas.calculateRMR(weight, height, calculateAge(), gender)
    }

    fun calculateAge(): Int {
        return Period.between(dateOfBirth, LocalDate.now()).years
    }

    fun calculateBMI(): Double {
        return Formulas.calculateBMI(weight, height)
    }

    fun maxHeartRate(): Int {
        return 220 - calculateAge()
    }

    private fun averageHeartRateDuringWorkouts(): Int {
        return Math.round(maxHeartRate() * 0.62).toInt()
    }

    private fun calculateMifflinStJeor(): Double {
        return Formulas.calculateMifflinStJeor(calculateAge(), weight, height, gender)
    }

    fun generateBodyStats(): BodyStats {
        return BodyStats(
                height = height,
                weight = weight,
                gender = gender,
                age = calculateAge(),
                RMR = calculateRMR(),
                BMI = calculateBMI(),
                maxHeartRate = maxHeartRate(),
                obesityCategory = ObesityCategory.getWithBmi(calculateBMI()),
                mifflinStJeorCalculation = calculateMifflinStJeor(),
                averageHRDuringWorkouts = averageHeartRateDuringWorkouts()

        )
    }
}