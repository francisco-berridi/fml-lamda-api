
import net.fitmylife.core.analytics.BodyAnalytics
import net.fitmylife.core.enums.Gender
import net.fitmylife.core.extensions.cm
import net.fitmylife.core.extensions.lb
import java.time.LocalDate
import java.time.Month

fun main(args: Array<String>) {

    val body = BodyAnalytics(Gender.Male, LocalDate.of(1982, Month.SEPTEMBER, 26), 181.cm, 200.lb)
    println(body.generateBodyStats())
}