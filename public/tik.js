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
                case "nowBoard":
                    tiktekto(dataValue)
                    break
                case "board":
                    if (dataValue.turn == "X") {
                        $("#"+dataValue.data).html("X");
                    } else {
                        $("#"+dataValue.data).html("O");
                    }
                    break
                case "user":
                    console.log(dataValue)
                    break
                case "O":
                    $("#winner").html("O Win")
                    break
                case "X":
                    $("#winner").html("X Win")
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
        $("#board tr").each(function() {
            $(this).find("td").text("");
        });
    });

    $("#play").click(function() {
        ws.send(JSON.stringify({
            type: "user",
            data: true
            }))
        });
    })



function tiktekto(data) {
    for (let index = 0; index < data.x.length; index++) {
        const element = data.x[index];
        $("#"+String(element)).html("X")
    }
    for (let index = 0; index < data.o.length; index++) {
        const element = data.o[index];
        $("#"+String(element)).html("O")
    }
}