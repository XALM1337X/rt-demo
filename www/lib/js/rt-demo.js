var RTDEMO = {
    test:function() {
        console.log("TEST_HIT");
    },
    submit:function() {
        var input = document.getElementById("t_fib_in").value
        var parse = parseInt(input);
        if (isNaN(parse)) {
            console.log("error");
            document.getElementById("t_fib_in").value =""
            document.getElementById("ret_block").innerHTML="Error must enter interger value.";
        } else {
            var payload = {
                    lookup: input
                };
            document.getElementById("t_fib_in").value =""
            var xhttp = new XMLHttpRequest();
            xhttp.open("POST",window.location.origin+"/fib_check", true);
            xhttp.onload = function (e) {
                if (xhttp.readyState === 4) {
                    if (xhttp.status === 200) {
                        var block = document.getElementById("ret_block");
                        block.innerHTML = xhttp.responseText;
                    } else {
                        console.error(xhttp.statusText);
                    }
                }
            };
            xhttp.onerror = function (e) {
                console.error(xhttp.statusText);
            };
            xhttp.send(JSON.stringify(payload));
        }
    }
};