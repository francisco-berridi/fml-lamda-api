package net.fitmylife.core.domain.activity

import net.fitmylife.core.enums.ExecutionType

data class Execution(
        val type: ExecutionType,
        val steps: ArrayList<ExecutionStep>,
        val repetition: Int
)