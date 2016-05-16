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
var peer = null;
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
  } else {
    // だれから送られてきたのか取得
    var sigFrom = signal.from;
    delete signal['from'];
    console.log("From: " + sigFrom);
    console.log(signal)
    receive(signal)
    if (peerReady == false) {
      answer(sigFrom);
      console.log("Answered.");
    }
    startCall();
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
  //setTimeout("open()", 3000);
}

function initialize() {
  initiator = true;
  peer = new SimplePeer({ initiator: initiator, stream: localStream })
  peer.on('signal', function (data) {
    data.to = "broadcast";
    var text = JSON.stringify(data);
    //console.log("peerSignal: " + text)
    socket.send(text);
  })
  peerReady = true;
}

function receive(signal) {
  peer = peer || new SimplePeer({ initiator: initiator, stream: localStream });
  peer.signal(signal)
}

function answer(to) {
  peer.on('signal', function (data) {
    data.to = to;
    var text = JSON.stringify(data);
    socket.send(text);
  })
  peerReady = true;
}

function startCall() {
  peer.on('stream', function (remoteStream) {
    $('#their-video').prop('src', window.URL.createObjectURL(remoteStream));
  })
  //displayTheirVideo(call);
}


function displayTheirVideo(call) {
  if (window.existingCall) {
      window.existingCall.close();
  }

  // 相手方とビデオ通信がされた時
  call.on('stream', function(stream){
      $('#their-video').prop('src', URL.createObjectURL(stream));
  });
}
