-- 数据流封装
json = require("json")
appx = require("appx")

function HandlerDataStream(mn)
    local regs, err = appx.getRegisters(mn)
    if err ~= nil then
        error(err)
    end
    local data = {}
    for i, v in ipairs(regs) do
        data[v.factor] = tonumber(v.value)
    end
    return json.encode({ datatime=os.date("%Y%m%d%H%M%S",os.time()), data = data, mn = mn})
end