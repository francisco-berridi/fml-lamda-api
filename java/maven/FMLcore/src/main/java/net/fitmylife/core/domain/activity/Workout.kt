package net.fitmylife.core.domain.activity

import net.fitmylife.core.enums.*
import java.time.Duration


data class Workout(
        val id: String,
        val name: String,
        val summary: String,
        val requiresGym: Boolean,
        val metValue: Int,
        val category: WorkoutCategory,
        val targetGoals: List<Goal>,
        val difficulty: Int,
        val aerobicCoefficient: Double?,
        val execution: Execution?,
        val targetExperience: Experience,
        val bmiRestriction: ObesityCategory,
        val duration: Duration,
        val cardioOption: List<CardioOption>,
        val intensity: Double,
        val resistanceProportion: Double,
        val cardioProportion: Double,
        val priority: Int?
)