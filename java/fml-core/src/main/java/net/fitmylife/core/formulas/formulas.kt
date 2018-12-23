package net.fitmylife.core.formulas

import net.fitmylife.core.enums.Gender

fun randomInt(min: Int, max: Int): Int = (Math.random() * (max - min) + min).toInt()


object Formulas {

    /** verified */
    fun calculateRMR(weight: Double, height: Double, age: Int, gender: Gender): Double = 9.99 * weight + 6.25 * height - 4.92 * age + gender.RMRFactor

    /** verified */
    fun calculateBMI(weight: Double, height: Double): Double = weight / Math.pow(height / 100, 2.0)

    /** verified */
    fun calculateMifflinStJeor(age: Int, weight: Double, height: Double, gender: Gender): Double =
            9.99 * weight + 6.25 * height - 4.92 * age + gender.RMRFactor

    // Using the Fitbit equation
    fun minuteOfActivity(weight: Double, heartRate: Int, gender: Gender): Double {
        return when (gender) {
            Gender.Female -> -5.92 + 0.0577 * heartRate - 0.0167 * weight + 0.00052 * heartRate * weight
            Gender.Male -> 3.56 - 0.0138 * heartRate - 0.1358 * weight + 0.00189 * heartRate * weight
        }
    }
}