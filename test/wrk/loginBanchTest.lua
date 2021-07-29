-- 启动阶段 （每个线程执行一次）
index = 0
request = function()
--     local user_index = math.random(200)
    local user_index = index % 200
    index = index + 1

    wrk.method = "POST"
    wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
    wrk.body   = "username=" .. user_index .. "u&password=1"
--     wrk.headers["Cookie"] = "token=" .. user_index .. "u_token"
    return wrk.format("POST", "/login")
end