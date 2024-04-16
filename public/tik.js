$(function(){
    if (!window.WebSocket) {
        alert("No WebSocket!")
        return
    }

    function tiktekto(data) {
        var dataData = data.data;
        var dataTurn = data.turn;

        if ($("#"+dataData).html().trim() === "") {
            if (dataTurn == "X") {
                $("#"+dataData).html("X");
                XorO = "O";
            } else {
                $("#"+dataData).html("O");
                XorO = "X";
            }
        }

        winner = checkWinner()
        if(winner != null) {
            $("#winner").html("Winner: "+winner);
        }
    }

    function checkWinner() {
        var board = [];
    
        $("#board tr").each(function() {
            var row = [];
            $(this).find("td").each(function() {
                row.push($(this).text().trim());
            });
            board.push(row);
        });
    
        for (var i = 0; i < 3; i++) {
            if (board[i][0] !== "" && board[i][0] === board[i][1] && board[i][0] === board[i][2]) {
                return board[i][0]; 
            }
            if (board[0][i] !== "" && board[0][i] === board[1][i] && board[0][i] === board[2][i]) {
                return board[0][i];
            }
        }
        if (board[0][0] !== "" && board[0][0] === board[1][1] && board[0][0] === board[2][2]) {
            return board[0][0];
        }
        if (board[0][2] !== "" && board[0][2] === board[1][1] && board[0][2] === board[2][0]) {
            return board[0][2];
        }
    
        return null;
    }

    var XorO = "X"

    connect = function() {
        ws = new WebSocket("ws://" + window.location.host + "/ws");
        ws.onmessage = function(e) {
            var dataValue = JSON.parse(e.data);
            switch (dataValue.type) {
                case "board":
                    tiktekto(dataValue)
                    break
                case "user":
                    console.log(dataValue)
                    break
            }
        };
    }

    connect();

    $("td").click(function() {
        var cellId = $(this).attr("id");
        if (ws.readyState === ws.OPEN) {
            ws.send(JSON.stringify({
                type: "board",
                data: cellId[-1],
                turn: XorO
            }));
        }
    });

    $("#re").click(function() {
        XorO = "X"
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