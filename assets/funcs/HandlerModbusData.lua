-- 计算传感器值 计算公式 = （（采集值-最小取样值）/（最大取样值-最小取样值）*（最大量程-最小量程） ）  + 最小量程
function HandlerModbusData(src)
	local minSpVal = 819 --  最小取样值 如 819
	local maxSpVal = 4095 --  最大取样值 如 4095
	local minRangeVal = 0 -- 最小量程 0
	local maxRangeVal = 300  -- 最大量程 300

	if (src < minSpVal )
	then
		error("原始值小于最小取样值")
	end

	local result = ((src-minSpVal)/(maxSpVal-minSpVal))*(maxRangeVal-minRangeVal) + (minRangeVal)

	if (result > maxRangeVal or result < minRangeVal )
	then
		error("计算结果超出范围")
	end

	return result*100 / 100
end

-- print(HandlerModbusData(1819))