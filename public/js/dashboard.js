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
  } else {
    // だれから送られてきたのか取得
    var sigFrom = signal.from;
    console.log("Call from " + sigFrom);
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
