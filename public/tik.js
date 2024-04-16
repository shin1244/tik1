$(function(){
    if (!window.WebSocket) {
        alert("No WebSocket!")
        return
    }

    connect = function() {
        ws = new WebSocket("ws://" + window.location.host + "/ws");
        ws.onmessage = function(e) {
            var dataValue = JSON.parse(e.data);
            switch (dataValue.type) {
                case "clear":
                    $("#play1").on("click");
                    $("#play1").html("플레이어1");
                    $("#play2").on("click");
                    $("#play2").html("플레이어2");
                    break
                case "player":
                    if (dataValue.data == "X"){
                        $("#play1").off("click");
                        $("#play1").html("가득참!");
                    } else if (dataValue.data == "O"){
                        $("#play2").off("click");
                        $("#play2").html("가득참!");
                    }
                    break
                case "nowBoard":
                    tiktekto(dataValue)
                    break
                case "userCount":
                    usercount(dataValue.data)
                    break;
                case "board":
                    if (dataValue.turn == "X") {
                        $("#"+dataValue.data).html("X");
                    } else {
                        $("#"+dataValue.data).html("O");
                    }
                    break
                case "O":
                    $("#winner").html("O Win")
                    disableClickEvent()
                    break
                case "X":
                    $("#winner").html("X Win")
                    disableClickEvent()
                    break
            }
        };
    }

    connect();

    $("td").click(function() {
        if ($(this).text().trim() === "") {
        var cellId = $(this).attr("id");
        if (ws.readyState === ws.OPEN) {
            ws.send(JSON.stringify({
                type: "board",
                data: cellId,
            }));
        }
    }
    });

    $("#re").click(function() {
        ableClickEvent()
        ws.send(JSON.stringify({
            type: "clear",
        }))
    });

    $("#play1").click(function() {
        ws.send(JSON.stringify({
            type: "user",
            data: "X"
            }))
        });

    $("#play2").click(function() {
        $("#play2").off("click");
        ws.send(JSON.stringify({
            type: "user",
            data: "O"
            }))
        });
    })



function tiktekto(data) {
    for (let index = 1; index < 10; index++) {
        if (data.x.indexOf(index) !== -1) {
            $("#"+String(index)).html("X")
        } else if (data.o.indexOf(index) !== -1) {
            $("#"+String(index)).html("O")
        } else {
            $("#"+String(index)).html("")
        }
    }
}

function disableClickEvent() {
    $("td").off("click");
}
function ableClickEvent() {
    $("td").on("click");
}
usercount = function(data) {
    $('#users').html('<b>유저수: '+data+'</b>');
}