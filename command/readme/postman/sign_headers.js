// 前置脚本，生成签名并植入headers
var headers = pm.request.headers; // 获取请求头
var product_id = headers.get('product_id');
var product_key = headers.get('product_key');
var timestamp = Math.floor(Date.now() / 1000);
var stringSignTemp = `product=${product_id}&time=${timestamp}&call_key=${product_key}`;
var sign = CryptoJS.MD5(stringSignTemp).toString(CryptoJS.enc.Hex);
headers.add({key: 'time', value: timestamp});
headers.add({key: 'sign', value: sign});
console.log("加密后请求headers:", pm.request.headers);