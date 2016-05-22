// ---- P2P ------
var roomID = $('#myroom').val();
var connectedIds = [];
var connectedPeers = [];
var initiator = false;
var peerReady = false;

function initialize() {
  initiator = true;
  var data = JSON.stringify({"type": "initialize", "roomID": roomID, "uid": uid});
  socket.send(data);
  peerReady = true;
}

function createConnection(sigId, initFlg) {
  console.log("Connect to: " + sigId);
  connectedIds.push(sigId);
  var newPeer = new SimplePeer({ initiator: initFlg, stream: localStream });
  connectedPeers.push(newPeer);
  sendSignal(newPeer, sigId);
  startCall(newPeer, sigId);
}

function sendSignal(peer, id) {
  peer.on('signal', function (data) {
    data.to = id;
    data.uid = uid;
    var text = JSON.stringify(data);
    socket.send(text);
    console.log("Signal send to: " + id + text);
  })
  peerReady = true;
}

function receive(signal, peer) {
  peer.signal(signal);
}

function startCall(peer, id) {
  // start dictation
  SpeechRec.start();
  // start video call
  $("#videos").append("<video id='video-"+id+"' autoplay style='width: 240px; height: 180px; border: 1px solid black;'></video>")
  peer.on('stream', function (remoteStream) {
    $("#video-"+id).prop('src', window.URL.createObjectURL(remoteStream));
  })
}

// ---- socket ------
var domain = $('#mydomain').val();
var uid = $('#myid').val();
var uname = $('#myname').val();
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
  var signal = JSON.parse(event.data);
  console.log("Received.");
  console.log(signal);

  if (signal.type == 'info') {
    // confirm to register socket client
    console.log("Socket is registerd. MyID is :" + signal.uid);
  } else if (signal.type == 'initialize') {
    var pool = signal.ids;
    for (var i=0; i < pool.length; i++) {
      var sigTo = pool[i];
      // まだ繋がっていなければ接続
      if (connectedIds.indexOf(sigTo) == -1) {
        createConnection(sigTo, true)
      }
    }
  } else if (signal.type == 'conversation') {
    // 会話のログ記述
    var u = signal.uname;
    var c = signal.content;
    $("#conversation").append("<p>"+u+": "+c+"</p>");
  } else {
    // だれから送られてきたのか取得
    var sigFrom = signal.from;
    delete signal['from'];
    console.log("Call From: " + sigFrom);
    // まだ繋がっていなければ接続
    if (connectedIds.indexOf(sigFrom) == -1) {
      createConnection(sigFrom, false)
    }
    receive(signal, connectedPeers[connectedIds.indexOf(sigFrom)]);
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

// ---- Video ----
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

// 音声認識
var skywayKey = $('#myskyway').val();

SpeechRec.config({
  'SkyWayKey':skywayKey,
  'OpusWorkerUrl':'/js/libopus.worker.js',
  'NrFlag':false,
  'SbmMode':1,
  'Recg:Nbest' : 1
});

if (SpeechRec.availability()) {
  console.log('Your browser supports SkyWay Speech Recognition.');
} else {
  console.error('Your browser does not support SkyWay Speech Recognition.');
  $("#start_speech").attr('disabled', true);
  $("#start_speech").text('お使いのブラウザでは音声認識機能はご利用になれません');
}

SpeechRec.on_error(function(e) {
  var x = { 'name': e.name, 'message': e.message };
  console.error('▲コールバック関数 on_error が実行されました at ' + (new Date).toLocaleString() + '\n' + JSON.stringify(x));
  //$(".mic").attr("src","./img/speak_now_0.png");
  $("#start_speech").attr('disabled', false);
});

SpeechRec.on_config(function(config){
  console.info('▲コールバック関数 on_config が実行されました( 音声認識を設定しました ): ' + JSON.stringify(config));
});
SpeechRec.on_start(function(){
  console.info('▲コールバック関数 on_start が実行されました( 音声認識を開始しました)');
});
SpeechRec.on_stop(function(){
  console.info('▲コールバック関数 on_stop が実行されました( 音声認識を停止しました / ユーザ操作 )');
});
SpeechRec.on_ask(function() {
  console.info('▲コールバック関数 on_ask が実行されました( マイクの使用可否を確認中です )');
});
SpeechRec.on_allow(function() {
  console.info('▲コールバック関数 on_allow が実行されました( マイクの使用が許可されました )');
});
SpeechRec.on_deny(function() {
  console.info('▲コールバック関数 on_deny が実行されました( マイクの使用が拒否されました )');
});
SpeechRec.on_voiceless(function(){
  console.warn('▲コールバック関数 on_voiceless が実行されました( 音声認識を停止しました / 始端が検出できませんでした )');
  SpeechRec.start();
});
SpeechRec.on_voice_begin(function(){
  console.warn('▲コールバック関数 on_voice_begin が実行されました( 始端が検出されました ) at ' + (new Date).toLocaleString());
});
SpeechRec.on_voice_too_long(function(){
  console.warn('▲コールバック関数 on_voice_too_long が実行されました( 音声認識を停止しました / 終端が検出できませんでした ) at ' + (new Date).toLocaleString());
  //$(".mic").attr("src","./img/speak_now_0.png");
  SpeechRec.start();
});
SpeechRec.on_voice_end(function(){
  console.warn('▲コールバック関数 on_voice_end が実行されました( 終端が検出されました ) at ' + (new Date).toLocaleString());
  //$(".mic").attr("src","./img/speak_now_0.png");
});
SpeechRec.on_no_result(function(){
  console.warn('▲コールバック関数 on_no_result が実行されました( 音声認識を停止しました / 認識結果が得られませんでした ) at ' + (new Date).toLocaleString());
  //$(".mic").attr("src","./img/speak_now_0.png");
  SpeechRec.start();
});
SpeechRec.on_mic_disabled(function(){
  console.info('▲コールバック関数 on_mic_disabled が実行されました( マイクが無効化されました )');
});
SpeechRec.on_mic_enabled(function(){
  console.info('▲コールバック関数 on_mic_enabled が実行されました( マイクが有効化されました )');
});

SpeechRec.on_result(function(result){
  console.log(result);
  var content = result.candidates[0].speech;
  $("#conversation").append("<p>"+uname+": "+content+"</p>");
  //$(".mic").attr("src","./img/speak_now_0.png");

  // send others
  for (var i=0; i < connectedIds.length; i++) {
    var data = JSON.stringify({"type": "conversation", "content": content, "uname": uname, "to": connectedIds[i]});
    console.log(data);
    socket.send(data);
    console.log(socket);
    console.log("send conversation");
  }
  SpeechRec.start();
});

SpeechRec.on_proc(function(info){
  /*
  var volume = info.volume;
  if (volume > -20) {
      $(".mic").attr("src","./img/speak_now_16.png");
  } else if (volume > -50) {
      $(".mic").attr("src","./img/speak_now_10.png");
  } else if (volume > -80) {
      $(".mic").attr("src","./img/speak_now_5.png");
  } else {
      $(".mic").attr("src","./img/speak_now_1.png");
  }
  */
});
