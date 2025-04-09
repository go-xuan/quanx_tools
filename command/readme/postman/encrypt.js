// 前置脚本-对请求体中password明文进行加密
var body = JSON.parse(pm.request.body.raw)  // 获取请求体
console.log("加密前请求body:",body);
var pwd = body.password
const req = {
    url: 'http://localhost:9001/user/api/v1/encrypt?text='+pwd,
    method: "GET"
};

pm.sendRequest(req, function (err, res) {
    if (err) {
        console.log(err);
    }else {
        var real_data = {
            username: body.username,
            password: res.json().data
        };
        console.log("加密后请求body:",real_data);
        pm.request.body.update(JSON.stringify(real_data));
    }
});