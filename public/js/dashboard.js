// ---- socket ------
var domain = $('#mydomain').val();
var uid = $('#myid').val();
var uri = 'wss://' + domain + ':443/ws';
var socket = null;

function socketOpen() {
  if (socket == null) {
    // WebSocket の初期化
    socket = new WebSocket(uri);
    socket.onopen = onOpen;
    socket.onmessage = onMessage;
    socket.onclose = onClose;
    socket.onerror = onError;
  }
}
socketOpen();

function onOpen(event) {
  console.log('socket opened.');
  registerSocket();
}

function onMessage(event) {
  console.log("sockert receive");
  var signal = JSON.parse(event.data);
  // 通信時の自分のIDを保存する
  if (signal.type == 'info') {
    // confirm to register socket client
    console.log("Socket is registerd.");
  } else if (signal.type == 'offer') {
    // どの部屋から電話がかかってきたか表示
    var roomID = signal.roomID;
    target = $('#room'+roomID);
    title = "<h6 id='calling'>You got a call from</h6>";
    $('#joined').before(title);
    $('#calling').after("<hr>");
    $('#calling').after(target);
    targetButton = $('#room'+roomID+"_button");
    targetButton.removeClass("btn-primary")
    targetButton.addClass("btn-success")
    $('#ringtone').get(0).play();
  }
}

function onError(event) {
  //console.log("Socket Error");
}

function onClose(event) {
  console.log("Socket is Closed.");
  socket = null;
  socketReady = false;
  // reopen after 3 seconds
  setTimeout("socketOpen()", 3000);
}

function registerSocket() {
  var data = JSON.stringify({ "uid": uid, "type": "register"});
  socket.send(data);
}
