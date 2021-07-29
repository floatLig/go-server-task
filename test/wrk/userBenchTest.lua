-- 启动阶段 （每个线程执行一次）
index = 0
request = function()
--     local user_index = math.random(200)
    local user_index = index % 200
    index = index + 1

    wrk.method = "POST"
    wrk.headers["Cookie"] = "token=" .. user_index .. "u_token"
    return wrk.format("POST", "/info")
end