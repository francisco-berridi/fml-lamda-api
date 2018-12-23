package net.fitmylife.core.services

import net.fitmylife.core.domain.activity.Workout
import net.fitmylife.core.domain.activity.WorkoutSchedule
import net.fitmylife.core.enums.*
import net.fitmylife.core.extensions.decimals
import org.json.JSONArray
import org.json.JSONObject
import java.lang.Exception
import java.time.Duration

object ParseService {

    fun getAllSchedules(obesityCategory: ObesityCategory, primaryGoal: Goal, secondaryGoal: Goal): ArrayList<WorkoutSchedule> {

        val all = arrayListOf<WorkoutSchedule>()

        val response = khttp.get(
                url = "https://dev.api.fitmylife.net/parse/classes/WorkoutSchedule",
                params = mapOf("limit" to "999"),
                headers = mapOf(
                        "X-Parse-Application-Id" to "fitmylifeAppId",
                        "X-Parse-REST-API-Key" to "flrcGF4lP4Lx6M2vq0rLo4MiAM2IN4BHRTr4QQKQ"))

        for (obj in response.jsonObject.get("results") as JSONArray) {
            println(obj)
            if (obj is JSONObject) {

            }
        }

        return all

    }

    fun getAllWorkouts(): ArrayList<Workout> {

        val all = arrayListOf<Workout>()

        val response = khttp.get(
                url = "https://dev.api.fitmylife.net/parse/classes/Workout",
                params = mapOf("limit" to "999"),
                headers = mapOf(
                        "X-Parse-Application-Id" to "fitmylifeAppId",
                        "X-Parse-REST-API-Key" to "flrcGF4lP4Lx6M2vq0rLo4MiAM2IN4BHRTr4QQKQ"))

        for (obj in response.jsonObject.get("results") as JSONArray) {
            if (obj is JSONObject) {

                // some workouts dont have sub category, exclude
                try {
                    obj.getString("subcategory")
                } catch (e: Exception) {
                    continue
                }


                all.add(parseWorkout(obj))
            }
        }

        return all
    }
}

private fun parseWorkout(obj: JSONObject): Workout {

    return Workout(
            id = obj.getString("objectId"),
            summary = try {
                obj.getString("summary")
            } catch (e: Exception) {
                ""
            },
            category = WorkoutCategory.fromParse(obj.getString("subcategory")),
            difficulty = obj.getInt("difficulty"),
            name = obj.getString("name"),
            duration = Duration.ofSeconds(obj.getLong("estimatedDuration")),
            targetExperience = parseTargetExperience(obj.getInt("difficulty")),
            metValue = obj.getInt("metValue"),
            resistanceProportion = obj.getDouble("resistanceComponent") / 100,
            cardioProportion = (1 - obj.getDouble("resistanceComponent") / 100).decimals(3),
            cardioOption = try {
                parseCarioOptions(obj.getJSONArray("cardioPreferences"))
            } catch (e: Exception) {
                arrayListOf<CardioOption>()
            },
            targetGoals = try {
                parseTargetGoals(obj.getJSONArray("goalRelated"))
            } catch (e: Exception) {
                arrayListOf<Goal>()
            },
            intensity = obj.getDouble("estimatedHeartRatePercentage") / 100,
            requiresGym = obj.getBoolean("requiresGym"),
            execution = null,
            priority = null,
            bmiRestriction = try {
                parseBmiRestriction(obj.getJSONArray("excludeObesityCategory"))
            } catch (e: Exception) {
                ObesityCategory.UnderWeight
            },
            aerobicCoefficient = null
    )
}


private fun parseTargetExperience(exp: Int): Experience {
    return when (exp) {
        in Int.MIN_VALUE until 10 -> Experience.None
        in 10..12 -> Experience.Beginner
        in 12..16 -> Experience.Intermediate
        in 17..Int.MAX_VALUE -> Experience.Advanced
        else -> Experience.None
    }
}

private fun parseCarioOptions(array: JSONArray): List<CardioOption> {

    val options = arrayListOf<CardioOption>()

    for (obj in array) {
        if (obj is String) {
            CardioOption.fromParse(obj)?.let {
                options.add(it)
            }
        }
    }

    return options
}

private fun parseTargetGoals(array: JSONArray): List<Goal> {
    val goals = arrayListOf<Goal>()

    for (obj in array) {
        if (obj is String) {
            Goal.fromParse(obj)?.let {
                goals.add(it)
            }
        }
    }

    return goals
}

private fun parseBmiRestriction(array: JSONArray): ObesityCategory {

    var restriction: ObesityCategory = ObesityCategory.UnderWeight

    for (obj in array) {
        if (obj is String) {

            ObesityCategory.fromParse(obj)?.let {
                if (it > restriction) {
                    restriction = it
                }
            }
        }
    }

    return restriction

}