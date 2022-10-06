function hasPerms(session, perms) {
    if (session.level === 'super') {
        return true;
    }
    flag = false;
    for (var i in perms) {
        if (session.perm_names.indexOf(perms[i]) !== -1) {
            flag = true;
            break;
        }
    }
    return flag;
}

var format_date = function (date, fmt) {
    var o = {
        "y+": date.getFullYear(),
        "M+": date.getMonth() + 1,                 //月份
        "d+": date.getDate(),                    //日
        "h+": date.getHours(),                   //小时
        "m+": date.getMinutes(),                 //分
        "s+": date.getSeconds(),                 //秒
        "q+": Math.floor((date.getMonth() + 3) / 3), //季度
        "S+": date.getMilliseconds()             //毫秒
    };
    for (var k in o) {
        if (new RegExp("(" + k + ")").test(fmt)) {
            if (k === "y+") {
                fmt = fmt.replace(RegExp.$1, ("" + o[k]).substr(4 - RegExp.$1.length));
            } else if (k === "S+") {
                var lens = RegExp.$1.length;
                lens = lens === 1 ? 3 : lens;
                fmt = fmt.replace(RegExp.$1, ("00" + o[k]).substr(("" + o[k]).length - 1, lens));
            } else {
                fmt = fmt.replace(RegExp.$1, (RegExp.$1.length === 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
            }
        }
    }
    return fmt;
};


function bytesToSize(_bytes) {
    if (!_bytes || _bytes === "0") {
        return 0
    }
    bytes = Number(_bytes)
    var sizes = ['bytes', 'K', 'M', 'G', 'T'];
    if (bytes === 0) return '0 Byte';
    var i = Number(Math.floor(Math.log(bytes) / Math.log(1024)));
    return Math.round(bytes / Math.pow(1024, i), 2) + ' ' + sizes[i];
}


function bpsToSize(_bps) {
    bps = Number(_bps)
    var sizes = ['bps', 'Kbps', 'Mbps', 'Gbps', 'Tbps'];
    if (bps === 0) return '0 bps';
    var i = Number(Math.floor(Math.log(bps) / Math.log(1000)));
    return Math.round(bps / Math.pow(1000, i), 2) + ' ' + sizes[i];
}

function doPost(action, formValues) {
    var form = document.createElement("FORM");
    document.body.appendChild(form);
    form.method = "post";
    form.action = action;
    form.style.display = "none";
    for (var k in formValues) {
        if (k !== 'submit') {
            var _input = document.createElement("input");
            _input.name = k;
            _input.value = formValues[k];
            form.appendChild(_input);
        }
    }
    form.submit();
}

/*
* 判断url是否合法
* */
function checkURL(url) {
    var str = url;
    //判断URL地址的正则表达式为:http(s)?://([\w-]+\.)+[\w-]+(/[\w- ./?%&=]*)?
    //下面的代码中应用了转义字符"\"输出一个字符"/"
    var Expression = /http(s)?:\/\/([\w-]+\.)+[\w-]+(\/[\w- .\/?%&=]*)?/;
    var objExp = new RegExp(Expression);
    return objExp.test(str) === true;
}

function isMobileDevice() {
    var sUserAgent = navigator.userAgent.toLowerCase();
    var bIsIpad = sUserAgent.match(/ipad/i) === 'ipad';
    var bIsIphone = sUserAgent.match(/iphone os/i) === 'iphone os';
    var bIsMidp = sUserAgent.match(/midp/i) === 'midp';
    var bIsUc7 = sUserAgent.match(/rv:1.2.3.4/i) === 'rv:1.2.3.4';
    var bIsUc = sUserAgent.match(/ucweb/i) === 'web';
    var bIsCE = sUserAgent.match(/windows ce/i) === 'windows ce';
    var bIsWM = sUserAgent.match(/windows mobile/i) === 'windows mobile';
    var bIsAndroid = sUserAgent.match(/android/i) === 'android';
    var pathname = location.pathname
    if (bIsIpad || bIsIphone || bIsMidp || bIsUc7 || bIsUc || bIsCE || bIsWM || bIsAndroid) {
        return true;
    }
    return false;
}


function concat2(a, b) {
    return a.concat(b.filter(function (v) {
        return !(a.indexOf(v) > -1)
    }))
}


/**
 * 秒转换友好的显示格式
 * 输出格式：21小时前
 * @param  {[type]} time [description]
 * @return {[type]}      [description]
 */
function second2Str(time) {
    //存储转换值
    var s;
    if ((time < 60 * 60)) {
        //超过十分钟少于1小时
        s = Math.floor(time / 60);
        return s + "分钟";
    } else if ((time < 60 * 60 * 24) && (time >= 60 * 60)) {
        //超过1小时少于24小时
        s = Math.floor(time / 60 / 60);
        return s + "小时";
    } else if (time >= 60 * 60 * 24) {
        //超过1天少于3天内
        s = Math.floor(time / 60 / 60 / 24);
        return s + "天";
    }
}


function date2str(date, y) {
    let f = function (n) {
        if (n < 10) {
            return "0" + n;
        } else {
            return n;
        }
    };
    let z = {
        y: date.getFullYear(),
        M: f(date.getMonth() + 1),
        d: f(date.getDate()),
        h: f(date.getHours()),
        m: f(date.getMinutes()),
        s: f(date.getSeconds()),
    };
    return y.replace(/(y+|M+|d+|h+|m+|s+)/g, function (v) {
        return ((v.length > 1 ? "0" : "") + eval('z.' + v.slice(-1))).slice(-(v.length > 2 ? v.length : 2))
    });

}

function toDecimal2(x) {
    let f = parseFloat(x);
    if (isNaN(f)) {
        return false;
    }
    f = Math.round(x * 100) / 100;
    let s = f.toString();
    let rs = s.indexOf('.');
    if (rs < 0) {
        rs = s.length;
        s += '.';
    }
    while (s.length <= rs + 2) {
        s += '0';
    }
    return s;
}

function fenToYuan(fen) {
    let num = fen;
    num = fen * 0.01;
    num += '';
    let reg = num.indexOf('.') > -1 ? /(\d{1,3})(?=(?:\d{3})+\.)/g : /(\d{1,3})(?=(?:\d{3})+$)/g;
    num = num.replace(reg, '$1');
    num = toDecimal2(num);
    return num
}

function yuanToFen(yuan, digit) {
    let m = 0,
        s1 = yuan.toString(),
        s2 = digit.toString();
    try {
        m += s1.split(".")[1].length
    } catch (e) {
    }
    try {
        m += s2.split(".")[1].length
    } catch (e) {
    }
    return Number(s1.replace(".", "")) * Number(s2.replace(".", "")) / Math.pow(10, m)
}

function checkTwoPointNum(inputNumber) {
    let partten = /^-?\d+\.?\d{0,2}$/;
    return partten.test(inputNumber)
}

function removeFromArray(arr, key) {
    for (var i = 0; i < arr.length; i++) {
        if (arr[i] === key) {
            arr.splice(i, 1);
        }
    }
}

function ifempty(src, defval) {
    if (src === undefined || src === null || src === "") {
        return defval;
    }
    return src;
}

function warpString(src, len) {
    if (!src) {
        return ""
    }
    if (src.length < len) {
        return src
    }
    if (src.length > len) {
        return src.substring(0, len) + "..."
    }
}


function warpRemark(obj) {
    return warpString(obj.remark, 16)
}

function GetRandom(min, max) {
    return Math.round(Math.random() * (max - min)) + min;
}