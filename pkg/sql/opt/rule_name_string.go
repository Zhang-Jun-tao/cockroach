// Code generated by "stringer -output=rule_name_string.go -type=RuleName rule_name.go rule_name.og.go"; DO NOT EDIT.

package opt

import "strconv"

const _RuleName_name = "InvalidRuleNameNumManualRuleNamesEliminateEmptyAndEliminateEmptyOrEliminateSingletonAndOrSimplifyAndSimplifyOrSimplifyFiltersFoldNullAndOrNegateComparisonEliminateNotNegateAndNegateOrCommuteVarInequalityCommuteConstInequalityNormalizeCmpPlusConstNormalizeCmpMinusConstNormalizeCmpConstMinusNormalizeTupleEqualityFoldNullComparisonLeftFoldNullComparisonRightFoldIsNullFoldNonNullIsNullFoldIsNotNullFoldNonNullIsNotNullCommuteNullIsDecorrelateJoinTryDecorrelateSelectHoistSelectExistsHoistSelectNotExistsHoistSelectSubqueryHoistProjectSubqueryHoistJoinSubqueryHoistValuesSubqueryNormalizeAnyFilterNormalizeNotAnyFilterEliminateDistinctEliminateGroupByProjectPushSelectIntoInlinableProjectInlineProjectInProjectEnsureJoinFiltersAndEnsureJoinFiltersPushFilterIntoJoinLeftPushFilterIntoJoinRightSimplifyLeftJoinSimplifyRightJoinEliminateSemiJoinEliminateAntiJoinEliminateJoinNoColsLeftEliminateJoinNoColsRightEliminateLimitPushLimitIntoProjectPushOffsetIntoProjectEliminateMax1RowFoldPlusZeroFoldZeroPlusFoldMinusZeroFoldMultOneFoldOneMultFoldDivOneInvertMinusEliminateUnaryMinusEliminateProjectEliminateProjectProjectPruneProjectColsPruneScanColsPruneSelectColsPruneLimitColsPruneOffsetColsPruneJoinLeftColsPruneJoinRightColsPruneAggColsPruneGroupByColsPruneValuesColsPruneRowNumberColsCommuteVarCommuteConstEliminateCoalesceSimplifyCoalesceEliminateCastFoldNullCastFoldNullUnaryFoldNullBinaryLeftFoldNullBinaryRightFoldNullInNonEmptyFoldNullInEmptyFoldNullNotInEmptyNormalizeInConstFoldInNullEliminateExistsProjectEliminateExistsGroupByEliminateSelectEnsureSelectFiltersAndEnsureSelectFiltersMergeSelectsPushSelectIntoProjectPushSelectIntoJoinLeftPushSelectIntoJoinRightMergeSelectInnerJoinPushSelectIntoGroupByPushLimitIntoScanGenerateIndexScansConstrainScanConstrainLookupJoinIndexScanNumRuleNames"

var _RuleName_index = [...]uint16{0, 15, 33, 50, 66, 89, 100, 110, 125, 138, 154, 166, 175, 183, 203, 225, 246, 268, 290, 312, 334, 357, 367, 384, 397, 417, 430, 445, 465, 482, 502, 521, 541, 558, 577, 595, 616, 633, 656, 686, 708, 728, 745, 767, 790, 806, 823, 840, 857, 880, 904, 918, 938, 959, 975, 987, 999, 1012, 1023, 1034, 1044, 1055, 1074, 1090, 1113, 1129, 1142, 1157, 1171, 1186, 1203, 1221, 1233, 1249, 1264, 1282, 1292, 1304, 1321, 1337, 1350, 1362, 1375, 1393, 1412, 1430, 1445, 1463, 1479, 1489, 1511, 1533, 1548, 1570, 1589, 1601, 1622, 1644, 1667, 1687, 1708, 1725, 1743, 1756, 1784, 1796}

func (i RuleName) String() string {
	if i >= RuleName(len(_RuleName_index)-1) {
		return "RuleName(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RuleName_name[_RuleName_index[i]:_RuleName_index[i+1]]
}