import net.fitmylife.core.enums.Goal
import net.fitmylife.core.enums.ObesityCategory
import net.fitmylife.core.services.ParseService

fun main(args: Array<String>) {


    for (schedule in ParseService.getAllSchedules(ObesityCategory.NormalWeight, Goal.WeightLoss, Goal.Circuit)) {
        println(schedule)
    }

    /*for (workout in  ParseService.getAllWorkouts()) {
        println(workout)
    }*/
}