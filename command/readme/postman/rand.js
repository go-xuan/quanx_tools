// 生成随机数并设置postman的环境变量值
const randRequest = {
    url: 'http://localhost:9999/tools/rand/struct',
    method: "POST",
    header: 'Content-Type: application/json',
    body: {
        mode: 'raw',
        raw: JSON.stringify([
            {name:"name", type: 'name', constraint: ''},
            {name:"no", type: 'int', constraint: 'min=10000&max=1000000'},
            {name:"year", type: 'int', constraint: 'min=2000&max=2023'},
            {name:"string", type: 'string', constraint: ''},
            {name:"id_card", type: 'id_card', constraint: ''},
            {name:"phone", type: 'phone', constraint: ''},
            {name:"uuid", type: 'uuid', constraint: 'old=-&new='},
            {name:"int", type: 'int', constraint: 'min=11&max=15'},
            {name:"date", type: 'date', constraint: 'min=2024-01-01'},
            {name:"time", type: 'time', constraint: 'min=2024-01-01'},
            {name:"gender", type: 'option', constraint: 'options=男,女,未知'},
            {name:"bool_cn", type: 'option', constraint: 'options=是,否'},
            {name:"bool", type: 'option', constraint: 'options=true,false'}
        ]),
    }
};

pm.sendRequest(randRequest, function (err, res) {
    if (err) {
        console.log(err);
    }else {
        console.log(res.json());
        randData = res.json().data
        pm.variables.set("name",randData.name);
        pm.variables.set("no",randData.no);
        pm.variables.set("year",randData.year);
        pm.variables.set("string",randData.string);
        pm.variables.set("id_card",randData.id_card);
        pm.variables.set("phone",randData.phone);
        pm.variables.set("uuid",randData.uuid);
        pm.variables.set("int",randData.int);
        pm.variables.set("date",randData.date);
        pm.variables.set("time",randData.time);
        pm.variables.set("gender",randData.gender);
        pm.variables.set("bool_cn",randData.bool_cn);
        pm.variables.set("bool",randData.bool);
    }
});