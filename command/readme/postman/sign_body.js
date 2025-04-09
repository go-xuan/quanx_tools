// 前置脚本，生成签名并植入body
var body = JSON.parse(pm.request.body.raw)  // 获取请求体
console.log("加密前请求body:",body);
var product_id = body.product_id
var product_key = body.product_key
var timestamp = Math.floor(Date.now() / 1000);
var stringSignTemp = `product=${product_id}&time=${timestamp}&call_key=${product_key}`;
var sign = CryptoJS.MD5(stringSignTemp).toString(CryptoJS.enc.Hex);
body.sign = sign
body.time = timestamp
console.log("加密后请求body:",body);
pm.request.body.update(JSON.stringify(body));