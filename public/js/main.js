var localStream = null;

window.onload = function(){
  navigator.getUserMedia  = navigator.getUserMedia ||
                            navigator.webkitGetUserMedia ||
                            navigator.mozGetUserMedia;
  window.URL = window.URL || window.webkitURL;

  var video = document.getElementById('my-video');
  navigator.getUserMedia({video: true, audio: true}, function(stream) { // for success case
    localStream = stream;
    video.src = window.URL.createObjectURL(stream);
   },
   function(err) { // for error case
    console.log(err);
   }
  );
}

// ---- socket ------
// create socket
var socketReady = false;
var port = 443;
var uri = 'wss://192.168.99.100:' + port + '/ws';
var socket = null;
var socketId = null;

function socketOpen() {
  if (socket == null) {
    // WebSocket の初期化
    socket = new WebSocket(uri);
    socket.onopen = onOpen;
    socket.onmessage = onMessage;
    socket.onclose = onClose;
    socket.onerror = onError;
    socketReady = true;
  }
}
socketOpen();

// P2P
var connectedIds = [];
var connectedPeers = [];
var initiator = false;
var peerReady = false;

// 接続イベント
function onOpen(event) {
  console.log('socket opened.');
}

// メッセージ受信イベント
function onMessage(event) {
  var signal = JSON.parse(event.data);
  console.log("Received.");
  // 通信時の自分のIDを保存する
  if (signal.type == 'config') {
    socketId = signal.id;
    console.log("MyId:" + socketId);
  } else if (signal.type == 'initialize') {
    var pool = signal.ids;
    for (var i=0; i < pool.length; i++) {
      var sigTo = pool[i];
      // まだ繋がっていなければ接続
      if (connectedIds.indexOf(sigTo) == -1) {
        createConnection(sigTo, true)
      }
    }
  } else {
    // だれから送られてきたのか取得
    var sigFrom = signal.from;
    delete signal['from'];
    console.log("Call From: " + sigFrom);
    // まだ繋がっていなければ接続
    if (connectedIds.indexOf(sigFrom) == -1) {
      createConnection(sigFrom, false)
    }
    receive(signal, connectedPeers[connectedIds.indexOf(sigFrom)])
  }
}

// エラーイベント
function onError(event) {
  //chat("エラーが発生しました。");
}

// 切断イベント
function onClose(event) {
  console.log("切断しました。");
  socket = null;
  socketReady = false;
  // 3秒後に再接続
  setTimeout("open()", 3000);
}

function initialize() {
  initiator = true;
  var data = { "to": "myself", "type": "initialize"};
  var text = JSON.stringify(data);
  socket.send(text);
  peerReady = true;
}

function createConnection(sigId, initFlg) {
  console.log("Connect to: " + sigId);
  connectedIds.push(sigId);
  var newPeer = new SimplePeer({ initiator: initFlg, stream: localStream });
  connectedPeers.push(newPeer);
  sendSignal(newPeer, sigId)
  startCall(newPeer, sigId);
}

function sendSignal(peer, id) {
  peer.on('signal', function (data) {
    data.to = id;
    var text = JSON.stringify(data);
    socket.send(text);
    console.log("Signal send to: " + id + text);
  })
  peerReady = true;
}

function receive(signal, peer) {
  peer.signal(signal)
}

function startCall(peer, id) {
  $("#their-video").append("<video id='video-"+id+"' autoplay style='width: 240px; height: 180px; border: 1px solid black;'></video>")
  peer.on('stream', function (remoteStream) {
    $("#video-"+id).prop('src', window.URL.createObjectURL(remoteStream));
  })
}
