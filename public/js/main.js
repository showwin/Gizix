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

// 音声認識
(function($) {
  "use strict";

  SpeechRec.config({
    'SkyWayKey':'aaec1d4e-8dc8-4aa5-84cf-47f6aac6596a',
    'OpusWorkerUrl':'/js/libopus.worker.js',
    'SbmMode':0,
    'Recg:Nbest' : 1
  });

  if (SpeechRec.availability()) {
    console.log('Your browser supports SkyWay Speech Recognition.');
  } else {
    console.error('Your browser does not support SkyWay Speech Recognition.');
    $("#start_speech").attr('disabled', true);
    $("#start_speech").text('お使いのブラウザでは音声認識機能はご利用になれません');
  }

  $("#start_speech").click(function(){
    console.log("音声認識を開始します");
    SpeechRec.start();
    $("#result").text("");
    $("#start_speech").attr('disabled', true);
  });

  SpeechRec.on_error(function(e) {
    var x = { 'name': e.name, 'message': e.message };
    console.error('▲コールバック関数 on_error が実行されました at ' + (new Date).toLocaleString() + '\n' + JSON.stringify(x));
    $(".mic").attr("src","./img/speak_now_0.png");
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
  });
  SpeechRec.on_voice_begin(function(){
    console.warn('▲コールバック関数 on_voice_begin が実行されました( 始端が検出されました ) at ' + (new Date).toLocaleString());
  });
  SpeechRec.on_voice_too_long(function(){
    console.warn('▲コールバック関数 on_voice_too_long が実行されました( 音声認識を停止しました / 終端が検出できませんでした ) at ' + (new Date).toLocaleString());
    $(".mic").attr("src","./img/speak_now_0.png");
    $("#start_speech").attr('disabled', false);
  });
  SpeechRec.on_voice_end(function(){
    console.warn('▲コールバック関数 on_voice_end が実行されました( 終端が検出されました ) at ' + (new Date).toLocaleString());
    $(".mic").attr("src","./img/speak_now_0.png");
    $("#start_speech").attr('disabled', false);
  });
  SpeechRec.on_no_result(function(){
    console.warn('▲コールバック関数 on_no_result が実行されました( 音声認識を停止しました / 認識結果が得られませんでした ) at ' + (new Date).toLocaleString());
    $(".mic").attr("src","./img/speak_now_0.png");
    $("#start_speech").attr('disabled', false);
  });
  SpeechRec.on_mic_disabled(function(){
    console.info('▲コールバック関数 on_mic_disabled が実行されました( マイクが無効化されました )');
  });
  SpeechRec.on_mic_enabled(function(){
    console.info('▲コールバック関数 on_mic_enabled が実行されました( マイクが有効化されました )');
  });

  SpeechRec.on_result(function(result){
    console.log(result);
    $("#conversation").append("<p>"+result.candidates[0].speech+"</p>");
    //$(".mic").attr("src","./img/speak_now_0.png");
    //$("#start_speech").attr('disabled', false);
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

})(jQuery);
